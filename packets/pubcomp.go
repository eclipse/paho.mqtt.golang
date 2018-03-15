package packets

import (
	"bytes"
	"io"
	"net"
)

// Pubcomp is the Variable Header definition for a Pubcomp control packet
type Pubcomp struct {
	PacketID   uint16
	ReasonCode byte
	Properties Properties
}

//Unpack is the implementation of the interface required function for a packet
func (p *Pubcomp) Unpack(r *bytes.Buffer) error {
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
func (p *Pubcomp) Buffers() net.Buffers {
	var b bytes.Buffer
	writeUint16(p.PacketID, &b)
	b.WriteByte(p.ReasonCode)
	idvp := p.Properties.Pack(PUBCOMP)
	propLen := encodeVBI(len(idvp))
	return net.Buffers{b.Bytes(), propLen, idvp}
}

func (p *Pubcomp) Send(w io.Writer) error {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: PUBCOMP}}
	cp.Content = p

	return cp.Send(w)
}

func NewPubcomp(opts ...func(p *Pubcomp)) *Pubcomp {
	p := &Pubcomp{
		Properties: Properties{
			User: make(map[string]string),
		},
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func PubcompFromPubrel(pr *Pubrel) func(*Pubcomp) {
	return func(pc *Pubcomp) {
		pc.PacketID = pr.PacketID
	}
}

func PubcompReasonCode(r byte) func(*Pubcomp) {
	return func(pc *Pubcomp) {
		pc.ReasonCode = r
	}
}

func PubcompProperties(p *Properties) func(*Pubcomp) {
	return func(pc *Pubcomp) {
		pc.Properties = *p
	}
}
