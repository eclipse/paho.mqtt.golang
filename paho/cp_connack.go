package paho

import "github.com/eclipse/paho.mqtt.golang/packets"

// Connack is a representation of the MQTT Connack packet
type Connack struct {
	SessionPresent bool
	ReasonCode     byte
	Properties     *ConnackProperties
}

// ConnackProperties is a struct of the properties that can be set
// for a Connack packet
type ConnackProperties struct {
	AssignedClientID     string
	ServerKeepAlive      *uint16
	WildcardSubAvailable bool
	SubIDAvailable       bool
	SharedSubAvailable   bool
	RetainAvailable      bool
	ResponseInfo         string
	AuthMethod           string
	AuthData             []byte
	ServerReference      string
	ReasonString         string
	ReceiveMaximum       *uint16
	TopicAliasMaximum    *uint16
	MaximumQoS           *byte
	MaximumPacketSize    *uint32
	User                 map[string]string
}

// PropertiesFromPacketProperties is a function that takes a lower level
// Properties struct and completes the properties of the Connack on
// which it is called
func (c *Connack) PropertiesFromPacketProperties(p *packets.Properties) {
	c.Properties = &ConnackProperties{
		AssignedClientID:     p.AssignedClientID,
		ServerKeepAlive:      p.ServerKeepAlive,
		WildcardSubAvailable: true,
		SubIDAvailable:       true,
		SharedSubAvailable:   true,
		RetainAvailable:      true,
		ResponseInfo:         p.ResponseInfo,
		AuthMethod:           p.AuthMethod,
		AuthData:             p.AuthData,
		ServerReference:      p.ServerReference,
		ReasonString:         p.ReasonString,
		ReceiveMaximum:       p.ReceiveMaximum,
		TopicAliasMaximum:    p.TopicAliasMaximum,
		MaximumQoS:           p.MaximumQOS,
		MaximumPacketSize:    p.MaximumPacketSize,
		User:                 p.User,
	}

	if p.WildcardSubAvailable != nil {
		c.Properties.WildcardSubAvailable = *p.WildcardSubAvailable == 1
	}
	if p.SubIDAvailable != nil {
		c.Properties.SubIDAvailable = *p.SubIDAvailable == 1
	}
	if p.SharedSubAvailable != nil {
		c.Properties.SharedSubAvailable = *p.SharedSubAvailable == 1
	}
	if p.RetainAvailable != nil {
		c.Properties.RetainAvailable = *p.RetainAvailable == 1
	}
}

// ConnackFromPacketConnack takes a packets library Connack and
// returns a paho library Connack
func ConnackFromPacketConnack(c *packets.Connack) *Connack {
	v := &Connack{
		SessionPresent: c.SessionPresent,
		ReasonCode:     c.ReasonCode,
	}
	v.PropertiesFromPacketProperties(c.Properties)

	return v
}
