package login

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/arukim/galaxy/network"

	"golang.org/x/net/websocket"
)

type Room struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type GameServer struct {
	pattern string
	rooms   []Room
	router  *network.Router
}

func NewGameServer(pattern string) *GameServer {
	s := GameServer{
		pattern: pattern,
	}

	s.rooms = []Room{{Name: "red room", ID: 0}}

	s.router = network.NewRouter([]network.CommandHandler{
		{
			Name: "join",
			Handler: func(Data *json.RawMessage) []byte {
				resp, _ := json.Marshal("joined the room")
				return resp
			},
		},
		{
			Name: "rooms",
			Handler: func(Data *json.RawMessage) []byte {
				resp, _ := json.Marshal(s.rooms)
				return resp
			},
		},
	})

	return &s
}

func (gs *GameServer) Listen() {
	log.Printf("GameServer is started on %v\n", gs.pattern)

	http.Handle(gs.pattern, websocket.Handler(gs.router.Listen))
}
