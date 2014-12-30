package packets

import (
	"bufio"
	"bytes"
	"code.google.com/p/go-uuid/uuid"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type ControlPacket interface {
	Write(io.Writer) error
	Unpack(io.Reader)
	String() string
	Details() Details
	UUID() uuid.UUID
}

var PacketNames = map[uint8]string{
	1:  "CONNECT",
	2:  "CONNACK",
	3:  "PUBLISH",
	4:  "PUBACK",
	5:  "PUBREC",
	6:  "PUBREL",
	7:  "PUBCOMP",
	8:  "SUBSCRIBE",
	9:  "SUBACK",
	10: "UNSUBSCRIBE",
	11: "UNSUBACK",
	12: "PINGREQ",
	13: "PINGRESP",
	14: "DISCONNECT",
}

const (
	Connect     = 1
	Connack     = 2
	Publish     = 3
	Puback      = 4
	Pubrec      = 5
	Pubrel      = 6
	Pubcomp     = 7
	Subscribe   = 8
	Suback      = 9
	Unsubscribe = 10
	Unsuback    = 11
	Pingreq     = 12
	Pingresp    = 13
	Disconnect  = 14
)

const (
	Accepted                     = 0x00
	RefusedBadProtocolVersion    = 0x01
	RefusedIDRejected            = 0x02
	RefusedServerUnavailable     = 0x03
	RefusedBadUsernameOrPassword = 0x04
	RefusedNotAuthorised         = 0x05
	NetworkError                 = 0xFE
	ProtocolViolation            = 0xFF
)

var ConnackReturnCodes = map[uint8]string{
	0:   "Connection Accepted",
	1:   "Connection Refused: Bad Protocol Version",
	2:   "Connection Refused: Client Identifier Rejected",
	3:   "Connection Refused: Server Unavailable",
	4:   "Connection Refused: Username or Password in unknown format",
	5:   "Connection Refused: Not Authorised",
	254: "Connection Error",
	255: "Connection Refused: Protocol Violation",
}

func ReadPacket(r io.Reader) (cp ControlPacket, err error) {
	var fh FixedHeader
	b := make([]byte, 1)

	_, err = io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}
	fh.unpack(b[0], r)
	cp = NewControlPacketWithHeader(fh)
	if cp == nil {
		return nil, errors.New("Bad data from client")
	}
	packetBytes := make([]byte, fh.RemainingLength)
	_, err = io.ReadFull(r, packetBytes)
	if err != nil {
		return nil, err
	}
	cp.Unpack(bytes.NewBuffer(packetBytes))
	return cp, nil
}

func NewControlPacket(packetType byte) (cp ControlPacket) {
	switch packetType {
	case Connect:
		cp = &ConnectPacket{FixedHeader: FixedHeader{MessageType: Connect}, uuid: uuid.NewUUID()}
	case Connack:
		cp = &ConnackPacket{FixedHeader: FixedHeader{MessageType: Connack}, uuid: uuid.NewUUID()}
	case Disconnect:
		cp = &DisconnectPacket{FixedHeader: FixedHeader{MessageType: Disconnect}, uuid: uuid.NewUUID()}
	case Publish:
		cp = &PublishPacket{FixedHeader: FixedHeader{MessageType: Publish}, uuid: uuid.NewUUID()}
	case Puback:
		cp = &PubackPacket{FixedHeader: FixedHeader{MessageType: Puback}, uuid: uuid.NewUUID()}
	case Pubrec:
		cp = &PubrecPacket{FixedHeader: FixedHeader{MessageType: Pubrec}, uuid: uuid.NewUUID()}
	case Pubrel:
		cp = &PubrelPacket{FixedHeader: FixedHeader{MessageType: Pubrel, Qos: 1}, uuid: uuid.NewUUID()}
	case Pubcomp:
		cp = &PubcompPacket{FixedHeader: FixedHeader{MessageType: Pubcomp}, uuid: uuid.NewUUID()}
	case Subscribe:
		cp = &SubscribePacket{FixedHeader: FixedHeader{MessageType: Subscribe, Qos: 1}, uuid: uuid.NewUUID()}
	case Suback:
		cp = &SubackPacket{FixedHeader: FixedHeader{MessageType: Suback}, uuid: uuid.NewUUID()}
	case Unsubscribe:
		cp = &UnsubscribePacket{FixedHeader: FixedHeader{MessageType: Unsubscribe}, uuid: uuid.NewUUID()}
	case Unsuback:
		cp = &UnsubackPacket{FixedHeader: FixedHeader{MessageType: Unsuback}, uuid: uuid.NewUUID()}
	case Pingreq:
		cp = &PingreqPacket{FixedHeader: FixedHeader{MessageType: Pingreq}, uuid: uuid.NewUUID()}
	case Pingresp:
		cp = &PingrespPacket{FixedHeader: FixedHeader{MessageType: Pingresp}, uuid: uuid.NewUUID()}
	default:
		return nil
	}
	return cp
}

