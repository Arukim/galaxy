package core

import (
	"golang.org/x/net/websocket"
)

// Client stoers websocket and current routing data
type Client struct {
	Name            string
	Router          *Router
	DefaultHandlers []*CommandHandler
}

// NewClient creates
func NewClient(ws *websocket.Conn, handlers []*CommandHandler) *Client {
	p := &Client{}

	p.DefaultHandlers = handlers
	p.Router = NewRouter(handlers, p)
	p.Router.Listen(ws)

	return p
}
