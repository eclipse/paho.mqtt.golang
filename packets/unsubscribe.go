package packets

import (
	"bytes"
	"io"
	"net"
)

// Unsubscribe is the Variable Header definition for a Unsubscribe control packet
type Unsubscribe struct {
	PacketID uint16
	Topics   []string
}

//Unpack is the implementation of the interface required function for a packet
func (u *Unsubscribe) Unpack(r *bytes.Buffer) error {
	for {
		t, err := readString(r)
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}
		u.Topics = append(u.Topics, t)
	}

	return nil
}

// Buffers is the implementation of the interface required function for a packet
func (u *Unsubscribe) Buffers() net.Buffers {
	var b bytes.Buffer
	writeUint16(u.PacketID, &b)
	for _, t := range u.Topics {
		writeString(t, &b)
	}
	return net.Buffers{b.Bytes()}
}

func (u *Unsubscribe) Send(w io.Writer) error {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: UNSUBSCRIBE}}
	cp.Content = u

	return cp.Send(w)
}
