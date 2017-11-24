package packets

import (
	"bytes"
	"io"
	"net"
)

// Subscribe is the Variable Header definition for a Subscribe control packet
type Subscribe struct {
	PacketID      uint16
	Properties    Properties
	Subscriptions map[string]SubOptions
}

type SubOptions struct {
	QoS               byte
	NoLocal           bool
	RetainAsPublished bool
	RetainHandling    byte
}

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

func NewSubscribe(opts ...func(c *Subscribe)) *Subscribe {
	s := &Subscribe{
		Subscriptions: make(map[string]SubOptions),
		Properties: Properties{
			User: make(map[string]string),
		},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func Sub(topic string, subOpts SubOptions) func(*Subscribe) {
	return func(s *Subscribe) {
		s.Subscriptions[topic] = subOpts
	}
}

func MultiSub(subs map[string]SubOptions) func(*Subscribe) {
	return func(s *Subscribe) {
		for k, v := range subs {
			s.Subscriptions[k] = v
		}
	}
}

func SubscribeProperties(p *Properties) func(*Subscribe) {
	return func(s *Subscribe) {
		s.Properties = *p
	}
}

//Unpack is the implementation of the interface required function for a packet
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

func (s *Subscribe) Send(w io.Writer) error {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: SUBSCRIBE, Flags: 2}}
	cp.Content = s

	return cp.Send(w)
}
