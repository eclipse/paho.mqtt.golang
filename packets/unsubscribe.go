package packets

import (
	"bytes"
	"io"
	"net"
)

// Unsubscribe is the Variable Header definition for a Unsubscribe control packet
type Unsubscribe struct {
	PacketID   uint16
	Properties Properties
	Topics     []string
}

// Unpack is the implementation of the interface required function for a packet
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

// NewUnsubscribe creates a new Unsubscribe packet and applies all the
// provided/listed option functions to configure the packet
func NewUnsubscribe(opts ...func(u *Unsubscribe)) *Unsubscribe {
	u := &Unsubscribe{
		Properties: Properties{
			User: make(map[string]string),
		},
	}

	for _, opt := range opts {
		opt(u)
	}

	return u
}

// UnsubscribeTopics is an Unsubscribe option function that takes a
// slice of strings being the topics that should be unsubscribed from
func UnsubscribeTopics(topics []string) func(*Unsubscribe) {
	return func(u *Unsubscribe) {
		u.Topics = topics
	}
}

// UnsubscribeProperties is an Unsubscribe option function that sets
// the Properties for the Unsubscribe packet
func UnsubscribeProperties(p *Properties) func(*Unsubscribe) {
	return func(u *Unsubscribe) {
		u.Properties = *p
	}
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

// Send is the implementation of the interface required function for a packet
func (u *Unsubscribe) Send(w io.Writer) error {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: UNSUBSCRIBE}}
	cp.Content = u

	return cp.Send(w)
}
