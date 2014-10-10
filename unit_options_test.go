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

import (
	"crypto/tls"
	"crypto/x509"
	"testing"
	"time"
)

func Test_NewClientOptions_default(t *testing.T) {
	o := NewClientOptions()

	if o.ClientId != "" {
		t.Fatalf("bad default client id")
	}

	if o.Username != "" {
		t.Fatalf("bad default username")
	}

	if o.Password != "" {
		t.Fatalf("bad default password")
	}

	if o.KeepAlive != 30*time.Second {
		t.Fatalf("bad default timeout")
	}
}

func Test_NewClientOptions_mix(t *testing.T) {
	o := NewClientOptions()
	o.AddBroker("tcp://192.168.1.2:9999")
	o.SetClientId("myclientid")
	o.SetUsername("myuser")
	o.SetPassword("mypassword")
	o.SetKeepAlive(88)

	if o.Servers[0].Scheme != "tcp" {
		t.Fatalf("bad scheme")
	}

	if o.Servers[0].Host != "192.168.1.2:9999" {
		t.Fatalf("bad host")
	}

	if o.ClientId != "myclientid" {
		t.Fatalf("bad set clientid")
	}

	if o.Username != "myuser" {
		t.Fatalf("bad set username")
	}

	if o.Password != "mypassword" {
		t.Fatalf("bad set password")
	}

	if o.KeepAlive != 88 {
		t.Fatalf("bad set timeout")
	}
}

func Test_ModifyOptions(t *testing.T) {
	o := NewClientOptions()
	o.AddBroker("tcp://3.3.3.3:12345")
	c := NewClient(o)
	o.AddBroker("ws://2.2.2.2:9999")
	o.SetOrderMatters(false)

	if c.options.Servers[0].Scheme != "tcp" {
		t.Fatalf("client options.server.Scheme was modified")
	}

	// if c.options.server.Host != "2.2.2.2:9999" {
	// 	t.Fatalf("client options.server.Host was modified")
	// }

	if o.Order != false {
		t.Fatalf("options.order was not modified")
	}
}

func Test_TlsConfig(t *testing.T) {
	o := NewClientOptions().SetTlsConfig(&tls.Config{
		RootCAs:            x509.NewCertPool(),
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          x509.NewCertPool(),
		InsecureSkipVerify: true})

	c := NewClient(o)

	if c.options.TlsConfig.ClientAuth != tls.NoClientCert {
		t.Fatalf("client options.tlsConfig ClientAuth incorrect")
	}

	if c.options.TlsConfig.InsecureSkipVerify != true {
		t.Fatalf("client options.tlsConfig InsecureSkipVerify incorrect")
	}
}

func Test_OnConnectionLost(t *testing.T) {
	onconnlost := func(client *MqttClient, err error) {
		panic(err)
	}
	o := NewClientOptions().SetConnectionLostHandler(onconnlost)

	c := NewClient(o)

	if c.options.OnConnectionLost == nil {
		t.Fatalf("client options.onconnlost was nil")
	}
}
