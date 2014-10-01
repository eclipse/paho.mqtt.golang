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

// Package mqtt provides an MQTT v3.1 client library.
package mqtt

import (
	"bufio"
	"errors"
	. "github.com/alsm/hrotti/packets"
	"net"
	"sync"
	"time"
)

type Client interface {
	IsConnected() bool
	Connect() Token
	Disconnect(uint)
	ForceDisconnect()
	disconnect()
	Publish(string, byte, bool, interface{}) Token
	Subscribe(string, byte, MessageHandler) Token
	SubscribeMultiple(map[string]byte, MessageHandler) Token
	Unsubscribe(...string) Token
}

// MqttClient is a lightweight MQTT v3.1 Client for communicating
// with an MQTT server using non-blocking methods that allow work
// to be done in the background.

// An application may connect to an MQTT server using:
//   A plain TCP socket
//   A secure SSL/TLS socket
//   A websocket

// To enable ensured message delivery at Quality of Service (QoS) levels
// described in the MQTT spec, a message persistence mechanism must be
// used. This is done by providing a type which implements the Store
// interface. For convenience, FileStore and MemoryStore are provided
// implementations that should be sufficient for most use cases. More
// information can be found in their respective documentation.

// Numerous connection options may be specified by configuring a
// and then supplying a ClientOptions type.
type MqttClient struct {
	sync.RWMutex
	messageIds
	callbacks       Callbacks
	conn            net.Conn
	bufferedConn    *bufio.ReadWriter
	ibound          chan ControlPacket
	obound          chan *PacketAndToken
	oboundP         chan *PacketAndToken
	errors          chan error
	stop            chan struct{}
	persist         Store
	options         ClientOptions
	lastContact     lastcontact
	pingOutstanding bool
	connected       bool
}

// NewClient will create an MQTT v3.1 client with all of the options specified
// in the provided ClientOptions. The client must have the Start method called
// on it before it may be used. This is to make sure resources (such as a net
// connection) are created before the application is actually ready.
func NewClient(ops *ClientOptions) *MqttClient {
	c := &MqttClient{}
	c.options = *ops

	if c.options.store == nil {
		c.options.store = NewMemoryStore()
	}
	c.persist = c.options.store
	c.connected = false
	c.messageIds = messageIds{index: make(map[uint16]Token)}
	return c
}

func (c *MqttClient) IsConnected() bool {
	c.RLock()
	defer c.RUnlock()
	return c.connected
}

