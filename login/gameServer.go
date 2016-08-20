package login

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/arukim/galaxy/network"

	"golang.org/x/net/websocket"
)

// Room each room runs one game
type Room struct {
	Name           string `json:"name"`
	ID             int    `json:"id"`
	PlayersPerGame int    `json:"playersPerGame"`
}

// GameServer is entry point for players
type GameServer struct {
	pattern string
	rooms   []Room
	clients []*network.WebClient
	router  *network.Router

	clientsLock *sync.Mutex
	roomsLock   *sync.Mutex
}

// NewGameServer creates new game server instance
func NewGameServer(pattern string) *GameServer {
	gs := GameServer{
		pattern: pattern,
	}

	gs.roomsLock = &sync.Mutex{}
	gs.clientsLock = &sync.Mutex{}
	gs.clients = make([]*network.WebClient, 1)

	gs.rooms = []Room{
		{
			Name:           "red room",
			ID:             0,
			PlayersPerGame: 1,
		},
		{
			Name:           "blue room",
			ID:             1,
			PlayersPerGame: 2,
		},
	}

	gs.router = network.NewRouter([]network.CommandHandler{
		{
			Name:    "join",
			Handler: gs.OnJoin,
		},
		{
			Name:    "rooms",
			Handler: gs.OnRooms,
		},
	})

	return &gs
}

// OnJoin - player is trying to enter room
func (gs *GameServer) OnJoin(d *json.RawMessage) []byte {

	resp, _ := json.Marshal("joined the room")
	return resp
}

// OnRooms - player is quering rooms
func (gs *GameServer) OnRooms(d *json.RawMessage) []byte {
	gs.roomsLock.Lock()
	defer gs.roomsLock.Unlock()

	resp, _ := json.Marshal(gs.rooms)
	return resp
}

// Listen starts the server
func (gs *GameServer) Listen() {
	log.Printf("GameServer is started on %v\n", gs.pattern)

	http.Handle(gs.pattern, websocket.Handler(func(ws *websocket.Conn) {
		client := network.NewWebClient(ws, r)

		gs.clientsLock.Lock()
		defer gs.clientsLock.Unlock()

		gs.clients = append(gs.clients, client)

		gs.router.Listen(ws)
	}))
}
