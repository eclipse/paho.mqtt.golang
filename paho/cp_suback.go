package paho

import "github.com/eclipse/paho.mqtt.golang/packets"

// Suback is a representation of an MQTT suback packet
type Suback struct {
	Reasons    []byte
	Properties *SubackProperties
}

// SubackProperties is a struct of the properties that can be set
// for a Suback packet
type SubackProperties struct {
	ReasonString string
	User         map[string]string
}

// Packet returns a packets library Suback from the paho Suback
// on which it is called
func (s *Suback) Packet() *packets.Suback {
	return &packets.Suback{
		Reasons: s.Reasons,
		Properties: &packets.Properties{
			User: s.Properties.User,
		},
	}
}

// SubackFromPacketSuback takes a packets library Suback and
// returns a paho library Suback
func SubackFromPacketSuback(s *packets.Suback) *Suback {
	return &Suback{
		Reasons: s.Reasons,
		Properties: &SubackProperties{
			ReasonString: s.Properties.ReasonString,
			User:         s.Properties.User,
		},
	}
}
