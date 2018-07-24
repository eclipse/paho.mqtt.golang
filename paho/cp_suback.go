package paho

import "github.com/eclipse/paho.mqtt.golang/packets"

// Suback is a representation of an MQTT suback packet
type Suback struct {
	Reasons    []byte
	Properties *SubackProperties
}

type SubackProperties struct {
	ReasonString string
	User         map[string]string
}

func (s *Suback) Packet() *packets.Suback {
	return &packets.Suback{
		Reasons: s.Reasons,
		Properties: &packets.Properties{
			User: s.Properties.User,
		},
	}
}

func SubackFromUnsubackPacket(s *packets.Suback) *Suback {
	return &Suback{
		Reasons: s.Reasons,
		Properties: &SubackProperties{
			ReasonString: s.Properties.ReasonString,
			User:         s.Properties.User,
		},
	}
}
