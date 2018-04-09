package packets

import (
	"bytes"
	"io"
	"net"
)

// Unsuback is the Variable Header definition for a Unsuback control packet
type Unsuback struct {
	PacketID   uint16
	Properties Properties
	Reasons    []byte
}

// UnsubackSuccess, etc are the list of valid unsuback reason codes.
const (
	UnsubackSuccess                     = 0x00
	UnsubackNoSubscriptionFound         = 0x11
	UnsubackUnspecifiedError            = 0x80
	UnsubackImplementationSpecificError = 0x83
	UnsubackNotAuthorized               = 0x87
	UnsubackTopicFilterInvalid          = 0x8F
	UnsubackPacketIdentifierInUse       = 0x91
)

// Unpack is the implementation of the interface required function for a packet
func (u *Unsuback) Unpack(r *bytes.Buffer) error {
	var err error
	u.PacketID, err = readUint16(r)
	if err != nil {
		return err
	}

	err = u.Properties.Unpack(r, UNSUBACK)
	if err != nil {
		return err
	}

	u.Reasons = r.Bytes()

	return nil
}

// Buffers is the implementation of the interface required function for a packet
func (u *Unsuback) Buffers() net.Buffers {
	var b bytes.Buffer
	writeUint16(u.PacketID, &b)
	idvp := u.Properties.Pack(UNSUBACK)
	propLen := encodeVBI(len(idvp))
	return net.Buffers{b.Bytes(), propLen, idvp, u.Reasons}
}

// Send is the implementation of the interface required function for a packet
func (u *Unsuback) Send(w io.Writer) error {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: UNSUBACK}}
	cp.Content = u

	return cp.Send(w)
}

// Reason returns a string representation of the meaning of the ReasonCode
func (u *Unsuback) Reason(index int) string {
	if index >= 0 && index < len(u.Reasons) {
		switch u.Reasons[index] {
		case 0x00:
			return "Success - The subscription is deleted"
		case 0x11:
			return "No subscription found - No matching Topic Filter is being used by the Client."
		case 0x80:
			return "Unspecified error - The unsubscribe could not be completed and the Server either does not wish to reveal the reason or none of the other Reason Codes apply."
		case 0x83:
			return "Implementation specific error - The UNSUBSCRIBE is valid but the Server does not accept it."
		case 0x87:
			return "Not authorized - The Client is not authorized to unsubscribe."
		case 0x8F:
			return "Topic Filter invalid - The Topic Filter is correctly formed but is not allowed for this Client."
		case 0x91:
			return "Packet Identifier in use - The specified Packet Identifier is already in use."
		}
	}
	return "Invalid Reason index"
}

// NewUnsuback creates a new Unsuback packet and applies all the
// provided/listed option functions to configure the packet
func NewUnsuback(opts ...func(u *Unsuback)) *Unsuback {
	u := &Unsuback{
		Properties: Properties{
			User: make(map[string]string),
		},
	}

	for _, opt := range opts {
		opt(u)
	}

	return u
}

// UnsubackReasons is a Unsuback option function that sets the
// reason codes for the Unsuback packet
func UnsubackReasons(r []byte) func(*Unsuback) {
	return func(u *Unsuback) {
		u.Reasons = r
	}
}

// UnsubackProperties is a Unsuback option function that sets
// the Properties for the Unsuback packet
func UnsubackProperties(p *Properties) func(*Unsuback) {
	return func(u *Unsuback) {
		u.Properties = *p
	}
}
