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

type SubscribeProperties struct {
	SubscriptionIdentifier *uint32
	User                   map[string]string
}

// PublishPropertiesFromPacketProperties is a function that takes a lower level
// Properties struct and returns a PublishProperties
func (s *Subscribe) PropertiesFromPacketProperties(prop *packets.Properties) {
	s.Properties = &SubscribeProperties{
		SubscriptionIdentifier: prop.SubscriptionIdentifier,
		User: prop.User,
	}
}

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

func (s *Subscribe) Packet() *packets.Subscribe {
	v := &packets.Subscribe{Subscriptions: s.PacketSubOptionsFromSubscribeOptions()}

	if s.Properties != nil {
		v.Properties = &packets.Properties{
			SubscriptionIdentifier: s.Properties.SubscriptionIdentifier,
			User: s.Properties.User,
		}
	}

	return v
}
