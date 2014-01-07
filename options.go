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
	"os"
)

// MessageHandler is a callback type which can be set to be
// executed upon the arrival of messages published to topics
// to which the client is subscribed.
type MessageHandler func(message Message)

// OnConnectionLost is a callback type which can be set to be
// executed upon an unintended disconnection from the MQTT broker.
// Disconnects caused by calling Disconnect or ForceDisconnect will
// not cause an OnConnectionLost callback to execute.
type OnConnectionLost func(reason error)

// ClientOptions contains configurable options for an MqttClient.
type ClientOptions struct {
	server        *url.URL
	server2       *url.URL
	clientId      string
	username      string
	password      string
	cleanses      bool
	order         bool
	will_enabled  bool
	will_topic    string
	will_payload  []byte
	will_qos      QoS
	will_retained bool
	maxinflight   uint
	tlsconfig     *tls.Config
	timeout       uint
	store         Store
	tracefile     *os.File
	tracelevel    tracelevel
	msgRouter     *router
	stopRouter    chan bool
	pubChanZero   chan *Message
	pubChanOne    chan *Message
	pubChanTwo    chan *Message
	onconnlost    OnConnectionLost
	mids          messageIds
}

// NewClientClientOptions will create a new ClientClientOptions type with some
// default values.
//   Port: 1883
//   CleanSession: True
//   Timeout: 30 (seconds)
//   Tracefile: os.Stdout
func NewClientOptions() *ClientOptions {
	o := &ClientOptions{
		server:        nil,
		server2:       nil,
		clientId:      "",
		username:      "",
		password:      "",
		cleanses:      true,
		order:         true,
		will_enabled:  false,
		will_topic:    "",
		will_payload:  nil,
		will_qos:      QOS_ZERO,
		will_retained: false,
		maxinflight:   10,
		tlsconfig:     nil,
		store:         nil,
		timeout:       30,
		tracefile:     os.Stdout,
		tracelevel:    Verbose,
		pubChanZero:   nil,
		pubChanOne:    nil,
		pubChanTwo:    nil,
		onconnlost:    DefaultErrorHandler,
		mids:          messageIds{index: make(map[MId]bool)},
	}
	o.msgRouter, o.stopRouter = newRouter()
	return o
}

// SetBroker will allow you to set the URI for your broker. The format should be
// scheme://host:port
// Where "scheme" is one of "tcp", "ssl", or "ws", "host" is the ip-address (or hostname)
// and "port" is the port on which the broker is accepting connections.
// For example, one could connect to tcp://test.mosquitto.org:1883
func (opts *ClientOptions) SetBroker(server string) *ClientOptions {
	opts.server, _ = url.Parse(server)
	return opts
}

// SetStandbyBroker will allow you to set a second URI to which the client will attempt
// to connect in the event of a connection failure. This is for use only in cases where
// two brokers are configured as a highly available pair. (For example, two IBM MessageSight
// appliances configured in High Availability mode).
func (opts *ClientOptions) SetStandbyBroker(server string) *ClientOptions {
	opts.server2, _ = url.Parse(server)
	return opts
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
	opts.cleanses = clean
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
func (opts *ClientOptions) SetTlsConfig(tlsconfig *tls.Config) *ClientOptions {
	opts.tlsconfig = tlsconfig
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

// SetTimeout will set the amount of time (in seconds) that the client
// should wait before sending a PING request to the broker. This will
// allow the client to know that a connection has not been lost with the
// server.
func (opts *ClientOptions) SetTimeout(timeout uint) *ClientOptions {
	opts.timeout = timeout
	return opts
}

// UnsetWill will cause any set will message to be disregarded.
func (opts *ClientOptions) UnsetWill() *ClientOptions {
	opts.will_enabled = false
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
	opts.will_enabled = true
	opts.will_topic = topic
	opts.will_payload = payload
	opts.will_qos = qos
	opts.will_retained = retained
	return opts
}

// SetTracefile will set the output for any trace statements that are generated
// by the client. By default, trace statements will be directed to os.Stdout.
func (opts *ClientOptions) SetTracefile(tracefile *os.File) *ClientOptions {
	opts.tracefile = tracefile
	return opts
}

// SetTraceLevel will set the trace level (verbosity) of the client.
// Options are:
//   Off
//   Critical
//   Warn
//   Verbose
func (opts *ClientOptions) SetTraceLevel(level tracelevel) *ClientOptions {
	opts.tracelevel = level
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
