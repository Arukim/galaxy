package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"

	"golang.org/x/net/websocket"
)

// Router is set of command handlers
type Router struct {
	handlers        map[string]*CommandHandler
	conn            *websocket.Conn
	client          *Client
	allowedCommands string
	exit            chan bool

	handlersLock *sync.Mutex
}

// NewRouter creates new router instance
func NewRouter(handlers []*CommandHandler, c *Client) *Router {
	r := &Router{}

	r.exit = make(chan bool)
	r.handlersLock = &sync.Mutex{}
	r.client = c

	r.SetHandlers(handlers)

	return r
}

// Send data to client
func (r *Router) Send(cmd string, data interface{}) {
	com := Command{
		Cmd:  cmd,
		Data: data,
	}

	resp, err := json.Marshal(com)

	if err != nil {
		panic("can't generate packet")
	}

	r.conn.Write(resp)
}

// SetHandlers set new handlers for the router
func (r *Router) SetHandlers(handlers []*CommandHandler) {
	r.handlersLock.Lock()
	defer r.handlersLock.Unlock()

	r.handlers = make(map[string]*CommandHandler)

	var buf bytes.Buffer

	for _, h := range handlers {
		r.handlers[h.Name] = h

		buf.WriteString(h.Name)
		buf.WriteString(", ")
	}
	r.allowedCommands = buf.String()
}

// Listen handles web socket read
func (r *Router) Listen(conn *websocket.Conn) {
	GameStats.OnClientConnected()
	log.Println("client connected")
	r.conn = conn

	for true {
		var res *Result
		var cmd CommandJSON
		// read command from socket
		err := websocket.JSON.Receive(r.conn, &cmd)
		// check for disconnect
		if err == io.EOF {
			log.Println("client disconnected")
			GameStats.OnClientDisconnected()
			return
		} else if err != nil {
			// handle know errors, e.g. bad json in request
			if strings.HasPrefix(err.Error(), "invalid character") {
				res = &Result{
					Status: 400,
					Data:   err.Error(),
				}

				r.reply(cmd.Cmd, res)
				continue
			}
			// Unknown error gonna crush the socket
			res = &Result{
				Status: 500,
				Data:   err.Error(),
			}

			r.reply(cmd.Cmd, res)
			log.Printf("Unhandled error %v", err.Error())
			continue
		}

		// log command
		log.Printf("received %v err %v\n", cmd.Cmd, err)

		// handlers can change over the time, so use lock
		r.handlersLock.Lock()
		handler, ok := r.handlers[cmd.Cmd]
		var allowedCommands = r.allowedCommands
		r.handlersLock.Unlock()

		// if we found handler - execute
		if ok {
			log.Printf("handling %v", cmd.Cmd)
			res = handler.Handle(cmd.Data, r.client)
			// else provide some data about problem
		} else {
			res = &Result{
				Status: 404,
				Data:   fmt.Sprintf(`Unknown command. Allowed command are: %v. Read examples at github.com/arukim/galaxy`, allowedCommands),
			}
		}

		r.reply(cmd.Cmd, res)
	}
}

// reply sends response into socket
func (r *Router) reply(cmd string, result *Result) {
	// create response struct
	response := ResponseJSON{
		Cmd:    cmd,
		Result: result,
	}

	// try to marshall
	resp, err := json.Marshal(response)
	if err != nil {
		// wtf something gone wrong, let's return 500
		response.Result.Status = 500
		response.Result.Data = "Fatal server error"
		if resp, err = json.Marshal(response); err != nil {
			// we can't even generate 500? time to panic
			panic("can't generate 500 error")
		}
	}

	// write the response
	r.conn.Write(resp)
}
