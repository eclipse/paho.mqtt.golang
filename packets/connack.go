package packets

import (
	"bufio"
	"bytes"
)

// Connack is the Variable Header definition for a connack control packet
type Connack struct {
	AssignedClientID  assignedClientID
	ServerKeepAlive   serverKeepAlive
	AuthMethod        authMethod
	AuthData          authData
	ReplyInfo         replyInfo
	ServerReference   serverReference
	ReasonString      reasonString
	ReceiveMaximum    receiveMaximum
	TopicAliasMaximum topicAliasMaximum
	MaximumQOS        maximumQOS
}

//Unpack is the implementation of the interface required function for a packet
func (c *Connack) Unpack(r bufio.Reader) (int, error) {
	return 0, nil
}

// Pack is the implementation of the interface required function for a packet
func (c *Connack) Pack(b bytes.Buffer) {
}
