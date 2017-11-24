package packets

import (
	"bytes"
	"io"
	"net"
)

// Auth is the Variable Header definition for a Auth control packet
type Auth struct {
	AuthReasonCode byte
	Properties     Properties
}

// Unpack is the implementation of the interface required function for a packet
func (a *Auth) Unpack(r *bytes.Buffer) error {
	var err error
	a.AuthReasonCode, err = r.ReadByte()
	if err != nil {
		return err
	}

	err = a.Properties.Unpack(r, AUTH)
	if err != nil {
		return err
	}

	return nil
}

// Buffers is the implementation of the interface required function for a packet
func (a *Auth) Buffers() net.Buffers {
	properties := a.Properties.Pack(AUTH)
	propLen := encodeVBI(len(properties))
	return net.Buffers{[]byte{a.AuthReasonCode}, propLen, properties}
}

func (a *Auth) Send(w io.Writer) error {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: AUTH}}
	cp.Content = a

	return cp.Send(w)
}
