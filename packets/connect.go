package packets

import (
	"bytes"
	"net"
)

// Connect is the Variable Header definition for a connect control packet
type Connect struct {
	passwordFlag    bool
	usernameFlag    bool
	ProtocolName    string
	ProtocolVersion byte
	WillTopic       string
	WillRetain      bool
	WillQOS         byte
	WillFlag        bool
	WillMessage     []byte
	CleanStart      bool
	Username        string
	Password        []byte
	KeepAlive       uint16
	ClientID        string
	IDVP            IDValuePair
}

// PackFlags takes the Connect flags and packs them into the single byte
// representation used on the wire by MQTT
func (c *Connect) PackFlags() (f byte) {
	if c.Username != "" {
		f |= 0x01 << 7
	}
	if c.Password != nil {
		f |= 0x01 << 6
	}
	if c.WillFlag {
		f |= 0x01 << 2
		f |= c.WillQOS << 3
		if c.WillRetain {
			f |= 0x01 << 5
		}
	}
	if c.CleanStart {
		f |= 0x01 << 1
	}
	return
}

// UnpackFlags takes the wire byte representing the connect options flags
// and fills out the appropriate variables in the struct
func (c *Connect) UnpackFlags(b byte) {
	c.CleanStart = 1&(b>>1) > 0
	c.WillFlag = 1&(b>>2) > 0
	c.WillQOS = 3 & (b >> 3)
	c.WillRetain = 1&(b>>5) > 0
	c.passwordFlag = 1&(b>>6) > 0
	c.usernameFlag = 1&(b>>7) > 0
}

//Unpack is the implementation of the interface required function for a packet
func (c *Connect) Unpack(r *bytes.Buffer) (int, error) {
	var length int
	var err error

	if c.ProtocolName, err = readString(r); err != nil {
		return 0, err
	}
	length += len(c.ProtocolName) + 2

	if c.ProtocolVersion, err = r.ReadByte(); err != nil {
		return 0, err
	}

	flags, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	c.UnpackFlags(flags)

	if c.KeepAlive, err = readUint16(r); err != nil {
		return 0, err
	}
	length += 4

	idvpLen, err := c.IDVP.Unpack(r, CONNECT)
	length += idvpLen
	if err != nil {
		return 0, err
	}

	c.ClientID, err = readString(r)
	length += len(c.ClientID) + 2
	if err != nil {
		return 0, err
	}

	if c.WillFlag {
		c.WillTopic, err = readString(r)
		length += len(c.WillTopic) + 2
		if err != nil {
			return 0, err
		}
		c.WillMessage, err = readBinary(r)
		length += len(c.WillMessage) + 2
		if err != nil {
			return 0, err
		}
	}

	if c.usernameFlag {
		c.Username, err = readString(r)
		length += len(c.Username) + 2
		if err != nil {
			return 0, err
		}
	}

	if c.passwordFlag {
		c.Password, err = readBinary(r)
		length += len(c.Password) + 2
		if err != nil {
			return 0, err
		}
	}

	return length, nil
}

// Buffers is the implementation of the interface required function for a packet
func (c *Connect) Buffers() net.Buffers {
	var header bytes.Buffer
	var body bytes.Buffer
	writeString(c.ProtocolName, &header)
	header.WriteByte(c.ProtocolVersion)
	header.WriteByte(c.PackFlags())
	writeUint16(c.KeepAlive, &header)
	idvp := c.IDVP.Pack(CONNECT)
	idvpLen := encodeVBI(len(idvp))

	writeString(c.ClientID, &body)
	if c.WillFlag {
		writeString(c.WillTopic, &body)
		writeBinary(c.WillMessage, &body)
	}
	if c.Username != "" {
		writeString(c.Username, &body)
	}
	if c.Password != nil {
		writeBinary(c.Password, &body)
	}

	return net.Buffers{header.Bytes(), idvpLen, idvp, body.Bytes()}
}
