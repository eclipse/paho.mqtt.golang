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
	Properties *Properties
}

// PubcompSuccess, etc are the list of valid pubcomp reason codes.
const (
	PubcompSuccess                  = 0x00
	PubcompPacketIdentifierNotFound = 0x92
)

//Unpack is the implementation of the interface required function for a packet
func (p *Pubcomp) Unpack(r *bytes.Buffer) error {
	var err error
	success := r.Len() == 2
	noProps := r.Len() == 3
	p.PacketID, err = readUint16(r)
	if err != nil {
		return err
	}
	if !success {
		p.ReasonCode, err = r.ReadByte()
		if err != nil {
			return err
		}

		if !noProps {
			err = p.Properties.Unpack(r, PUBACK)
			if err != nil {
				return err
			}
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

// WriteTo is the implementation of the interface required function for a packet
func (p *Pubcomp) WriteTo(w io.Writer) (int64, error) {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: PUBCOMP}}
	cp.Content = p

	return cp.WriteTo(w)
}

// Reason returns a string representation of the meaning of the ReasonCode
func (p *Pubcomp) Reason() string {
	switch p.ReasonCode {
	case 0:
		return "Success - Packet Identifier released. Publication of QoS 2 message is complete."
	case 146:
		return "Packet Identifier not found - The Packet Identifier is not known. This is not an error during recovery, but at other times indicates a mismatch between the Session State on the Client and Server."
	}

	return ""
}

// NewPubcomp creates a new Pubcomp packet and applies all the
// provided/listed option functions to configure the packet
func NewPubcomp(opts ...func(p *Pubcomp)) *Pubcomp {
	p := &Pubcomp{
		Properties: &Properties{
			User: make(map[string]string),
		},
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// PubcompReason is a Pubcomp option function that sets the
// reason code for the Pubcomp packet
func PubcompReason(r byte) func(*Pubcomp) {
	return func(p *Pubcomp) {
		p.ReasonCode = r
	}
}

// PubcompFromPubrel reads the PacketID from the provided Pubrel
// and creates a Pubcomp packet with the same PacketID, this is used
// in the QoS2 flow
func PubcompFromPubrel(pr *Pubrel) func(*Pubcomp) {
	return func(pc *Pubcomp) {
		pc.PacketID = pr.PacketID
	}
}

// PubcompReasonCode is a Pubcomp option function that sets the
// reason code for the Pubcomp packet
func PubcompReasonCode(r byte) func(*Pubcomp) {
	return func(pc *Pubcomp) {
		pc.ReasonCode = r
	}
}

// PubcompProperties is a Pubcomp option function that sets
// the Properties for the Pubcomp packet
func PubcompProperties(p *Properties) func(*Pubcomp) {
	return func(pc *Pubcomp) {
		pc.Properties = p
	}
}
