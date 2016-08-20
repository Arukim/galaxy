package network

import (
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/net/websocket"
)

// Command is model for network commands
type Command struct {
	Cmd  string           `json:"cmd"`
	Data *json.RawMessage `json:"data"`
}

// CommandHandler handles one command
type CommandHandler struct {
	Name    string
	Handler func(Data *json.RawMessage) []byte
}

// Router is set of command handlers
type Router struct {
	handlers map[string]CommandHandler
}

func NewRouter(handlers []CommandHandler) *Router {
	r := &Router{}
	r.handlers = make(map[string]CommandHandler)

	for _, h := range handlers {
		r.handlers[h.Name] = h
	}

	return r
}

func (r *Router) Listen(ws *websocket.Conn) {
	defer ws.Close()

	fmt.Println("client connected")
	for true {
		var cmd Command
		var resp []byte
		err := websocket.JSON.Receive(ws, &cmd)

		if err == io.EOF {
			fmt.Println("client disconnected")
			return
		}

		fmt.Println(cmd)
		if handler, ok := r.handlers[cmd.Cmd]; ok {
			resp = handler.Handler(cmd.Data)
		} else {
			resp, _ = json.Marshal("Unknown command. Allowed commands: rooms, info. Use json.")
		}

		ws.Write(resp)
	}
}
