/*
 * Copyright (c) 2013 IBM Corp.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v1.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-v10.html
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 */

package mqtt

import "fmt"
import "time"
import "bytes"

import "io/ioutil"
import "crypto/tls"
import "crypto/x509"
import "testing"

func Test_Start(t *testing.T) {
	ops := NewClientOptions().SetClientId("Start").
		SetBroker(FVT_TCP).
		SetStore(NewFileStore("/tmp/fvt/Start"))
	c := NewClient(ops)

	_, err := c.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	c.Disconnect(250)
}

func Test_InvalidConnRc(t *testing.T) {
	ops := NewClientOptions().SetClientId("FailClientID").
		SetBroker("tcp://" + FVT_IP + ":17003").
		SetStore(NewFileStore("/tmp/fvt/InvalidConnRc"))

	c := NewClient(ops)
	_, err := c.Start()
	if err != ErrNotAuthorized {
		t.Fatalf("Did not receive error as expected, got %v", err)
	}
	c.Disconnect(250)
}

// Helper function for Test_Start_Ssl
func NewTlsConfig() *tls.Config {
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile("samples/samplecerts/CAfile.pem")
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	cert, err := tls.LoadX509KeyPair("samples/samplecerts/client-crt.pem", "samples/samplecerts/client-key.pem")
	if err != nil {
		panic(err)
	}

	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}
}

// temporarily disabled b/c Seth doesn't know
//   how to configure IMA for SSL
func Test_Start_Ssl(t *testing.T) {
	tlsconfig := NewTlsConfig()
	ops := NewClientOptions().SetClientId("StartSsl").
		SetBroker(FVT_SSL).
		SetStore(NewFileStore("/tmp/fvt/Start_Ssl")).
		SetTlsConfig(tlsconfig)

	c := NewClient(ops)

	_, err := c.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	c.Disconnect(250)
}

func Test_Publish_1(t *testing.T) {
	ops := NewClientOptions()
	ops.SetBroker(FVT_TCP)
	ops.SetClientId("Publish_1")
	ops.SetStore(NewFileStore("/tmp/fvt/Publish_1"))

	c := NewClient(ops)
	_, err := c.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	c.Publish(QOS_ZERO, "/test/Publish", "Publish qo0")

	c.Disconnect(250)
}

func Test_Publish_2(t *testing.T) {
	ops := NewClientOptions()
	ops.SetBroker(FVT_TCP)
	ops.SetClientId("Publish_2")
	ops.SetStore(NewFileStore("/tmp/fvt/Publish_2"))

	c := NewClient(ops)
	_, err := c.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	c.Publish(QOS_ZERO, "/test/Publish", "Publish1 qos0")
	c.Publish(QOS_ZERO, "/test/Publish", "Publish2 qos0")

	c.Disconnect(250)
}

func Test_Publish_3(t *testing.T) {
	ops := NewClientOptions()
	ops.SetBroker(FVT_TCP)
	ops.SetClientId("Publish_3")
	ops.SetStore(NewFileStore("/tmp/fvt/Publish_3"))

	c := NewClient(ops)
	_, err := c.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	c.Publish(QOS_ZERO, "/test/Publish", "Publish1 qos0")
	c.Publish(QOS_ONE, "/test/Publish", "Publish2 qos1")
	c.Publish(QOS_TWO, "/test/Publish", "Publish2 qos2")

	c.Disconnect(250)
}

func Test_Subscribe(t *testing.T) {
	pops := NewClientOptions()
	pops.SetBroker(FVT_TCP)
	pops.SetClientId("Subscribe_tx")
	pops.SetStore(NewFileStore("/tmp/fvt/Subscribe/p"))
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.SetBroker(FVT_TCP)
	sops.SetClientId("Subscribe_rx")
	sops.SetStore(NewFileStore("/tmp/fvt/Subscribe/s"))
	var f MessageHandler = func(msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	}
	sops.SetDefaultPublishHandler(f)
	s := NewClient(sops)

	_, err := s.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	filter, _ := NewTopicFilter("/test/sub", 0)
	s.StartSubscription(nil, filter)

	_, err = p.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	p.Publish(QOS_ZERO, "/test/sub", "Publish qos0")

	p.Disconnect(250)
	s.Disconnect(250)
}

