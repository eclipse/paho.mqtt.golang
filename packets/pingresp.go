package packets

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"io"
)

//PINGRESP packet

type PingrespPacket struct {
	FixedHeader
	uuid uuid.UUID
}

func (pr *PingrespPacket) String() string {
	str := fmt.Sprintf("%s", pr.FixedHeader)
	return str
}

func (pr *PingrespPacket) Write(w io.Writer) error {
	packet := pr.FixedHeader.pack()
	_, err := packet.WriteTo(w)

	return err
}

func (pr *PingrespPacket) Unpack(b io.Reader) {
}

func (pr *PingrespPacket) Details() Details {
	return Details{Qos: 0, MessageID: 0}
}

func (pr *PingrespPacket) UUID() uuid.UUID {
	return pr.uuid
}
