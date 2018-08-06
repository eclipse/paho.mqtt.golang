package packets

import (
	"bytes"
	"io"
	"net"
)

// Pingreq is the Variable Header definition for a Pingreq control packet
type Pingreq struct {
}

//Unpack is the implementation of the interface required function for a packet
func (p *Pingreq) Unpack(r *bytes.Buffer) error {
	return nil
}

// Buffers is the implementation of the interface required function for a packet
func (p *Pingreq) Buffers() net.Buffers {
	return nil
}

// WriteTo is the implementation of the interface required function for a packet
func (p *Pingreq) WriteTo(w io.Writer) (int64, error) {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: PINGREQ}}
	cp.Content = p

	return cp.WriteTo(w)
}
