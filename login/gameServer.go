package login

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/arukim/galaxy/core"

	q "github.com/ahmetalpbalkan/go-linq"
	"golang.org/x/net/websocket"
)

// GameServer is entry point for players
type GameServer struct {
	pattern string
	rooms   []*Room
	players []*core.Player

	clientsLock *sync.Mutex
}

// NewGameServer creates new game server instance
func NewGameServer(pattern string) *GameServer {
	gs := GameServer{
		pattern: pattern,
	}

	gs.clientsLock = &sync.Mutex{}
	gs.players = []*core.Player{}

	gs.rooms = []*Room{
		NewRoom(RoomSettings{Name: "red room", PlayersPerGame: 1}),
		NewRoom(RoomSettings{Name: "blue room", PlayersPerGame: 2}),
	}

	return &gs
}

// OnJoin - player is trying to enter room
func (gs *GameServer) OnJoin(d *json.RawMessage, p *core.Player) *core.Result {
	var resp *core.Result
	var roomName string
	json.Unmarshal(*d, &roomName)

	log.Printf("Request to join room %v", roomName)

	r, found, err := q.From(gs.rooms).
		FirstBy(func(i q.T) (bool, error) {
			return i.(*Room).Name == roomName, nil
		})

	if err == nil && found {
		room := r.(*Room)
		room.AddPlayer(p)
		resp = core.NewSuccessResult("joined the room " + room.Name)
	} else {
		resp = core.NewErrorResult(fmt.Sprintf("room %s not found", roomName))
	}

	return resp
}

// OnRooms - player is quering rooms
func (gs *GameServer) OnRooms(d *json.RawMessage, p *core.Player) *core.Result {
	return core.NewSuccessResult(gs.rooms)
}

// Listen starts the server
func (gs *GameServer) Listen() {
	log.Printf("GameServer is started on %v\n", gs.pattern)

	http.Handle(gs.pattern, websocket.Handler(func(ws *websocket.Conn) {
		player := core.NewPlayer(ws, []*core.CommandHandler{
			{
				Name:   "join",
				Handle: gs.OnJoin,
			},
			{
				Name:   "rooms",
				Handle: gs.OnRooms,
			},
		})

		gs.clientsLock.Lock()
		defer gs.clientsLock.Unlock()

		gs.players = append(gs.players, player)
	}))
}
