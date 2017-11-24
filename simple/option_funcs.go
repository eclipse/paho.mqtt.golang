package simple

import (
	"net"
)

func WithConn(conn net.Conn) func(*Client) error {
	return func(c *Client) error {
		c.Conn = conn
		return nil
	}
}

func OpenConn(network, server string) func(*Client) error {
	return func(c *Client) error {
		conn, err := net.Dial(network, server)
		if err != nil {
			return err
		}
		c.Conn = conn
		return nil
	}
}
