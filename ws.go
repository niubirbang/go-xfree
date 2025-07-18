package goxfree

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

type ws struct {
	dialer *websocket.Dialer
}

func (w *ws) conn(path string) (*websocket.Conn, error) {
	path = strings.TrimLeft(path, "/")
	u := url.URL{Scheme: "ws", Host: "unix", Path: fmt.Sprintf("/%s", path)}
	conn, resp, err := w.dialer.Dial(u.String(), http.Header{})
	if err != nil {
		if resp != nil {
			return nil, fmt.Errorf("ws response status: %s", resp.Status)
		}
		return nil, err
	}
	return conn, nil
}
