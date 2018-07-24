package packets

import (
	"bytes"
	"io"
	"net"
)

// Pubrec is the Variable Header definition for a Pubrec control packet
type Pubrec struct {
	PacketID   uint16
	ReasonCode byte
	Properties *Properties
}

// PubrecSuccess, etc are the list of valid Pubrec reason codes
const (
	PubrecSuccess                     = 0x00
	PubrecNoMatchingSubscribers       = 0x10
	PubrecUnspecifiedError            = 0x80
	PubrecImplementationSpecificError = 0x83
	PubrecNotAuthorized               = 0x87
	PubrecTopicNameInvalid            = 0x90
	PubrecPacketIdentifierInUse       = 0x91
	PubrecQuotaExceeded               = 0x97
	PubrecPayloadFormatInvalid        = 0x99
)

//Unpack is the implementation of the interface required function for a packet
func (p *Pubrec) Unpack(r *bytes.Buffer) error {
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
func (p *Pubrec) Buffers() net.Buffers {
	var b bytes.Buffer
	writeUint16(p.PacketID, &b)
	b.WriteByte(p.ReasonCode)
	idvp := p.Properties.Pack(PUBACK)
	propLen := encodeVBI(len(idvp))
	return net.Buffers{b.Bytes(), propLen, idvp}
}

// WriteTo is the implementation of the interface required function for a packet
func (p *Pubrec) WriteTo(w io.Writer) (int64, error) {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: PUBREC}}
	cp.Content = p

	return cp.WriteTo(w)
}

// Reason returns a string representation of the meaning of the ReasonCode
func (p *Pubrec) Reason() string {
	switch p.ReasonCode {
	case 0:
		return "Success - The message is accepted. Publication of the QoS 2 message proceeds."
	case 16:
		return "No matching subscribers. - The message is accepted but there are no subscribers. This is sent only by the Server. If the Server knows that case there are no matching subscribers, it MAY use this Reason Code instead of 0x00 (Success)"
	case 128:
		return "Unspecified error - The receiver does not accept the publish but either does not want to reveal the reason, or it does not match one of the other values."
	case 131:
		return "Implementation specific error - The PUBLISH is valid but the receiver is not willing to accept it."
	case 135:
		return "Not authorized - The PUBLISH is not authorized."
	case 144:
		return "Topic Name invalid - The Topic Name is not malformed, but is not accepted by this Client or Server."
	case 145:
		return "Packet Identifier in use - The Packet Identifier is already in use. This might indicate a mismatch in the Session State between the Client and Server."
	case 151:
		return "Quota exceeded - An implementation or administrative imposed limit has been exceeded."
	case 153:
		return "Payload format invalid - The payload format does not match the one specified in the Payload Format Indicator."
	}

	return ""
}

// NewPubrec creates a new Pubrec packet and applies all the
// provided/listed option functions to configure the packet
func NewPubrec(opts ...func(p *Pubrec)) *Pubrec {
	p := &Pubrec{
		Properties: &Properties{
			User: make(map[string]string),
		},
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// PubrecReason is a Pubrec option function that sets the
// reason code for the Pubrec packet
func PubrecReason(r byte) func(*Pubrec) {
	return func(p *Pubrec) {
		p.ReasonCode = r
	}
}

// PubrecFromPublish reads the PacketID from the provided Publish
// and creates a Pubrec packet with the same PacketID, this is used
// in the QoS2 flow
func PubrecFromPublish(pb *Publish) func(*Pubrec) {
	return func(pr *Pubrec) {
		pr.PacketID = pb.PacketID
	}
}
