//go:build !js
// +build !js

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
	"time"

	"nhooyr.io/websocket"
)

// NewWebsocket returns a new websocket and returns a net.Conn compatible interface using the gorilla/websocket package
func NewWebsocket(host string, tlsc *tls.Config, timeout time.Duration, requestHeader http.Header, options *WebsocketOptions) (net.Conn, error) {
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	if options == nil {
		// Apply default options
		options = &WebsocketOptions{}
	}

	httpTransport := &http.Transport{TLSClientConfig: tlsc, Proxy: options.Proxy}
	httpClient := &http.Client{Transport: httpTransport, Timeout: timeout}

	dialOptions := websocket.DialOptions{
		HTTPClient:   httpClient,
		HTTPHeader:   requestHeader,
		Subprotocols: []string{"mqtt"},
	}

	ctx := context.Background()

	ws, resp, err := websocket.Dial(ctx, host, &dialOptions)

	if err != nil {
		if resp != nil {
			WARN.Println(CLI, fmt.Sprintf("Websocket handshake failure. StatusCode: %d. Body: %s", resp.StatusCode, resp.Body))
		}
		return nil, err
	}

	return websocket.NetConn(ctx, ws, websocket.MessageBinary), nil
}
