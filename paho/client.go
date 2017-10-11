package paho

import (
	"net"
	"sync"
	"time"

	p "github.com/eclipse/paho.mqtt.golang/packets"
)

//Client provides the implementation of an MQTT v5 client
//
type Client struct {
	sync.Mutex
	conn               net.Conn
	PingTimeout        time.Duration
	Pinger             Pinger
	OnConnect          func(p.Connack)
	OnDisconnect       func(p.Disconnect)
	PropertiesModifier func(*Properties)
}

func NewClient(c net.Conn, ph Pinger) *Client {
	cl := &Client{
		conn:   c,
		Pinger: ph,
	}

	return cl
}

func (c *Client) Connect(p *Properties) (*p.Connack, error) {

}
