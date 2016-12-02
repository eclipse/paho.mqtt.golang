package packets

import (
	"bufio"
	"bytes"
)

// Publish is the Variable Header definition for a publish control packet
type Publish struct {
	PacketID         packetID
	PayloadFormat    payloadFormat
	PubExpiry        pubExpiry
	ReplyTopic       replyTopic
	CorrelationData  correlationData
	TopicAlias       topicAlias
	UserDefinedPairs userDefinedPair
}

//Unpack is the implementation of the interface required function for a packet
func (p *Publish) Unpack(r bufio.Reader) (int, error) {
	return 0, nil
}

// Pack is the implementation of the interface required function for a packet
func (p *Publish) Pack(b bytes.Buffer) {
}
