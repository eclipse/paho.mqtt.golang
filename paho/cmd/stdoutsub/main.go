package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	pk "github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/eclipse/paho.mqtt.golang/paho"
)

func main() {
	server := flag.String("server", "127.0.0.1:1883", "The MQTT server to connect to ex: 127.0.0.1:1883")
	topic := flag.String("topic", "#", "Topic to subscribe to")
	qos := flag.Int("qos", 0, "The QoS to subscribe to messages at")
	clientid := flag.String("clientid", "", "A clientid for the connection")
	username := flag.String("username", "", "A username to authenticate to the MQTT server")
	password := flag.String("password", "", "Password to match username")
	flag.Parse()

	paho.SetDebugLogger(log.New(os.Stderr, "SUB: ", log.LstdFlags))
	msgChan := make(chan paho.Message)

	c, err := paho.NewClient(
		paho.OpenTCPConn(*server),
		paho.DefaultMessageHandler(func(m paho.Message) {
			msgChan <- m
		}))

	cp := &pk.Connect{
		KeepAlive:  30,
		ClientID:   *clientid,
		CleanStart: true,
	}

	if *username != "" {
		cp.Username = *username
	}
	if *password != "" {
		cp.Password = []byte(*password)
	}

	ca, err := c.Connect(cp)
	if err != nil {
		log.Fatalln(err)
	}
	if ca.ReasonCode != 0 {
		log.Fatalf("Failed to connect to %s : %d - %s", *server, ca.ReasonCode, ca.Reason())
	}

	fmt.Printf("Connected to %s\n", *server)

	ic := make(chan os.Signal, 1)
	signal.Notify(ic, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ic
		fmt.Println("signal received, exiting")
		if c != nil {
			c.Disconnect(&pk.Disconnect{ReasonCode: pk.DisconnectNormalDisconnection})
		}
		os.Exit(0)
	}()

	sa, err := c.Subscribe(&pk.Subscribe{
		Subscriptions: map[string]pk.SubOptions{
			*topic: pk.SubOptions{QoS: byte(*qos)},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	if sa.Reasons[0] != byte(*qos) {
		log.Fatalf("Failed to subscribe to %s : %s", *topic, sa.Reason(0))
	}
	log.Printf("Subscribed to %s", *topic)

	for m := range msgChan {
		log.Println("Received message:", string(m.Payload))
	}
}
