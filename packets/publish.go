package packets

import (
	"bytes"
	"net"
)

// Publish is the Variable Header definition for a publish control packet
type Publish struct {
	Duplicate bool
	QoS       byte
	Retain    bool
	Topic     string
	PacketID  uint16
	IDVP      IDValuePair
}

//Unpack is the implementation of the interface required function for a packet
func (p *Publish) Unpack(r *bytes.Buffer) (int, error) {
	var err error
	p.Topic, err = readString(r)
	if err != nil {
		return 0, err
	}
	p.PacketID, err = readUint16(r)
	if err != nil {
		return 0, err
	}

	idvpLen, err := p.IDVP.Unpack(r, PUBLISH)
	if err != nil {
		return 0, err
	}

	return idvpLen + 4 + len(p.Topic), nil
}

// Buffers is the implementation of the interface required function for a packet
func (p *Publish) Buffers() net.Buffers {
	var b bytes.Buffer
	writeString(p.Topic, &b)
	writeUint16(p.PacketID, &b)
	idvp := p.IDVP.Pack(PUBLISH)
	idvpLen := encodeVBI(len(idvp))
	return net.Buffers{b.Bytes(), idvpLen, idvp}

}
