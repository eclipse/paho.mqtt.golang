/*
 * Copyright (c) 2026 IBM Corp and others.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v2.0
 * and Eclipse Distribution License v1.0 which accompany this distribution.
 *
 * The Eclipse Public License is available at
 *    https://www.eclipse.org/legal/epl-2.0/
 * and the Eclipse Distribution License is available at
 *   http://www.eclipse.org/org/documents/edl-v10.php.
 */

package mqtt

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
)

// Check that the library handles a case where a malicious server sends a SUBACK packet with a different number of return codes than
// the matching SUBSCRIBE requested subscriptions.
// See [MQTT-3.8.4-5] in the spec for the relevant rule.
func Test_startIncomingComms_subackReturnCodeSubscriptionMismatch(t *testing.T) {
	const messageID = 1

	cases := []struct {
		subs          []string
		returnCodes   []byte
		expectedError bool
	}{
		{
			subs:          []string{"topic/a"},
			returnCodes:   []byte{0},
			expectedError: false,
		},
		{
			subs:          []string{"topic/a", "topic/b"},
			returnCodes:   []byte{0, 0},
			expectedError: false,
		},
		{
			subs:          []string{"topic/a"},
			returnCodes:   []byte{0, 1},
			expectedError: true,
		},
		{
			subs:          []string{"topic/a", "topic/b"},
			returnCodes:   []byte{0},
			expectedError: true,
		},
		{
			subs:          []string{"topic/a"},
			returnCodes:   []byte{},
			expectedError: true,
		},
	}

	for _, c := range cases {
		// Set the topics requested in the matching SUBSCRIBE.
		token := newToken(packets.Subscribe).(*SubscribeToken)
		token.subs = c.subs

		// Build the SUBACK to process.
		suback := packets.NewControlPacket(packets.Suback).(*packets.SubackPacket)
		suback.MessageID = messageID
		suback.ReturnCodes = c.returnCodes
		var conn bytes.Buffer
		if err := suback.Write(&conn); err != nil {
			t.Fatalf("failed to write suback: %v", err)
		}

		inboundFromStore := make(chan packets.ControlPacket) // Store unused in this test
		close(inboundFromStore)

		// Start the incoming processor, the SUBACK is already on conn
		output := startIncomingComms(&conn, &testCommsFns{token: token}, inboundFromStore, noopSLogger)

		// Regardless of the result, the token should be done
		select {
		case <-token.Done():
		case <-time.After(time.Second):
			t.Fatalf("subscribe token was not completed")
		}

		// capture everything from the channel to ensure it's fully drained
		var received []incomingComms
	drainOutput:
		for {
			select {
			case msg, ok := <-output:
				if !ok {
					break drainOutput
				}
				received = append(received, msg)
			case <-time.After(time.Second):
				t.Fatalf("startIncomingComms did not complete")
			}
		}

		if c.expectedError {
			malformedSubackErrors := 0
			for _, msg := range received {
				if errors.Is(msg.err, ErrMalformedSuback) {
					malformedSubackErrors++
				}
			}
			if malformedSubackErrors != 1 {
				t.Errorf("expected ErrMalformedSuback once on chan (sub: %v, codes: %v), got %d in %v", c.subs, c.returnCodes, malformedSubackErrors, received)
			}
			if !errors.Is(token.Error(), ErrMalformedSuback) {
				t.Errorf("expected ErrMalformedSuback (sub: %v, codes: %v), got %v", c.subs, c.returnCodes, token.Error())
			}
		} else {
			if len(received) != 1 || !errors.Is(received[0].err, io.EOF) {
				t.Errorf("expected normal closure, got %v", received)
			}
			if token.Error() != nil {
				t.Errorf("expected successful SUBACK (sub: %v, codes: %v), got %v", c.subs, c.returnCodes, token.Error())
			}
		}
	}
}

// testCommsFns is a basic implementation of commsFns for use with startIncomingComms
type testCommsFns struct {
	token tokenCompletor
}

func (c *testCommsFns) getToken(uint16) tokenCompletor {
	return c.token
}

func (c *testCommsFns) freeID(uint16) {}

func (c *testCommsFns) UpdateLastReceived() {}

func (c *testCommsFns) UpdateLastSent() {}

func (c *testCommsFns) getWriteTimeOut() time.Duration {
	return 0
}

func (c *testCommsFns) persistOutbound(packets.ControlPacket) {}

func (c *testCommsFns) persistInbound(packets.ControlPacket) {}

func (c *testCommsFns) pingRespReceived() {}
