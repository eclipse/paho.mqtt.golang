package paho

import "github.com/eclipse/paho.mqtt.golang/packets"

// Unsubscribe is a representation of an MQTT unsubscribe packet
type Unsubscribe struct {
	Topics     []string
	Properties *UnsubscribeProperties
}

type UnsubscribeProperties struct {
	User map[string]string
}

func (u *Unsubscribe) Packet() *packets.Unsubscribe {
	return &packets.Unsubscribe{
		Topics: u.Topics,
		Properties: &packets.Properties{
			User: u.Properties.User,
		},
	}
}
