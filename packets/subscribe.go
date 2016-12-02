package packets

import (
	"bufio"
	"bytes"
)

// Subscribe is the Variable Header definition for a Subscribe control packet
type Subscribe struct {
	packetID packetID
}

//Unpack is the implementation of the interface required function for a packet
func (s *Subscribe) Unpack(r bufio.Reader) (int, error) {
	return 0, nil
}

// Pack is the implementation of the interface required function for a packet
func (s *Subscribe) Pack(b bytes.Buffer) {
}
