package paho

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
)

// Client is the struct representing an MQTT client
type Client struct {
	sync.Mutex
	caCtx         *caContext
	raCtx         *CPContext
	stop          chan struct{}
	ClientID      string
	workers       sync.WaitGroup
	Conn          net.Conn
	MIDs          MIDService
	AuthHandler   Auther
	PingHandler   Pinger
	Router        Router
	Persistence   Persistence
	PacketTimeout time.Duration
	OnDisconnect  func(packets.Disconnect)
	serverProps   CommsProperties
}

// CommsProperties is a struct of the communication properties that may
// be set by the server in the Connack and that the client needs to be
// aware of for future subscribes/publishes
type CommsProperties struct {
	ReceiveMaximum       uint16
	MaximumQoS           byte
	MaximumPacketSize    uint32
	TopicAliasMaximum    uint16
	RetainAvailable      bool
	WildcardSubAvailable bool
	SubIDAvailable       bool
	SharedSubAvailable   bool
}

type caContext struct {
	Return  chan *packets.Connack
	Context context.Context
}

// NewClient is used to create a new default instance of an MQTT client.
// It returns a pointer to the new client instance.
// The default client uses the provided PingHandler, MessageID and
// StandardRouter implementations, and a noop Persistence.
// These should be replaced if desired before the client is connected.
// client.Conn *MUST* be set to an already connected net.Conn before
// Connect() is called.
func NewClient() *Client {
	debug.Println("Creating new client")
	c := &Client{
		stop: make(chan struct{}),
		serverProps: CommsProperties{
			ReceiveMaximum:       65535,
			MaximumQoS:           2,
			MaximumPacketSize:    4294967295,
			TopicAliasMaximum:    0,
			RetainAvailable:      true,
			WildcardSubAvailable: true,
			SubIDAvailable:       true,
			SharedSubAvailable:   true,
		},
		Persistence:   &noopPersistence{},
		MIDs:          &MIDs{index: make(map[uint16]*CPContext)},
		PacketTimeout: 10 * time.Second,
		Router: &StandardRouter{
			subscriptions: make(map[string][]MessageHandler),
		},
	}

	c.PingHandler = &PingHandler{
		pingFailHandler: func(e error) {
			c.Error(e)
		},
	}

	return c
}

// Connect is used to connect the client to a server. It presumes that
// the Client instance already has a working network connection.
// The function takes a pre-prepared Connect packet, and uses that to
// establish an MQTT connection. Assuming the connection completes
// successfully the rest of the client is initiated and the Connack
// returned. Otherwise the failure Connack (if there is one) is returned
// along with an error indicating the reason for the failure to connect.
func (c *Client) Connect(ctx context.Context, cp *Connect) (*Connack, error) {
	debug.Println("Connecting")
	c.Lock()
	defer c.Unlock()

	keepalive := cp.KeepAlive

	c.ClientID = cp.ClientID

	if c.Conn == nil {
		return nil, fmt.Errorf("Client has no connection")
	}

	debug.Println("Starting Incoming")
	c.workers.Add(1)
	go func() {
		defer c.workers.Done()
		c.Incoming()
	}()

	connCtx, cf := context.WithTimeout(ctx, c.PacketTimeout)
	defer cf()
	c.caCtx = &caContext{make(chan *packets.Connack, 1), connCtx}
	defer func() {
		c.caCtx = nil
	}()

	ccp := cp.Packet()

	ccp.ProtocolName = "MQTT"
	ccp.ProtocolVersion = 5

	debug.Println("Sending CONNECT")
	if _, err := ccp.WriteTo(c.Conn); err != nil {
		return nil, err
	}

	debug.Println("Waiting for CONNACK")
	var cap *packets.Connack
	select {
	case <-connCtx.Done():
		if e := connCtx.Err(); e == context.DeadlineExceeded {
			debug.Println("Timeout waiting for CONNACK")
			return nil, e
		}
	case cap = <-c.caCtx.Return:
	}

	ca := ConnackFromPacketConnack(cap)

	if ca.ReasonCode >= 0x80 {
		debug.Println("Received an error code in Connack:", ca.ReasonCode)
		return ca, fmt.Errorf("Failed to connect to server: %s", ca.Properties.ReasonString)
	}

	if ca.Properties.ServerKeepAlive != nil {
		keepalive = *ca.Properties.ServerKeepAlive
	}
	if ca.Properties.AssignedClientID != "" {
		c.ClientID = ca.Properties.AssignedClientID
	}
	if ca.Properties.ReceiveMaximum != nil {
		c.serverProps.ReceiveMaximum = *ca.Properties.ReceiveMaximum
	}
	if ca.Properties.MaximumQoS != nil {
		c.serverProps.MaximumQoS = *ca.Properties.MaximumQoS
	}
	if ca.Properties.MaximumPacketSize != nil {
		c.serverProps.MaximumPacketSize = *ca.Properties.MaximumPacketSize
	}
	if ca.Properties.TopicAliasMaximum != nil {
		c.serverProps.TopicAliasMaximum = *ca.Properties.TopicAliasMaximum
	}
	c.serverProps.RetainAvailable = ca.Properties.RetainAvailable
	c.serverProps.WildcardSubAvailable = ca.Properties.WildcardSubAvailable
	c.serverProps.SubIDAvailable = ca.Properties.SubIDAvailable
	c.serverProps.SharedSubAvailable = ca.Properties.SharedSubAvailable

	debug.Println("Received CONNACK, starting PingHandler")
	c.workers.Add(1)
	go func() {
		defer c.workers.Done()
		c.PingHandler.Start(c.Conn, time.Duration(keepalive)*time.Second)
	}()

	return ca, nil
}

