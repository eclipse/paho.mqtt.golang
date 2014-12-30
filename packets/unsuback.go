package packets

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"io"
)

//UNSUBACK packet

type UnsubackPacket struct {
	FixedHeader
	MessageID uint16
	uuid      uuid.UUID
}

func (ua *UnsubackPacket) String() string {
	str := fmt.Sprintf("%s\n", ua.FixedHeader)
	str += fmt.Sprintf("MessageID: %d", ua.MessageID)
	return str
}

func (ua *UnsubackPacket) Write(w io.Writer) error {
	var err error
	ua.FixedHeader.RemainingLength = 2
	packet := ua.FixedHeader.pack()
	packet.Write(encodeUint16(ua.MessageID))
	_, err = packet.WriteTo(w)

	return err
}

func (ua *UnsubackPacket) Unpack(b io.Reader) {
	ua.MessageID = decodeUint16(b)
}

func (ua *UnsubackPacket) Details() Details {
	return Details{Qos: 0, MessageID: ua.MessageID}
}

func (ua *UnsubackPacket) UUID() uuid.UUID {
	return ua.uuid
}