func Test_Will(t *testing.T) {
	willmsgc := make(chan string)

	sops := NewClientOptions().SetBroker(FVT_TCP)
	sops.SetClientId("will-giver")
	sops.SetWill("/wills", "good-byte!", QOS_ZERO, false)
	sops.SetOnConnectionLost(func(err error) {
	})
	c := NewClient(sops)

	wops := NewClientOptions()
	wops.SetBroker(FVT_TCP)
	wops.SetClientId("will-subscriber")
	wops.SetStore(NewFileStore("/tmp/fvt/Will"))
	wops.SetDefaultPublishHandler(func(msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		willmsgc <- string(msg.Payload())
	})
	wsub := NewClient(wops)

	_, err := wsub.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	filter, _ := NewTopicFilter("/wills", 0)
	wsub.StartSubscription(nil, filter)

	_, err = c.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}
	time.Sleep(time.Duration(1) * time.Second)

	c.conn.Close() // Force kill the client with a will

	wsub.Disconnect(250)

	if <-willmsgc != "good-byte!" {
		t.Fatalf("will message did not have correct payload")
	}
}

func Test_Binary_Will(t *testing.T) {
	willmsgc := make(chan []byte)
	will := []byte{
		0xDE,
		0xAD,
		0xBE,
		0xEF,
	}

	sops := NewClientOptions().SetBroker(FVT_TCP)
	sops.SetClientId("will-giver")
	sops.SetBinaryWill("/wills", will, QOS_ZERO, false)
	sops.SetOnConnectionLost(func(err error) {
	})
	c := NewClient(sops)

	wops := NewClientOptions().SetBroker(FVT_TCP)
	wops.SetClientId("will-subscriber")
	wops.SetStore(NewFileStore("/tmp/fvt/Binary_Will"))
	wops.SetDefaultPublishHandler(func(msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %v\n", msg.Payload())
		willmsgc <- msg.Payload()
	})
	wsub := NewClient(wops)

	_, err := wsub.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	filter, _ := NewTopicFilter("/wills", 0)
	wsub.StartSubscription(nil, filter)

	_, err = c.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}
	time.Sleep(time.Duration(1) * time.Second)

	c.conn.Close() // Force kill the client with a will

	wsub.Disconnect(250)

	if !bytes.Equal(<-willmsgc, will) {
		t.Fatalf("will message did not have correct payload")
	}
}

/**
"[...] a publisher is responsible for determining the maximum QoS a
message can be delivered at, but a subscriber is able to downgrade
the QoS to one more suitable for its usage.
The QoS of a message is never upgraded."
**/

/***********************************
 * Tests to cover the 9 QoS combos *
 ***********************************/

func wait(c chan bool) {
	fmt.Println("choke is waiting")
	<-c
}

// Pub 0, Sub 0

func Test_p0s0(t *testing.T) {
	store := "/tmp/fvt/p0s0"
	topic := "/test/p0s0"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.SetBroker(FVT_TCP)
	pops.SetClientId("p0s0-pub")
	pops.SetStore(NewFileStore(store + "/p"))
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.SetBroker(FVT_TCP)
	sops.SetClientId("p0s0-sub")
	sops.SetStore(NewFileStore(store + "/s"))
	var f MessageHandler = func(msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	_, err := s.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	filter, _ := NewTopicFilter(topic, 0)
	receipt, err := s.StartSubscription(nil, filter)
	if err != nil {
		t.Fatalf("Error on MqttClient.StartSubscription(): ", err)
	}
	<-receipt

	_, err = p.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}
	receipt = p.Publish(QOS_ZERO, topic, "p0s0 payload 1")
	<-receipt
	receipt = p.Publish(QOS_ZERO, topic, "p0s0 payload 2")
	<-receipt

	wait(choke)
	wait(choke)

	receipt = p.Publish(QOS_ZERO, topic, "p0s0 payload 3")
	<-receipt
	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)

	chkcond(isemptydir(store + "/p"))
	chkcond(isemptydir(store + "/s"))
}

