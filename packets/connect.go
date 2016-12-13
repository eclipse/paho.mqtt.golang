package packets

import (
	"bufio"
	"bytes"
)

// Connect is the Variable Header definition for a connect control packet
type Connect struct {
	ProtocolName              []byte
	ProtocolVersion           byte
	UsernameFlag              bool
	PasswordFlag              bool
	WillTopic                 string
	WillRetain                bool
	WillQOS                   byte
	WillFlag                  bool
	WillMessage               []byte
	CleanStart                bool
	Username                  string
	Password                  []byte
	KeepAlive                 uint16
	ClientID                  string
	IdvpLen                   int
	SessionExpiryIntervalFlag bool
	SessionExpiryInterval     sessionExpiryInterval
	WillDelayIntervalFlag     bool
	WillDelayInterval         willDelayInterval
	ReceiveMaximum            receiveMaximum
	TopicAliasMaximum         topicAliasMaximum
	RequestReplyInfo          requestReplyInfo
	RequestProblemInfo        requestProblemInfo
	UserDefinedPairs          userDefinedPair
	AuthMethod                authMethod
	AuthData                  authData
	MaximumQOS                maximumQOS
}

func (c *Connect) connectFlags() (f byte) {
	if c.UsernameFlag {
		f |= 0x01 << 7
	}
	if c.PasswordFlag {
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

func (c *Connect) packIDVP() bytes.Buffer {
	var b bytes.Buffer

	if c.SessionExpiryIntervalFlag {
		b.WriteByte(idvpSessionExpiryInterval)
		writeUint32(uint32(c.SessionExpiryInterval), b)
	}

	if c.WillDelayIntervalFlag {
		b.WriteByte(idvpWillDelayInterval)
		writeUint32(uint32(c.WillDelayInterval), b)
	}

	if c.ReceiveMaximum > 0 {
		b.WriteByte(idvpReceiveMaximum)
		writeUint16(uint16(c.ReceiveMaximum), b)
	}

	b.WriteByte(idvpTopicAliasMaximum)
	writeUint16(uint16(c.TopicAliasMaximum), b)

	b.WriteByte(idvpRequestReplyInfo)
	b.WriteByte(byte(c.RequestReplyInfo))

	b.WriteByte(idvpRequestProblemInfo)
	b.WriteByte(byte(c.RequestProblemInfo))

	if c.AuthMethod != "" {
		b.WriteByte(idvpAuthMethod)
		writeString(string(c.AuthMethod), b)
		b.WriteByte(idvpAuthData)
		b.Write(c.AuthData)
	}

	for k, v := range c.UserDefinedPairs {
		b.WriteByte(idvpUserDefinedPair)
		writeString(k, b)
		writeString(v, b)
	}

	return b
}

//Unpack is the implementation of the interface required function for a packet
func (c *Connect) Unpack(r bufio.Reader) (int, error) {
	return 0, nil
}

// Pack is the implementation of the interface required function for a packet
func (c *Connect) Pack(b bytes.Buffer) {
	b.Write(c.ProtocolName)
	b.WriteByte(c.ProtocolVersion)
	b.WriteByte(c.connectFlags())
	writeUint16(c.KeepAlive, b)
	idvp := c.packIDVP()
	encodeVBI(idvp.Len(), b)
	b.Write(idvp.Bytes())
	writeString(c.ClientID, b)
	if c.WillFlag {
		writeString(c.WillTopic, b)
		writeBinary(c.WillMessage, b)
	}
	if c.UsernameFlag {
		writeString(c.Username, b)
	}
	if c.PasswordFlag {
		writeBinary(c.Password, b)
	}
}
