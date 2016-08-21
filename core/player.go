package core

import (
	"golang.org/x/net/websocket"
)

// WebClient stoers websocket and current routing data
type Player struct {
	router *Router
	ws     *websocket.Conn
}

// NewWebClient creates
func NewPlayer(ws *websocket.Conn, handlers []*CommandHandler) *Player {
	c := Player{
		ws: ws,
	}
	c.router = NewRouter(handlers, &c)
	c.router.Listen(ws)
	return &c
}
