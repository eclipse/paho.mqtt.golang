package paho

import "github.com/eclipse/paho.mqtt.golang/packets"

// Unsuback is a representation of an MQTT Unsuback packet
type Unsuback struct {
	Reasons    []byte
	Properties *UnsubackProperties
}

type UnsubackProperties struct {
	ReasonString string
	User         map[string]string
}

func (u *Unsuback) Packet() *packets.Unsuback {
	return &packets.Unsuback{
		Reasons: u.Reasons,
		Properties: &packets.Properties{
			User: u.Properties.User,
		},
	}
}

func UnsubackFromUnsubackPacket(u *packets.Unsuback) *Unsuback {
	return &Unsuback{
		Reasons: u.Reasons,
		Properties: &UnsubackProperties{
			ReasonString: u.Properties.ReasonString,
			User:         u.Properties.User,
		},
	}
}
