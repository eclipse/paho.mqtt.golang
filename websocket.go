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
	"net/http"
	"net/url"
)

// WebsocketOptions are config options for a websocket dialer
type WebsocketOptions struct {
	ReadBufferSize  int
	WriteBufferSize int
	Proxy           ProxyFunction
}

type ProxyFunction func(req *http.Request) (*url.URL, error)
