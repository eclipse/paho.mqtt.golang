package main

import (
	"log"

	pk "github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/eclipse/paho.mqtt.golang/simple"
)

func main() {
	c, err := simple.NewClient(simple.OpenConn("tcp", "127.0.0.1:1883"))

	cp := pk.NewConnect(
		pk.KeepAlive(30),
		pk.ClientID("testGo"),
	)

	ca, err := c.Connect(cp)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(ca.Reason())

	s := pk.NewSubscribe(
		pk.Sub("test/1", pk.SubOptions{QoS: 2}),
	)
	s.PacketID = 1
	sa, err := c.Subscribe(s)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(sa.Reason(0))
}
