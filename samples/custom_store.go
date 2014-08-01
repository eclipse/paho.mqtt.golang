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

// This demonstrates how to implement your own Store interface and provide
// it to the go-mqtt client.

package main

import (
	"fmt"
	"time"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

// This NoOpStore type implements the go-mqtt/Store interface, which
// allows it to be used by the go-mqtt client library. However, it is
// highly recommended that you do not use this NoOpStore in production,
// because it will NOT provide any sort of guaruntee of message delivery.
type NoOpStore struct {
	// Contain nothing
}

func (store *NoOpStore) Open() {
	// Do nothing
}

func (store *NoOpStore) Put(string, *MQTT.Message) {
	// Do nothing
}

func (store *NoOpStore) Get(string) *MQTT.Message {
	// Do nothing
	return nil
}

func (store *NoOpStore) Del(string) {
	// Do nothing
}

func (store *NoOpStore) All() []string {
	return nil
}

func (store *NoOpStore) Close() {
	// Do Nothing
}

func (store *NoOpStore) Reset() {
	// Do Nothing
}

func (store *NoOpStore) SetTracer(tracer *MQTT.Tracer) {
	// Do Nothing
}

func main() {
	myNoOpStore := &NoOpStore{}

	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://test.mosquitto.org:1883")
	opts.SetClientId("custom-store")
	opts.SetStore(myNoOpStore)

	var callback MQTT.MessageHandler = func(client *MQTT.MqttClient, msg MQTT.Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	}

	c := MQTT.NewClient(opts)
	_, err := c.Start()
	if err != nil {
		panic(err)
	}

	filter, _ := MQTT.NewTopicFilter("/go-mqtt/sample", 0)
	c.StartSubscription(callback, filter)

	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("this is msg #%d!", i)
		c.Publish(MQTT.QOS_ONE, "/go-mqtt/sample", []byte(text))
	}

	for i := 1; i < 5; i++ {
		time.Sleep(1 * time.Second)
	}

	c.Disconnect(250)
}
