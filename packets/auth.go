package packets

import (
	"bytes"
	"io"
	"net"
)

// Auth is the Variable Header definition for a Auth control packet
type Auth struct {
	ReasonCode byte
	Properties *Properties
}

// Unpack is the implementation of the interface required function for a packet
func (a *Auth) Unpack(r *bytes.Buffer) error {
	var err error

	success := r.Len() == 0
	noProps := r.Len() == 1
	if !success {
		a.ReasonCode, err = r.ReadByte()
		if err != nil {
			return err
		}

		if !noProps {
			err = a.Properties.Unpack(r, AUTH)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Buffers is the implementation of the interface required function for a packet
func (a *Auth) Buffers() net.Buffers {
	properties := a.Properties.Pack(AUTH)
	propLen := encodeVBI(len(properties))
	return net.Buffers{[]byte{a.ReasonCode}, propLen, properties}
}

// WriteTo is the implementation of the interface required function for a packet
func (a *Auth) WriteTo(w io.Writer) (int64, error) {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: AUTH}}
	cp.Content = a

	return cp.WriteTo(w)
}
