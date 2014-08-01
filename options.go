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
	"crypto/tls"
	"net/url"
	"time"
)

// MessageHandler is a callback type which can be set to be
// executed upon the arrival of messages published to topics
// to which the client is subscribed.
type MessageHandler func(client *MqttClient, message Message)

// OnConnectionLost is a callback type which can be set to be
// executed upon an unintended disconnection from the MQTT broker.
// Disconnects caused by calling Disconnect or ForceDisconnect will
// not cause an OnConnectionLost callback to execute.
type OnConnectionLost func(client *MqttClient, reason error)

// ClientOptions contains configurable options for an MqttClient.
type ClientOptions struct {
	servers         []*url.URL
	clientId        string
	username        string
	password        string
	cleanSession    bool
	order           bool
	willEnabled     bool
	willTopic       string
	willPayload     []byte
	willQos         QoS
	willRetained    bool
	maxInflight     uint
	tlsConfig       *tls.Config
	keepAlive       uint
	store           Store
	msgRouter       *router
	stopRouter      chan bool
	incomingPubChan chan *Message
	onconnlost      OnConnectionLost
	mids            messageIds
	writeTimeout    time.Duration
}

// NewClientClientOptions will create a new ClientClientOptions type with some
// default values.
//   Port: 1883
//   CleanSession: True
//   Timeout: 30 (seconds)
//   Tracefile: os.Stdout
func NewClientOptions() *ClientOptions {
	o := &ClientOptions{
		servers:         nil,
		clientId:        "",
		username:        "",
		password:        "",
		cleanSession:    true,
		order:           true,
		willEnabled:     false,
		willTopic:       "",
		willPayload:     nil,
		willQos:         QOS_ZERO,
		willRetained:    false,
		maxInflight:     10,
		tlsConfig:       nil,
		store:           nil,
		keepAlive:       30,
		incomingPubChan: nil,
		onconnlost:      DefaultErrorHandler,
		mids:            messageIds{index: make(map[MId]bool)},
		writeTimeout:    0, // 0 represents timeout disabled
	}
	o.msgRouter, o.stopRouter = newRouter()
	return o
}

// AddBroker adds a broker URI to the list of brokers to be used. The format should be
// scheme://host:port
// Where "scheme" is one of "tcp", "ssl", or "ws", "host" is the ip-address (or hostname)
// and "port" is the port on which the broker is accepting connections.
func (o *ClientOptions) AddBroker(server string) *ClientOptions {
	brokerURI, _ := url.Parse(server)
	o.servers = append(o.servers, brokerURI)
	return o
}

// SetClientId will set the client id to be used by this client when
// connecting to the MQTT broker. According to the MQTT v3.1 specification,
// a client id mus be no longer than 23 characters.
func (opts *ClientOptions) SetClientId(clientid string) *ClientOptions {
	opts.clientId = clientid
	return opts
}

// SetUsername will set the username to be used by this client when connecting
// to the MQTT broker. Note: without the use of SSL/TLS, this information will
// be sent in plaintext accross the wire.
func (opts *ClientOptions) SetUsername(username string) *ClientOptions {
	opts.username = username
	return opts
}

// SetPassword will set the password to be used by this client when connecting
// to the MQTT broker. Note: without the use of SSL/TLS, this information will
// be sent in plaintext accross the wire.
func (opts *ClientOptions) SetPassword(password string) *ClientOptions {
	opts.password = password
	return opts
}

// SetCleanSession will set the "clean session" flag in the connect message
// when this client connects to an MQTT broker. By setting this flag, you are
// indicating that no messages saved by the broker for this client should be
// delivered. Any messages that were going to be sent by this client before
// diconnecting previously but didn't will not be sent upon connecting to the
// broker.
func (opts *ClientOptions) SetCleanSession(clean bool) *ClientOptions {
	opts.cleanSession = clean
	return opts
}

// SetOrderMatters will set the message routing to guarantee order within
// each QoS level. By default, this value is true. If set to false,
// this flag indicates that messages can be delivered asynchronously
// from the client to the application and possibly arrive out of order.
func (opts *ClientOptions) SetOrderMatters(order bool) *ClientOptions {
	opts.order = order
	return opts
}

// SetMaxInFlight will set a limit on the maximum number of "in-flight" messages
// going from the client to the server. This setting is currently ignored.
// func (opts *ClientOptions) SetMaxInFlight(max uint) *ClientOptions {
// 	opts.maxinflight = max
// 	return opts
// }

// SetTlsConfig will set an SSL/TLS configuration to be used when connecting
// to an MQTT broker. Please read the official Go documentation for more
// information.
func (opts *ClientOptions) SetTlsConfig(tlsConfig *tls.Config) *ClientOptions {
	opts.tlsConfig = tlsConfig
	return opts
}

// SetStore will set the implementation of the Store interface
// used to provide message persistence in cases where QoS levels
// QoS_ONE or QoS_TWO are used. If no store is provided, then the
// client will use MemoryStore by default.
func (opts *ClientOptions) SetStore(store Store) *ClientOptions {
	opts.store = store
	return opts
}

// SetKeepAlive will set the amount of time (in seconds) that the client
// should wait before sending a PING request to the broker. This will
// allow the client to know that a connection has not been lost with the
// server.
func (opts *ClientOptions) SetKeepAlive(keepAlive uint) *ClientOptions {
	opts.keepAlive = keepAlive
	return opts
}

// UnsetWill will cause any set will message to be disregarded.
func (opts *ClientOptions) UnsetWill() *ClientOptions {
	opts.willEnabled = false
	return opts
}

// SetWill accepts a string will message to be set. When the client connects,
// it will give this will message to the broker, which will then publish the
// provided payload (the will) to any clients that are subscribed to the provided
// topic.
func (opts *ClientOptions) SetWill(topic string, payload string, qos QoS, retained bool) *ClientOptions {
	opts.SetBinaryWill(topic, []byte(payload), qos, retained)
	return opts
}

// SetBinaryWill accepts a []byte will message to be set. When the client connects,
// it will give this will message to the broker, which will then publish the
// provided payload (the will) to any clients that are subscribed to the provided
// topic.
func (opts *ClientOptions) SetBinaryWill(topic string, payload []byte, qos QoS, retained bool) *ClientOptions {
	opts.willEnabled = true
	opts.willTopic = topic
	opts.willPayload = payload
	opts.willQos = qos
	opts.willRetained = retained
	return opts
}

// SetDefaultPublishHandler
func (opts *ClientOptions) SetDefaultPublishHandler(defaultHandler MessageHandler) *ClientOptions {
	opts.msgRouter.setDefaultHandler(defaultHandler)
	return opts
}

// SetOnConnectionLost will set the OnConnectionLost callback to be executed
// in the case where the client unexpectedly loses connection with the MQTT broker.
func (opts *ClientOptions) SetOnConnectionLost(onLost OnConnectionLost) *ClientOptions {
	opts.onconnlost = onLost
	return opts
}

// SetWriteTimeout puts a limit on how long a mqtt publish should block until it unblocks with a
// timeout error. A duration of 0 never times out.
func (opts *ClientOptions) SetWriteTimeout(t time.Duration) {
	opts.writeTimeout = t
}
