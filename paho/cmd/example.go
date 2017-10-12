package main

import (
	"log"
	"net"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/eclipse/paho.mqtt.golang/paho"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:1883")
	if err != nil {
		log.Fatalln(err)
	}
	c := paho.NewClient(conn, paho.NewPingHandler(conn, paho.PFH))
	ca, err := c.Connect(paho.Properties{}, packets.Connect{
		KeepAlive: 30,
		ClientID:  "testGo",
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(ca.Reason())
}
