package packets

import (
	"bufio"
	"bytes"
)

// Pubcomp is the Variable Header definition for a Pubcomp control packet
type Pubcomp struct {
	packetID packetID
}

//Unpack is the implementation of the interface required function for a packet
func (p *Pubcomp) Unpack(r bufio.Reader) (int, error) {
	return 0, nil
}

// Pack is the implementation of the interface required function for a packet
func (p *Pubcomp) Pack(b bytes.Buffer) {
}
