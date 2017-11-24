package packets

import (
	"bytes"
	"io"
	"net"
)

// Unsuback is the Variable Header definition for a Unsuback control packet
type Unsuback struct {
	PacketID   uint16
	Properties Properties
}

//Unpack is the implementation of the interface required function for a packet
func (u *Unsuback) Unpack(r *bytes.Buffer) error {
	var err error
	u.PacketID, err = readUint16(r)
	if err != nil {
		return err
	}

	err = u.Properties.Unpack(r, UNSUBACK)
	if err != nil {
		return err
	}

	return nil
}

// Buffers is the implementation of the interface required function for a packet
func (u *Unsuback) Buffers() net.Buffers {
	var b bytes.Buffer
	writeUint16(u.PacketID, &b)
	idvp := u.Properties.Pack(UNSUBACK)
	propLen := encodeVBI(len(idvp))
	return net.Buffers{b.Bytes(), propLen, idvp}
}

func (u *Unsuback) Send(w io.Writer) error {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: UNSUBACK}}
	cp.Content = u

	return cp.Send(w)
}
