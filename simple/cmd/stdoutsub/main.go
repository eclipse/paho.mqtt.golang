package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
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

	conn, err := net.Dial("tcp", *server)
	if err != nil {
		log.Fatalln(err)
	}
	c := simple.NewClient(conn)

	cp := packets.NewControlPacket(packets.CONNECT)
	connect := cp.Content.(*packets.Connect)
	connect.KeepAlive = 30
	connect.ClientID = *clientid
	connect.CleanStart = true
	if *username != "" {
		connect.UsernameFlag = true
		connect.Username = *username
		connect.PasswordFlag = true
		connect.Password = []byte(*password)
	}

	ca, err := c.Connect(cp)
	if err != nil {
		log.Fatalln(err)
	}
	if ca.Content.(*packets.Connack).ReasonCode != 0 {
		log.Fatalf("Failed to connect to %s : %d - %s", *server, ca.Content.(*packets.Connack).ReasonCode, ca.Content.(*packets.Connack).Reason())
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
					d := packets.NewControlPacket(packets.DISCONNECT)
					d.Content.(*packets.Disconnect).DisconnectReasonCode = 0

					c.Disconnect(d)
					if conn != nil {
						conn.Close()
					}
				}
				os.Exit(0)
			}
		}
	}()

	s := packets.NewSubscribe(map[string]packets.SubOptions{*topic: packets.SubOptions{QoS: byte(*qos)}})
	s.Content.(*packets.Subscribe).PacketID = 1
	sa, err := c.Subscribe(s)
	if err != nil {
		log.Fatalln(err)
	}
	if sa.Content.(*packets.Suback).Reasons[0] != 0 {
		log.Fatalf("Failed to subscribe to %s : %s", *topic, sa.Content.(*packets.Suback).Reason(0))
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
