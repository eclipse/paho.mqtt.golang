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

package main

import (
	"fmt"
	"os"
	"time"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

var f MQTT.MessageHandler = func(client *MQTT.MqttClient, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	opts := MQTT.NewClientOptions().AddBroker("tcp://test.mosquitto.org:1883").SetClientId("trivial")
	opts.SetDefaultPublishHandler(f)

	c := MQTT.NewClient(opts)
	_, err := c.Start()
	if err != nil {
		panic(err)
	}

	filter, _ := MQTT.NewTopicFilter("/go-mqtt/sample", 0)
	if receipt, err := c.StartSubscription(nil, filter); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		<-receipt
	}

	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("this is msg #%d!", i)
		receipt := c.Publish(MQTT.QOS_ONE, "/go-mqtt/sample", []byte(text))
		<-receipt
	}

	time.Sleep(3 * time.Second)

	if receipt, err := c.EndSubscription("/go-mqtt/sample"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		<-receipt
	}

	c.Disconnect(250)
}
