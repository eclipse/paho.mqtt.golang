package paho

import (
	"fmt"

	p "github.com/eclipse/paho.mqtt.golang/packets"
)

// Message is a struct to represent the contents of a received
// Publish in an easier way to handle, these are passed to
// the MessageHandlers related to the subscriptions of the client
type Message struct {
	Topic      string
	QoS        byte
	Retain     bool
	Properties p.Properties
	Payload    []byte
}

// MessageFromPublish takes a Publish packet and returns a Message
func MessageFromPublish(pb *p.Publish) Message {
	return Message{
		Topic:      pb.Topic,
		QoS:        pb.QoS,
		Retain:     pb.Retain,
		Properties: pb.Properties,
		Payload:    pb.Payload,
	}
}

// String returns a string displaying the contents of a Message
func (m *Message) String() string {
	ret := fmt.Sprintf("topic: %s  qos: %d  retain: %t\n", m.Topic, m.QoS, m.Retain)
	if m.Properties.PayloadFormat != nil {
		ret += fmt.Sprintf("PayloadFormat: %v\n", m.Properties.PayloadFormat)
	}
	if m.Properties.PubExpiry != nil {
		ret += fmt.Sprintf("PubExpiry: %v\n", m.Properties.PubExpiry)
	}
	if m.Properties.ContentType != "" {
		ret += fmt.Sprintf("ContentType: %v\n", m.Properties.ContentType)
	}
	if m.Properties.ReplyTopic != "" {
		ret += fmt.Sprintf("ReplyTopic: %v\n", m.Properties.ReplyTopic)
	}
	if m.Properties.CorrelationData != nil {
		ret += fmt.Sprintf("CorrelationData: %v\n", m.Properties.CorrelationData)
	}
	if m.Properties.SubscriptionIdentifier != nil {
		ret += fmt.Sprintf("SubscriptionIdentifier: %v\n", m.Properties.SubscriptionIdentifier)
	}
	if m.Properties.SessionExpiryInterval != nil {
		ret += fmt.Sprintf("SessionExpiryInterval: %v\n", m.Properties.SessionExpiryInterval)
	}
	if m.Properties.AssignedClientID != "" {
		ret += fmt.Sprintf("AssignedClientID: %v\n", m.Properties.AssignedClientID)
	}
	if m.Properties.ServerKeepAlive != nil {
		ret += fmt.Sprintf("ServerKeepAlive: %v\n", m.Properties.ServerKeepAlive)
	}
	if m.Properties.AuthMethod != "" {
		ret += fmt.Sprintf("AuthMethod: %v\n", m.Properties.AuthMethod)
	}
	if m.Properties.AuthData != nil {
		ret += fmt.Sprintf("AuthData: %v\n", m.Properties.AuthData)
	}
	if m.Properties.RequestProblemInfo != nil {
		ret += fmt.Sprintf("RequestProblemInfo: %v\n", m.Properties.RequestProblemInfo)
	}
	if m.Properties.WillDelayInterval != nil {
		ret += fmt.Sprintf("WillDelayInterval: %v\n", m.Properties.WillDelayInterval)
	}
	if m.Properties.RequestResponseInfo != nil {
		ret += fmt.Sprintf("RequestResponseInfo: %v\n", m.Properties.RequestResponseInfo)
	}
	if m.Properties.ResponseInfo != "" {
		ret += fmt.Sprintf("ResponseInfo: %v\n", m.Properties.ResponseInfo)
	}
	if m.Properties.ServerReference != "" {
		ret += fmt.Sprintf("ServerReference: %v\n", m.Properties.ServerReference)
	}
	if m.Properties.ReasonString != "" {
		ret += fmt.Sprintf("ReasonString: %v\n", m.Properties.ReasonString)
	}
	if m.Properties.ReceiveMaximum != nil {
		ret += fmt.Sprintf("ReceiveMaximum: %v\n", m.Properties.ReceiveMaximum)
	}
	if m.Properties.TopicAliasMaximum != nil {
		ret += fmt.Sprintf("TopicAliasMaximum: %v\n", m.Properties.TopicAliasMaximum)
	}
	if m.Properties.TopicAlias != nil {
		ret += fmt.Sprintf("TopicAlias: %v\n", m.Properties.TopicAlias)
	}
	if m.Properties.MaximumQOS != nil {
		ret += fmt.Sprintf("MaximumQOS: %v\n", m.Properties.MaximumQOS)
	}
	if m.Properties.RetainAvailable != nil {
		ret += fmt.Sprintf("RetainAvailable: %v\n", m.Properties.RetainAvailable)
	}
	if m.Properties.MaximumPacketSize != nil {
		ret += fmt.Sprintf("MaximumPacketSize: %v\n", m.Properties.MaximumPacketSize)
	}
	if m.Properties.WildcardSubAvailable != nil {
		ret += fmt.Sprintf("WildcardSubAvailable: %v\n", m.Properties.WildcardSubAvailable)
	}
	if m.Properties.SubIDAvailable != nil {
		ret += fmt.Sprintf("SubIDAvailable: %v\n", m.Properties.SubIDAvailable)
	}
	if m.Properties.SharedSubAvailable != nil {
		ret += fmt.Sprintf("SharedSubAvailable: %v\n", m.Properties.SharedSubAvailable)
	}
	for k, v := range m.Properties.User {
		ret += fmt.Sprintf("User: %s : %s\n", k, v)
	}
	ret += string(m.Payload)

	return ret
}
