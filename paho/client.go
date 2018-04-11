package paho

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	p "github.com/eclipse/paho.mqtt.golang/packets"
)

// Client is the struct representing an MQTT client
type Client struct {
	sync.Mutex
	Stop          chan struct{}
	Workers       sync.WaitGroup
	Conn          net.Conn
	MIDs          MIDService
	PingHandler   Pinger
	Router        Router
	Persistence   Persistence
	PacketTimeout time.Duration
	Disconnected  func(p.Disconnect)
}

// NewClient is used to create a new instance of an MQTT client. It
// returns a pointer to the new client instance and an error.
// opts is a variadic of functions that take a pointer to a Client and
// returns an error. These functions are used to modify the properties
// of the client (such as setting the PingHandler, Router, etc), a
// selection of such functions are provided in this library.
func NewClient(opts ...func(*Client) error) (*Client, error) {
	debug.Println("Creating new client")
	c := &Client{
		Stop: make(chan struct{}),
	}

	debug.Println("Setting options")
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	if c.PingHandler == nil {
		c.PingHandler = &PingHandler{
			pingFailHandler: func(e error) {
				c.Error(e)
			},
		}
	}

	if c.Persistence == nil {
		c.Persistence = &noopPersistence{}
	}

	if c.MIDs == nil {
		c.MIDs = &MIDs{index: make(map[uint16]*CPContext)}
	}

	if c.PacketTimeout == 0 {
		c.PacketTimeout = 10 * time.Second
	}

	return c, nil
}

// Connect is used to connect the client to a server. It presumes that
// the Client instance already has a working network connection.
// The function takes a pre-prepared Connect packet, and uses that to
// establish an MQTT connection. Assuming the connection completes
// successfully the rest of the client is initiated and the Connack
// returned. Otherwise the failure Connack (if there is one) is returned
// along with an error indicating the reason for the failure to connect.
func (c *Client) Connect(cp *p.Connect) (*p.Connack, error) {
	debug.Println("Connecting")
	c.Lock()
	defer c.Unlock()

	if c.Conn == nil {
		return nil, fmt.Errorf("Client has no connection")
	}

	cp.ProtocolName = "MQTT"
	cp.ProtocolVersion = 5

	debug.Println("Sending CONNECT")
	if err := cp.Send(c.Conn); err != nil {
		return nil, err
	}

	debug.Println("Waiting for CONNACK")
	ca, err := p.ReadPacket(c.Conn)
	if err != nil {
		return nil, err
	}

	if ca.Type != p.CONNACK {
		return nil, fmt.Errorf("Received %d instead of Connack", ca.Type)
	}

	if ca.Content.(*p.Connack).ReasonCode >= 0x80 {
		debug.Println("Received an error code in Connack:", ca.Content.(*p.Connack).ReasonCode)
		return ca.Content.(*p.Connack), fmt.Errorf("Failed to connect to server: %s", ca.Content.(*p.Connack).Reason())
	}

	debug.Println("Received CONNACK, starting PingHandler and Incoming")
	c.Workers.Add(2)
	go c.PingHandler.Start(c.Conn, time.Duration(cp.KeepAlive)*time.Second)
	go c.Incoming()

	return ca.Content.(*p.Connack), nil
}

