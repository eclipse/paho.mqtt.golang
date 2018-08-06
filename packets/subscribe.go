package packets

import (
	"bytes"
	"io"
	"net"
)

// Subscribe is the Variable Header definition for a Subscribe control packet
type Subscribe struct {
	PacketID      uint16
	Properties    *Properties
	Subscriptions map[string]SubOptions
}

// SubOptions is the struct representing the options for a subscription
type SubOptions struct {
	QoS               byte
	NoLocal           bool
	RetainAsPublished bool
	RetainHandling    byte
}

// Pack is the implementation of the interface required function for a packet
func (s *SubOptions) Pack() byte {
	var ret byte
	ret |= s.QoS & 0x03
	if s.NoLocal {
		ret |= 1 << 2
	}
	if s.RetainAsPublished {
		ret |= 1 << 3
	}
	ret |= s.RetainHandling & 0x30

	return ret
}

// NewSubscribe creates a new Subscribe packet and applies all the
// provided/listed option functions to configure the packet
func NewSubscribe(opts ...func(c *Subscribe)) *Subscribe {
	s := &Subscribe{
		Subscriptions: make(map[string]SubOptions),
		Properties: &Properties{
			User: make(map[string]string),
		},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// SubscribeSingle is a Subscribe option function that adds a single subscription
// to the subscribe packet for the topic given and with the
// associated SubOptions
func SubscribeSingle(topic string, subOpts SubOptions) func(*Subscribe) {
	return func(s *Subscribe) {
		s.Subscriptions[topic] = subOpts
	}
}

// SubscribeMulti is a Subscribe option function that adds a map of string
// to SubOptions, the string keys of the map are topics to subscribe to
// the SubOption values are the subscription options to be applied to
// the subscription
func SubscribeMulti(subs map[string]SubOptions) func(*Subscribe) {
	return func(s *Subscribe) {
		for k, v := range subs {
			s.Subscriptions[k] = v
		}
	}
}

// SubscribeProperties is a Subscribe option function that sets
// the Properties for the Subscribe packet
func SubscribeProperties(p *Properties) func(*Subscribe) {
	return func(s *Subscribe) {
		s.Properties = p
	}
}

// Unpack is the implementation of the interface required function for a packet
func (s *Subscribe) Unpack(r *bytes.Buffer) error {
	var err error
	s.PacketID, err = readUint16(r)
	if err != nil {
		return err
	}

	err = s.Properties.Unpack(r, SUBSCRIBE)
	if err != nil {
		return err
	}

	return nil
}

// Buffers is the implementation of the interface required function for a packet
func (s *Subscribe) Buffers() net.Buffers {
	var b bytes.Buffer
	writeUint16(s.PacketID, &b)
	var subs bytes.Buffer
	for t, o := range s.Subscriptions {
		writeString(t, &subs)
		subs.WriteByte(o.Pack())
	}
	idvp := s.Properties.Pack(SUBSCRIBE)
	propLen := encodeVBI(len(idvp))
	return net.Buffers{b.Bytes(), propLen, idvp, subs.Bytes()}
}

// WriteTo is the implementation of the interface required function for a packet
func (s *Subscribe) WriteTo(w io.Writer) (int64, error) {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: SUBSCRIBE, Flags: 2}}
	cp.Content = s

	return cp.WriteTo(w)
}
