package packets

import (
	"bufio"
	"bytes"
)

// Suback is the Variable Header definition for a Suback control packet
type Suback struct {
	packetID packetID
}

//Unpack is the implementation of the interface required function for a packet
func (s *Suback) Unpack(r bufio.Reader) (int, error) {
	return 0, nil
}

// Pack is the implementation of the interface required function for a packet
func (s *Suback) Pack(b bytes.Buffer) {
}