// Pub 0, Sub 1

func Test_p0s1(t *testing.T) {
	store := "/tmp/fvt/p0s1"
	topic := "/test/p0s1"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.SetBroker(FVT_TCP)
	pops.SetClientId("p0s1-pub")
	pops.SetStore(NewFileStore(store + "/p"))
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.SetBroker(FVT_TCP)
	sops.SetClientId("p0s1-sub")
	sops.SetStore(NewFileStore(store + "/s"))
	var f MessageHandler = func(msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	_, err := s.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	filter, _ := NewTopicFilter(topic, 0)
	receipt, err := s.StartSubscription(nil, filter)
	if err != nil {
		t.Fatalf("Error on MqttClient.StartSubscription(): ", err)
	}
	<-receipt

	_, err = p.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}
	receipt = p.Publish(QOS_ZERO, topic, "p0s1 payload 1")
	<-receipt
	receipt = p.Publish(QOS_ZERO, topic, "p0s1 payload 2")
	<-receipt
	wait(choke)
	wait(choke)

	receipt = p.Publish(QOS_ZERO, topic, "p0s1 payload 3")
	<-receipt
	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)

	chkcond(isemptydir(store + "/p"))
	chkcond(isemptydir(store + "/s"))
}

// Pub 0, Sub 2

func Test_p0s2(t *testing.T) {
	store := "/tmp/fvt/p0s2"
	topic := "/test/p0s2"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.SetBroker(FVT_TCP)
	pops.SetClientId("p0s2-pub")
	pops.SetStore(NewFileStore(store + "/p"))
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.SetBroker(FVT_TCP)
	sops.SetClientId("p0s2-sub")
	sops.SetStore(NewFileStore(store + "/s"))
	var f MessageHandler = func(msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	_, err := s.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	filter, _ := NewTopicFilter(topic, 2)
	receipt, err := s.StartSubscription(nil, filter)
	if err != nil {
		t.Fatalf("Error on MqttClient.StartSubscription(): ", err)
	}
	<-receipt

	_, err = p.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}
	receipt = p.Publish(QOS_ZERO, topic, "p0s2 payload 1")
	<-receipt
	receipt = p.Publish(QOS_ZERO, topic, "p0s2 payload 2")
	<-receipt
	wait(choke)
	wait(choke)

	receipt = p.Publish(QOS_ZERO, topic, "p0s2 payload 3")
	<-receipt
	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)

	chkcond(isemptydir(store + "/p"))
	chkcond(isemptydir(store + "/s"))
}

// Pub 1, Sub 0

func Test_p1s0(t *testing.T) {
	store := "/tmp/fvt/p1s0"
	topic := "/test/p1s0"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.SetBroker(FVT_TCP)
	pops.SetClientId("p1s0-pub")
	pops.SetStore(NewFileStore(store + "/p"))
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.SetBroker(FVT_TCP)
	sops.SetClientId("p1s0-sub")
	sops.SetStore(NewFileStore(store + "/s"))
	var f MessageHandler = func(msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	_, err := s.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	filter, _ := NewTopicFilter(topic, 0)
	receipt, err := s.StartSubscription(nil, filter)
	if err != nil {
		t.Fatalf("Error on MqttClient.StartSubscription(): ", err)
	}
	<-receipt

	_, err = p.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}
	receipt = p.Publish(QOS_ONE, topic, "p1s0 payload 1")
	<-receipt
	receipt = p.Publish(QOS_ONE, topic, "p1s0 payload 2")
	<-receipt
	wait(choke)
	wait(choke)

	receipt = p.Publish(QOS_ONE, topic, "p1s0 payload 3")
	<-receipt
	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)

	chkcond(isemptydir(store + "/p"))
	chkcond(isemptydir(store + "/s"))
}

