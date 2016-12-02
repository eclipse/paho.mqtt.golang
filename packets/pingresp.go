package packets

import (
	"bufio"
	"bytes"
)

// Pingresp is the Variable Header definition for a Pingresp control packet
type Pingresp struct {
}

//Unpack is the implementation of the interface required function for a packet
func (p *Pingresp) Unpack(r bufio.Reader) (int, error) {
	return 0, nil
}

// Pack is the implementation of the interface required function for a packet
func (p *Pingresp) Pack(b bytes.Buffer) {
}
