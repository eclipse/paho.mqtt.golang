package packets

import (
	"bytes"
	"io"
	"net"
)

// Pubrel is the Variable Header definition for a Pubrel control packet
type Pubrel struct {
	PacketID   uint16
	ReasonCode byte
	Properties Properties
}

//Unpack is the implementation of the interface required function for a packet
func (p *Pubrel) Unpack(r *bytes.Buffer) error {
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
func (p *Pubrel) Buffers() net.Buffers {
	var b bytes.Buffer
	writeUint16(p.PacketID, &b)
	b.WriteByte(p.ReasonCode)
	idvp := p.Properties.Pack(PUBREL)
	propLen := encodeVBI(len(idvp))
	return net.Buffers{b.Bytes(), propLen, idvp}
}

// Send is the implementation of the interface required function for a packet
func (p *Pubrel) Send(w io.Writer) error {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: PUBREL, Flags: 2}}
	cp.Content = p

	return cp.Send(w)
}

// NewPubrel creates a new Pubrel packet and applies all the
// provided/listed option functions to configure the packet
func NewPubrel(opts ...func(p *Pubrel)) *Pubrel {
	p := &Pubrel{
		Properties: Properties{
			User: make(map[string]string),
		},
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// PubrelFromPubrec reads the PacketID from the provided Pubrec
// and creates a Pubrel packet with the same PacketID, this is used
// in the QoS2 flow
func PubrelFromPubrec(p *Pubrec) func(*Pubrel) {
	return func(pr *Pubrel) {
		pr.PacketID = p.PacketID
	}
}

// PubrelReasonCode is a Pubrel option function that sets the
// reason code for the Pubrel packet
func PubrelReasonCode(r byte) func(*Pubrel) {
	return func(pr *Pubrel) {
		pr.ReasonCode = r
	}
}

// PubrelProperties is a Pubrel option function that sets
// the Properties for the Pubrel packet
func PubrelProperties(p *Properties) func(*Pubrel) {
	return func(pr *Pubrel) {
		pr.Properties = *p
	}
}
