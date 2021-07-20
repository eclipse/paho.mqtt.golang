/*
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v2.0
 * and Eclipse Distribution License v1.0 which accompany this distribution.
 *
 * The Eclipse Public License is available at
 *    https://www.eclipse.org/legal/epl-2.0/
 * and the Eclipse Distribution License is available at
 *   http://www.eclipse.org/org/documents/edl-v10.php.
 */

package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
)

// httpProxy is a HTTP/HTTPS connect capable proxy.
type httpProxy struct {
	host     string
	haveAuth bool
	username string
	password string
	forward  proxy.Dialer
}

func (s httpProxy) String() string {
	return fmt.Sprintf("HTTP proxy dialer for %s", s.host)
}

func newHTTPProxy(uri *url.URL, forward proxy.Dialer) (proxy.Dialer, error) {
	s := new(httpProxy)
	s.host = uri.Host
	s.forward = forward
	if uri.User != nil {
		s.haveAuth = true
		s.username = uri.User.Username()
		s.password, _ = uri.User.Password()
	}

	return s, nil
}

func (s *httpProxy) Dial(_, addr string) (net.Conn, error) {
	reqURL := url.URL{
		Scheme: "https",
		Host:   addr,
	}

	req, err := http.NewRequest("CONNECT", reqURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Close = false
	if s.haveAuth {
		req.SetBasicAuth(s.username, s.password)
	}
	req.Header.Set("User-Agent", "paho.mqtt")

	// Dial and create the client connection.
	c, err := s.forward.Dial("tcp", s.host)
	if err != nil {
		return nil, err
	}

	err = req.Write(c)
	if err != nil {
		_ = c.Close()
		return nil, err
	}

	resp, err := http.ReadResponse(bufio.NewReader(c), req)
	if err != nil {
		_ = c.Close()
		return nil, err
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		_ = c.Close()
		return nil, fmt.Errorf("proxied connection returned an error: %v", resp.Status)
	}

	return c, nil
}
