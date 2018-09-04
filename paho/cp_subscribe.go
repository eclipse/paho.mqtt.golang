package paho

import "github.com/eclipse/paho.mqtt.golang/packets"

// Subscribe is a representation of a MQTT subscribe packet
type Subscribe struct {
	Subscriptions map[string]SubscribeOptions
	Properties    *SubscribeProperties
}

// SubscribeOptions is the struct representing the options for a subscription
type SubscribeOptions struct {
	QoS               byte
	NoLocal           bool
	RetainAsPublished bool
	RetainHandling    byte
}

// SubscribeProperties is a struct of the properties that can be set
// for a Subscribe packet
type SubscribeProperties struct {
	SubscriptionIdentifier *uint32
	User                   map[string]string
}

// InitProperties is a function that takes a packet library
// Properties struct and completes the properties of the Subscribe on
// which it is called
func (s *Subscribe) InitProperties(prop *packets.Properties) {
	s.Properties = &SubscribeProperties{
		SubscriptionIdentifier: prop.SubscriptionIdentifier,
		User:                   prop.User,
	}
}

// PacketSubOptionsFromSubscribeOptions returns a map of string to packet
// library SubOptions for the paho Subscribe on which it is called
func (s *Subscribe) PacketSubOptionsFromSubscribeOptions() map[string]packets.SubOptions {
	r := make(map[string]packets.SubOptions)
	for k, v := range s.Subscriptions {
		r[k] = packets.SubOptions{
			QoS:               v.QoS,
			NoLocal:           v.NoLocal,
			RetainAsPublished: v.RetainAsPublished,
			RetainHandling:    v.RetainHandling,
		}
	}

	return r
}

// Packet returns a packets library Subscribe from the paho Subscribe
// on which it is called
func (s *Subscribe) Packet() *packets.Subscribe {
	v := &packets.Subscribe{Subscriptions: s.PacketSubOptionsFromSubscribeOptions()}

	if s.Properties != nil {
		v.Properties = &packets.Properties{
			SubscriptionIdentifier: s.Properties.SubscriptionIdentifier,
			User:                   s.Properties.User,
		}
	}

	return v
}
