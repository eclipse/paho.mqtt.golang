package packets

import (
	"bytes"
	"io"
	"net"
)

// Puback is the Variable Header definition for a Puback control packet
type Puback struct {
	PacketID   uint16
	ReasonCode byte
	Properties Properties
}

//Unpack is the implementation of the interface required function for a packet
func (p *Puback) Unpack(r *bytes.Buffer) error {
	var err error
	success := r.Len() == 2
	p.PacketID, err = readUint16(r)
	if err != nil {
		return err
	}
	if !success {
		p.ReasonCode, err = r.ReadByte()
		if err != nil {
			return err
		}

		err = p.Properties.Unpack(r, PUBACK)
		if err != nil {
			return err
		}
	}
	return nil
}

// Buffers is the implementation of the interface required function for a packet
func (p *Puback) Buffers() net.Buffers {
	var b bytes.Buffer
	writeUint16(p.PacketID, &b)
	b.WriteByte(p.ReasonCode)
	idvp := p.Properties.Pack(PUBACK)
	propLen := encodeVBI(len(idvp))
	return net.Buffers{b.Bytes(), propLen, idvp}
}

// Send is the implementation of the interface required function for a packet
func (p *Puback) Send(w io.Writer) error {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: PUBACK}}
	cp.Content = p

	return cp.Send(w)
}

// NewPuback creates a new Puback packet and applies all the
// provided/listed option functions to configure the packet
func NewPuback(opts ...func(p *Puback)) *Puback {
	p := &Puback{
		Properties: Properties{
			User: make(map[string]string),
		},
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// PubackFromPublish reads the PacketID from the provided Publish
// and creates a Puback packet with the same PacketID, this is used
// to respond to a QoS1 publish
func PubackFromPublish(pb *Publish) func(*Puback) {
	return func(pa *Puback) {
		pa.PacketID = pb.PacketID
	}
}
