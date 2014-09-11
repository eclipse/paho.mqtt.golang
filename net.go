/*
 * Copyright (c) 2013 IBM Corp.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v1.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-v10.html
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 */

package mqtt

import (
	"code.google.com/p/go.net/websocket"
	"crypto/tls"
	"io"
	"net"
	"net/url"
	"time"
)

func openConnection(uri *url.URL, tlsc *tls.Config) (conn net.Conn, err error) {
	switch uri.Scheme {
	case "ws":
		conn, err = websocket.Dial(uri.String(), "mqtt", "ws://localhost")
		if err != nil {
			return
		}
		conn.(*websocket.Conn).PayloadType = websocket.BinaryFrame
	case "tcp":
		conn, err = net.Dial("tcp", uri.Host)
	case "ssl":
		fallthrough
	case "tls":
		fallthrough
	case "tcps":
		conn, err = tls.Dial("tcp", uri.Host, tlsc)
	}
	return
}

// This function is only used for receiving a connack
// when the connection is first started.
// This prevents receiving incoming data while resume
// is in progress if clean session is false.
func connect(c *MqttClient) (rc ConnRC) {
	rc = CONN_FAILURE
	DEBUG.Println(NET, "connect started")

	//connack is always 4 bytes
	ca := make([]byte, 4)
	_, err := io.ReadFull(c.bufferedConn, ca)
	if err != nil {
		ERROR.Println(NET, "connect got error")
		select {
		case c.errors <- err:
		default:
			// c.errors is a buffer of one, so there must already be an error closing this connection.
		}
		return
	}
	msg := decode(ca)

	if msg == nil || msg.msgType() != CONNACK {
		ERROR.Println(NET, "received msg that was nil or not CONNACK")
		return
	}

	DEBUG.Println(NET, "received connack")
	return msg.connRC()
}

// actually read incoming messages off the wire
// send Message object into ibound channel
func incoming(c *MqttClient) {

	var err error

	DEBUG.Println(NET, "incoming started")

	for {
		var rerr error
		var msg *Message
		msgType := make([]byte, 1)
		DEBUG.Println(NET, "incoming waiting for network data")
		msgType[0], rerr = c.bufferedConn.ReadByte()
		if rerr != nil {
			err = rerr
			break
		}
		bytes, remLen := decodeRemlenFromNetwork(c.bufferedConn)
		fixedHeader := make([]byte, len(bytes)+1)
		copy(fixedHeader, append(msgType, bytes...))
		if remLen > 0 {
			data := make([]byte, remLen)
			DEBUG.Println(NET, remLen, "more incoming bytes to read")
			_, rerr = io.ReadFull(c.bufferedConn, data)
			if rerr != nil {
				err = rerr
				break
			}
			DEBUG.Println(NET, "data:", data)
			msg = decode(append(fixedHeader, data...))
		} else {
			msg = decode(fixedHeader)
		}
		if msg != nil {
			DEBUG.Println(NET, "incoming received inbound message, type", msg.msgType())
			c.ibound <- msg
		} else {
			CRITICAL.Println(NET, "incoming msg was nil")
		}
	}
	// We received an error on read.
	// If disconnect is in progress, swallow error and return
	select {
	case <-c.stop:
		DEBUG.Println(NET, "incoming stopped")
		return
		// Not trying to disconnect, send the error to the errors channel
	default:
		ERROR.Println(NET, "incoming stopped with error")
		select {
		case c.errors <- err:
		default:
			// c.errors is a buffer of one, so there must already be an error closing this connection.
		}
		return
	}
}

// receive a Message object on obound, and then
// actually send outgoing message to the wire
func outgoing(c *MqttClient) {

	DEBUG.Println(NET, "outgoing started")

	for {
		DEBUG.Println(NET, "outgoing waiting for an outbound message")
		select {
		case out := <-c.obound:
			msg := out.m
			msgtype := msg.msgType()
			DEBUG.Println(NET, "obound got msg to write, type:", msgtype)
			if msg.QoS() != QOS_ZERO && msg.MsgId() == 0 {
				msg.setMsgId(c.options.mids.getId())
			}
			if out.r != nil {
				c.receipts.put(msg.MsgId(), out.r)
			}
			msg.setTime()
			persist_obound(c.persist, msg)

			if c.options.writeTimeout > 0 {
				c.conn.SetWriteDeadline(time.Now().Add(c.options.writeTimeout))
			}

			if _, err := c.conn.Write(msg.Bytes()); err != nil {
				ERROR.Println(NET, "outgoing stopped with error")
				select {
				case c.errors <- err:
				default:
					// c.errors is a buffer of one, so there must already be an error closing this connection.
				}
				return
			}

			if c.options.writeTimeout > 0 {
				// If we successfully wrote, we don't want the timeout to happen during an idle period
				// so we reset it to infinite.
				c.conn.SetWriteDeadline(time.Time{})
			}

			if (msg.QoS() == QOS_ZERO) &&
				(msgtype == PUBLISH || msgtype == SUBSCRIBE || msgtype == UNSUBSCRIBE) {
				c.receipts.get(msg.MsgId()) <- Receipt{}
				c.receipts.end(msg.MsgId())
			}
			c.lastContact.update()
			DEBUG.Println(NET, "obound wrote msg, id:", msg.MsgId())
		case msg := <-c.oboundP:
			msgtype := msg.msgType()
			DEBUG.Println(NET, "obound priority msg to write, type", msgtype)
			_, err := c.conn.Write(msg.Bytes())
			if err != nil {
				ERROR.Println(NET, "outgoing stopped with error")
				select {
				case c.errors <- err:
				default:
					// c.errors is a buffer of one, so there must already be an error closing this connection.
				}
				return
			}
			c.lastContact.update()
			if msgtype == DISCONNECT {
				DEBUG.Println(NET, "outbound wrote disconnect, now closing connection")
				c.conn.Close()
				return
			}
		}
	}
}

