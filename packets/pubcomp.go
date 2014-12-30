package packets

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"io"
)

//PUBCOMP packet

type PubcompPacket struct {
	FixedHeader
	MessageID uint16
	uuid      uuid.UUID
}

func (pc *PubcompPacket) String() string {
	str := fmt.Sprintf("%s\n", pc.FixedHeader)
	str += fmt.Sprintf("MessageID: %d", pc.MessageID)
	return str
}

func (pc *PubcompPacket) Write(w io.Writer) error {
	var err error
	pc.FixedHeader.RemainingLength = 2
	packet := pc.FixedHeader.pack()
	packet.Write(encodeUint16(pc.MessageID))
	_, err = packet.WriteTo(w)

	return err
}

func (pc *PubcompPacket) Unpack(b io.Reader) {
	pc.MessageID = decodeUint16(b)
}

func (pc *PubcompPacket) Details() Details {
	return Details{Qos: pc.Qos, MessageID: pc.MessageID}
}

func (pc *PubcompPacket) UUID() uuid.UUID {
	return pc.uuid
}
