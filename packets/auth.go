package packets

import (
	"bufio"
	"bytes"
)

// Auth is the Variable Header definition for a Auth control packet
type Auth struct {
	AuthMethod authMethod
	AuthData   authData
}

// Unpack is the implementation of the interface required function for a packet
func (a *Auth) Unpack(r bufio.Reader) (int, error) {
	return 0, nil
}

// Pack is the implementation of the interface required function for a packet
func (a *Auth) Pack(b bytes.Buffer) {
}
