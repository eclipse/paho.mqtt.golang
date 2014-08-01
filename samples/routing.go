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

/*----------------------------------------------------------------------
This sample is designed to demonstrate the ability to set individual
callbacks on a per-subscription basis. There are three handlers in use:
 brokerLoadHandler -        $SYS/broker/load/#
 brokerConnectionHandler -  $SYS/broker/connection/#
 brokerClientHandler -      $SYS/broker/clients/#
The client will receive 100 messages total from those subscriptions,
and then print the total number of messages received from each.
It may take a few moments for the sample to complete running, as it
must wait for messages to be published.
-----------------------------------------------------------------------*/

package main

import (
	"fmt"
	"os"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

var broker_load = make(chan bool)
var broker_connection = make(chan bool)
var broker_clients = make(chan bool)

var brokerLoadHandler MQTT.MessageHandler = func(client *MQTT.MqttClient, msg MQTT.Message) {
	broker_load <- true
	fmt.Printf("BrokerLoadHandler         ")
	fmt.Printf("[%s]  ", msg.Topic())
	fmt.Printf("%s\n", msg.Payload())
}

var brokerConnectionHandler MQTT.MessageHandler = func(client *MQTT.MqttClient, msg MQTT.Message) {
	broker_connection <- true
	fmt.Printf("BrokerConnectionHandler   ")
	fmt.Printf("[%s]  ", msg.Topic())
	fmt.Printf("%s\n", msg.Payload())
}

var brokerClientsHandler MQTT.MessageHandler = func(client *MQTT.MqttClient, msg MQTT.Message) {
	broker_clients <- true
	fmt.Printf("BrokerClientsHandler      ")
	fmt.Printf("[%s]  ", msg.Topic())
	fmt.Printf("%s\n", msg.Payload())
}

func main() {
	opts := MQTT.NewClientOptions().AddBroker("tcp://test.mosquitto.org:1883").SetClientId("router-sample")
	opts.SetCleanSession(true)

	c := MQTT.NewClient(opts)
	_, err := c.Start()
	if err != nil {
		panic(err)
	}

	loadFilter, _ := MQTT.NewTopicFilter("$SYS/broker/load/#", 0)
	if receipt, err := c.StartSubscription(brokerLoadHandler, loadFilter); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		<-receipt
	}

	connectionFilter, _ := MQTT.NewTopicFilter("$SYS/broker/connection/#", 0)
	if receipt, err := c.StartSubscription(brokerConnectionHandler, connectionFilter); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		<-receipt
	}

	clientsFilter, _ := MQTT.NewTopicFilter("$SYS/broker/clients/#", 0)
	if receipt, err := c.StartSubscription(brokerClientsHandler, clientsFilter); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		<-receipt
	}

	num_bload := 0
	num_bconns := 0
	num_bclients := 0

	for i := 0; i < 100; i++ {
		select {
		case <-broker_load:
			num_bload++
		case <-broker_connection:
			num_bconns++
		case <-broker_clients:
			num_bclients++
		}
	}

	fmt.Printf("Received %3d Broker Load messages\n", num_bload)
	fmt.Printf("Received %3d Broker Connection messages\n", num_bconns)
	fmt.Printf("Received %3d Broker Clients messages\n", num_bclients)

	c.Disconnect(250)
}
