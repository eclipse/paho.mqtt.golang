package packets

import (
	"bytes"
	"io"
	"net"
)

// Disconnect is the Variable Header definition for a Disconnect control packet
type Disconnect struct {
	DisconnectReasonCode byte
	Properties           Properties
}

func NewDisconnect(opts ...func(*Disconnect)) *Disconnect {
	d := &Disconnect{}

	for _, opt := range opts {
		opt(d)
	}

	return d
}

func DisconnectReason(r byte) func(*Disconnect) {
	return func(d *Disconnect) {
		d.DisconnectReasonCode = r
	}
}

func DisconnectProperties(p Properties) func(*Disconnect) {
	return func(d *Disconnect) {
		d.Properties = p
	}
}

//Unpack is the implementation of the interface required function for a packet
func (d *Disconnect) Unpack(r *bytes.Buffer) error {
	var err error
	d.DisconnectReasonCode, err = r.ReadByte()
	if err != nil {
		return err
	}

	err = d.Properties.Unpack(r, DISCONNECT)
	if err != nil {
		return err
	}

	return nil
}

// Buffers is the implementation of the interface required function for a packet
func (d *Disconnect) Buffers() net.Buffers {
	idvp := d.Properties.Pack(DISCONNECT)
	propLen := encodeVBI(len(idvp))
	return net.Buffers{[]byte{d.DisconnectReasonCode}, propLen, idvp}
}

func (d *Disconnect) Send(w io.Writer) error {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: DISCONNECT}}
	cp.Content = d

	return cp.Send(w)
}
