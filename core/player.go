package core

import (
	"golang.org/x/net/websocket"
)

// Player stoers websocket and current routing data
type Player struct {
	Name            string
	Router          *Router
	DefaultHandlers []*CommandHandler
}

// NewPlayer creates
func NewPlayer(ws *websocket.Conn, handlers []*CommandHandler) *Player {
	p := &Player{}

	p.DefaultHandlers = handlers
	p.Router = NewRouter(handlers, p)
	p.Router.Listen(ws)

	return p
}
