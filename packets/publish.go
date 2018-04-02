package packets

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
)

// Publish is the Variable Header definition for a publish control packet
type Publish struct {
	Duplicate  bool
	QoS        byte
	Retain     bool
	Topic      string
	PacketID   uint16
	Properties Properties
	Payload    []byte
}

//Unpack is the implementation of the interface required function for a packet
func (p *Publish) Unpack(r *bytes.Buffer) error {
	var err error
	p.Topic, err = readString(r)
	if err != nil {
		return err
	}
	if p.QoS > 0 {
		p.PacketID, err = readUint16(r)
		if err != nil {
			return err
		}
	}

	err = p.Properties.Unpack(r, PUBLISH)
	if err != nil {
		return err
	}

	p.Payload, err = ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return nil
}

// Buffers is the implementation of the interface required function for a packet
func (p *Publish) Buffers() net.Buffers {
	var b bytes.Buffer
	writeString(p.Topic, &b)
	if p.QoS > 0 {
		writeUint16(p.PacketID, &b)
	}
	properties := p.Properties.Pack(PUBLISH)
	propLen := encodeVBI(len(properties))
	return net.Buffers{b.Bytes(), propLen, properties, p.Payload}

}

// Send is the implementation of the interface required function for a packet
func (p *Publish) Send(w io.Writer) error {
	f := p.QoS << 1
	if p.Duplicate {
		f |= 1 << 3
	}
	if p.Retain {
		f |= 1
	}

	cp := &ControlPacket{FixedHeader: FixedHeader{Type: PUBLISH, Flags: f}}
	cp.Content = p

	return cp.Send(w)
}

// NewPublish creates a new Publish packet and applies all the
// provided/listed option functions to configure the packet
func NewPublish(opts ...func(p *Publish)) *Publish {
	p := &Publish{
		Properties: Properties{
			User: make(map[string]string),
		},
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Message is a publish option function that sets the topic,
// qos, retain and payload for the Publish packet to the
// values provided
func Message(topic string, qos byte, retain bool, payload []byte) func(*Publish) {
	return func(p *Publish) {
		p.Topic = topic
		p.QoS = qos
		p.Retain = retain
		p.Payload = payload
	}
}

// PublishProperties is a Publish option function that sets
// the Properties for the Publish packet
func PublishProperties(p *Properties) func(*Publish) {
	return func(pp *Publish) {
		pp.Properties = *p
	}
}
