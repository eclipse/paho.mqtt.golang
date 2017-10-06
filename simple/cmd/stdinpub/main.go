package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
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

	for {
		message, err := stdin.ReadString('\n')
		if err == io.EOF {
			os.Exit(0)
		}

		if _, err = c.SendMessage(*topic, byte(*qos), *retained, nil, []byte(message)); err != nil {
			log.Println(err)
		}
	}
}
