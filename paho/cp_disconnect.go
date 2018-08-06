package paho

import "github.com/eclipse/paho.mqtt.golang/packets"

// Disconnect is a representation of the MQTT Disconnect packet
type Disconnect struct {
	ReasonCode byte
	Properties *DisconnectProperties
}

// DisconnectProperties is a struct of the properties that can be set
// for a Disconnect packet
type DisconnectProperties struct {
	SessionExpiryInterval *uint32
	ServerReference       string
	ReasonString          string
	User                  map[string]string
}

// PropertiesFromPacketProperties is a function that takes a lower level
// Properties struct and completes the properties of the Disconnect on
// which it is called
func (d *Disconnect) PropertiesFromPacketProperties(p *packets.Properties) {
	d.Properties = &DisconnectProperties{
		SessionExpiryInterval: p.SessionExpiryInterval,
		ServerReference:       p.ServerReference,
		ReasonString:          p.ReasonString,
		User:                  p.User,
	}
}

// DisconnectFromPacketDisconnect takes a packets library Disconnect and
// returns a paho library Disconnect
func DisconnectFromPacketDisconnect(p *packets.Disconnect) *Disconnect {
	v := &Disconnect{ReasonCode: p.ReasonCode}
	v.PropertiesFromPacketProperties(p.Properties)

	return v
}

// Packet returns a packets library Disconnect from the paho Disconnect
// on which it is called
func (d *Disconnect) Packet() *packets.Disconnect {
	v := &packets.Disconnect{ReasonCode: d.ReasonCode}

	if d.Properties != nil {
		v.Properties = &packets.Properties{
			SessionExpiryInterval: d.Properties.SessionExpiryInterval,
			ServerReference:       d.Properties.ServerReference,
			ReasonString:          d.Properties.ReasonString,
			User:                  d.Properties.User,
		}
	}

	return v
}
