package main

import (
	"log"
	"os"

	pk "github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/eclipse/paho.mqtt.golang/paho"
)

func main() {
	paho.SetDebugLogger(log.New(os.Stdout, "DEBUG: ", log.LstdFlags))
	c, err := paho.NewClient(
		paho.OpenTCPConn("127.0.0.1:1883"),
	)

	cp := pk.NewConnect(
		pk.KeepAlive(30),
		pk.ClientID("testGo"),
	)

	log.Println("Connecting")
	ca, err := c.Connect(cp)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Connected")

	log.Println(ca.Reason())

	s := pk.NewSubscribe(
		pk.Sub("test/1", pk.SubOptions{QoS: 2}),
	)
	sa, err := c.Subscribe(s)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(sa.Reason(0))
}
