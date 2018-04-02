package packets

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

// PacketType is a type alias to byte representing the different
// MQTT control packet types
type PacketType byte

// The following consts are the packet type number for each of the
// different control packets in MQTT
const (
	_ PacketType = iota
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

// Packet is the interface defining the unique parts of a controlpacket
type Packet interface {
	Unpack(*bytes.Buffer) error
	Buffers() net.Buffers
	Send(io.Writer) error
}

// FixedHeader is the definition of a control packet fixed header
type FixedHeader struct {
	Flags           byte
	Type            PacketType
	remainingLength int
}

// Pack operates on a FixedHeader and takes the option values and produces
// the wire format byte that represents these.
func (f *FixedHeader) Pack() []byte {
	var b bytes.Buffer

	b.WriteByte(byte(f.Type)<<4 | f.Flags)
	b.Write(encodeVBI(f.remainingLength))

	return b.Bytes()
}

// ControlPacket is the definition of a control packet
type ControlPacket struct {
	FixedHeader
	Content Packet
}

// PacketID is a helper function that returns the value of the PacketID
// field from any kind of mqtt packet in the Content element
func (c *ControlPacket) PacketID() uint16 {
	switch r := c.Content.(type) {
	case *Publish:
		return r.PacketID
	case *Puback:
		return r.PacketID
	case *Pubrec:
		return r.PacketID
	case *Pubrel:
		return r.PacketID
	case *Pubcomp:
		return r.PacketID
	case *Subscribe:
		return r.PacketID
	case *Suback:
		return r.PacketID
	case *Unsubscribe:
		return r.PacketID
	case *Unsuback:
		return r.PacketID
	default:
		return 0
	}
}

// NewControlPacket takes a packetType and returns a pointer to a
// ControlPacket where the VariableHeader field is a pointer to an
// instance of a VariableHeader definition for that packetType
func NewControlPacket(t PacketType) *ControlPacket {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: t}}
	switch t {
	case CONNECT:
		cp.Content = &Connect{
			ProtocolName:    "MQTT",
			ProtocolVersion: 5,
			Properties:      Properties{User: make(map[string]string)},
		}
	case CONNACK:
		cp.Content = &Connack{Properties: Properties{User: make(map[string]string)}}
	case PUBLISH:
		cp.Content = &Publish{Properties: Properties{User: make(map[string]string)}}
	case PUBACK:
		cp.Content = &Puback{Properties: Properties{User: make(map[string]string)}}
	case PUBREC:
		cp.Content = &Pubrec{Properties: Properties{User: make(map[string]string)}}
	case PUBREL:
		cp.Flags = 2
		cp.Content = &Pubrel{Properties: Properties{User: make(map[string]string)}}
	case PUBCOMP:
		cp.Content = &Pubcomp{Properties: Properties{User: make(map[string]string)}}
	case SUBSCRIBE:
		cp.Flags = 2
		cp.Content = &Subscribe{
			Subscriptions: make(map[string]SubOptions),
			Properties:    Properties{User: make(map[string]string)},
		}
	case SUBACK:
		cp.Content = &Suback{Properties: Properties{User: make(map[string]string)}}
	case UNSUBSCRIBE:
		cp.Flags = 2
		cp.Content = &Unsubscribe{Properties: Properties{User: make(map[string]string)}}
	case UNSUBACK:
		cp.Content = &Unsuback{Properties: Properties{User: make(map[string]string)}}
	case PINGREQ:
		cp.Content = &Pingreq{}
	case PINGRESP:
		cp.Content = &Pingresp{}
	case DISCONNECT:
		cp.Content = &Disconnect{Properties: Properties{User: make(map[string]string)}}
	case AUTH:
		cp.Flags = 1
		cp.Content = &Auth{Properties: Properties{User: make(map[string]string)}}
	default:
		return nil
	}

	return cp
}

