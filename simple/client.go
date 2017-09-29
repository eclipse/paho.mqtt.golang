package simple

import (
	"fmt"
	"io"
	"sync"
	"time"

	p "github.com/eclipse/paho.mqtt.golang/packets"
)

type Client struct {
	sync.Mutex
	conn           io.ReadWriteCloser
	ReceiveTimeout time.Duration
	PingTimeout    time.Duration
}

func NewClient(conn io.ReadWriteCloser) *Client {
	c := &Client{
		conn:           conn,
		ReceiveTimeout: 10 * time.Second,
	}

	return c
}

func (c *Client) Connect(cp *p.ControlPacket) (*p.ControlPacket, error) {
	c.Lock()
	defer c.Unlock()

	if cp.Type != p.CONNECT {
		return nil, fmt.Errorf("Connect requires a Connect packet")
	}

	c.PingTimeout = time.Duration(cp.Content.(*p.Connect).KeepAlive) * time.Second

	if err := cp.Send(c.conn); err != nil {
		return nil, err
	}

	ca, err := p.ReadPacket(c.conn)
	if err != nil {
		return nil, err
	}

	if ca.Type != p.CONNACK {
		return nil, fmt.Errorf("Received %d instead of Connack", ca.Type)
	}

	return ca, nil
}

func (c *Client) Ping() error {
	c.Lock()
	defer c.Unlock()

	if err := p.NewControlPacket(p.PINGREQ).Send(c.conn); err != nil {
		return err
	}

	prChan := func() chan error {
		r := make(chan error)
		go func() {
			if pr, err := p.ReadPacket(c.conn); err != nil {
				r <- err
			} else {
				if pr.Type == p.PINGRESP {
					r <- nil
				} else {
					r <- fmt.Errorf("Received %d instead of PingResp", pr.Type)
				}
			}
		}()
		return r
	}()

	select {
	case <-time.After(c.PingTimeout):
		c.conn.Close()
		return fmt.Errorf("Failed to receive a PingReponse")
	case err := <-prChan:
		return err
	}
}

func (c *Client) Subscribe(s *p.ControlPacket) (*p.ControlPacket, error) {
	c.Lock()
	defer c.Unlock()

	if s.Type != p.SUBSCRIBE {
		return nil, fmt.Errorf("Subscribe requires a Subscribe packet")
	}

	if err := s.Send(c.conn); err != nil {
		return nil, err
	}

	sa, err := p.ReadPacket(c.conn)
	if err != nil {
		return nil, err
	}

	if sa.Type != p.SUBACK {
		return nil, fmt.Errorf("Received %d instead of Suback", sa.Type)
	}

	return sa, nil
}

func (c *Client) Unsubscribe(u *p.ControlPacket) (*p.ControlPacket, error) {
	c.Lock()
	defer c.Unlock()

	if u.Type != p.UNSUBSCRIBE {
		return nil, fmt.Errorf("Unsubscribe requires a Unsubscribe packet")
	}

	if err := u.Send(c.conn); err != nil {
		return nil, err
	}

	ua, err := p.ReadPacket(c.conn)
	if err != nil {
		return nil, err
	}

	if ua.Type != p.SUBACK {
		return nil, fmt.Errorf("Received %d instead of Unsuback", ua.Type)
	}

	return ua, nil
}

func (c *Client) Publish(pb *p.ControlPacket) (*p.ControlPacket, error) {
	c.Lock()
	defer c.Unlock()

	if pb.Type != p.PUBLISH {
		return nil, fmt.Errorf("Unsubscribe requires a Unsubscribe packet")
	}

	if err := pb.Send(c.conn); err != nil {
		return nil, err
	}

	switch pb.Content.(*p.Publish).QoS {
	case 0:
		return nil, nil
	case 1:
		pa, err := p.ReadPacket(c.conn)
		if err != nil {
			return nil, err
		}

		if pa.Type != p.PUBACK {
			return nil, fmt.Errorf("Received %d instead of Puback", pa.Type)
		}

		return pa, nil
	case 2:
		prec, err := p.ReadPacket(c.conn)
		if err != nil {
			return nil, err
		}

		if prec.Type != p.PUBREC {
			return nil, fmt.Errorf("Received %d instead of Pubrec", prec.Type)
		}

		prel := p.NewControlPacket(p.PUBREL)
		prel.Content.(*p.Publish).PacketID = pb.Content.(*p.Publish).PacketID

		if err := prel.Send(c.conn); err != nil {
			return nil, err
		}

		pcomp, err := p.ReadPacket(c.conn)
		if err != nil {
			return nil, err
		}

		if pcomp.Type != p.PUBCOMP {
			return nil, fmt.Errorf("Received %d instead of Pubcomp", pcomp.Type)
		}

		return pcomp, nil
	}

	return nil, fmt.Errorf("oops")
}
