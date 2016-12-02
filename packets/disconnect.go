package packets

import (
	"bufio"
	"bytes"
)

// Disconnect is the Variable Header definition for a Disconnect control packet
type Disconnect struct {
	sessionExpiryInterval sessionExpiryInterval
	ServerReference       serverReference
	ReasonString          reasonString
}

//Unpack is the implementation of the interface required function for a packet
func (d *Disconnect) Unpack(r bufio.Reader) (int, error) {
	return 0, nil
}

// Pack is the implementation of the interface required function for a packet
func (d *Disconnect) Pack(b bytes.Buffer) {
}
