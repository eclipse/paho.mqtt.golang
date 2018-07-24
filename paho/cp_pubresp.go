package paho

import "github.com/eclipse/paho.mqtt.golang/packets"

type PublishResponse struct {
	ReasonCode byte
	Properties *PublishResponseProperties
}

type PublishResponseProperties struct {
	ReasonString string
	User         map[string]string
}

func PublishResponseFromPuback(pa *packets.Puback) *PublishResponse {
	return &PublishResponse{
		ReasonCode: pa.ReasonCode,
		Properties: &PublishResponseProperties{
			ReasonString: pa.Properties.ReasonString,
			User:         pa.Properties.User,
		},
	}
}

func PublishResponseFromPubcomp(pc *packets.Pubcomp) *PublishResponse {
	return &PublishResponse{
		ReasonCode: pc.ReasonCode,
		Properties: &PublishResponseProperties{
			ReasonString: pc.Properties.ReasonString,
			User:         pc.Properties.User,
		},
	}
}