// Connect will create a connection to the message broker
// If clean session is false, then a slice will
// be returned containing Receipts for all messages
// that were in-flight at the last disconnect.
// If clean session is true, then any existing client
// state will be removed.
func (c *MqttClient) Connect() Token {
	var err error
	t := newToken(CONNECT).(*ConnectToken)
	DEBUG.Println(CLI, "Connect()")

	go func() {
		var rc byte
		cm := newConnectMsgFromOptions(c.options)

		for _, broker := range c.options.servers {
		CONN:
			DEBUG.Println(CLI, "about to write new connect msg")
			c.conn, err = openConnection(broker, c.options.tlsConfig)
			if err == nil {
				DEBUG.Println(CLI, "socket connected to broker")
				switch c.options.protocolVersion {
				case 3:
					DEBUG.Println(CLI, "Using MQTT 3.1 protocol")
					cm.ProtocolName = "MQIsdp"
					cm.ProtocolVersion = 3
				default:
					DEBUG.Println(CLI, "Using MQTT 3.1.1 protocol")
					c.options.protocolVersion = 4
					cm.ProtocolName = "MQTT"
					cm.ProtocolVersion = 4
				}
				cm.Write(c.conn)

				rc = c.connect()
				if rc != CONN_ACCEPTED {
					c.conn.Close()
					c.conn = nil
					//if the protocol version was explicitly set don't do any fallback
					if c.options.protocolVersionExplicit {
						ERROR.Println(CLI, "Connecting to", broker, "CONNACK was not CONN_ACCEPTED, but rather", ConnackReturnCodes[rc])
						continue
					}
					if c.options.protocolVersion == 4 {
						DEBUG.Println(CLI, "Trying reconnect using MQTT 3.1 protocol")
						c.options.protocolVersion = 3
						goto CONN
					}
				}
				break
			} else {
				ERROR.Println(CLI, err.Error())
				WARN.Println(CLI, "failed to connect to broker, trying next")
				rc = CONN_NETWORK_ERROR
			}
		}

		if c.conn == nil {
			ERROR.Println(CLI, "Failed to connect to a broker")
			t.returnCode = rc
			if rc != CONN_NETWORK_ERROR {
				t.err = connErrors[rc]
			} else {
				t.err = errors.New(connErrors[rc].Error() + " : " + err.Error())
			}
			t.flowComplete()
			return
		}
		c.bufferedConn = bufio.NewReadWriter(bufio.NewReader(c.conn), bufio.NewWriter(c.conn))

		c.persist.Open()

		c.obound = make(chan *PacketAndToken)
		c.oboundP = make(chan *PacketAndToken)
		c.ibound = make(chan ControlPacket)
		c.errors = make(chan error)
		c.stop = make(chan struct{})

		go outgoing(c)
		go alllogic(c)

		c.options.incomingPubChan = make(chan *PublishPacket, 100)
		c.options.msgRouter.matchAndDispatch(c.options.incomingPubChan, c.options.order, c)

		c.connected = true
		DEBUG.Println(CLI, "client is connected")

		if c.options.keepAlive != 0 {
			go keepalive(c)
		}

		// Take care of any messages in the store
		//var leftovers []Receipt
		if c.options.cleanSession == false {
			//leftovers = c.resume()
		} else {
			c.persist.Reset()
		}

		// Do not start incoming until resume has completed
		go incoming(c)

		DEBUG.Println(CLI, "exit startMqttClient")
		t.flowComplete()
	}()
	return t
}

// This function is only used for receiving a connack
// when the connection is first started.
// This prevents receiving incoming data while resume
// is in progress if clean session is false.
func (c *MqttClient) connect() byte {
	DEBUG.Println(NET, "connect started")

	ca, err := ReadPacket(c.conn)
	if err != nil {
		ERROR.Println(NET, "connect got error", err)
		//c.errors <- err
		return CONN_NETWORK_ERROR
	}
	msg := ca.(*ConnackPacket)

	if msg == nil || msg.FixedHeader.MessageType != CONNACK {
		ERROR.Println(NET, "received msg that was nil or not CONNACK")
	} else {
		DEBUG.Println(NET, "received connack")
	}
	return msg.ReturnCode
}

// Disconnect will end the connection with the server, but not before waiting
// the specified number of milliseconds to wait for existing work to be
// completed.
func (c *MqttClient) Disconnect(quiesce uint) {
	if !c.IsConnected() {
		WARN.Println(CLI, "already disconnected")
		return
	}
	DEBUG.Println(CLI, "disconnecting")
	c.connected = false

	// wait for work to finish, or quiesce time consumed
	end := time.After(time.Duration(quiesce) * time.Millisecond)

	// for now we just wait for the time specified and hope the work is done
	select {
	case <-end:
		DEBUG.Println(CLI, "quiesce expired, forcing disconnect")
		// case <- other:
		// 	DEBUG.Println(CLI, "finished processing work, graceful disconnect")
	}
	c.disconnect()
}

// ForceDisconnect will end the connection with the mqtt broker immediately.
func (c *MqttClient) ForceDisconnect() {
	if !c.IsConnected() {
		WARN.Println(CLI, "already disconnected")
		return
	}
	DEBUG.Println(CLI, "forcefully disconnecting")
	c.disconnect()
}