// Incoming is the Client function that reads and handles incoming
// packets from the server. The function is started as a goroutine
// from Connect(), it exits when it receives a server initiated
// Disconnect, the Stop channel is closed or there is an error reading
// a packet from the network connection
func (c *Client) Incoming() {
	for {
		select {
		case <-c.stop:
			debug.Println("Client stopping, Incoming stopping")
			return
		default:
			recv, err := packets.ReadPacket(c.Conn)
			if err != nil {
				c.Error(err)
				return
			}
			debug.Println("Received a control packet:", recv.Type)
			switch recv.Type {
			case packets.CONNACK:
				cap := recv.Content.(*packets.Connack)
				if c.caCtx != nil {
					c.caCtx.Return <- cap
				}
			case packets.AUTH:
				ap := recv.Content.(*packets.Auth)
				switch ap.ReasonCode {
				case 0x0:
					if c.AuthHandler != nil {
						go c.AuthHandler.Authenticated()
					}
					if c.raCtx != nil {
						c.raCtx.Return <- *recv
					}
				case 0x18:
					if c.AuthHandler != nil {
						if _, err := c.AuthHandler.Authenticate(AuthFromPacketAuth(ap)).Packet().WriteTo(c.Conn); err != nil {
							c.Error(err)
							return
						}
					}
				}
			case packets.PUBLISH:
				pb := recv.Content.(*packets.Publish)
				go c.Router.Route(pb)
				switch pb.QoS {
				case 1:
					packets.NewPuback(packets.PubackFromPublish(pb)).WriteTo(c.Conn)
				case 2:
					packets.NewPubrec(packets.PubrecFromPublish(pb)).WriteTo(c.Conn)
				}
			case packets.PUBACK, packets.PUBCOMP, packets.SUBACK, packets.UNSUBACK:
				if cpCtx := c.MIDs.Get(recv.PacketID()); cpCtx != nil {
					cpCtx.Return <- *recv
				} else {
					debug.Println("Received a response for a message ID we don't know:", recv.PacketID())
				}
			case packets.PUBREC:
				if cpCtx := c.MIDs.Get(recv.PacketID()); cpCtx == nil {
					debug.Println("Received a PUBREC for a message ID we don't know:", recv.PacketID())
					packets.NewPubrel(
						packets.PubrelFromPubrec(recv.Content.(*packets.Pubrec)),
						packets.PubrelReasonCode(0x92),
					).WriteTo(c.Conn)
				} else {
					pr := recv.Content.(*packets.Pubrec)
					if pr.ReasonCode >= 0x80 {
						//Received a failure code, shortcut and return
						cpCtx.Return <- *recv
					} else {
						packets.NewPubrel(packets.PubrelFromPubrec(pr)).WriteTo(c.Conn)
					}
				}
			case packets.PUBREL:
				//Auto respond to pubrels unless failure code
				pr := recv.Content.(*packets.Pubrel)
				if pr.ReasonCode < 0x80 {
					//Received a failure code, continue
					continue
				} else {
					packets.NewPubcomp(packets.PubcompFromPubrel(pr)).WriteTo(c.Conn)
				}
			case packets.DISCONNECT:
				if c.OnDisconnect != nil {
					go c.OnDisconnect(*recv.Content.(*packets.Disconnect))
				}
				if c.raCtx != nil {
					c.raCtx.Return <- *recv
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
	debug.Println("Error called:", e)
	c.Lock()
	select {
	case <-c.stop:
		//already shutting down, do nothing
	default:
		close(c.stop)
	}
	c.Conn.Close()
	c.Unlock()
}

// Authenticate is used to initiate a reauthentication of credentials with the
// server. This function sends the initial Auth packet to start the reauthentication
// then relies on the client AuthHandler managing any further requests from the
// server until either a successful Auth packet is passed back, or a Disconnect
// is received.
func (c *Client) Authenticate(ctx context.Context, a *Auth) (*AuthResponse, error) {
	debug.Println("Client initiated reauthentication")
	c.Lock()
	defer c.Unlock()

	c.raCtx = &CPContext{make(chan packets.ControlPacket, 1), ctx}
	defer func() {
		c.raCtx = nil
	}()

	debug.Println("Sending AUTH")
	if _, err := a.Packet().WriteTo(c.Conn); err != nil {
		return nil, err
	}

	var rp packets.ControlPacket
	select {
	case <-ctx.Done():
		if e := ctx.Err(); e == context.DeadlineExceeded {
			debug.Println("Timeout waiting for Auth to complete")
			return nil, e
		}
	case rp = <-c.raCtx.Return:
	}

	switch rp.Type {
	case packets.AUTH:
		//If we've received one here it must be successful, the only way
		//to abort a reauth is a server initiated disconnect
		return AuthResponseFromPacketAuth(rp.Content.(*packets.Auth)), nil
	case packets.DISCONNECT:
		return AuthResponseFromPacketDisconnect(rp.Content.(*packets.Disconnect)), nil
	}

	return nil, fmt.Errorf("Error with Auth, didn't receive Auth or Disconnect")
}

// Subscribe is used to send a Subscription request to the MQTT server.
// It is passed a pre-prepared Subscribe packet and blocks waiting for
// a response Suback, or for the timeout to fire. Any reponse Suback
// is returned from the function, along with any errors.
func (c *Client) Subscribe(ctx context.Context, s *Subscribe) (*Suback, error) {
	debug.Printf("Subscribing to %+v", s.Subscriptions)
	subCtx, cf := context.WithTimeout(ctx, c.PacketTimeout)
	defer cf()
	cpCtx := &CPContext{make(chan packets.ControlPacket, 1), subCtx}

	sp := s.Packet()

	sp.PacketID = c.MIDs.Request(cpCtx)
	debug.Println("Sending SUBSCRIBE")
	if _, err := sp.WriteTo(c.Conn); err != nil {
		return nil, err
	}
	debug.Println("Waiting for SUBACK")
	var sap packets.ControlPacket

	select {
	case <-subCtx.Done():
		if e := subCtx.Err(); e == context.DeadlineExceeded {
			debug.Println("Timeout waiting for SUBACK")
			return nil, e
		}
	case sap = <-cpCtx.Return:
	}

	if sap.Type != packets.SUBACK {
		return nil, fmt.Errorf("Received %d instead of Suback", sap.Type)
	}
	debug.Println("Received SUBACK")

	sa := SubackFromPacketSuback(sap.Content.(*packets.Suback))
	switch {
	case len(sa.Reasons) == 1:
		if sa.Reasons[0] >= 0x80 {
			debug.Println("Received an error code in Suback:", sa.Reasons[0])
			return sa, fmt.Errorf("Failed to subscribe to topic: %s", sa.Properties.ReasonString)
		}
	default:
		for _, code := range sa.Reasons {
			if code >= 0x80 {
				debug.Println("Received an error code in Suback:", code)
				return sa, fmt.Errorf("At least one requested subscription failed")
			}
		}
	}

	return sa, nil
}

// Unsubscribe is used to send an Unsubscribe request to the MQTT server.
// It is passed a pre-prepared Unsubscribe packet and blocks waiting for
// a response Unsuback, or for the timeout to fire. Any reponse Unsuback
// is returned from the function, along with any errors.
func (c *Client) Unsubscribe(ctx context.Context, u *Unsubscribe) (*Unsuback, error) {
	debug.Printf("Unsubscribing from %+v", u.Topics)
	unsubCtx, cf := context.WithTimeout(ctx, c.PacketTimeout)
	defer cf()
	cpCtx := &CPContext{make(chan packets.ControlPacket, 1), unsubCtx}

	up := u.Packet()

	up.PacketID = c.MIDs.Request(cpCtx)
	debug.Println("Sending UNSUBSCRIBE")
	if _, err := up.WriteTo(c.Conn); err != nil {
		return nil, err
	}
	debug.Println("Waiting for UNSUBACK")
	var uap packets.ControlPacket

	select {
	case <-unsubCtx.Done():
		if e := unsubCtx.Err(); e == context.DeadlineExceeded {
			debug.Println("Timeout waiting for UNSUBACK")
			return nil, e
		}
	case uap = <-cpCtx.Return:
	}

	if uap.Type != packets.UNSUBACK {
		return nil, fmt.Errorf("Received %d instead of Unsuback", uap.Type)
	}
	debug.Println("Received SUBACK")

	ua := UnsubackFromPacketUnsuback(uap.Content.(*packets.Unsuback))
	switch {
	case len(ua.Reasons) == 1:
		if ua.Reasons[0] >= 0x80 {
			debug.Println("Received an error code in Suback:", ua.Reasons[0])
			return ua, fmt.Errorf("Failed to unsubscribe from topic: %s", ua.Properties.ReasonString)
		}
	default:
		for _, code := range ua.Reasons {
			if code >= 0x80 {
				debug.Println("Received an error code in Suback:", code)
				return ua, fmt.Errorf("At least one requested unsubscribe failed")
			}
		}
	}

	return ua, nil
}

// Publish is used to send a publication to the MQTT server.
// It is passed a pre-prepared Publish packet and blocks waiting for
// the appropriate response, or for the timeout to fire.
// Any reponse message is returned from the function, along with any errors.
func (c *Client) Publish(ctx context.Context, p *Publish) (*PublishResponse, error) {
	debug.Printf("Sending message to %s", p.Topic)
	c.Lock()
	defer c.Unlock()

	pb := p.Packet()

	switch p.QoS {
	case 0:
		debug.Println("Sending QoS0 message")
		if _, err := pb.WriteTo(c.Conn); err != nil {
			return nil, err
		}
		return nil, nil
	case 1, 2:
		return c.publishQoS12(ctx, pb)
	}

	return nil, fmt.Errorf("oops")
}

func (c *Client) publishQoS12(ctx context.Context, pb *packets.Publish) (*PublishResponse, error) {
	debug.Println("Sending QoS1 message")
	pubCtx, cf := context.WithTimeout(ctx, c.PacketTimeout)
	defer cf()
	cpCtx := &CPContext{make(chan packets.ControlPacket, 1), pubCtx}

	pb.PacketID = c.MIDs.Request(cpCtx)
	if _, err := pb.WriteTo(c.Conn); err != nil {
		return nil, err
	}
	var resp packets.ControlPacket

	select {
	case <-pubCtx.Done():
		if e := pubCtx.Err(); e == context.DeadlineExceeded {
			debug.Println("Timeout waiting for Publish response")
			return nil, e
		}
	case resp = <-cpCtx.Return:
	}

	switch pb.QoS {
	case 1:
		if resp.Type != packets.PUBACK {
			return nil, fmt.Errorf("Received %d instead of PUBACK", resp.Type)
		}
		debug.Println("Received PUBACK for", pb.PacketID)

		pr := PublishResponseFromPuback(resp.Content.(*packets.Puback))
		if pr.ReasonCode >= 0x80 {
			debug.Println("Received an error code in Puback:", pr.ReasonCode)
			return pr, fmt.Errorf("Error publishing: %s", resp.Content.(*packets.Puback).Reason())
		}
		return pr, nil
	case 2:
		switch resp.Type {
		case packets.PUBCOMP:
			debug.Println("Received PUBCOMP for", pb.PacketID)
			pr := PublishResponseFromPubcomp(resp.Content.(*packets.Pubcomp))
			return pr, nil
		case packets.PUBREC:
			debug.Printf("Received PUBREC for %s (must have errored)", pb.PacketID)
			pr := PublishResponseFromPubrec(resp.Content.(*packets.Pubrec))
			return pr, nil
		default:
			return nil, fmt.Errorf("Received %d instead of PUBCOMP", resp.Type)
		}
	}

	debug.Println("Ended up with a non QoS1/2 message:", pb.QoS)
	return nil, fmt.Errorf("Ended up with a non QoS1/2 message: %d", pb.QoS)
}

// Disconnect is used to send a Disconnect packet to the MQTT server
// Whether or not the attempt to send the Disconnect packet fails
// (and if it does this function returns any error) the network connection
// is closed.
func (c *Client) Disconnect(d *Disconnect) error {
	debug.Println("Disconnecting")
	c.Lock()
	defer c.Unlock()
	defer c.Conn.Close()

	_, err := d.Packet().WriteTo(c.Conn)

	return err
}
