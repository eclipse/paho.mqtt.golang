package packets

import (
	"bytes"
	"io"
	"net"
)

// Suback is the Variable Header definition for a Suback control packet
type Suback struct {
	PacketID   uint16
	Properties Properties
	Reasons    []byte
}

// SubackGrantedQoS0, etc are the list of valid suback reason codes.
const (
	SubackGrantedQoS0                         = 0x00
	SubackGrantedQoS1                         = 0x01
	SubackGrantedQoS2                         = 0x02
	SubackUnspecifiederror                    = 0x80
	SubackImplementationspecificerror         = 0x83
	SubackNotauthorized                       = 0x87
	SubackTopicFilterinvalid                  = 0x8F
	SubackPacketIdentifierinuse               = 0x91
	SubackQuotaexceeded                       = 0x97
	SubackSharedSubscriptionnotsupported      = 0x9E
	SubackSubscriptionIdentifiersnotsupported = 0xA1
	SubackWildcardsubscriptionsnotsupported   = 0xA2
)

//Unpack is the implementation of the interface required function for a packet
func (s *Suback) Unpack(r *bytes.Buffer) error {
	var err error
	s.PacketID, err = readUint16(r)
	if err != nil {
		return err
	}

	err = s.Properties.Unpack(r, SUBACK)
	if err != nil {
		return err
	}

	s.Reasons = r.Bytes()

	return nil
}

// Buffers is the implementation of the interface required function for a packet
func (s *Suback) Buffers() net.Buffers {
	var b bytes.Buffer
	writeUint16(s.PacketID, &b)
	idvp := s.Properties.Pack(SUBACK)
	propLen := encodeVBI(len(idvp))
	return net.Buffers{b.Bytes(), propLen, idvp, s.Reasons}
}

// Send is the implementation of the interface required function for a packet
func (s *Suback) Send(w io.Writer) error {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: SUBACK}}
	cp.Content = s

	return cp.Send(w)
}

// Reason returns a string representation of the meaning of the ReasonCode
func (s *Suback) Reason(index int) string {
	if index >= 0 && index < len(s.Reasons) {
		switch s.Reasons[index] {
		case 0:
			return "Granted QoS 0 - The subscription is accepted and the maximum QoS sent will be QoS 0. This might be a lower QoS than was requested."
		case 1:
			return "Granted QoS 1 - The subscription is accepted and the maximum QoS sent will be QoS 1. This might be a lower QoS than was requested."
		case 2:
			return "Granted QoS 2 - The subscription is accepted and any received QoS will be sent to this subscription."
		case 128:
			return "Unspecified error - The subscription is not accepted and the Server either does not wish to reveal the reason or none of the other Reason Codes apply."
		case 131:
			return "Implementation specific error - The SUBSCRIBE is valid but the Server does not accept it."
		case 135:
			return "Not authorized - The Client is not authorized to make this subscription."
		case 143:
			return "Topic Filter invalid - The Topic Filter is correctly formed but is not allowed for this Client."
		case 145:
			return "Packet Identifier in use - The specified Packet Identifier is already in use."
		case 151:
			return "Quota exceeded - An implementation or administrative imposed limit has been exceeded."
		case 158:
			return "Shared Subscription not supported - The Server does not support Shared Subscriptions for this Client."
		case 161:
			return "Subscription Identifiers not supported - The Server does not support Subscription Identifiers; the subscription is not accepted."
		case 162:
			return "Wildcard subscriptions not supported - The Server does not support Wildcard subscription; the subscription is not accepted."
		}
	}
	return "Invalid Reason index"
}

// NewSuback creates a new Suback packet and applies all the
// provided/listed option functions to configure the packet
func NewSuback(opts ...func(sa *Suback)) *Suback {
	sa := &Suback{
		Properties: Properties{
			User: make(map[string]string),
		},
	}

	for _, opt := range opts {
		opt(sa)
	}

	return sa
}

// SubackReasons is a Suback option function that sets the
// reason codes for the Suback packet
func SubackReasons(r []byte) func(*Suback) {
	return func(sa *Suback) {
		sa.Reasons = r
	}
}

// SubackProperties is a Suback option function that sets
// the Properties for the Suback packet
func SubackProperties(p *Properties) func(*Suback) {
	return func(sa *Suback) {
		sa.Properties = *p
	}
}
