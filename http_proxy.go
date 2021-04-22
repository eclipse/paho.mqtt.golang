package mqtt

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"
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
	if resp.StatusCode != 200 {
		_ = c.Close()
		err = fmt.Errorf("connect server using proxy error, StatusCode [%d]", resp.StatusCode)
		return nil, err
	}

	return c, nil
}

func init() {
	proxy.RegisterDialerType("http", newHTTPProxy)
	proxy.RegisterDialerType("https", newHTTPProxy)
}
