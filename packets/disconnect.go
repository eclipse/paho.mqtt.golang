package packets

import (
	"bytes"
	"net"
)

// Disconnect is the Variable Header definition for a Disconnect control packet
type Disconnect struct {
	DisconnectReasonCode byte
	IDVP                 IDValuePair
}

//Unpack is the implementation of the interface required function for a packet
func (d *Disconnect) Unpack(r *bytes.Buffer) (int, error) {
	var err error
	d.DisconnectReasonCode, err = r.ReadByte()
	if err != nil {
		return 0, err
	}

	idvpLen, err := d.IDVP.Unpack(r, DISCONNECT)
	if err != nil {
		return 0, err
	}

	return idvpLen + 1, nil
}

// Buffers is the implementation of the interface required function for a packet
func (d *Disconnect) Buffers() net.Buffers {
	idvp := d.IDVP.Pack(DISCONNECT)
	idvpLen := encodeVBI(len(idvp))
	return net.Buffers{[]byte{d.DisconnectReasonCode}, idvpLen, idvp}
}
