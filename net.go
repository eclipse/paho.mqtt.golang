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
	. "github.com/alsm/hrotti/packets"
	"net"
	"net/url"
	"reflect"
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

// actually read incoming messages off the wire
// send Message object into ibound channel
func incoming(c *MqttClient) {

	var err error
	var cp ControlPacket

	DEBUG.Println(NET, "incoming started")

	for {
		if cp, err = ReadPacket(c.conn); err != nil {
			break
		}
		DEBUG.Println(NET, "Received Message")
		c.ibound <- cp
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
		c.errors <- err
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
		case pub := <-c.obound:
			msg := pub.p.(*PublishPacket)
			if msg.Qos != 0 && msg.MessageID == 0 {
				msg.MessageID = c.getId(pub.t)
				pub.t.(*PublishToken).messageId = msg.MessageID
			} else {
				pub.t.flowComplete()
			}
			//persist_obound(c.persist, msg)

			if c.options.writeTimeout > 0 {
				c.conn.SetWriteDeadline(time.Now().Add(c.options.writeTimeout))
			}

			if err := msg.Write(c.conn); err != nil {
				ERROR.Println(NET, "outgoing stopped with error")
				c.errors <- err
				return
			}

			if c.options.writeTimeout > 0 {
				// If we successfully wrote, we don't want the timeout to happen during an idle period
				// so we reset it to infinite.
				c.conn.SetWriteDeadline(time.Time{})
			}

			c.lastContact.update()
			DEBUG.Println(NET, "obound wrote msg, id:", msg.MessageID)
		case msg := <-c.oboundP:
			msgtype := reflect.TypeOf(msg.p)
			switch msg.p.(type) {
			case *SubscribePacket:
				msg.p.(*SubscribePacket).MessageID = c.getId(msg.t)
			case *UnsubscribePacket:
				msg.p.(*UnsubscribePacket).MessageID = c.getId(msg.t)
			}
			DEBUG.Println(NET, "obound priority msg to write, type", msgtype)
			if err := msg.p.Write(c.conn); err != nil {
				ERROR.Println(NET, "outgoing stopped with error")
				c.errors <- err
				return
			}
			c.lastContact.update()
			switch msg.p.(type) {
			case *DisconnectPacket:
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
			DEBUG.Println(NET, "logic got msg on ibound")
			//persist_ibound(c.persist, msg)
			switch msg.(type) {
			case *PingrespPacket:
				DEBUG.Println(NET, "received pingresp")
				c.pingOutstanding = false
			case *SubackPacket:
				sa := msg.(*SubackPacket)
				DEBUG.Println(NET, "received suback, id:", sa.MessageID)
				token := c.getToken(sa.MessageID).(*SubscribeToken)
				DEBUG.Println(NET, "granted qoss", sa.GrantedQoss)
				for i, qos := range sa.GrantedQoss {
					token.subResult[token.subs[i]] = qos
				}
				token.flowComplete()
				go c.freeId(sa.MessageID)
			case *UnsubackPacket:
				ua := msg.(*UnsubackPacket)
				DEBUG.Println(NET, "received unsuback, id:", ua.MessageID)
				// c.receipts.get(msg.MsgId()) <- Receipt{}
				// c.receipts.end(msg.MsgId())
				go c.freeId(ua.MessageID)
			case *PublishPacket:
				pp := msg.(*PublishPacket)
				DEBUG.Println(NET, "received publish, msgId:", pp.MessageID)
				DEBUG.Println(NET, "putting msg on onPubChan")
				switch pp.Qos {
				case 2:
					c.options.incomingPubChan <- pp
					DEBUG.Println(NET, "done putting msg on incomingPubChan")
					pr := NewControlPacket(PUBREC).(*PubrecPacket)
					pr.MessageID = pp.MessageID
					DEBUG.Println(NET, "putting pubrec msg on obound")
					c.oboundP <- &PacketAndToken{p: pr, t: nil}
					DEBUG.Println(NET, "done putting pubrec msg on obound")
				case 1:
					c.options.incomingPubChan <- pp
					DEBUG.Println(NET, "done putting msg on incomingPubChan")
					pa := NewControlPacket(PUBACK).(*PubackPacket)
					pa.MessageID = pp.MessageID
					DEBUG.Println(NET, "putting puback msg on obound")
					c.oboundP <- &PacketAndToken{p: pa, t: nil}
					DEBUG.Println(NET, "done putting puback msg on obound")
				case 0:
					select {
					case c.options.incomingPubChan <- pp:
						DEBUG.Println(NET, "done putting msg on incomingPubChan")
					case err, ok := <-c.errors:
						DEBUG.Println(NET, "error while putting msg on pubChanZero")
						// We are unblocked, but need to put the error back on so the outer
						// select can handle it appropriately.
						if ok {
							go func(errVal error, errChan chan error) {
								errChan <- errVal
							}(err, c.errors)
						}
					}
				}
			case *PubackPacket:
				pa := msg.(*PubackPacket)
				DEBUG.Println(NET, "received puback, id:", pa.MessageID)
				// c.receipts.get(msg.MsgId()) <- Receipt{}
				// c.receipts.end(msg.MsgId())
				c.getToken(pa.MessageID).flowComplete()
				c.freeId(pa.MessageID)
			case *PubrecPacket:
				prec := msg.(*PubrecPacket)
				DEBUG.Println(NET, "received pubrec, id:", prec.MessageID)
				prel := NewControlPacket(PUBREL).(*PubrelPacket)
				prel.MessageID = prec.MessageID
				select {
				case c.oboundP <- &PacketAndToken{p: prel, t: nil}:
				case <-time.After(time.Second):
				}
			case *PubrelPacket:
				pr := msg.(*PubrelPacket)
				DEBUG.Println(NET, "received pubrel, id:", pr.MessageID)
				pc := NewControlPacket(PUBCOMP).(*PubcompPacket)
				pc.MessageID = pr.MessageID
				select {
				case c.oboundP <- &PacketAndToken{p: pc, t: nil}:
				case <-time.After(time.Second):
				}
			case *PubcompPacket:
				pc := msg.(*PubcompPacket)
				DEBUG.Println(NET, "received pubcomp, id:", pc.MessageID)
				// c.receipts.get(msg.MsgId()) <- Receipt{}
				// c.receipts.end(msg.MsgId())
				c.getToken(pc.MessageID).flowComplete()
				c.freeId(pc.MessageID)
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
