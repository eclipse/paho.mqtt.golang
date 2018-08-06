package paho

import "github.com/eclipse/paho.mqtt.golang/packets"

// PublishResponse is a generic representation of a response
// to a QoS1 or QoS2 Publish
type PublishResponse struct {
	ReasonCode byte
	Properties *PublishResponseProperties
}

// PublishResponseProperties is the properties associated with
// a response to a QoS1 or QoS2 Publish
type PublishResponseProperties struct {
	ReasonString string
	User         map[string]string
}

// PublishResponseFromPuback takes a packets library Puback and
// returns a paho library PublishResponse
func PublishResponseFromPuback(pa *packets.Puback) *PublishResponse {
	return &PublishResponse{
		ReasonCode: pa.ReasonCode,
		Properties: &PublishResponseProperties{
			ReasonString: pa.Properties.ReasonString,
			User:         pa.Properties.User,
		},
	}
}

// PublishResponseFromPubcomp takes a packets library Pubcomp and
// returns a paho library PublishResponse
func PublishResponseFromPubcomp(pc *packets.Pubcomp) *PublishResponse {
	return &PublishResponse{
		ReasonCode: pc.ReasonCode,
		Properties: &PublishResponseProperties{
			ReasonString: pc.Properties.ReasonString,
			User:         pc.Properties.User,
		},
	}
}

// PublishResponseFromPubrec takes a packets library Pubrec and
// returns a paho library PublishResponse
func PublishResponseFromPubrec(pr *packets.Pubrec) *PublishResponse {
	return &PublishResponse{
		ReasonCode: pr.ReasonCode,
		Properties: &PublishResponseProperties{
			ReasonString: pr.Properties.ReasonString,
			User:         pr.Properties.User,
		},
	}
}
