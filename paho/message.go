package paho

import (
	"fmt"

	p "github.com/eclipse/paho.mqtt.golang/packets"
)

type Message struct {
	Topic   string
	QoS     byte
	Retain  bool
	IDVP    p.IDValuePair
	Payload []byte
}

func (m *Message) String() string {
	ret := fmt.Sprintf("topic: %s  qos: %d  retain: %t\n", m.Topic, m.QoS, m.Retain)
	if m.IDVP.PayloadFormat != nil {
		ret += fmt.Sprintf("PayloadFormat: %v\n", m.IDVP.PayloadFormat)
	}
	if m.IDVP.PubExpiry != nil {
		ret += fmt.Sprintf("PubExpiry: %v\n", m.IDVP.PubExpiry)
	}
	if m.IDVP.ContentType != "" {
		ret += fmt.Sprintf("ContentType: %v\n", m.IDVP.ContentType)
	}
	if m.IDVP.ReplyTopic != "" {
		ret += fmt.Sprintf("ReplyTopic: %v\n", m.IDVP.ReplyTopic)
	}
	if m.IDVP.CorrelationData != nil {
		ret += fmt.Sprintf("CorrelationData: %v\n", m.IDVP.CorrelationData)
	}
	if m.IDVP.SubscriptionIdentifier != nil {
		ret += fmt.Sprintf("SubscriptionIdentifier: %v\n", m.IDVP.SubscriptionIdentifier)
	}
	if m.IDVP.SessionExpiryInterval != nil {
		ret += fmt.Sprintf("SessionExpiryInterval: %v\n", m.IDVP.SessionExpiryInterval)
	}
	if m.IDVP.AssignedClientID != "" {
		ret += fmt.Sprintf("AssignedClientID: %v\n", m.IDVP.AssignedClientID)
	}
	if m.IDVP.ServerKeepAlive != nil {
		ret += fmt.Sprintf("ServerKeepAlive: %v\n", m.IDVP.ServerKeepAlive)
	}
	if m.IDVP.AuthMethod != "" {
		ret += fmt.Sprintf("AuthMethod: %v\n", m.IDVP.AuthMethod)
	}
	if m.IDVP.AuthData != nil {
		ret += fmt.Sprintf("AuthData: %v\n", m.IDVP.AuthData)
	}
	if m.IDVP.RequestProblemInfo != nil {
		ret += fmt.Sprintf("RequestProblemInfo: %v\n", m.IDVP.RequestProblemInfo)
	}
	if m.IDVP.WillDelayInterval != nil {
		ret += fmt.Sprintf("WillDelayInterval: %v\n", m.IDVP.WillDelayInterval)
	}
	if m.IDVP.RequestResponseInfo != nil {
		ret += fmt.Sprintf("RequestResponseInfo: %v\n", m.IDVP.RequestResponseInfo)
	}
	if m.IDVP.ResponseInfo != "" {
		ret += fmt.Sprintf("ResponseInfo: %v\n", m.IDVP.ResponseInfo)
	}
	if m.IDVP.ServerReference != "" {
		ret += fmt.Sprintf("ServerReference: %v\n", m.IDVP.ServerReference)
	}
	if m.IDVP.ReasonString != "" {
		ret += fmt.Sprintf("ReasonString: %v\n", m.IDVP.ReasonString)
	}
	if m.IDVP.ReceiveMaximum != nil {
		ret += fmt.Sprintf("ReceiveMaximum: %v\n", m.IDVP.ReceiveMaximum)
	}
	if m.IDVP.TopicAliasMaximum != nil {
		ret += fmt.Sprintf("TopicAliasMaximum: %v\n", m.IDVP.TopicAliasMaximum)
	}
	if m.IDVP.TopicAlias != nil {
		ret += fmt.Sprintf("TopicAlias: %v\n", m.IDVP.TopicAlias)
	}
	if m.IDVP.MaximumQOS != nil {
		ret += fmt.Sprintf("MaximumQOS: %v\n", m.IDVP.MaximumQOS)
	}
	if m.IDVP.RetainAvailable != nil {
		ret += fmt.Sprintf("RetainAvailable: %v\n", m.IDVP.RetainAvailable)
	}
	if m.IDVP.MaximumPacketSize != nil {
		ret += fmt.Sprintf("MaximumPacketSize: %v\n", m.IDVP.MaximumPacketSize)
	}
	if m.IDVP.WildcardSubAvailable != nil {
		ret += fmt.Sprintf("WildcardSubAvailable: %v\n", m.IDVP.WildcardSubAvailable)
	}
	if m.IDVP.SubIDAvailable != nil {
		ret += fmt.Sprintf("SubIDAvailable: %v\n", m.IDVP.SubIDAvailable)
	}
	if m.IDVP.SharedSubAvailable != nil {
		ret += fmt.Sprintf("SharedSubAvailable: %v\n", m.IDVP.SharedSubAvailable)
	}
	for k, v := range m.IDVP.UserProperty {
		ret += fmt.Sprintf("UserProperty: %s : %s\n", k, v)
	}
	ret += string(m.Payload)

	return ret
}
