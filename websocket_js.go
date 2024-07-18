//go:build js

/*
 * This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v2.0
 * and Eclipse Distribution License v1.0 which accompany this distribution.
 *
 * The Eclipse Public License is available at
 *    https://www.eclipse.org/legal/epl-2.0/
 * and the Eclipse Distribution License is available at
 *   http://www.eclipse.org/org/documents/edl-v10.php.
 *
 * Contributors:
 */

package mqtt

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"nhooyr.io/websocket"
)

// WebsocketOptions are config options for a websocket dialer
type WebsocketOptions struct {
	ReadBufferSize  int
	WriteBufferSize int
	Proxy           ProxyFunction
}

type ProxyFunction func(req *http.Request) (*url.URL, error)

// NewWebsocket returns a new websocket and returns a net.Conn compatible interface using the nhooyr.io/websocket package
func NewWebsocket(host string, tlsc *tls.Config, _ time.Duration, requestHeader http.Header, _ *WebsocketOptions) (net.Conn, error) {
	opts := websocket.DialOptions{
		Subprotocols: []string{"mqtt"},
	}
	ctx := context.Background()
	ws, resp, err := websocket.Dial(ctx, host, &opts)
	if err != nil {
		if resp != nil {
			WARN.Println(CLI, fmt.Sprintf("Websocket handshake failure. StatusCode: %d. Body: %s", resp.StatusCode, resp.Body))
		}
		return nil, err
	}

	return websocket.NetConn(ctx, ws, websocket.MessageBinary), err
}
