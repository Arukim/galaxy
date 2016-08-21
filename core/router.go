package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/net/websocket"
)

// Router is set of command handlers
type Router struct {
	handlers        map[string]*CommandHandler
	player          *Player
	allowedCommands string
}

// NewRouter creates new router instance
func NewRouter(handlers []*CommandHandler, p *Player) *Router {
	r := &Router{}
	r.handlers = make(map[string]*CommandHandler)

	var buf bytes.Buffer

	for _, h := range handlers {
		r.handlers[h.Name] = h

		buf.WriteString(h.Name)
		buf.WriteString(", ")
	}
	r.allowedCommands = buf.String()

	return r
}

// Listen handles web socket read-write
func (r *Router) Listen(ws *websocket.Conn) {
	defer ws.Close()

	fmt.Println("client connected")
	for true {
		var cmd Command
		var resp []byte
		var res *Result
		err := websocket.JSON.Receive(ws, &cmd)

		if err == io.EOF {
			fmt.Println("client disconnected")
			return
		}

		fmt.Println(cmd)
		if handler, ok := r.handlers[cmd.Cmd]; ok {
			res = handler.Handle(cmd.Data, r.player)

			//resp, _ = json.Marshal("joined the room " + room.Name)
		} else {
			res = &Result{
				Status: 404,
				Data:   fmt.Sprintf(`Unknown command. Allowed command are: %v. Read examples at github.com/arukim/galaxy`, r.allowedCommands),
			}
		}

		response := ResponseJSON{
			Cmd:    cmd.Cmd,
			Result: res,
		}

		resp, err = json.Marshal(response)
		if err != nil {
			response.Result.Status = 500
			response.Result.Data = "Fatal server error"
			if resp, err = json.Marshal(response); err != nil {
				panic("can't generate 500 error")
			}
		}

		ws.Write(resp)
	}
}
