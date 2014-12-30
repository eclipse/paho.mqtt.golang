package packets

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"io"
)

//DISCONNECT packet

type DisconnectPacket struct {
	FixedHeader
	uuid uuid.UUID
}

func (d *DisconnectPacket) String() string {
	str := fmt.Sprintf("%s\n", d.FixedHeader)
	return str
}

func (d *DisconnectPacket) Write(w io.Writer) error {
	packet := d.FixedHeader.pack()
	_, err := packet.WriteTo(w)

	return err
}

func (d *DisconnectPacket) Unpack(b io.Reader) {
}

func (d *DisconnectPacket) Details() Details {
	return Details{Qos: 0, MessageID: 0}
}

func (d *DisconnectPacket) UUID() uuid.UUID {
	return d.uuid
}
