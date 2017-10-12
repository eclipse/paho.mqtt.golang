package paho

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	p "github.com/eclipse/paho.mqtt.golang/packets"
)

//Client provides the implementation of an MQTT v5 client
//
type Client struct {
	sync.Mutex
	conn               net.Conn
	settings           ClientSettings
	PingTimeout        time.Duration
	Pinger             Pinger
	OnConnect          func(p.Connack)
	OnDisconnect       func(p.Disconnect)
	PropertiesModifier func(*Properties)
}

type ClientSettings struct {
	PingTimeout       time.Duration
	TopicAliasMaximum *uint16
	ReceiveMaximum    *uint16
	MaximumPacketSize *uint32
}

func NewClient(c net.Conn, ph Pinger) *Client {
	cl := &Client{
		conn:   c,
		Pinger: ph,
	}

	return cl
}

func (c *Client) Connect(pt Properties, cn p.Connect) (*p.Connack, error) {
	if v, iv := pt.Validate(p.CONNECT); !v {
		return nil, fmt.Errorf("Invalid properties for a Connect packet: %s", strings.Join(iv, ", "))
	}

	cn.IDVP = p.IDValuePair(pt)
	cn.ProtocolName = "MQTT"
	cn.ProtocolVersion = 5
	connectCP := p.NewControlPacket(p.CONNECT)
	connectCP.Content = &cn

	if err := connectCP.Send(c.conn); err != nil {
		return nil, err
	}

	connackCP, err := p.ReadPacket(c.conn)
	if err != nil {
		return nil, err
	}
	ca, ok := connackCP.Content.(*p.Connack)
	if !ok {
		return nil, fmt.Errorf("Response packet to Connect was not Connack")
	}
	if ca.ReasonCode != 0 {
		return ca, fmt.Errorf("%s", ca.Reason())
	}
	if ca.IDVP.TopicAliasMaximum != nil {
		c.settings.TopicAliasMaximum = ca.IDVP.TopicAliasMaximum
	}

	return ca, nil
}
