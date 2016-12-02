package packets

import (
	"bufio"
	"bytes"
)

// Unsuback is the Variable Header definition for a Unsuback control packet
type Unsuback struct {
	packetID packetID
}

//Unpack is the implementation of the interface required function for a packet
func (u *Unsuback) Unpack(r bufio.Reader) (int, error) {
	return 0, nil
}

// Pack is the implementation of the interface required function for a packet
func (u *Unsuback) Pack(b bytes.Buffer) {
}
