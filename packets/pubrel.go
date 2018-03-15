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

func (p *Pubrel) Send(w io.Writer) error {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: PUBREL, Flags: 2}}
	cp.Content = p

	return cp.Send(w)
}

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

func PubrelFromPubrec(p *Pubrec) func(*Pubrel) {
	return func(pr *Pubrel) {
		pr.PacketID = p.PacketID
	}
}

func PubrelReasonCode(r byte) func(*Pubrel) {
	return func(pr *Pubrel) {
		pr.ReasonCode = r
	}
}

func PubrelProperties(p *Properties) func(*Pubrel) {
	return func(pr *Pubrel) {
		pr.Properties = *p
	}
}
