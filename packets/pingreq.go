package packets

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"io"
)

//PINGREQ packet

type PingreqPacket struct {
	FixedHeader
	uuid uuid.UUID
}

func (pr *PingreqPacket) String() string {
	str := fmt.Sprintf("%s", pr.FixedHeader)
	return str
}

func (pr *PingreqPacket) Write(w io.Writer) error {
	packet := pr.FixedHeader.pack()
	_, err := packet.WriteTo(w)

	return err
}

func (pr *PingreqPacket) Unpack(b io.Reader) {
}

func (pr *PingreqPacket) Details() Details {
	return Details{Qos: 0, MessageID: 0}
}

func (pr *PingreqPacket) UUID() uuid.UUID {
	return pr.uuid
}
