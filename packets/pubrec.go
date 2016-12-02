package packets

import (
	"bufio"
	"bytes"
)

// Pubrec is the Variable Header definition for a Pubrec control packet
type Pubrec struct {
	packetID packetID
}

//Unpack is the implementation of the interface required function for a packet
func (p *Pubrec) Unpack(r bufio.Reader) (int, error) {
	return 0, nil
}

// Pack is the implementation of the interface required function for a packet
func (p *Pubrec) Pack(b bytes.Buffer) {
}
