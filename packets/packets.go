package packets

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type packetType byte
type packetID uint16

const (
	_ packetType = iota
	CONNECT
	CONNACK
	PUBLISH
	PUBACK
	PUBREC
	PUBREL
	PUBCOMP
	SUBSCRIBE
	SUBACK
	UNSUBSCRIBE
	UNSUBACK
	PINGREQ
	PINGRESP
	DISCONNECT
	AUTH
)

type packet interface {
	Unpack(bufio.Reader) (int, error)
	Pack(bytes.Buffer)
}

// FixedHeader is the definition of a control packet fixed header
type FixedHeader struct {
	cpType          packetType
	flags           byte
	remainingLength int
}

// ControlPacket is the definition of a control packet
type ControlPacket struct {
	FixedHeader
	VariableHeader packet
	Payload        []byte
}

// NewControlPacket takes a packetType and returns a pointer to a
// ControlPacket where the VariableHeader field is a pointer to an
// instance of a VariableHeader definition for that packetType
func NewControlPacket(t packetType) *ControlPacket {
	cp := &ControlPacket{FixedHeader: FixedHeader{cpType: t}}
	switch t {
	case CONNECT:
		cp.VariableHeader = &Connect{
			ProtocolName:    []byte{0x00, 0x04, 'M', 'Q', 'T', 'T'},
			ProtocolVersion: 5}
	case CONNACK:
		cp.VariableHeader = &Connack{}
	case PUBLISH:
		cp.VariableHeader = &Publish{}
	case PUBACK:
		cp.VariableHeader = &Puback{}
	case PUBREC:
		cp.VariableHeader = &Pubrec{}
	case PUBREL:
		cp.flags = 2
		cp.VariableHeader = &Pubrel{}
	case PUBCOMP:
		cp.VariableHeader = &Pubcomp{}
	case SUBSCRIBE:
		cp.flags = 2
		cp.VariableHeader = &Subscribe{}
	case SUBACK:
		cp.VariableHeader = &Suback{}
	case UNSUBSCRIBE:
		cp.flags = 2
		cp.VariableHeader = &Unsubscribe{}
	case UNSUBACK:
		cp.VariableHeader = &Unsuback{}
	case PINGREQ:
		cp.VariableHeader = &Pingreq{}
	case PINGRESP:
		cp.VariableHeader = &Pingresp{}
	case DISCONNECT:
		cp.VariableHeader = &Disconnect{}
	case AUTH:
		cp.flags = 1
		cp.VariableHeader = &Auth{}
	default:
		return nil
	}

	return cp
}

// ReadPacket reads a control packet from a bufio.Reader and returns a completed
// struct with the appropriate data
func ReadPacket(r bufio.Reader) (*ControlPacket, error) {
	t, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	cp := NewControlPacket(packetType(t >> 4))
	if cp == nil {
		return nil, fmt.Errorf("Invalid packet type requested, %d", t)
	}
	cp.flags = t & 0xF
	cp.remainingLength, err = decodeVBI(r)
	if err != nil {
		return nil, err
	}
	length, err := cp.VariableHeader.Unpack(r)
	if err != nil {
		return nil, err
	}
	payloadLength := cp.remainingLength - length
	if payloadLength > 0 {
		cp.Payload = make([]byte, payloadLength)
		n, err := r.Read(cp.Payload)
		if err != nil {
			return nil, err
		}
		if n != payloadLength {
			return nil, fmt.Errorf("Failed to read payload, expected %d bytes, read %d", payloadLength, n)
		}
	}

	return cp, nil
}

// Send writes a packet to an io.Writer, handling packing all the parts of
// a control packet.
func (c *ControlPacket) Send(w io.Writer) error {
	var b bytes.Buffer

	b.WriteByte(byte(c.cpType)<<4 | c.flags)
	encodeVBI(c.remainingLength, b)
	c.VariableHeader.Pack(b)
	b.Write(c.Payload)

	_, err := b.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}

func encodeVBI(length int, b bytes.Buffer) {
	for {
		digit := byte(length % 128)
		length /= 128
		if length > 0 {
			digit |= 0x80
		}
		b.WriteByte(digit)
		if length == 0 {
			break
		}
	}
}

func decodeVBI(r bufio.Reader) (int, error) {
	var vbi uint32
	var multiplier uint32
	for {
		digit, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		vbi |= uint32(digit&127) << multiplier
		if (digit & 128) == 0 {
			break
		}
		multiplier += 7
	}
	return int(vbi), nil
}

func writeUint16(u uint16, b bytes.Buffer) {
	b.WriteByte(byte(u >> 8))
	b.WriteByte(byte(u))
}

func writeUint32(u uint32, b bytes.Buffer) {
	b.WriteByte(byte(u >> 24))
	b.WriteByte(byte(u >> 16))
	b.WriteByte(byte(u >> 8))
	b.WriteByte(byte(u))
}

func writeString(s string, b bytes.Buffer) {
	writeUint16(uint16(len(s)), b)
	b.WriteString(s)
}

func writeBinary(d []byte, b bytes.Buffer) {
	writeUint16(uint16(len(d)), b)
	b.Write(d)
}
