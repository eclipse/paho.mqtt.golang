package packets

import (
	"bytes"
	"net"
)

// Pingreq is the Variable Header definition for a Pingreq control packet
type Pingreq struct {
}

//Unpack is the implementation of the interface required function for a packet
func (p *Pingreq) Unpack(r *bytes.Buffer) (int, error) {
	return 0, nil
}

// Buffers is the implementation of the interface required function for a packet
func (p *Pingreq) Buffers() net.Buffers {
	return nil
}
