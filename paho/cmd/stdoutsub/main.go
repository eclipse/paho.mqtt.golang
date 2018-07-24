package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	msgChan := make(chan *paho.Publish)

	c, err := paho.NewClient(
		paho.OpenTCPConn(*server),
		paho.DefaultMessageHandler(func(m *paho.Publish) {
			msgChan <- m
		}))

	cp := &paho.Connect{
		KeepAlive:  30,
		ClientID:   *clientid,
		CleanStart: true,
		Username:   *username,
		Password:   []byte(*password),
	}

	if *username != "" {
		cp.UsernameFlag = true
	}
	if *password != "" {
		cp.PasswordFlag = true
	}

	ca, err := c.Connect(cp)
	if err != nil {
		log.Fatalln(err)
	}
	if ca.ReasonCode != 0 {
		log.Fatalf("Failed to connect to %s : %d - %s", *server, ca.ReasonCode, ca.Properties.ReasonString)
	}

	fmt.Printf("Connected to %s\n", *server)

	ic := make(chan os.Signal, 1)
	signal.Notify(ic, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ic
		fmt.Println("signal received, exiting")
		if c != nil {
			d := &paho.Disconnect{ReasonCode: 0}
			c.Disconnect(d)
		}
		os.Exit(0)
	}()

	sa, err := c.Subscribe(context.Background(), &paho.Subscribe{
		Subscriptions: map[string]paho.SubscribeOptions{
			*topic: paho.SubscribeOptions{QoS: byte(*qos)},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	if sa.Reasons[0] != byte(*qos) {
		log.Fatalf("Failed to subscribe to %s : %d", *topic, sa.Reasons[0])
	}
	log.Printf("Subscribed to %s", *topic)

	for m := range msgChan {
		log.Println("Received message:", string(m.Payload))
	}
}
