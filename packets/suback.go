package packets

import (
	"bytes"
	"net"
)

// Suback is the Variable Header definition for a Suback control packet
type Suback struct {
	PacketID uint16
	IDVP     IDValuePair
}

//Unpack is the implementation of the interface required function for a packet
func (s *Suback) Unpack(r *bytes.Buffer) (int, error) {
	var err error
	s.PacketID, err = readUint16(r)
	if err != nil {
		return 0, err
	}

	idvpLen, err := s.IDVP.Unpack(r, SUBACK)
	if err != nil {
		return 0, err
	}

	return idvpLen + 2, nil
}

// Buffers is the implementation of the interface required function for a packet
func (s *Suback) Buffers() net.Buffers {
	var b bytes.Buffer
	writeUint16(s.PacketID, &b)
	idvp := s.IDVP.Pack(SUBACK)
	idvpLen := encodeVBI(len(idvp))
	return net.Buffers{b.Bytes(), idvpLen, idvp}
}
