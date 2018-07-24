package packets

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

// Puback is the Variable Header definition for a Puback control packet
type Puback struct {
	PacketID   uint16
	ReasonCode byte
	Properties *Properties
}

// PubackSuccess, etc are the list of valid puback reason codes.
const (
	PubackSuccess                     = 0x00
	PubackNoMatchingSubscribers       = 0x10
	PubackUnspecifiedError            = 0x80
	PubackImplementationSpecificError = 0x83
	PubackNotAuthorized               = 0x87
	PubackTopicNameInvalid            = 0x90
	PubackPacketIdentifierInUse       = 0x91
	PubackQuotaExceeded               = 0x97
	PubackPayloadFormatInvalid        = 0x99
)

//Unpack is the implementation of the interface required function for a packet
func (p *Puback) Unpack(r *bytes.Buffer) error {
	var err error
	success := r.Len() == 2
	noProps := r.Len() == 3
	fmt.Println("length", r.Len())
	p.PacketID, err = readUint16(r)
	if err != nil {
		fmt.Println("Error in readUint16")
		return err
	}
	if !success {
		p.ReasonCode, err = r.ReadByte()
		if err != nil {
			fmt.Println("Error at readbyte")
			return err
		}

		if !noProps {
			err = p.Properties.Unpack(r, PUBACK)
			if err != nil {
				fmt.Println("Error at properties unpack")
				return err
			}
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

// WriteTo is the implementation of the interface required function for a packet
func (p *Puback) WriteTo(w io.Writer) (int64, error) {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: PUBACK}}
	cp.Content = p

	return cp.WriteTo(w)
}

// Reason returns a string representation of the meaning of the ReasonCode
func (p *Puback) Reason() string {
	switch p.ReasonCode {
	case 0:
		return "The message is accepted. Publication of the QoS 1 message proceeds."
	case 16:
		return "The message is accepted but there are no subscribers. This is sent only by the Server. If the Server knows that there are no matching subscribers, it MAY use this Reason Code instead of 0x00 (Success)."
	case 128:
		return "The receiver does not accept the publish but either does not want to reveal the reason, or it does not match one of the other values."
	case 131:
		return "The PUBLISH is valid but the receiver is not willing to accept it."
	case 135:
		return "The PUBLISH is not authorized."
	case 144:
		return "The Topic Name is not malformed, but is not accepted by this Client or Server."
	case 145:
		return "The Packet Identifier is already in use. This might indicate a mismatch in the Session State between the Client and Server."
	case 151:
		return "An implementation or administrative imposed limit has been exceeded."
	case 153:
		return "The payload format does not match the specified Payload Format Indicator."
	}

	return ""
}

// NewPuback creates a new Puback packet and applies all the
// provided/listed option functions to configure the packet
func NewPuback(opts ...func(pa *Puback)) *Puback {
	pa := &Puback{
		Properties: &Properties{
			User: make(map[string]string),
		},
	}

	for _, opt := range opts {
		opt(pa)
	}

	return pa
}

// PubackReasonCode is a Puback option function that sets the
// reason code for the Puback packet
func PubackReasonCode(r byte) func(*Puback) {
	return func(pa *Puback) {
		pa.ReasonCode = r
	}
}

// PubackProperties is a Puback option function that sets
// the Properties for the Puback packet
func PubackProperties(p *Properties) func(*Puback) {
	return func(pa *Puback) {
		pa.Properties = p
	}
}

// PubackFromPublish reads the PacketID from the provided Publish
// and creates a Puback packet with the same PacketID, this is used
// to respond to a QoS1 publish
func PubackFromPublish(pb *Publish) func(*Puback) {
	return func(pa *Puback) {
		pa.PacketID = pb.PacketID
	}
}
