package packets

import (
	"bufio"
	"bytes"
)

// Pingreq is the Variable Header definition for a Pingreq control packet
type Pingreq struct {
}

//Unpack is the implementation of the interface required function for a packet
func (p *Pingreq) Unpack(r bufio.Reader) (int, error) {
	return 0, nil
}

// Pack is the implementation of the interface required function for a packet
func (p *Pingreq) Pack(b bytes.Buffer) {
}