// Pub 1, Sub 1

func Test_p1s1(t *testing.T) {
	store := "/tmp/fvt/p1s1"
	topic := "/test/p1s1"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.SetBroker(FVT_TCP)
	pops.SetClientId("p1s1-pub")
	pops.SetStore(NewFileStore(store + "/p"))
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.SetBroker(FVT_TCP)
	sops.SetClientId("p1s1-sub")
	sops.SetStore(NewFileStore(store + "/s"))
	var f MessageHandler = func(msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	_, err := s.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start()", err)
	}

	filter, _ := NewTopicFilter(topic, 1)
	receipt, err := s.StartSubscription(nil, filter)
	if err != nil {
		t.Fatalf("Error on MqttClient.StartSubscription()", err)
	}
	<-receipt

	_, err = p.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}
	receipt = p.Publish(QOS_ONE, topic, "p1s1 payload 1")
	<-receipt
	receipt = p.Publish(QOS_ONE, topic, "p1s1 payload 2")
	<-receipt
	wait(choke)
	wait(choke)

	receipt = p.Publish(QOS_ONE, topic, "p1s1 payload 3")
	<-receipt
	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)

	chkcond(isemptydir(store + "/p"))
	chkcond(isemptydir(store + "/s"))
}

// Pub 1, Sub 2

func Test_p1s2(t *testing.T) {
	store := "/tmp/fvt/p1s2"
	topic := "/test/p1s2"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.SetBroker(FVT_TCP)
	pops.SetClientId("p1s2-pub")
	pops.SetStore(NewFileStore(store + "/p"))
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.SetBroker(FVT_TCP)
	sops.SetClientId("p1s2-sub")
	sops.SetStore(NewFileStore(store + "/s"))
	var f MessageHandler = func(msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	_, err := s.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}
	filter, _ := NewTopicFilter(topic, 2)
	receipt, err := s.StartSubscription(nil, filter)
	if err != nil {
		t.Fatalf("Error on MqttClient.StartSubscription(): ", err)
	}
	<-receipt

	_, err = p.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}
	receipt = p.Publish(QOS_ONE, topic, "p1s2 payload 1")
	<-receipt
	receipt = p.Publish(QOS_ONE, topic, "p1s2 payload 2")
	<-receipt
	wait(choke)
	wait(choke)

	receipt = p.Publish(QOS_ONE, topic, "p1s2 payload 3")
	<-receipt
	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)

	chkcond(isemptydir(store + "/p"))
	chkcond(isemptydir(store + "/s"))
}

// Pub 2, Sub 0

func Test_p2s0(t *testing.T) {
	store := "/tmp/fvt/p2s0"
	topic := "/test/p2s0"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.SetBroker(FVT_TCP)
	pops.SetClientId("p2s0-pub")
	pops.SetStore(NewFileStore(store + "/p"))
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.SetBroker(FVT_TCP)
	sops.SetClientId("p2s0-sub")
	sops.SetStore(NewFileStore(store + "/s"))
	var f MessageHandler = func(msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	_, err := s.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	filter, _ := NewTopicFilter(topic, 0)
	receipt, err := s.StartSubscription(nil, filter)
	if err != nil {
		t.Fatalf("Error on MqttClient.StartSubscription(): ", err)
	}
	<-receipt

	_, err = p.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}
	p.Publish(QOS_TWO, topic, "p2s0 payload 1")
	p.Publish(QOS_TWO, topic, "p2s0 payload 2")
	wait(choke)
	wait(choke)

	p.Publish(QOS_TWO, topic, "p2s0 payload 3")
	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)

	chkcond(isemptydir(store + "/p"))
	chkcond(isemptydir(store + "/s"))
}

