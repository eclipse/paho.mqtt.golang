package packets

import (
	"bytes"
	"net"
)

// Auth is the Variable Header definition for a Auth control packet
type Auth struct {
	AuthReasonCode byte
	IDVP           IDValuePair
}

// Unpack is the implementation of the interface required function for a packet
func (a *Auth) Unpack(r *bytes.Buffer) (int, error) {
	var err error
	a.AuthReasonCode, err = r.ReadByte()
	if err != nil {
		return 0, err
	}

	idvpLen, err := a.IDVP.Unpack(r, AUTH)
	if err != nil {
		return 0, err
	}

	return idvpLen + 1, nil
}

// Buffers is the implementation of the interface required function for a packet
func (a *Auth) Buffers() net.Buffers {
	idvp := a.IDVP.Pack(AUTH)
	idvpLen := encodeVBI(len(idvp))
	return net.Buffers{[]byte{a.AuthReasonCode}, idvpLen, idvp}
}
