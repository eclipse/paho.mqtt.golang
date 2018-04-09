package packets

import (
	"bytes"
	"io"
	"net"
)

// Pingresp is the Variable Header definition for a Pingresp control packet
type Pingresp struct {
}

//Unpack is the implementation of the interface required function for a packet
func (p *Pingresp) Unpack(r *bytes.Buffer) error {
	return nil
}

// Buffers is the implementation of the interface required function for a packet
func (p *Pingresp) Buffers() net.Buffers {
	return nil
}

// Send is the implementation of the interface required function for a packet
func (p *Pingresp) Send(w io.Writer) error {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: PINGRESP}}
	cp.Content = p

	return cp.Send(w)
}