// Pub 2, Sub 1

func Test_p2s1(t *testing.T) {
	store := "/tmp/fvt/p2s1"
	topic := "/test/p2s1"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.SetBroker(FVT_TCP)
	pops.SetClientId("p2s1-pub")
	pops.SetStore(NewFileStore(store + "/p"))
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.SetBroker(FVT_TCP)
	sops.SetClientId("p2s1-sub")
	sops.SetStore(NewFileStore(store + "/s"))
	var f MessageHandler = func(msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	_, err := s.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	filter, _ := NewTopicFilter(topic, 1)
	receipt, err := s.StartSubscription(nil, filter)
	if err != nil {
		t.Fatalf("Error on MqttClient.StartSubscription(): ", err)
	}
	<-receipt

	_, err = p.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}
	receipt = p.Publish(QOS_TWO, topic, "p2s1 payload 1")
	<-receipt
	receipt = p.Publish(QOS_TWO, topic, "p2s1 payload 2")
	<-receipt
	wait(choke)
	wait(choke)

	receipt = p.Publish(QOS_TWO, topic, "p2s1 payload 3")
	<-receipt
	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)

	chkcond(isemptydir(store + "/p"))
	chkcond(isemptydir(store + "/s"))
}

// Pub 2, Sub 2

func Test_p2s2(t *testing.T) {
	store := "/tmp/fvt/p2s2"
	topic := "/test/p2s2"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.SetBroker(FVT_TCP)
	pops.SetClientId("p2s2-pub")
	pops.SetStore(NewFileStore(store + "/p"))
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.SetBroker(FVT_TCP)
	sops.SetClientId("p2s2-sub")
	sops.SetStore(NewFileStore(store + "/s"))
	var f MessageHandler = func(msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	_, err := s.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	filter, _ := NewTopicFilter(topic, 2)
	receipt, err := s.StartSubscription(nil, filter)
	if err != nil {
		t.Fatalf("Error on MqttClient.StartSubscription(): ", err)
	}
	<-receipt

	_, err = p.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}
	receipt = p.Publish(QOS_TWO, topic, "p2s2 payload 1")
	<-receipt
	receipt = p.Publish(QOS_TWO, topic, "p2s2 payload 2")
	<-receipt
	wait(choke)
	wait(choke)

	receipt = p.Publish(QOS_TWO, topic, "p2s2 payload 3")
	<-receipt
	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)

	chkcond(isemptydir(store + "/p"))
	chkcond(isemptydir(store + "/s"))
}