// receive Message objects on ibound
// store messages if necessary
// send replies on obound
// delete messages from store if necessary
func alllogic(c *MqttClient) {

	DEBUG.Println(NET, "logic started")

	for {
		DEBUG.Println(NET, "logic waiting for msg on ibound")

		select {
		case msg := <-c.ibound:
			DEBUG.Println(NET, "logic got msg on ibound, type", msg.msgType())
			persist_ibound(c.persist, msg)
			switch msg.msgType() {
			case PINGRESP:
				DEBUG.Println(NET, "received pingresp")
				c.pingOutstanding = false
			case SUBACK:
				DEBUG.Println(NET, "received suback, id:", msg.MsgId())
				c.receipts.get(msg.MsgId()) <- Receipt{}
				c.receipts.end(msg.MsgId())
				go c.options.mids.freeId(msg.MsgId())
			case UNSUBACK:
				DEBUG.Println(NET, "received unsuback, id:", msg.MsgId())
				c.receipts.get(msg.MsgId()) <- Receipt{}
				c.receipts.end(msg.MsgId())
				go c.options.mids.freeId(msg.MsgId())
			case PUBLISH:
				DEBUG.Println(NET, "received publish, msgId:", msg.MsgId())
				DEBUG.Println(NET, "putting msg on onPubChan")
				switch msg.QoS() {
				case QOS_TWO:
					c.options.incomingPubChan <- msg
					DEBUG.Println(NET, "done putting msg on incomingPubChan")
					pubrecMsg := newPubRecMsg()
					pubrecMsg.setMsgId(msg.MsgId())
					DEBUG.Println(NET, "putting pubrec msg on obound")
					c.obound <- sendable{pubrecMsg, nil}
					DEBUG.Println(NET, "done putting pubrec msg on obound")
				case QOS_ONE:
					c.options.incomingPubChan <- msg
					DEBUG.Println(NET, "done putting msg on incomingPubChan")
					pubackMsg := newPubAckMsg()
					pubackMsg.setMsgId(msg.MsgId())
					DEBUG.Println(NET, "putting puback msg on obound")
					c.obound <- sendable{pubackMsg, nil}
					DEBUG.Println(NET, "done putting puback msg on obound")
				case QOS_ZERO:
					select {
					case c.options.incomingPubChan <- msg:
						DEBUG.Println(NET, "done putting msg on incomingPubChan")
					case err, ok := <-c.errors:
						DEBUG.Println(NET, "error while putting msg on pubChanZero")
						// We are unblocked, but need to put the error back on so the outer
						// select can handle it appropriately.
						if ok {
							go func(errVal error, errChan chan error) {
								select {
								case errChan <- errVal:
								default:
									// c.errors is a buffer of one, so there must already be an error closing this connection.
								}
							}(err, c.errors)
						}
					}
				}
			case PUBACK:
				DEBUG.Println(NET, "received puback, id:", msg.MsgId())
				c.receipts.get(msg.MsgId()) <- Receipt{}
				c.receipts.end(msg.MsgId())
				go c.options.mids.freeId(msg.MsgId())
			case PUBREC:
				DEBUG.Println(NET, "received pubrec, id:", msg.MsgId())
				id := msg.MsgId()
				pubrelMsg := newPubRelMsg()
				pubrelMsg.setMsgId(id)
				select {
				case c.obound <- sendable{pubrelMsg, nil}:
				case <-time.After(time.Second):
				}
			case PUBREL:
				DEBUG.Println(NET, "received pubrel, id:", msg.MsgId())
				pubcompMsg := newPubCompMsg()
				pubcompMsg.setMsgId(msg.MsgId())
				select {
				case c.obound <- sendable{pubcompMsg, nil}:
				case <-time.After(time.Second):
				}
			case PUBCOMP:
				DEBUG.Println(NET, "received pubcomp, id:", msg.MsgId())
				c.receipts.get(msg.MsgId()) <- Receipt{}
				c.receipts.end(msg.MsgId())
				go c.options.mids.freeId(msg.MsgId())
			}
		case <-c.stop:
			WARN.Println(NET, "logic stopped")
			return
		case err := <-c.errors:
			c.connected = false
			ERROR.Println(NET, "logic got error")
			// clean up go routines
			// incoming most likely stopped if outgoing stopped,
			// but let it know to stop anyways.
			close(c.options.stopRouter)
			close(c.stop)
			c.conn.Close()

			// Call onConnectionLost or default error handler
			go c.options.onconnlost(c, err)
			return
		}
	}
}
