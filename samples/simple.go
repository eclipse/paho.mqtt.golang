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
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"os"
	"time"
)

var f MQTT.MessageHandler = func(msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	opts := MQTT.NewClientOptions().SetBroker("tcp://test.mosquitto.org:1883").SetClientId("trivial")
	opts.SetTraceLevel(MQTT.Off)
	opts.SetDefaultPublishHandler(f)

	c := MQTT.NewClient(opts)
	_, err := c.Start()
	if err != nil {
		panic(err)
	}

	if receipt, err := c.StartSubscription(nil, "/go-mqtt/sample", MQTT.QOS_ZERO); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		<-receipt
	}

	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("this is msg #%d!", i)
		receipt := c.Publish(MQTT.QOS_ONE, "/go-mqtt/sample", text)
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
