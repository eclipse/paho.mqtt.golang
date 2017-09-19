package packets

import (
	"bytes"
	"net"
)

// Subscribe is the Variable Header definition for a Subscribe control packet
type Subscribe struct {
	PacketID      uint16
	IDVP          IDValuePair
	Subscriptions map[string]byte
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
		subs.WriteByte(o)
	}
	idvp := s.IDVP.Pack(SUBSCRIBE)
	idvpLen := encodeVBI(len(idvp))
	return net.Buffers{b.Bytes(), idvpLen, idvp, subs.Bytes()}
}
