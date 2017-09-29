package main

import (
	"log"
	"net"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/eclipse/paho.mqtt.golang/simple"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:1883")
	if err != nil {
		log.Fatalln(err)
	}
	c := simple.NewClient(conn)

	cp := packets.NewControlPacket(packets.CONNECT)
	connect := cp.Content.(*packets.Connect)
	connect.KeepAlive = 30
	connect.ClientID = "testGo"

	ca, err := c.Connect(cp)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(ca.Content.(*packets.Connack).Reason())

	s := packets.NewSubscribe(map[string]byte{"test/1": 2})
	s.Content.(*packets.Subscribe).PacketID = 1
	sa, err := c.Subscribe(s)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(sa.Content.(*packets.Suback).Reason(0))
}
