package packets

import (
	"bufio"
	"bytes"
)

// Puback is the Variable Header definition for a Puback control packet
type Puback struct {
	packetID packetID
}

//Unpack is the implementation of the interface required function for a packet
func (p *Puback) Unpack(r bufio.Reader) (int, error) {
	return 0, nil
}

// Pack is the implementation of the interface required function for a packet
func (p *Puback) Pack(b bytes.Buffer) {
}
