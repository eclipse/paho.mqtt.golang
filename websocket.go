package mqtt

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// NewWebsocket creates a new websocket and returns a net.Conn compatiable interface.
func NewWebsocket(host string, tlsc *tls.Config, timeout time.Duration, requestHeader http.Header) (net.Conn, error) {
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	dialer := &websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		HandshakeTimeout:  timeout,
		EnableCompression: false,
		TLSClientConfig:   tlsc,
	}
	host = strings.Replace(host, "wss2", "wss", -1)
	ws, _, err := dialer.Dial(host, requestHeader)

	wrapper := &websocketConnector{
		Conn: ws,
	}
	return wrapper, err
}

/*
	Supporting the net.Conn interface for the gorilla/websocket library
	as suggested in issue: https://github.com/gorilla/websocket/issues/282
*/
type websocketConnector struct {
	*websocket.Conn
	r   io.Reader
	rio sync.Mutex
	wio sync.Mutex
}

func (c *websocketConnector) SetDeadline(t time.Time) error {
	if err := c.SetReadDeadline(t); err != nil {
		return err
	}
	err := c.SetWriteDeadline(t)
	return err
}

func (c *websocketConnector) Write(p []byte) (int, error) {
	c.wio.Lock()
	defer c.wio.Unlock()

	err := c.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (c *websocketConnector) Read(p []byte) (int, error) {
	c.rio.Lock()
	defer c.rio.Unlock()
	for {
		if c.r == nil {
			// Advance to next message.
			var err error
			_, c.r, err = c.NextReader()
			if err != nil {
				return 0, err
			}
		}
		n, err := c.r.Read(p)
		if err == io.EOF {
			// At end of message.
			c.r = nil
			if n > 0 {
				return n, nil
			}
			// No data read, continue to next message.
			continue
		}
		return n, err
	}
}
