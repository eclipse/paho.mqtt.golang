package packets

import (
	"bytes"
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"io"
)

//CONNACK packet

type ConnackPacket struct {
	FixedHeader
	TopicNameCompression byte
	ReturnCode           byte
	uuid                 uuid.UUID
}

func (ca *ConnackPacket) String() string {
	str := fmt.Sprintf("%s\n", ca.FixedHeader)
	str += fmt.Sprintf("returncode: %d", ca.ReturnCode)
	return str
}

func (ca *ConnackPacket) Write(w io.Writer) error {
	var body bytes.Buffer
	var err error

	body.WriteByte(ca.TopicNameCompression)
	body.WriteByte(ca.ReturnCode)
	ca.FixedHeader.RemainingLength = 2
	packet := ca.FixedHeader.pack()
	packet.Write(body.Bytes())
	_, err = packet.WriteTo(w)

	return err
}

func (ca *ConnackPacket) Unpack(b io.Reader) {
	ca.TopicNameCompression = decodeByte(b)
	ca.ReturnCode = decodeByte(b)
}

func (ca *ConnackPacket) Details() Details {
	return Details{Qos: 0, MessageID: 0}
}

func (ca *ConnackPacket) UUID() uuid.UUID {
	return ca.uuid
}
