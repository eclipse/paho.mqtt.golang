/*
 * Copyright (c) 2021 IBM Corp and others.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v2.0
 * and Eclipse Distribution License v1.0 which accompany this distribution.
 *
 * The Eclipse Public License is available at
 *    https://www.eclipse.org/legal/epl-2.0/
 * and the Eclipse Distribution License is available at
 *   http://www.eclipse.org/org/documents/edl-v10.php.
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 *    Matt Brittan
 */
package mqtt

import (
	"net"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestCustomConnectionFunction(t *testing.T) {
	// Set netpipe to emulate a connection of a different type
	netClient, netServer := net.Pipe()
	defer netClient.Close()
	defer netServer.Close()
	var firstMessage = ""
	go func() {
		// read first message only
		bytes := make([]byte, 1024)
		n, err := netServer.Read(bytes)
		if err != nil {
			t.Errorf("%v", err)
		}
		firstMessage = string(bytes[:n])
	}()
	// Set custom network connection function and client connect
	var customConnectionFunc OpenConnectionFunc = func(uri *url.URL, options ClientOptions) (net.Conn, error) {
		return netClient, nil
	}
	options := &ClientOptions{
		CustomOpenConnectionFn: customConnectionFunc,
	}
	brokerAddr := netServer.LocalAddr().Network()
	options.AddBroker(brokerAddr)
	client := NewClient(options)

	// Try to connect using custom function, wait for 2 seconds, to pass MQTT first message
	if token := client.Connect(); token.WaitTimeout(2*time.Second) && token.Error() != nil {
		t.Errorf("%v", token.Error())
	}

	// Analyze first message sent by client and received by the server
	if len(firstMessage) <= 0 || !strings.Contains(firstMessage, "MQTT") {
		t.Error("no message received on connect")
	}
}
