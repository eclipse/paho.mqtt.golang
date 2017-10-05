package paho

import (
	"log"
	"net"
	"sync"
	"time"
)

type Client struct {
	sync.Mutex
	conn        net.Conn
	PingTimeout time.Duration
	Pinger      Pinger
}

func NewClient(c net.Conn, pt time.Duration) *Client {
	cl := &Client{
		conn:   c,
		Pinger: NewPingHander(c, pt, PFH),
	}

	return cl
}

func PFH(err error) {
	log.Fatalln(err)
}