// Incoming is the Client function that reads and handles incoming
// packets from the server. The function is started as a goroutine
// from Connect(), it exits when it receives a server initiated
// Disconnect, the Stop channel is closed or there is an error reading
// a packet from the network connection
func (c *Client) Incoming() {
	defer c.Workers.Done()
	for {
		select {
		case <-c.Stop:
			debug.Println("Client stopping, Incoming stopping")
			return
		default:
			recv, err := p.ReadPacket(c.Conn)
			if err != nil {
				c.Error(err)
				return
			}
			debug.Println("Received a control packet:", recv.Type)
			switch recv.Type {
			case p.PUBLISH:
				pb := recv.Content.(*p.Publish)
				go c.Router.Route(pb)
				switch pb.QoS {
				case 1:
					p.NewPuback(p.PubackFromPublish(pb)).Send(c.Conn)
				case 2:
					p.NewPubrec(p.PubrecFromPublish(pb)).Send(c.Conn)
				}
			case p.PUBACK, p.PUBCOMP, p.SUBACK, p.UNSUBACK:
				if cpCtx := c.MIDs.Get(recv.PacketID()); cpCtx != nil {
					cpCtx.Return <- *recv
				} else {
					debug.Println("Received a response for a message ID we don't know:", recv.PacketID())
				}
			case p.PUBREC:
				if cpCtx := c.MIDs.Get(recv.PacketID()); cpCtx == nil {
					debug.Println("Received a PUBREC for a message ID we don't know:", recv.PacketID())
					p.NewPubrel(
						p.PubrelFromPubrec(recv.Content.(*p.Pubrec)),
						p.PubrelReasonCode(0x92),
					).Send(c.Conn)
				} else {
					pr := recv.Content.(*p.Pubrec)
					if pr.ReasonCode >= 0x80 {
						//Received a failure code, shortcut and return
						cpCtx.Return <- *recv
					} else {
						p.NewPubrel(p.PubrelFromPubrec(pr)).Send(c.Conn)
					}
				}
			case p.PUBREL:
				//Auto respond to pubrels unless failure code
				pr := recv.Content.(*p.Pubrel)
				if pr.ReasonCode < 0x80 {
					//Received a failure code, continue
					continue
				} else {
					p.NewPubcomp(p.PubcompFromPubrel(pr)).Send(c.Conn)
				}
			case p.DISCONNECT:
				if c.Disconnected != nil {
					go c.Disconnected(*recv.Content.(*p.Disconnect))
				}
				c.Error(fmt.Errorf("Received server initiated disconnect"))
			}
		}
	}
}

// Error is called to signify that an error situation has occurred, this
// causes the client's Stop channel to be closed (if it hasn't already been)
// which results in the other client goroutines terminating.
// It also closes the client network connection.
func (c *Client) Error(e error) {
	c.Lock()
	debug.Println("Error called:", e)
	select {
	case <-c.Stop:
		//already shutting down, do nothing
	default:
		close(c.Stop)
	}
	c.Conn.Close()
	c.Unlock()
}

// Subscribe is used to send a Subscription request to the MQTT server.
// It is passed a pre-prepared Subscribe packet and blocks waiting for
// a response Suback, or for the timeout to fire. Any reponse Suback
// is returned from the function, along with any errors.
func (c *Client) Subscribe(s *p.Subscribe) (*p.Suback, error) {
	debug.Printf("Subscribing to %+v", s.Subscriptions)
	ctx, cf := context.WithTimeout(context.Background(), c.PacketTimeout)
	defer cf()
	cpCtx := &CPContext{make(chan p.ControlPacket, 1), ctx}

	s.PacketID = c.MIDs.Request(cpCtx)
	debug.Println("Sending SUBSCRIBE")
	if err := s.Send(c.Conn); err != nil {
		return nil, err
	}
	debug.Println("Waiting for SUBACK")
	var resp p.ControlPacket

	select {
	case <-ctx.Done():
		if e := ctx.Err(); e == context.DeadlineExceeded {
			debug.Println("Timeout waiting for SUBACK")
			return nil, e
		}
	case resp = <-cpCtx.Return:
	}

	if resp.Type != p.SUBACK {
		return nil, fmt.Errorf("Received %d instead of Suback", resp.Type)
	}
	debug.Println("Received SUBACK")

	r := resp.Content.(*p.Suback).Reasons
	switch {
	case len(r) == 1:
		if r[0] >= 0x80 {
			debug.Println("Received an error code in Suback:", r[0])
			return resp.Content.(*p.Suback), fmt.Errorf("Failed to subscribe to topic: %s", resp.Content.(*p.Suback).Reason(0))
		}
	default:
		for _, code := range r {
			if code >= 0x80 {
				debug.Println("Received an error code in Suback:", code)
				return resp.Content.(*p.Suback), fmt.Errorf("At least one requested subscription failed")
			}
		}
	}

	return resp.Content.(*p.Suback), nil
}

