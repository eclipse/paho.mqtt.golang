package paho

import "github.com/eclipse/paho.mqtt.golang/packets"

// Connect is a representation of the MQTT Connect packet
type Connect struct {
	Username       string
	Password       []byte
	ClientID       string
	CleanStart     bool
	KeepAlive      uint16
	Properties     *ConnectProperties
	WillMessage    *WillMessage
	WillProperties *WillProperties
	UsernameFlag   bool
	PasswordFlag   bool
}

// ConnectProperties is a struct of the properties that can be set
// for a Connect packet
type ConnectProperties struct {
	SessionExpiryInterval *uint32
	AuthMethod            string
	AuthData              []byte
	WillDelayInterval     *uint32
	RequestProblemInfo    bool
	RequestResponseInfo   bool
	ReceiveMaximum        *uint16
	TopicAliasMaximum     *uint16
	MaximumQOS            *byte
	MaximumPacketSize     *uint32
	User                  map[string]string
}

// InitProperties is a function that takes a lower level
// Properties struct and completes the properties of the Connect on
// which it is called
func (c *Connect) InitProperties(p *packets.Properties) {
	c.Properties = &ConnectProperties{
		SessionExpiryInterval: p.SessionExpiryInterval,
		AuthMethod:            p.AuthMethod,
		AuthData:              p.AuthData,
		WillDelayInterval:     p.WillDelayInterval,
		RequestResponseInfo:   false,
		RequestProblemInfo:    true,
		ReceiveMaximum:        p.ReceiveMaximum,
		TopicAliasMaximum:     p.TopicAliasMaximum,
		MaximumQOS:            p.MaximumQOS,
		MaximumPacketSize:     p.MaximumPacketSize,
		User:                  p.User,
	}

	if p.RequestResponseInfo != nil {
		c.Properties.RequestResponseInfo = *p.RequestProblemInfo == 1
	}
	if p.RequestProblemInfo != nil {
		c.Properties.RequestProblemInfo = *p.RequestProblemInfo == 1
	}
}

// InitWillProperties is a function that takes a lower level
// Properties struct and completes the properties of the Will in the Connect on
// which it is called
func (c *Connect) InitWillProperties(p *packets.Properties) {
	c.WillProperties = &WillProperties{
		WillDelayInterval: p.WillDelayInterval,
		PayloadFormat:     p.PayloadFormat,
		MessageExpiry:     p.MessageExpiry,
		ContentType:       p.ContentType,
		ResponseTopic:     p.ResponseTopic,
		CorrelationData:   p.CorrelationData,
		User:              p.User,
	}
}

// ConnectFromPacketConnect takes a packets library Connect and
// returns a paho library Connect
func ConnectFromPacketConnect(p *packets.Connect) *Connect {
	v := &Connect{
		UsernameFlag: p.UsernameFlag,
		Username:     p.Username,
		PasswordFlag: p.PasswordFlag,
		Password:     p.Password,
		ClientID:     p.ClientID,
		CleanStart:   p.CleanStart,
		KeepAlive:    p.KeepAlive,
	}
	v.InitProperties(p.Properties)
	if p.WillFlag {
		v.WillMessage = &WillMessage{
			Retain:  p.WillRetain,
			QoS:     p.WillQOS,
			Topic:   p.WillTopic,
			Payload: p.WillMessage,
		}
		v.InitWillProperties(p.WillProperties)
	}

	return v
}

// Packet returns a packets library Connect from the paho Connect
// on which it is called
func (c *Connect) Packet() *packets.Connect {
	v := &packets.Connect{
		UsernameFlag: c.UsernameFlag,
		Username:     c.Username,
		PasswordFlag: c.PasswordFlag,
		Password:     c.Password,
		ClientID:     c.ClientID,
		CleanStart:   c.CleanStart,
		KeepAlive:    c.KeepAlive,
	}

	if c.Properties != nil {
		v.Properties = &packets.Properties{
			SessionExpiryInterval: c.Properties.SessionExpiryInterval,
			AuthMethod:            c.Properties.AuthMethod,
			AuthData:              c.Properties.AuthData,
			WillDelayInterval:     c.Properties.WillDelayInterval,
			ReceiveMaximum:        c.Properties.ReceiveMaximum,
			TopicAliasMaximum:     c.Properties.TopicAliasMaximum,
			MaximumQOS:            c.Properties.MaximumQOS,
			MaximumPacketSize:     c.Properties.MaximumPacketSize,
			User:                  c.Properties.User,
		}
		if c.Properties.RequestResponseInfo {
			v.Properties.RequestResponseInfo = Byte(1)
		}
		if !c.Properties.RequestProblemInfo {
			v.Properties.RequestProblemInfo = Byte(0)
		}
	}

	if c.WillMessage != nil {
		v.WillFlag = true
		v.WillQOS = c.WillMessage.QoS
		v.WillTopic = c.WillMessage.Topic
		v.WillRetain = c.WillMessage.Retain
		v.WillMessage = c.WillMessage.Payload
		if c.WillProperties != nil {
			v.WillProperties = &packets.Properties{
				WillDelayInterval: c.WillProperties.WillDelayInterval,
				PayloadFormat:     c.WillProperties.PayloadFormat,
				MessageExpiry:     c.WillProperties.MessageExpiry,
				ContentType:       c.WillProperties.ContentType,
				ResponseTopic:     c.WillProperties.ResponseTopic,
				CorrelationData:   c.WillProperties.CorrelationData,
				User:              c.WillProperties.User,
			}
		}
	}

	return v
}

// WillMessage is a representation of the LWT message that can
// be sent with the Connect packet
type WillMessage struct {
	Retain  bool
	QoS     byte
	Topic   string
	Payload []byte
}

// WillProperties is a struct of the properties that can be set
// for a Will in a Connect packet
type WillProperties struct {
	WillDelayInterval *uint32
	PayloadFormat     *byte
	MessageExpiry     *uint32
	ContentType       string
	ResponseTopic     string
	CorrelationData   []byte
	User              map[string]string
}
