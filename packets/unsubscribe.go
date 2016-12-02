package packets

import (
	"bufio"
	"bytes"
)

// Unsubscribe is the Variable Header definition for a Unsubscribe control packet
type Unsubscribe struct {
	packetID packetID
}

//Unpack is the implementation of the interface required function for a packet
func (u *Unsubscribe) Unpack(r bufio.Reader) (int, error) {
	return 0, nil
}

// Pack is the implementation of the interface required function for a packet
func (u *Unsubscribe) Pack(b bytes.Buffer) {
}
