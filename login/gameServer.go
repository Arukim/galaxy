package login

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type Room struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type GameServer struct {
	pattern string
	rooms   []Room
}

func NewGameServer(pattern string) *GameServer {
	s := GameServer{
		pattern: pattern,
	}
	s.rooms = []Room{{Name: "red room", ID: 0}}
	return &s
}

func (gs *GameServer) onConnected(ws *websocket.Conn) {
	defer ws.Close()

	fmt.Println("client connected")
	for true {
		var msg string
		var resp []byte
		err := websocket.JSON.Receive(ws, &msg)

		if err == io.EOF {
			fmt.Println("client disconnected")
			return
		}

		fmt.Println(msg)
		switch msg {
		case "rooms":
			resp, _ = json.Marshal(gs.rooms)
		case "join":
			resp, _ = json.Marshal("connected")
		case "info":
			resp, _ = json.Marshal("Welcome galaxy game server")
		default:
			resp, _ = json.Marshal("Unknown command. Allowed commands: rooms, info. Use json.")
		}

		ws.Write(resp)
	}
}

func (gs *GameServer) Listen() {
	log.Printf("GameServer is started on %v\n", gs.pattern)

	http.Handle(gs.pattern, websocket.Handler(gs.onConnected))
}
