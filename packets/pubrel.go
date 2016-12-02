package packets

import (
	"bufio"
	"bytes"
)

// Pubrel is the Variable Header definition for a Pubrel control packet
type Pubrel struct {
	packetID packetID
}

//Unpack is the implementation of the interface required function for a packet
func (p *Pubrel) Unpack(r bufio.Reader) (int, error) {
	return 0, nil
}

// Pack is the implementation of the interface required function for a packet
func (p *Pubrel) Pack(b bytes.Buffer) {
}
