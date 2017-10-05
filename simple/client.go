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
	conn            net.Conn
	ReceiveTimeout  time.Duration
	PingTimeout     time.Duration
	LastPing        time.Time
	PingOutstanding bool
}

func NewClient(conn net.Conn) *Client {
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

	if c.PingOutstanding && time.Now().Sub(c.LastPing) > (c.PingTimeout+c.PingTimeout/2) {
		return fmt.Errorf("Failed to receive a PingReponse")
	}

	return p.NewControlPacket(p.PINGREQ).Send(c.conn)

	// prChan := make(chan error)
	// go func() {
	// 	pr, err := p.ReadPacket(c.conn)
	// 	if err != nil {
	// 		log.Println(err)
	// 		prChan <- err
	// 	} else {
	// 		if pr.Type == p.PINGRESP {
	// 			prChan <- nil
	// 		} else {
	// 			prChan <- fmt.Errorf("Received %d instead of PingResp", pr.Type)
	// 		}
	// 	}
	// }()

	// select {
	// case <-time.After(c.PingTimeout):
	// 	return fmt.Errorf("Failed to receive a PingReponse")
	// case err := <-prChan:
	// 	return err
	// }
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

func (c *Client) SendMessage(topic string, qos byte, retain bool, idvp *p.IDValuePair, payload []byte) (*p.ControlPacket, error) {
	pb := p.NewControlPacket(p.PUBLISH)
	ct := pb.Content.(*p.Publish)
	ct.Topic = topic
	ct.QoS = qos
	ct.Retain = retain
	if qos > 0 {
		ct.PacketID = 1
	}
	if idvp != nil {
		ct.IDVP = *idvp
	}
	pb.Payload = payload

	return c.Publish(pb)
}

func (c *Client) Disconnect(d *p.ControlPacket) error {
	c.Lock()
	defer c.Unlock()

	if d.Type != p.DISCONNECT {
		return fmt.Errorf("Unsubscribe requires a Unsubscribe packet")
	}

	return d.Send(c.conn)
}

func (c *Client) Receive() (*Message, error) {
	c.Lock()
	defer c.Unlock()

	var m Message

	c.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	cp, err := p.ReadPacket(c.conn)
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
		m.IDVP = cp.Content.(*p.Publish).IDVP
		m.Payload = cp.Payload
	}

	return &m, nil
}
