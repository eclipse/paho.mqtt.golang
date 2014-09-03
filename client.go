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
	Start() ([]Receipt, error)
	Disconnect(uint)
	ForceDisconnect()
	disconnect()
	Publish(string, byte, bool, interface{})
	Subscribe(string, byte, MessageHandler) error
	SubscribeMultiple(map[string]byte, MessageHandler) error
	Unsubscribe(...string) error
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
	conn         net.Conn
	bufferedConn *bufio.ReadWriter
	ibound       chan ControlPacket
	obound       chan *PublishPacket
	oboundP      chan ControlPacket
	begin        chan byte
	errors       chan error
	stop         chan struct{}
	//receipts        *receiptMap
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
	return c
}

func (c *MqttClient) IsConnected() bool {
	c.RLock()
	defer c.RUnlock()
	return c.connected
}

// Start will create a connection to the message broker
// If clean session is false, then a slice will
// be returned containing Receipts for all messages
// that were in-flight at the last disconnect.
// If clean session is true, then any existing client
// state will be removed.
func (c *MqttClient) Start() ([]Receipt, error) {
	DEBUG.Println(CLI, "Start()")

	for _, broker := range c.options.servers {
		conn, err := openConnection(broker, c.options.tlsConfig)
		if err == nil {
			c.conn = conn
			DEBUG.Println(CLI, "connected to broker")
			break
		} else {
			WARN.Println(CLI, "failed to connect to broker, trying next")
		}
	}

	if c.conn == nil {
		ERROR.Println(CLI, "Failed to connect to a broker")
		return nil, errors.New("Failed to connect to a broker")
	}
	c.bufferedConn = bufio.NewReadWriter(bufio.NewReader(c.conn), bufio.NewWriter(c.conn))

	c.persist.Open()
	//c.receipts = newReceiptMap()

	DEBUG.Println(CLI, "about to start generateMsgIds")
	c.options.mids.generateMsgIds()

	c.obound = make(chan *PublishPacket)
	c.ibound = make(chan ControlPacket)
	c.oboundP = make(chan ControlPacket)
	c.errors = make(chan error)
	c.stop = make(chan struct{})

	go outgoing(c)
	go alllogic(c)

	cm := newConnectMsgFromOptions(c.options)
	cm.ProtocolName = "MQIsdp"
	cm.ProtocolVersion = 3
	DEBUG.Println(CLI, "about to write new connect msg")
	c.oboundP <- cm

	rc := connect(c)
	if rc != CONN_ACCEPTED {
		CRITICAL.Println(CLI, "CONNACK was not CONN_ACCEPTED, but rather", ConnackReturnCodes[rc])
		// Stop all go routines except outgoing
		close(c.stop)
		c.conn.Close()
		return nil, connErrors[rc]
	}

	c.options.incomingPubChan = make(chan *PublishPacket, 100)
	c.options.msgRouter.matchAndDispatch(c.options.incomingPubChan, c.options.order, c)

	c.connected = true
	DEBUG.Println(CLI, "client is connected")

	if c.options.keepAlive != 0 {
		go keepalive(c)
	}

	// Take care of any messages in the store
	var leftovers []Receipt
	if c.options.cleanSession == false {
		//leftovers = c.resume()
	} else {
		c.persist.Reset()
	}

	// Do not start incoming until resume has completed
	go incoming(c)

	DEBUG.Println(CLI, "exit startMqttClient")
	if err := connErrors[rc]; err != nil {
		// Cleanup before returning.
		close(c.stop)
		c.conn.Close()
	}
	return leftovers, connErrors[rc]
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
	c.oboundP <- dm
	// Stop all go routines
	close(c.stop)

	DEBUG.Println(CLI, "disconnected")
	c.persist.Close()
}

// Publish will publish a message with the specified QoS
// and content to the specified topic.
// Returns a read only channel used to track
// the delivery of the message.
func (c *MqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) {
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
	c.obound <- pub
}

// Start a new subscription. Provide a MessageHandler to be executed when
// a message is published on the topic provided.
func (c *MqttClient) Subscribe(topic string, qos byte, callback MessageHandler) error {
	var err error
	DEBUG.Println(CLI, "enter Subscribe")
	if !c.IsConnected() {
		return ErrNotConnected
	}
	s := NewControlPacket(SUBSCRIBE).(*SubscribePacket)
	DEBUG.Println(s.String())
	if err = validateTopicAndQos(topic, qos); err != nil {
		return err
	}
	s.Topics = append(s.Topics, topic)
	s.Qoss = append(s.Qoss, qos)

	if callback != nil {
		c.options.msgRouter.addRoute(topic, callback)
	}

	c.oboundP <- s
	DEBUG.Println(CLI, "exit Subscribe")
	return nil
}

// Start a new subscription for multiple topics. Provide a MessageHandler to
// be executed when a message is published on one of the topics provided.
func (c *MqttClient) SubscribeMultiple(filters map[string]byte, callback MessageHandler) error {
	var err error
	DEBUG.Println(CLI, "enter SubscribeMultiple")
	if !c.IsConnected() {
		return ErrNotConnected
	}
	s := NewControlPacket(SUBSCRIBE).(*SubscribePacket)
	if s.Topics, s.Qoss, err = validateSubscribeMap(filters); err != nil {
		return err
	}

	if callback != nil {
		for topic, _ := range filters {
			c.options.msgRouter.addRoute(topic, callback)
		}
	}

	//r := make(chan Receipt, 1)

	c.oboundP <- s

	DEBUG.Println(CLI, "exit SubscribeMultiple")
	return nil
}

// Unsubscribe will end the subscription from each of the topics provided.
// Messages published to those topics from other clients will no longer be
// received.
func (c *MqttClient) Unsubscribe(topics ...string) error {
	DEBUG.Println(CLI, "enter Unsubscribe")
	if !c.IsConnected() {
		return ErrNotConnected
	}
	u := NewControlPacket(UNSUBSCRIBE).(*UnsubscribePacket)
	u.Topics = make([]string, len(topics))
	copy(u.Topics, topics)

	c.oboundP <- u
	for _, topic := range topics {
		c.options.msgRouter.deleteRoute(topic)
	}

	DEBUG.Println(CLI, "exit Unsubscribe")
	return nil
}
