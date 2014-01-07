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

import "os"
import "crypto/tls"
import "crypto/x509"
import "testing"

func Test_NewClientOptions_default(t *testing.T) {
	o := NewClientOptions()

	if o.clientId != "" {
		t.Fatalf("bad default client id")
	}

	if o.username != "" {
		t.Fatalf("bad default username")
	}

	if o.password != "" {
		t.Fatalf("bad default password")
	}

	if o.tlsconfig != nil {
		t.Fatalf("bad default tlsconfig")
	}

	if o.timeout != 30 {
		t.Fatalf("bad default timeout")
	}

	if o.tracefile != os.Stdout {
		t.Fatalf("bad default tracefile")
	}

	if o.msgRouter == nil {
		t.Fatalf("bad default msgRouter")
	}

	if o.stopRouter == nil {
		t.Fatalf("bad stopRouter")
	}

	if o.pubChanZero != nil {
		t.Fatalf("bad default pubChanZero")
	}

	if o.pubChanOne != nil {
		t.Fatalf("bad default pubChanOne")
	}

	if o.pubChanTwo != nil {
		t.Fatalf("bad default pubChanTwo")
	}
}

func Test_NewClientOptions_mix(t *testing.T) {
	o := NewClientOptions()
	o.SetBroker("tcp://192.168.1.2:9999")
	o.SetClientId("myclientid")
	o.SetUsername("myuser")
	o.SetPassword("mypassword")
	o.SetTimeout(88)
	o.SetTracefile(os.Stderr)
	o.SetTraceLevel(Warn)

	if o.server.Scheme != "tcp" {
		t.Fatalf("bad scheme")
	}

	if o.server.Host != "192.168.1.2:9999" {
		t.Fatalf("bad host")
	}

	if o.clientId != "myclientid" {
		t.Fatalf("bad set clientid")
	}

	if o.username != "myuser" {
		t.Fatalf("bad set username")
	}

	if o.password != "mypassword" {
		t.Fatalf("bad set password")
	}

	if o.timeout != 88 {
		t.Fatalf("bad set timeout")
	}

	if o.tracefile != os.Stderr {
		t.Fatalf("bad set tracefile")
	}

	if o.tracelevel != Warn {
		t.Fatalf("bad set tracelevel")
	}
}

func Test_ModifyOptions(t *testing.T) {
	o := NewClientOptions()
	o.SetBroker("tcp://3.3.3.3:12345")
	c := NewClient(o)
	o.SetBroker("ws://2.2.2.2:9999")
	o.SetOrderMatters(false)

	if c.options.server.Scheme != "tcp" {
		t.Fatalf("client options.server.Scheme was modified")
	}

	// if c.options.server.Host != "2.2.2.2:9999" {
	// 	t.Fatalf("client options.server.Host was modified")
	// }

	if o.order != false {
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

	if c.options.tlsconfig == nil {
		t.Fatalf("client options.tlsconfig was nil")
	}

	if c.options.tlsconfig.ClientAuth != tls.NoClientCert {
		t.Fatalf("client options.tlsconfig ClientAuth incorrect")
	}

	if c.options.tlsconfig.InsecureSkipVerify != true {
		t.Fatalf("client options.tlsconfig InsecureSkipVerify incorrect")
	}
}

func Test_OnConnectionLost(t *testing.T) {
	onconnlost := func(err error) {
		panic(err)
	}
	o := NewClientOptions().SetOnConnectionLost(onconnlost)

	c := NewClient(o)

	if c.options.onconnlost == nil {
		t.Fatalf("client options.onconnlost was nil")
	}
}

// func Test_MaxInFlight(t *testing.T) {
// 	dflt := NewClientOptions()
// 	if dflt.maxinflight != 10 {
// 		t.Fatalf("options.maxinflight was not defaulted to 10")
// 	}

// 	set := NewClientOptions().SetMaxInFlight(1)
// 	if set.maxinflight != 1 {
// 		t.Fatalf("options.maxinflight was not set to 1")
// 	}
// }
