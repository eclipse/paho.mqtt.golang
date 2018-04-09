package paho

import (
	"crypto/tls"
	"net"
	"time"

	p "github.com/eclipse/paho.mqtt.golang/packets"
)

// WithConn sets the internal connection for the Client to the net.Conn passed in
func WithConn(conn net.Conn) func(*Client) error {
	return func(c *Client) error {
		c.Conn = conn
		return nil
	}
}

// OpenTCPConn will try and create a TCP connection to be used by the client
// the 'server' string is in the go style of "<server>:<port>"
func OpenTCPConn(server string) func(*Client) error {
	return func(c *Client) error {
		conn, err := net.Dial("tcp", server)
		if err != nil {
			return err
		}
		c.Conn = conn
		return nil
	}
}

// OpenTCP6Conn will try and create a TCP6 connection to be used by the client
// the 'server' string is in the go style of "<server>:<port>"
func OpenTCP6Conn(server string) func(*Client) error {
	return func(c *Client) error {
		conn, err := net.Dial("tcp6", server)
		if err != nil {
			return err
		}
		c.Conn = conn
		return nil
	}
}

// OpenTLSConn will try and create a TCP connection to be used by the client,
// and use the TLS config from 'config' to negotiate a secure connection,
// the 'server' string is in the go style of "<server>:<port>"
func OpenTLSConn(server string, config *tls.Config) func(*Client) error {
	return func(c *Client) error {
		conn, err := tls.Dial("tcp", server, config)
		if err != nil {
			return err
		}
		c.Conn = conn
		return nil
	}
}

// OpenTLS6Conn will try and create a TCP6 connection to be used by the client,
// and use the TLS config from 'config' to negotiate a secure connection,
// the 'server' string is in the go style of "<server>:<port>"
func OpenTLS6Conn(server string, config *tls.Config) func(*Client) error {
	return func(c *Client) error {
		conn, err := tls.Dial("tcp6", server, config)
		if err != nil {
			return err
		}
		c.Conn = conn
		return nil
	}
}

// WithPinger will set the client's internal ping mechanism to be the Pinger
// passed into this function, if this is not used the default Pinger provided
// in this library will be used.
func WithPinger(p Pinger) func(*Client) error {
	return func(c *Client) error {
		c.PingHandler = p
		return nil
	}
}

// WithRouter will set the client's internal message router to be the Router
// passed into this function, if this is not used the default Router provided
// in this library will be used.
func WithRouter(r Router) func(*Client) error {
	return func(c *Client) error {
		c.Router = r
		return nil
	}
}

// DefaultMessageHandler sets the client's internal message routing to always
// invoke the MessageHandler passed to this function when the client receives
// a message. This cannot be combined with other Routers, the last function
// to set a Router will be the configuration used by the client.
func DefaultMessageHandler(h MessageHandler) func(*Client) error {
	return func(c *Client) error {
		c.Router = &SingleHandlerRouter{
			messageHandler: h,
		}
		return nil
	}
}

// WithMIDService sets the client's internal message ID service to be the
// MIDService passed into this function. If this option is not used the default
// MIDService provided by this library will be used.
func WithMIDService(m MIDService) func(*Client) error {
	return func(c *Client) error {
		c.MIDs = m
		c.MIDs.Clear()
		return nil
	}
}

// WithPersistence sets the client's internal message persistence mechanism
// to be the Persistence passed into this function. If this option is not used
// *no* persistence mechanism is used. This library currently also provides a
// memory based Persistence.
func WithPersistence(p Persistence) func(*Client) error {
	return func(c *Client) error {
		c.Persistence = p
		return nil
	}
}

// PacketTimeout set's the client's internal timeout value for sending a message
// this value applies to the time to both send and receive/confirm a message.
// ie; this same value applies to sending a QoS1 message as sending and completing
// the fourway handshake for a QoS2 message.
// There is an exception to this that the timeout does not apply to the initial
// Connect call.
func PacketTimeout(d time.Duration) func(*Client) error {
	return func(c *Client) error {
		c.PacketTimeout = d
		return nil
	}
}

// DisconnectFunc sets a function to be called by the client when it is sent
// a DISCONNECT packet by the server.
func DisconnectFunc(d func(p.Disconnect)) func(*Client) error {
	return func(c *Client) error {
		c.Disconnected = d
		return nil
	}
}
