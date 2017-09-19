package packets

import (
	"bytes"
	"net"
)

// Connack is the Variable Header definition for a connack control packet
type Connack struct {
	SessionPresent bool
	ReasonCode     byte
	IDVP           IDValuePair
}

//Unpack is the implementation of the interface required function for a packet
func (c *Connack) Unpack(r *bytes.Buffer) (int, error) {
	connackFlags, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	c.SessionPresent = connackFlags&0x01 > 0

	c.ReasonCode, err = r.ReadByte()
	if err != nil {
		return 0, err
	}

	idvpLen, err := c.IDVP.Unpack(r, CONNECT)
	if err != nil {
		return 0, err
	}

	return idvpLen + 2, nil
}

// Buffers is the implementation of the interface required function for a packet
func (c *Connack) Buffers() net.Buffers {
	return nil
}
