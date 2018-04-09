package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	pk "github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/eclipse/paho.mqtt.golang/paho"
)

func main() {
	stdin := bufio.NewReader(os.Stdin)
	hostname, _ := os.Hostname()

	server := flag.String("server", "127.0.0.1:1883", "The full URL of the MQTT server to connect to")
	topic := flag.String("topic", hostname, "Topic to publish and receive the messages on")
	qos := flag.Int("qos", 0, "The QoS to send the messages at")
	name := flag.String("chatname", hostname, "The name to attach to your messages")
	clientid := flag.String("clientid", "", "A clientid for the connection")
	username := flag.String("username", "", "A username to authenticate to the MQTT server")
	password := flag.String("password", "", "Password to match username")
	flag.Parse()

	c, err := paho.NewClient(
		paho.OpenTCPConn(*server),
		paho.DefaultMessageHandler(func(m paho.Message) {
			log.Printf("%s : %s", m.Properties.User["chatname"], string(m.Payload))
		}),
	)

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

	_, err = c.Connect(cp)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Connected to %s\n", *server)

	ic := make(chan os.Signal, 1)
	signal.Notify(ic, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ic
		fmt.Println("signal received, exiting")
		if c != nil {
			d := pk.NewDisconnect(pk.DisconnectReason(0))
			c.Disconnect(d)
		}
		os.Exit(0)
	}()

	_, err = c.Subscribe(
		pk.NewSubscribe(
			pk.Sub(*topic, pk.SubOptions{QoS: byte(*qos), NoLocal: true}),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Subscribed to %s", *topic)

	for {
		message, err := stdin.ReadString('\n')
		if err == io.EOF {
			os.Exit(0)
		}
		pb := pk.NewPublish(
			pk.Message(*topic, byte(*qos), false, []byte(message)),
			pk.PublishProperties(pk.NewProperties(
				pk.UserSingle("chatname", *name),
			)),
		)

		if _, err = c.Publish(pb); err != nil {
			log.Println(err)
		}
	}
}
