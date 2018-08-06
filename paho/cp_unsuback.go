package paho

import "github.com/eclipse/paho.mqtt.golang/packets"

// Unsuback is a representation of an MQTT Unsuback packet
type Unsuback struct {
	Reasons    []byte
	Properties *UnsubackProperties
}

// UnsubackProperties is a struct of the properties that can be set
// for a Unsuback packet
type UnsubackProperties struct {
	ReasonString string
	User         map[string]string
}

// Packet returns a packets library Unsuback from the paho Unsuback
// on which it is called
func (u *Unsuback) Packet() *packets.Unsuback {
	return &packets.Unsuback{
		Reasons: u.Reasons,
		Properties: &packets.Properties{
			User: u.Properties.User,
		},
	}
}

// UnsubackFromPacketUnsuback takes a packets library Unsuback and
// returns a paho library Unsuback
func UnsubackFromPacketUnsuback(u *packets.Unsuback) *Unsuback {
	return &Unsuback{
		Reasons: u.Reasons,
		Properties: &UnsubackProperties{
			ReasonString: u.Properties.ReasonString,
			User:         u.Properties.User,
		},
	}
}
