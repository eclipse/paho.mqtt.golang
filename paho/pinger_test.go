package paho

import (
	"bytes"
	"net"
	"time"
)

type DummyConn struct {
	bytes.Buffer
}

func (d *DummyConn) Close() error {
	return nil
}

func (d *DummyConn) LocalAddr() net.Addr {
	return &net.IPAddr{IP: net.IPv4(127, 0, 0, 1)}
}

func (d *DummyConn) RemoteAddr() net.Addr {
	return &net.IPAddr{IP: net.IPv4(127, 0, 0, 1)}
}

func (d *DummyConn) SetDeadline(t time.Time) error {
	return nil
}

func (d *DummyConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (d *DummyConn) SetWriteDeadline(t time.Time) error {
	return nil
}
