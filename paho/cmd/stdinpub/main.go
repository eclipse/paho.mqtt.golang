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
	topic := flag.String("topic", hostname, "Topic to publish the messages on")
	qos := flag.Int("qos", 0, "The QoS to send the messages at")
	retained := flag.Bool("retained", false, "Are the messages sent with the retained flag")
	clientid := flag.String("clientid", "", "A clientid for the connection")
	username := flag.String("username", "", "A username to authenticate to the MQTT server")
	password := flag.String("password", "", "Password to match username")
	flag.Parse()

	c, err := paho.NewClient(paho.OpenTCPConn(*server))

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
			d := pk.NewDisconnect(pk.DisconnectReason(0))
			c.Disconnect(d)
		}
		os.Exit(0)
	}()

	for {
		message, err := stdin.ReadString('\n')
		if err == io.EOF {
			os.Exit(0)
		}

		if _, err = c.Publish(&pk.Publish{
			Topic:   *topic,
			QoS:     byte(*qos),
			Retain:  *retained,
			Payload: []byte(message),
		}); err != nil {
			log.Println(err)
		}
	}
}
