package network

import "golang.org/x/net/websocket"

// WebClient stoers websocket and current routing data
type WebClient struct {
	router Router
	ws     *websocket.Conn
}

// NewWebClient creates
func NewWebClient(ws *websocket.Conn, r Router) *WebClient {
	c := WebClient{
		ws:     ws,
		router: r,
	}
	return &c
}
