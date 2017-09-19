package packets

import (
	"bytes"
	"net"
)

// Unsuback is the Variable Header definition for a Unsuback control packet
type Unsuback struct {
	PacketID uint16
	IDVP     IDValuePair
}

//Unpack is the implementation of the interface required function for a packet
func (u *Unsuback) Unpack(r *bytes.Buffer) (int, error) {
	var err error
	u.PacketID, err = readUint16(r)
	if err != nil {
		return 0, err
	}

	idvpLen, err := u.IDVP.Unpack(r, UNSUBACK)
	if err != nil {
		return 0, err
	}

	return idvpLen + 2, nil
}

// Buffers is the implementation of the interface required function for a packet
func (u *Unsuback) Buffers() net.Buffers {
	var b bytes.Buffer
	writeUint16(u.PacketID, &b)
	idvp := u.IDVP.Pack(UNSUBACK)
	idvpLen := encodeVBI(len(idvp))
	return net.Buffers{b.Bytes(), idvpLen, idvp}
}