// ReadPacket reads a control packet from a io.Reader and returns a completed
// struct with the appropriate data
func ReadPacket(r io.Reader) (*ControlPacket, error) {
	t := make([]byte, 1)
	_, err := io.ReadFull(r, t)
	if err != nil {
		return nil, err
	}
	cp := NewControlPacket(PacketType(t[0] >> 4))
	if cp == nil {
		return nil, fmt.Errorf("Invalid packet type requested, %d", t[0]>>4)
	}
	cp.Flags = t[0] & 0xF
	if cp.Type == PUBLISH {
		cp.Content.(*Publish).QoS = (cp.Flags & 0x6) >> 1
	}
	vbi, err := getVBI(r)
	if err != nil {
		return nil, err
	}
	cp.remainingLength, err = decodeVBI(vbi)
	if err != nil {
		return nil, err
	}
	content := make([]byte, cp.remainingLength)
	n, err := io.ReadFull(r, content)
	if err != nil {
		return nil, err
	}
	if n != cp.remainingLength {
		return nil, fmt.Errorf("Failed to read packet, expected %d bytes, read %d", cp.remainingLength, n)
	}
	err = cp.Content.Unpack(bytes.NewBuffer(content))
	if err != nil {
		return nil, err
	}
	// payloadLength := cp.remainingLength - length
	// if payloadLength > 0 {
	// 	cp.Payload = content[length:]
	// }
	return cp, nil
}

// Send writes a packet to an io.Writer, handling packing all the parts of
// a control packet.
func (c *ControlPacket) Send(w io.Writer) error {
	var packet net.Buffers

	buffers := c.Content.Buffers()
	for _, b := range buffers {
		c.remainingLength += len(b)
	}

	packet = append(packet, c.FixedHeader.Pack())
	packet = append(packet, buffers...)

	_, err := packet.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}

func encodeVBI(length int) []byte {
	var x int
	b := make([]byte, 4)
	for {
		digit := byte(length % 128)
		length /= 128
		if length > 0 {
			digit |= 0x80
		}
		b[x] = digit
		x++
		if length == 0 {
			return b[:x]
		}
	}
}

func getVBI(r io.Reader) (*bytes.Buffer, error) {
	var ret bytes.Buffer
	digit := make([]byte, 1)
	for {
		_, err := io.ReadFull(r, digit)
		if err != nil {
			return nil, err
		}
		ret.WriteByte(digit[0])
		if digit[0] <= 0x7f {
			return &ret, nil
		}
	}
}

func decodeVBI(r *bytes.Buffer) (int, error) {
	var vbi uint32
	var multiplier uint32
	for {
		digit, err := r.ReadByte()
		if err != nil && err != io.EOF {
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

func writeUint16(u uint16, b *bytes.Buffer) error {
	if err := b.WriteByte(byte(u >> 8)); err != nil {
		return err
	}
	return b.WriteByte(byte(u))
}

func writeUint32(u uint32, b *bytes.Buffer) error {
	if err := b.WriteByte(byte(u >> 24)); err != nil {
		return err
	}
	if err := b.WriteByte(byte(u >> 16)); err != nil {
		return err
	}
	if err := b.WriteByte(byte(u >> 8)); err != nil {
		return err
	}
	return b.WriteByte(byte(u))
}

func writeString(s string, b *bytes.Buffer) {
	writeUint16(uint16(len(s)), b)
	b.WriteString(s)
}

func writeBinary(d []byte, b *bytes.Buffer) {
	writeUint16(uint16(len(d)), b)
	b.Write(d)
}

func readUint16(b *bytes.Buffer) (uint16, error) {
	b1, err := b.ReadByte()
	if err != nil {
		return 0, err
	}
	b2, err := b.ReadByte()
	if err != nil {
		return 0, err
	}
	return (uint16(b1) << 8) | uint16(b2), nil
}

func readUint32(b *bytes.Buffer) (uint32, error) {
	b1, err := b.ReadByte()
	if err != nil {
		return 0, err
	}
	b2, err := b.ReadByte()
	if err != nil {
		return 0, err
	}
	b3, err := b.ReadByte()
	if err != nil {
		return 0, err
	}
	b4, err := b.ReadByte()
	if err != nil {
		return 0, err
	}
	return (uint32(b1) << 24) | (uint32(b2) << 16) | (uint32(b3) << 8) | uint32(b4), nil
}

func readBinary(b *bytes.Buffer) ([]byte, error) {
	size, err := readUint16(b)
	if err != nil {
		return nil, err
	}
	s := make([]byte, size)
	if _, err := b.Read(s); err != nil {
		return nil, err
	}

	return s, nil
}

func readString(b *bytes.Buffer) (string, error) {
	s, err := readBinary(b)
	return string(s), err
}
