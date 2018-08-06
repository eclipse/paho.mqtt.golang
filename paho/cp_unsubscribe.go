package paho

import "github.com/eclipse/paho.mqtt.golang/packets"

// Unsubscribe is a representation of an MQTT unsubscribe packet
type Unsubscribe struct {
	Topics     []string
	Properties *UnsubscribeProperties
}

// UnsubscribeProperties is a struct of the properties that can be set
// for a Unsubscribe packet
type UnsubscribeProperties struct {
	User map[string]string
}

// Packet returns a packets library Unsubscribe from the paho Unsubscribe
// on which it is called
func (u *Unsubscribe) Packet() *packets.Unsubscribe {
	return &packets.Unsubscribe{
		Topics: u.Topics,
		Properties: &packets.Properties{
			User: u.Properties.User,
		},
	}
}
