package packets

import (
	"bytes"
	"net"
)

// Pingresp is the Variable Header definition for a Pingresp control packet
type Pingresp struct {
}

//Unpack is the implementation of the interface required function for a packet
func (p *Pingresp) Unpack(r *bytes.Buffer) (int, error) {
	return 0, nil
}

// Buffers is the implementation of the interface required function for a packet
func (p *Pingresp) Buffers() net.Buffers {
	return nil
}
