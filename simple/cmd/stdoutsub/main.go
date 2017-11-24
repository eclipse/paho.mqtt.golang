package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	pk "github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/eclipse/paho.mqtt.golang/simple"
)

func main() {
	server := flag.String("server", "127.0.0.1:1883", "The MQTT server to connect to ex: 127.0.0.1:1883")
	topic := flag.String("topic", "#", "Topic to subscribe to")
	qos := flag.Int("qos", 0, "The QoS to subscribe to messages at")
	clientid := flag.String("clientid", "", "A clientid for the connection")
	username := flag.String("username", "", "A username to authenticate to the MQTT server")
	password := flag.String("password", "", "Password to match username")
	flag.Parse()

	c, err := simple.NewClient(simple.OpenConn("tcp", *server))

	cp := pk.NewConnect(
		pk.KeepAlive(30),
		pk.ClientID(*clientid),
		pk.CleanStart(true),
	)

	if *username != "" {
		pk.Username(*username)(cp)
	}
	if *password != "" {
		pk.Password([]byte(*password))(cp)
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
		for {
			select {
			case <-time.After(30 * time.Second):
				if err := c.Ping(); err != nil {
					log.Fatalln(err)
				}
			case <-ic:
				fmt.Println("signal received, exiting")
				if c != nil {
					d := pk.NewDisconnect(pk.DisconnectReason(0))
					c.Disconnect(d)
				}
				os.Exit(0)
			}
		}
	}()

	s := pk.NewSubscribe(
		pk.Sub(*topic, pk.SubOptions{QoS: byte(*qos)}),
	)
	s.PacketID = 1

	sa, err := c.Subscribe(s)
	if err != nil {
		log.Fatalln(err)
	}
	if sa.Reasons[0] != 0 {
		log.Fatalf("Failed to subscribe to %s : %s", *topic, sa.Reason(0))
	}
	log.Printf("Subscribed to %s", *topic)

	for {
		pb, err := c.Receive()
		if err != nil {
			log.Fatalln(err)
		}

		if pb != nil {
			log.Printf("Received message %s", pb)
		}
	}
}
