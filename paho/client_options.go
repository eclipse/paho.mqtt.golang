package paho

import (
	"crypto/tls"
	"net"
	"time"

	p "github.com/eclipse/paho.mqtt.golang/packets"
)

func WithConn(conn net.Conn) func(*Client) error {
	return func(c *Client) error {
		c.Conn = conn
		return nil
	}
}

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

func WithPinger(p Pinger) func(*Client) error {
	return func(c *Client) error {
		c.PingHandler = p
		return nil
	}
}

func WithRouter(r Router) func(*Client) error {
	return func(c *Client) error {
		c.Router = r
		return nil
	}
}

func DefaultMessageHandler(h MessageHandler) func(*Client) error {
	return func(c *Client) error {
		c.Router = &SingleHandlerRouter{
			messageHandler: h,
		}
		return nil
	}
}

func WithMIDService(m MIDService) func(*Client) error {
	return func(c *Client) error {
		c.MIDs = m
		c.MIDs.Clear()
		return nil
	}
}

func WithMemoryPersistence() func(*Client) error {
	return func(c *Client) error {
		c.Persistence = &MemoryPersistence{}
		return nil
	}
}

func PacketTimeout(d time.Duration) func(*Client) error {
	return func(c *Client) error {
		c.PacketTimeout = d
		return nil
	}
}

func DisconnectFunc(d func(p.Disconnect)) func(*Client) error {
	return func(c *Client) error {
		c.Disconnected = d
		return nil
	}
}
