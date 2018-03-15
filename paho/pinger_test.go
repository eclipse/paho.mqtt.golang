package paho

import (
	"bytes"
	"net"
	"testing"
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

func Test_pingHandler_Start(t *testing.T) {
	type fields struct {
		lastPing        time.Time
		pingOutstanding int32
		pingFailHandler PingFailHandler
	}
	type args struct {
		c  net.Conn
		pt time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"1", fields{time.Now(), 0, nil}, args{&DummyConn{}, 10 * time.Second}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pingHandler{
				stop:            make(chan struct{}),
				lastPing:        tt.fields.lastPing,
				pingOutstanding: tt.fields.pingOutstanding,
			}
			p.Start(tt.args.c, tt.args.pt)
		})
	}
}
