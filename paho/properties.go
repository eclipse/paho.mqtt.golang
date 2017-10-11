package paho

import "github.com/eclipse/paho.mqtt.golang/packets"

type Properties packets.IDValuePair

func NewProperties() *Properties {
	return &Properties{}
}

func (p *Properties) SetPayloadFormat(v byte) *Properties {
	p.PayloadFormat = &v
	return p
}

func (p *Properties) SetPubExpiry(v uint32) *Properties {
	p.PubExpiry = &v
	return p
}

func (p *Properties) SetContentType(v string) *Properties {
	p.ContentType = v
	return p
}

func (p *Properties) SetReplyTopic(v string) *Properties {
	p.ReplyTopic = v
	return p
}

func (p *Properties) SetCorrelationData(v []byte) *Properties {
	p.CorrelationData = v
	return p
}

func (p *Properties) SetSubscriptionIdentifier(v *uint32) *Properties {
	p.SubscriptionIdentifier = v
	return p
}

func (p *Properties) SetSessionExpiryInterval(v uint32) *Properties {
	p.SessionExpiryInterval = &v
	return p
}

func (p *Properties) SetAssignedClientID(v string) *Properties {
	p.AssignedClientID = v
	return p
}

func (p *Properties) SetServerKeepAlive(v uint16) *Properties {
	p.ServerKeepAlive = &v
	return p
}

func (p *Properties) SetAuthMethod(v string) *Properties {
	p.AuthMethod = v
	return p
}

func (p *Properties) SetAuthData(v []byte) *Properties {
	p.AuthData = v
	return p
}

func (p *Properties) SetRequestProblemInfo(v byte) *Properties {
	p.RequestProblemInfo = &v
	return p
}

func (p *Properties) SetWillDelayInterval(v uint32) *Properties {
	p.WillDelayInterval = &v
	return p
}

func (p *Properties) SetRequestResponseInfo(v byte) *Properties {
	p.RequestResponseInfo = &v
	return p
}

func (p *Properties) SetResponseInfo(v string) *Properties {
	p.ResponseInfo = v
	return p
}

func (p *Properties) SetServerReference(v string) *Properties {
	p.ServerReference = v
	return p
}

func (p *Properties) SetReasonString(v string) *Properties {
	p.ReasonString = v
	return p
}

func (p *Properties) SetReceiveMaximum(v uint16) *Properties {
	p.ReceiveMaximum = &v
	return p
}

func (p *Properties) SetTopicAliasMaximum(v uint16) *Properties {
	p.TopicAliasMaximum = &v
	return p
}

func (p *Properties) SetTopicAlias(v uint16) *Properties {
	p.TopicAlias = &v
	return p
}

func (p *Properties) SetMaximumQOS(v byte) *Properties {
	p.MaximumQOS = &v
	return p
}

func (p *Properties) SetRetainAvailable(v byte) *Properties {
	p.RetainAvailable = &v
	return p
}

func (p *Properties) SetUserProperty(k, v string) *Properties {
	p.UserProperty[k] = v
	return p
}

func (p *Properties) SetMaximumPacketSize(v uint32) *Properties {
	p.MaximumPacketSize = &v
	return p
}

func (p *Properties) SetWildcardSubAvailable(v byte) *Properties {
	p.WildcardSubAvailable = &v
	return p
}

func (p *Properties) SetSubIDAvailable(v byte) *Properties {
	p.SubIDAvailable = &v
	return p
}

func (p *Properties) SetSharedSubAvailable(v byte) *Properties {
	p.SharedSubAvailable = &v
	return p
}

func (p *Properties) Validate() (bool, []string) {

}
