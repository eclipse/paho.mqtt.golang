package packets

import (
	"bytes"
	"net"
)

// Subscribe is the Variable Header definition for a Subscribe control packet
type Subscribe struct {
	PacketID      uint16
	IDVP          IDValuePair
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

func NewSubscribe(subs map[string]SubOptions) *ControlPacket {
	s := NewControlPacket(SUBSCRIBE)
	s.Content.(*Subscribe).Subscriptions = subs

	return s
}

//Unpack is the implementation of the interface required function for a packet
func (s *Subscribe) Unpack(r *bytes.Buffer) (int, error) {
	var err error
	s.PacketID, err = readUint16(r)
	if err != nil {
		return 0, err
	}

	idvpLen, err := s.IDVP.Unpack(r, SUBSCRIBE)
	if err != nil {
		return 0, err
	}

	return idvpLen + 2, nil
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
	idvp := s.IDVP.Pack(SUBSCRIBE)
	idvpLen := encodeVBI(len(idvp))
	return net.Buffers{b.Bytes(), idvpLen, idvp, subs.Bytes()}
}