func (c *MqttClient) disconnect() {
	c.connected = false
	dm := NewControlPacket(DISCONNECT).(*DisconnectPacket)

	// Send disconnect message and stop outgoing
	c.oboundP <- &PacketAndToken{p: dm, t: nil}
	// Stop all go routines
	close(c.stop)

	DEBUG.Println(CLI, "disconnected")
	c.persist.Close()
}

// Publish will publish a message with the specified QoS
// and content to the specified topic.
// Returns a read only channel used to track
// the delivery of the message.
func (c *MqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) Token {
	pub := NewControlPacket(PUBLISH).(*PublishPacket)
	pub.Qos = qos
	pub.TopicName = topic
	pub.Retain = retained
	switch payload.(type) {
	case string:
		pub.Payload = []byte(payload.(string))
	case []byte:
		pub.Payload = payload.([]byte)
	default:
	}

	DEBUG.Println(CLI, "sending publish message, topic:", topic)
	token := newToken(PUBLISH)
	c.obound <- &PacketAndToken{p: pub, t: token}
	return token
}

// Start a new subscription. Provide a MessageHandler to be executed when
// a message is published on the topic provided.
func (c *MqttClient) Subscribe(topic string, qos byte, callback MessageHandler) Token {
	token := newToken(SUBSCRIBE).(*SubscribeToken)
	DEBUG.Println(CLI, "enter Subscribe")
	if !c.IsConnected() {
		token.err = ErrNotConnected
		return token
	}
	sub := NewControlPacket(SUBSCRIBE).(*SubscribePacket)
	if err := validateTopicAndQos(topic, qos); err != nil {
		token.err = err
		return token
	}
	sub.Topics = append(sub.Topics, topic)
	sub.Qoss = append(sub.Qoss, qos)
	DEBUG.Println(sub.String())

	if callback != nil {
		c.options.msgRouter.addRoute(topic, callback)
	}

	token.subs = append(token.subs, topic)
	c.oboundP <- &PacketAndToken{p: sub, t: token}
	DEBUG.Println(CLI, "exit Subscribe")
	return token
}

// Start a new subscription for multiple topics. Provide a MessageHandler to
// be executed when a message is published on one of the topics provided.
func (c *MqttClient) SubscribeMultiple(filters map[string]byte, callback MessageHandler) Token {
	var err error
	token := newToken(SUBSCRIBE).(*SubscribeToken)
	DEBUG.Println(CLI, "enter SubscribeMultiple")
	if !c.IsConnected() {
		token.err = ErrNotConnected
		return token
	}
	sub := NewControlPacket(SUBSCRIBE).(*SubscribePacket)
	if sub.Topics, sub.Qoss, err = validateSubscribeMap(filters); err != nil {
		token.err = err
		return token
	}

	if callback != nil {
		for topic, _ := range filters {
			c.options.msgRouter.addRoute(topic, callback)
		}
	}
	token.subs = make([]string, len(sub.Topics))
	copy(token.subs, sub.Topics)
	c.oboundP <- &PacketAndToken{p: sub, t: token}
	DEBUG.Println(CLI, "exit SubscribeMultiple")
	return token
}

// Unsubscribe will end the subscription from each of the topics provided.
// Messages published to those topics from other clients will no longer be
// received.
func (c *MqttClient) Unsubscribe(topics ...string) Token {
	token := newToken(UNSUBSCRIBE).(*UnsubscribeToken)
	DEBUG.Println(CLI, "enter Unsubscribe")
	if !c.IsConnected() {
		token.err = ErrNotConnected
		return token
	}
	unsub := NewControlPacket(UNSUBSCRIBE).(*UnsubscribePacket)
	unsub.Topics = make([]string, len(topics))
	copy(unsub.Topics, topics)

	c.oboundP <- &PacketAndToken{p: unsub, t: token}
	for _, topic := range topics {
		c.options.msgRouter.deleteRoute(topic)
	}

	DEBUG.Println(CLI, "exit Unsubscribe")
	return token
}