func Test_PublishMessage(t *testing.T) {
	store := "/tmp/fvt/PublishMessage"
	topic := "/test/pubmsg"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.SetBroker(FVT_TCP)
	pops.SetClientId("pubmsg-pub")
	pops.SetStore(NewFileStore(store + "/p"))
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.SetBroker(FVT_TCP)
	sops.SetClientId("pubmsg-sub")
	sops.SetStore(NewFileStore(store + "/s"))
	var f MessageHandler = func(msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		if string(msg.Payload()) != "pubmsg payload" {
			t.Fatalf("Message payload incorrect")
		}
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	_, err := s.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	filter, _ := NewTopicFilter(topic, 2)
	receipt, err := s.StartSubscription(nil, filter)
	if err != nil {
		t.Fatalf("Error on MqttClient.StartSubscription(): ", err)
	}
	<-receipt

	_, err = p.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	text := "pubmsg payload"
	m := NewMessage(text)
	p.PublishMessage(topic, m)
	p.PublishMessage(topic, m)
	wait(choke)
	wait(choke)

	p.PublishMessage(topic, m)
	wait(choke)

	p.Disconnect(250)
	s.Disconnect(250)

	chkcond(isemptydir(store + "/p"))
	chkcond(isemptydir(store + "/s"))
}

func Test_PublishEmptyMessage(t *testing.T) {
	store := "/tmp/fvt/PublishEmptyMessage"
	topic := "/test/pubmsgempty"
	choke := make(chan bool)

	pops := NewClientOptions()
	pops.SetBroker(FVT_TCP)
	pops.SetClientId("pubmsgempty-pub")
	pops.SetStore(NewFileStore(store + "/p"))
	p := NewClient(pops)

	sops := NewClientOptions()
	sops.SetBroker(FVT_TCP)
	sops.SetClientId("pubmsgempty-sub")
	sops.SetStore(NewFileStore(store + "/s"))
	var f MessageHandler = func(msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		if string(msg.Payload()) != "" {
			t.Fatalf("Message payload incorrect")
		}
		choke <- true
	}
	sops.SetDefaultPublishHandler(f)

	s := NewClient(sops)
	_, err := s.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	filter, _ := NewTopicFilter(topic, 2)
	receipt, err := s.StartSubscription(nil, filter)
	if err != nil {
		t.Fatalf("Error on MqttClient.StartSubscription(): ", err)
	}
	<-receipt

	_, err = p.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	m := NewMessage(nil)
	p.PublishMessage(topic, m)
	p.PublishMessage(topic, m)
	wait(choke)
	wait(choke)

	p.PublishMessage(topic, m)
	wait(choke)

	p.Disconnect(250)
}

func Test_Cleanstore(t *testing.T) {
	store := "/tmp/fvt/cleanstore"
	topic := "/test/cleanstore"

	pops := NewClientOptions()
	pops.SetBroker(FVT_TCP)
	pops.SetClientId("cleanstore-pub")
	pops.SetStore(NewFileStore(store + "/p"))
	p := NewClient(pops)

	var s *MqttClient
	sops := NewClientOptions()
	sops.SetBroker(FVT_TCP)
	sops.SetClientId("cleanstore-sub")
	sops.SetCleanSession(false)
	sops.SetStore(NewFileStore(store + "/s"))
	var f MessageHandler = func(msg Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
		// Close the connection after receiving
		// the first message so that hopefully
		// there is something in the store to be
		// cleaned.
		s.conn.Close()
	}
	sops.SetDefaultPublishHandler(f)
	var ocl OnConnectionLost = func(reason error) {
		fmt.Printf("OnConnectionLost\n")
	}
	sops.SetOnConnectionLost(ocl)

	s = NewClient(sops)
	_, err := s.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	filter, _ := NewTopicFilter(topic, 2)
	receipt, err := s.StartSubscription(nil, filter)
	if err != nil {
		t.Fatalf("Error on MqttClient.StartSubscription(): ", err)
	}
	<-receipt

	_, err = p.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	m := NewMessage("test message")
	p.PublishMessage(topic, m)
	p.PublishMessage(topic, m)
	p.PublishMessage(topic, m)

	p.Disconnect(250)

	sops = NewClientOptions()
	sops.SetBroker(FVT_TCP)
	sops.SetClientId("cleanstore-sub")
	sops.SetCleanSession(true)
	sops.SetStore(NewFileStore(store + "/s"))
	sops.SetDefaultPublishHandler(f)

	s2 := NewClient(sops)
	_, err = s2.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}

	// at this point existing state should be cleared...
	// how to check?
}

/*
// A test to make sure ping mechanism is working
// This test can be left commented out because it's annoying to wait for
func Test_ping3_idle10(t *testing.T) {
	ops := NewClientOptions()
	ops.SetBroker(FVT_TCP)
	//ops.SetBroker("tcp://test.mosquitto.org:1883")
	ops.SetClientId("p3i10")
	ops.SetTimeout(4)

	c := NewClient(ops)
	_, err := c.Start()
	if err != nil {
		t.Fatalf("Error on MqttClient.Start(): ", err)
	}
	time.Sleep(time.Duration(10) * time.Second)
	c.Disconnect(250)
}
*/
