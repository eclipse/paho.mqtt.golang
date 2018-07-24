package paho

import (
	"bytes"
	"fmt"

	"github.com/eclipse/paho.mqtt.golang/packets"
)

// Publish is a reporesentation of the MQTT Publish packet
type Publish struct {
	QoS        byte
	Retain     bool
	Topic      string
	Properties *PublishProperties
	Payload    []byte
}

// PublishProperties is a struct of the properties that can be set
// for a publish packet
type PublishProperties struct {
	PayloadFormat          *byte
	MessageExpiry          *uint32
	ContentType            string
	ResponseTopic          string
	CorrelationData        []byte
	TopicAlias             *uint16
	SubscriptionIdentifier *uint32
	User                   map[string]string
}

// PublishPropertiesFromPacketProperties is a function that takes a lower level
// Properties struct and returns a PublishProperties
func (p *Publish) PropertiesFromPacketProperties(prop *packets.Properties) {
	p.Properties = &PublishProperties{
		PayloadFormat:          prop.PayloadFormat,
		MessageExpiry:          prop.MessageExpiry,
		ContentType:            prop.ContentType,
		ResponseTopic:          prop.ResponseTopic,
		CorrelationData:        prop.CorrelationData,
		TopicAlias:             prop.TopicAlias,
		SubscriptionIdentifier: prop.SubscriptionIdentifier,
		User: prop.User,
	}
}

func PublishFromPacketPublish(p *packets.Publish) *Publish {
	v := &Publish{
		QoS:     p.QoS,
		Retain:  p.Retain,
		Topic:   p.Topic,
		Payload: p.Payload,
	}
	v.PropertiesFromPacketProperties(p.Properties)

	return v
}

func (p *Publish) Packet() *packets.Publish {
	v := &packets.Publish{
		QoS:     p.QoS,
		Retain:  p.Retain,
		Topic:   p.Topic,
		Payload: p.Payload,
	}
	if p.Properties != nil {
		v.Properties = &packets.Properties{
			PayloadFormat:          p.Properties.PayloadFormat,
			MessageExpiry:          p.Properties.MessageExpiry,
			ContentType:            p.Properties.ContentType,
			ResponseTopic:          p.Properties.ResponseTopic,
			CorrelationData:        p.Properties.CorrelationData,
			TopicAlias:             p.Properties.TopicAlias,
			SubscriptionIdentifier: p.Properties.SubscriptionIdentifier,
			User: p.Properties.User,
		}
	}

	return v
}

func (p *Publish) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "topic: %s  qos: %d  retain: %t\n", p.Topic, p.QoS, p.Retain)
	if p.Properties.PayloadFormat != nil {
		fmt.Fprintf(&b, "PayloadFormat: %v\n", p.Properties.PayloadFormat)
	}
	if p.Properties.MessageExpiry != nil {
		fmt.Fprintf(&b, "MessageExpiry: %v\n", p.Properties.MessageExpiry)
	}
	if p.Properties.ContentType != "" {
		fmt.Fprintf(&b, "ContentType: %v\n", p.Properties.ContentType)
	}
	if p.Properties.ResponseTopic != "" {
		fmt.Fprintf(&b, "ResponseTopic: %v\n", p.Properties.ResponseTopic)
	}
	if p.Properties.CorrelationData != nil {
		fmt.Fprintf(&b, "CorrelationData: %v\n", p.Properties.CorrelationData)
	}
	if p.Properties.SubscriptionIdentifier != nil {
		fmt.Fprintf(&b, "SubscriptionIdentifier: %v\n", p.Properties.SubscriptionIdentifier)
	}
	for k, v := range p.Properties.User {
		fmt.Fprintf(&b, "User: %s : %s\n", k, v)
	}
	b.WriteString(string(p.Payload))

	return b.String()
}