// Unsubscribe is used to send an Unsubscribe request to the MQTT server.
// It is passed a pre-prepared Unsubscribe packet and blocks waiting for
// a response Unsuback, or for the timeout to fire. Any reponse Unsuback
// is returned from the function, along with any errors.
func (c *Client) Unsubscribe(u *p.Unsubscribe) (*p.Unsuback, error) {
	debug.Printf("Unsubscribing from %+v", u.Topics)
	ctx, cf := context.WithTimeout(context.Background(), c.PacketTimeout)
	defer cf()
	cpCtx := &CPContext{make(chan p.ControlPacket, 1), ctx}

	u.PacketID = c.MIDs.Request(cpCtx)
	debug.Println("Sending UNSUBSCRIBE")
	if err := u.Send(c.Conn); err != nil {
		return nil, err
	}
	debug.Println("Waiting for UNSUBACK")
	var resp p.ControlPacket

	select {
	case <-ctx.Done():
		if e := ctx.Err(); e == context.DeadlineExceeded {
			debug.Println("Timeout waiting for UNSUBACK")
			return nil, e
		}
	case resp = <-cpCtx.Return:
	}

	if resp.Type != p.UNSUBACK {
		return nil, fmt.Errorf("Received %d instead of Unsuback", resp.Type)
	}
	debug.Println("Received SUBACK")

	r := resp.Content.(*p.Unsuback).Reasons
	switch {
	case len(r) == 1:
		if r[0] >= 0x80 {
			debug.Println("Received an error code in Suback:", r[0])
			return resp.Content.(*p.Unsuback), fmt.Errorf("Failed to unsubscribe from topic: %s", resp.Content.(*p.Unsuback).Reason(0))
		}
	default:
		for _, code := range r {
			if code >= 0x80 {
				debug.Println("Received an error code in Suback:", code)
				return resp.Content.(*p.Unsuback), fmt.Errorf("At least one requested unsubscribe failed")
			}
		}
	}

	return resp.Content.(*p.Unsuback), nil
}

// Publish is used to send a publication to the MQTT server.
// It is passed a pre-prepared Publish packet and blocks waiting for
// the appropriate response, or for the timeout to fire.
// Any reponse message is returned from the function, along with any errors.
func (c *Client) Publish(pb *p.Publish) (p.Packet, error) {
	debug.Printf("Sending message to %s", pb.Topic)
	c.Lock()
	defer c.Unlock()

	switch pb.QoS {
	case 0:
		debug.Println("Sending QoS0 message")
		if err := pb.Send(c.Conn); err != nil {
			return nil, err
		}
		return nil, nil
	case 1, 2:
		return c.publishQoS12(pb)
	}

	return nil, fmt.Errorf("oops")
}

func (c *Client) publishQoS12(pb *p.Publish) (p.Packet, error) {
	debug.Println("Sending QoS1 message")
	ctx, cf := context.WithTimeout(context.Background(), c.PacketTimeout)
	defer cf()
	cpCtx := &CPContext{make(chan p.ControlPacket, 1), ctx}

	pb.PacketID = c.MIDs.Request(cpCtx)
	if err := pb.Send(c.Conn); err != nil {
		return nil, err
	}
	var resp p.ControlPacket

	select {
	case <-ctx.Done():
		if e := ctx.Err(); e == context.DeadlineExceeded {
			debug.Println("Timeout waiting for SUBACK")
			return nil, e
		}
	case resp = <-cpCtx.Return:
	}

	switch pb.QoS {
	case 1:
		if resp.Type != p.PUBACK {
			return nil, fmt.Errorf("Received %d instead of PUBACK", resp.Type)
		}
		debug.Println("Received PUBACK for", pb.PacketID)
		return resp.Content.(*p.Puback), nil
	case 2:
		if resp.Type != p.PUBCOMP {
			return nil, fmt.Errorf("Received %d instead of PUBCOMP", resp.Type)
		}
		debug.Println("Received PUBCOMP for", pb.PacketID)
		return resp.Content.(*p.Pubcomp), nil
	}

	debug.Println("Ended up with a non QoS1/2 message:", pb.QoS)
	return nil, fmt.Errorf("Ended up with a non QoS1/2 message: %d", pb.QoS)
}

// Disconnect is used to send a Disconnect packet to the MQTT server
// Whether or not the attempt to send the Disconnect packet fails
// (and if it does this function returns any error) the network connection
// is closed.
func (c *Client) Disconnect(d *p.Disconnect) error {
	c.Lock()
	defer c.Unlock()
	defer c.Conn.Close()

	return d.Send(c.Conn)
}
