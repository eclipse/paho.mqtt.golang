package simple

import (
	"fmt"
	"net"
	"sync"
	"time"

	p "github.com/eclipse/paho.mqtt.golang/packets"
)

type Client struct {
	sync.Mutex
	Conn            net.Conn
	PingTimeout     time.Duration
	LastPing        time.Time
	PingOutstanding bool
}

func NewClient(opts ...func(*Client) error) (*Client, error) {
	c := &Client{
		PingOutstanding: false,
		PingTimeout:     30 * time.Second,
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Client) Connect(cp *p.Connect) (*p.Connack, error) {
	c.Lock()
	defer c.Unlock()

	c.PingTimeout = time.Duration(cp.KeepAlive) * time.Second

	if err := cp.Send(c.Conn); err != nil {
		return nil, err
	}

	ca, err := p.ReadPacket(c.Conn)
	if err != nil {
		return nil, err
	}

	if ca.Type != p.CONNACK {
		return nil, fmt.Errorf("Received %d instead of Connack", ca.Type)
	}

	return ca.Content.(*p.Connack), nil
}

func (c *Client) Ping() error {
	c.Lock()
	defer c.Unlock()

	if c.PingOutstanding && time.Now().Sub(c.LastPing) > (c.PingTimeout+c.PingTimeout/2) {
		return fmt.Errorf("Failed to receive a PingReponse")
	}

	return p.NewControlPacket(p.PINGREQ).Send(c.Conn)
}

func (c *Client) Subscribe(s *p.Subscribe) (*p.Suback, error) {
	c.Lock()
	defer c.Unlock()

	if err := s.Send(c.Conn); err != nil {
		return nil, err
	}

	sa, err := p.ReadPacket(c.Conn)
	if err != nil {
		return nil, err
	}

	if sa.Type != p.SUBACK {
		return nil, fmt.Errorf("Received %d instead of Suback", sa.Type)
	}

	return sa.Content.(*p.Suback), nil
}

func (c *Client) Unsubscribe(u *p.Unsubscribe) (*p.Unsuback, error) {
	c.Lock()
	defer c.Unlock()

	if err := u.Send(c.Conn); err != nil {
		return nil, err
	}

	ua, err := p.ReadPacket(c.Conn)
	if err != nil {
		return nil, err
	}

	if ua.Type != p.SUBACK {
		return nil, fmt.Errorf("Received %d instead of Unsuback", ua.Type)
	}

	return ua.Content.(*p.Unsuback), nil
}

func (c *Client) Publish(pb *p.Publish) (p.Packet, error) {
	c.Lock()
	defer c.Unlock()

	if err := pb.Send(c.Conn); err != nil {
		return nil, err
	}

	switch pb.QoS {
	case 0:
		return nil, nil
	case 1:
		pa, err := p.ReadPacket(c.Conn)
		if err != nil {
			return nil, err
		}

		if pa.Type != p.PUBACK {
			return nil, fmt.Errorf("Received %d instead of Puback", pa.Type)
		}

		return pa.Content.(*p.Puback), nil
	case 2:
		prec, err := p.ReadPacket(c.Conn)
		if err != nil {
			return nil, err
		}

		if prec.Type != p.PUBREC {
			return nil, fmt.Errorf("Received %d instead of Pubrec", prec.Type)
		}

		prel := p.NewControlPacket(p.PUBREL)
		prel.Content.(*p.Pubrel).PacketID = pb.PacketID

		if err := prel.Send(c.Conn); err != nil {
			return nil, err
		}

		pcomp, err := p.ReadPacket(c.Conn)
		if err != nil {
			return nil, err
		}

		if pcomp.Type != p.PUBCOMP {
			return nil, fmt.Errorf("Received %d instead of Pubcomp", pcomp.Type)
		}

		return pcomp.Content.(*p.Pubcomp), nil
	}

	return nil, fmt.Errorf("oops")
}

func (c *Client) Disconnect(d *p.Disconnect) error {
	c.Lock()
	defer c.Unlock()
	defer c.Conn.Close()

	return d.Send(c.Conn)
}

func (c *Client) Receive() (*Message, error) {
	c.Lock()
	defer c.Unlock()

	var m Message

	c.Conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	cp, err := p.ReadPacket(c.Conn)
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	switch cp.Type {
	case p.PINGRESP:
		c.PingOutstanding = false
		return nil, nil
	case p.PUBLISH:
		m.Topic = cp.Content.(*p.Publish).Topic
		m.QoS = cp.Content.(*p.Publish).QoS
		m.Retain = cp.Content.(*p.Publish).Retain
		m.Properties = cp.Content.(*p.Publish).Properties
		m.Payload = cp.Content.(*p.Publish).Payload
	}

	return &m, nil
}
