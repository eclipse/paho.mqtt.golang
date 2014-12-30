package packets

import (
	"bytes"
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"io"
)

//SUBSCRIBE packet

type SubscribePacket struct {
	FixedHeader
	MessageID uint16
	Topics    []string
	Qoss      []byte
	uuid      uuid.UUID
}

func (s *SubscribePacket) String() string {
	str := fmt.Sprintf("%s\n", s.FixedHeader)
	str += fmt.Sprintf("MessageID: %d topics: %s", s.MessageID, s.Topics)
	return str
}

func (s *SubscribePacket) Write(w io.Writer) error {
	var body bytes.Buffer
	var err error

	body.Write(encodeUint16(s.MessageID))
	for i, topic := range s.Topics {
		body.Write(encodeString(topic))
		body.WriteByte(s.Qoss[i])
	}
	s.FixedHeader.RemainingLength = body.Len()
	packet := s.FixedHeader.pack()
	packet.Write(body.Bytes())
	_, err = packet.WriteTo(w)

	return err
}

func (s *SubscribePacket) Unpack(b io.Reader) {
	s.MessageID = decodeUint16(b)
	payloadLength := s.FixedHeader.RemainingLength - 2
	for payloadLength > 0 {
		topic := decodeString(b)
		s.Topics = append(s.Topics, topic)
		qos := decodeByte(b)
		s.Qoss = append(s.Qoss, qos)
		payloadLength -= 2 + len(topic) + 1 //2 bytes of string length, plus string, plus 1 byte for Qos
	}
}

func (s *SubscribePacket) Details() Details {
	return Details{Qos: 1, MessageID: s.MessageID}
}

func (s *SubscribePacket) UUID() uuid.UUID {
	return s.uuid
}