func NewControlPacketWithHeader(fh FixedHeader) (cp ControlPacket) {
	switch fh.MessageType {
	case Connect:
		cp = &ConnectPacket{FixedHeader: fh, uuid: uuid.NewUUID()}
	case Connack:
		cp = &ConnackPacket{FixedHeader: fh, uuid: uuid.NewUUID()}
	case Disconnect:
		cp = &DisconnectPacket{FixedHeader: fh, uuid: uuid.NewUUID()}
	case Publish:
		cp = &PublishPacket{FixedHeader: fh, uuid: uuid.NewUUID()}
	case Puback:
		cp = &PubackPacket{FixedHeader: fh, uuid: uuid.NewUUID()}
	case Pubrec:
		cp = &PubrecPacket{FixedHeader: fh, uuid: uuid.NewUUID()}
	case Pubrel:
		cp = &PubrelPacket{FixedHeader: fh, uuid: uuid.NewUUID()}
	case Pubcomp:
		cp = &PubcompPacket{FixedHeader: fh, uuid: uuid.NewUUID()}
	case Subscribe:
		cp = &SubscribePacket{FixedHeader: fh, uuid: uuid.NewUUID()}
	case Suback:
		cp = &SubackPacket{FixedHeader: fh, uuid: uuid.NewUUID()}
	case Unsubscribe:
		cp = &UnsubscribePacket{FixedHeader: fh, uuid: uuid.NewUUID()}
	case Unsuback:
		cp = &UnsubackPacket{FixedHeader: fh, uuid: uuid.NewUUID()}
	case Pingreq:
		cp = &PingreqPacket{FixedHeader: fh, uuid: uuid.NewUUID()}
	case Pingresp:
		cp = &PingrespPacket{FixedHeader: fh, uuid: uuid.NewUUID()}
	default:
		return nil
	}
	return cp
}

type Details struct {
	Qos       byte
	MessageID uint16
}

type FixedHeader struct {
	MessageType     byte
	Dup             bool
	Qos             byte
	Retain          bool
	RemainingLength int
}

func (fh FixedHeader) String() string {
	return fmt.Sprintf("%s: dup: %t qos: %d retain: %t rLength: %d", PacketNames[fh.MessageType], fh.Dup, fh.Qos, fh.Retain, fh.RemainingLength)
}

func boolToByte(b bool) byte {
	switch b {
	case true:
		return 1
	default:
		return 0
	}
}

func (fh *FixedHeader) pack() bytes.Buffer {
	var header bytes.Buffer
	header.WriteByte(fh.MessageType<<4 | boolToByte(fh.Dup)<<3 | fh.Qos<<1 | boolToByte(fh.Retain))
	header.Write(encodeLength(fh.RemainingLength))
	return header
}

func (fh *FixedHeader) unpack(typeAndFlags byte, r io.Reader) {
	fh.MessageType = typeAndFlags >> 4
	fh.Dup = (typeAndFlags>>3)&0x01 > 0
	fh.Qos = (typeAndFlags >> 1) & 0x03
	fh.Retain = typeAndFlags&0x01 > 0
	fh.RemainingLength = decodeLength(r)
}

func decodeByte(b io.Reader) byte {
	num := make([]byte, 1)
	b.Read(num)
	return num[0]
}

func decodeUint16(b io.Reader) uint16 {
	num := make([]byte, 2)
	b.Read(num)
	return binary.BigEndian.Uint16(num)
}

func encodeUint16(num uint16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, num)
	return bytes
}

func encodeString(field string) []byte {
	fieldLength := make([]byte, 2)
	binary.BigEndian.PutUint16(fieldLength, uint16(len(field)))
	return append(fieldLength, []byte(field)...)
}

func decodeString(b io.Reader) string {
	fieldLength := decodeUint16(b)
	field := make([]byte, fieldLength)
	b.Read(field)
	return string(field)
}

func decodeBytes(b io.Reader) []byte {
	fieldLength := decodeUint16(b)
	field := make([]byte, fieldLength)
	b.Read(field)
	return field
}

func encodeBytes(field []byte) []byte {
	fieldLength := make([]byte, 2)
	binary.BigEndian.PutUint16(fieldLength, uint16(len(field)))
	return append(fieldLength, field...)
}

func encodeLength(length int) []byte {
	var encLength []byte
	for {
		digit := byte(length % 128)
		length /= 128
		if length > 0 {
			digit |= 0x80
		}
		encLength = append(encLength, digit)
		if length == 0 {
			break
		}
	}
	return encLength
}

func decodeLength(r io.Reader) int {
	var rLength uint32
	var multiplier uint32
	b := make([]byte, 1)
	for {
		io.ReadFull(r, b)
		digit := b[0]
		rLength |= uint32(digit&127) << multiplier
		if (digit & 128) == 0 {
			break
		}
		multiplier += 7
	}
	return int(rLength)
}

func DecodeLength(r *bufio.Reader) int {
	var rLength uint32
	var multiplier uint32 = 0
	//b := make([]byte, 1)
	for {
		//io.ReadFull(r, b)
		digit, _ := r.ReadByte()
		//digit := b[0]
		rLength |= uint32(digit&127) << multiplier
		if (digit & 128) == 0 {
			break
		}
		multiplier += 7
	}
	return int(rLength)
}
