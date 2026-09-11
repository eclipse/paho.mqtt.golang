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

func Test_connectMQTT_rejectsOversizedConnack(t *testing.T) {
	// The broker response advertises a two-byte body but does not send it. A
	// one-byte limit must reject the packet before attempting to read the body.
	conn := bytes.NewBuffer([]byte{0x20, 0x02})
	connectPacket := packets.NewControlPacket(packets.Connect).(*packets.ConnectPacket)

	_, _, err := connectMQTT(conn, connectPacket, 4, noopSLogger, 1)
	if !errors.Is(err, packets.ErrPacketTooLarge) {
		t.Fatalf("expected ErrPacketTooLarge, got %v", err)
	}
}

func Test_startIncomingComms_rejectsOversizedPacket(t *testing.T) {
	conn := bytes.NewBuffer([]byte{0x30, 0xff, 0xff, 0xff, 0x7f})
	inboundFromStore := make(chan packets.ControlPacket)
	close(inboundFromStore)

	output := startIncomingComms(
		conn,
		&testCommsFns{maxIncomingPacketSize: 1024},
		inboundFromStore,
		noopSLogger,
	)

	select {
	case result := <-output:
		if !errors.Is(result.err, packets.ErrPacketTooLarge) {
			t.Fatalf("expected ErrPacketTooLarge, got %v", result.err)
		}
	case <-time.After(time.Second):
		t.Fatal("startIncomingComms did not report the oversized packet")
	}
}

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
	token                 tokenCompletor
	maxIncomingPacketSize uint32
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

func (c *testCommsFns) getMaxIncomingPacketSize() uint32 {
	return c.maxIncomingPacketSize
}

func (c *testCommsFns) persistOutbound(packets.ControlPacket) {}

func (c *testCommsFns) persistInbound(packets.ControlPacket) {}

func (c *testCommsFns) pingRespReceived() {}

func Test_verifyCONNACK_rejectsReservedReturnCode(t *testing.T) {
	// Fixed header 0x20, remaining length 2, flags 0, reserved return code 0x06.
	conn := bytes.NewBuffer([]byte{0x20, 0x02, 0x00, 0x06})
	rc, _, err := verifyCONNACK(conn, noopSLogger, 0)
	if !errors.Is(err, packets.ErrorProtocolViolation) {
		t.Fatalf("expected ErrorProtocolViolation, got rc=%d err=%v", rc, err)
	}
	if rc != packets.ErrNetworkError {
		t.Fatalf("expected ErrNetworkError sentinel on decode failure, got 0x%02x", rc)
	}
}

func Test_verifyCONNACK_acceptsStandardReturnCode(t *testing.T) {
	conn := bytes.NewBuffer([]byte{0x20, 0x02, 0x01, 0x00})
	rc, sessionPresent, err := verifyCONNACK(conn, noopSLogger, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rc != packets.Accepted {
		t.Fatalf("expected Accepted, got 0x%02x", rc)
	}
	if !sessionPresent {
		t.Fatal("expected SessionPresent")
	}
}
