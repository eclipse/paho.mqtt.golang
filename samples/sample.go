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
	"flag"
	"fmt"
	"os"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

/*
Options:
 [-help]                      Display help
 [-a pub|sub]                 Action pub (publish) or sub (subscribe)
 [-m <message>]               Payload to send
 [-n <number>]                Number of messages to send or receive
 [-q 0|1|2]                   Quality of Service
 [-clean]                     CleanSession (true if -clean is present)
 [-id <clientid>]             CliendID
 [-user <user>]               User
 [-password <password>]       Password
 [-broker <uri>]              Broker URI
 [-topic <topic>]             Topic
 [-store <path>]              Store Directory

*/

func main() {
	topic := flag.String("topic", "", "The topic name to/from which to publish/subscribe")
	broker := flag.String("broker", "", "The broker URI. ex: tcp://10.10.1.1:1883")
	password := flag.String("password", "", "The password (optional)")
	user := flag.String("user", "", "The User (optional)")
	id := flag.String("id", "", "The ClientID (optional)")
	cleansess := flag.Bool("clean", false, "Set Clean Session (default false)")
	qos := flag.Int("qos", 0, "The Quality of Service 0,1,2 (default 0)")
	num := flag.Int("num", 1, "The number of messages to publish or subscribe (default 1)")
	payload := flag.String("message", "", "The message text to publish (default empty)")
	action := flag.String("action", "", "Action publish or subscribe (required)")
	store := flag.String("store", ":memory:", "The Store Directory (default use memory store)")
	flag.Parse()

	if *broker == "" {
		fmt.Println("Invalid setting for -broker")
		return
	}

	if *action != "pub" && *action != "sub" {
		fmt.Println("Invalid setting for -action, must be pub or sub")
		return
	}

	if *topic == "" {
		fmt.Println("Invalid setting for -topic, must not be empty")
		return
	}

	fmt.Printf("Sample Info:\n")
	fmt.Printf("\taction:    %s\n", *action)
	fmt.Printf("\tbroker:    %s\n", *broker)
	fmt.Printf("\tclientid:  %s\n", *id)
	fmt.Printf("\tuser:      %s\n", *user)
	fmt.Printf("\tpassword:  %s\n", *password)
	fmt.Printf("\ttopic:     %s\n", *topic)
	fmt.Printf("\tmessage:   %s\n", *payload)
	fmt.Printf("\tqos:       %d\n", *qos)
	fmt.Printf("\tcleansess: %v\n", *cleansess)
	fmt.Printf("\tnum:       %d\n", *num)
	fmt.Printf("\tstore:     %s\n", *store)

	opts := MQTT.NewClientOptions()
	opts.AddBroker(*broker)
	opts.SetClientId(*id)
	opts.SetUsername(*user)
	opts.SetPassword(*password)
	opts.SetCleanSession(*cleansess)
	if *store != ":memory:" {
		opts.SetStore(MQTT.NewFileStore(*store))
	}

	if *action == "pub" {
		client := MQTT.NewClient(opts)
		_, err := client.Start()
		gotareceipt := make(chan bool)
		if err != nil {
			panic(err)
		}
		fmt.Println("Sample Publisher Started")
		for i := 0; i < *num; i++ {
			fmt.Println("---- doing publish ----")
			receipt := client.Publish(MQTT.QoS(*qos), *topic, []byte(*payload))

			go func() {
				<-receipt
				fmt.Println("  message delivered!")
				gotareceipt <- true
			}()
		}

		for i := 0; i < *num; i++ {
			<-gotareceipt
		}

		client.Disconnect(250)
		fmt.Println("Sample Publisher Disconnected")
	} else {
		num_received := 0
		choke := make(chan [2]string)

		opts.SetDefaultPublishHandler(func(client *MQTT.MqttClient, msg MQTT.Message) {
			choke <- [2]string{msg.Topic(), string(msg.Payload())}
		})

		client := MQTT.NewClient(opts)
		_, err := client.Start()
		if err != nil {
			panic(err)
		}

		filter, e := MQTT.NewTopicFilter(*topic, byte(*qos))
		if e != nil {
			fmt.Println(e)
			os.Exit(1)
		}

		client.StartSubscription(nil, filter)

		for num_received < *num {
			incoming := <-choke
			fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
			num_received++
		}

		client.Disconnect(250)
		fmt.Println("Sample Subscriber Disconnected")
	}
}
