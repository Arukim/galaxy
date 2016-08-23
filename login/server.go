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

// Server (login) is entry point for players
type Server struct {
	pattern string
	rooms   []*Room
	players []*core.Player

	clientsLock *sync.Mutex
}

// NewServer creates new game server instance
func NewServer(pattern string) *Server {
	s := Server{
		pattern: pattern,
	}

	s.clientsLock = &sync.Mutex{}
	s.players = []*core.Player{}

	s.rooms = []*Room{
		NewRoom(RoomSettings{Name: "red room", PlayersPerGame: 1}),
		NewRoom(RoomSettings{Name: "blue room", PlayersPerGame: 2}),
	}

	return &s
}

// OnJoin - player is trying to enter room
func (s *Server) OnJoin(d *json.RawMessage, p *core.Player) *core.Result {
	var resp *core.Result
	var roomName string
	json.Unmarshal(*d, &roomName)

	log.Printf("Request to join room %v", roomName)

	r, found, err := q.From(s.rooms).
		FirstBy(func(i q.T) (bool, error) {
			return i.(*Room).Name == roomName, nil
		})

	if err == nil && found {
		room := r.(*Room)
		go room.AddPlayer(p)
		resp = core.NewSuccessResult("joined the room " + room.Name)
	} else {
		resp = core.NewErrorResult(fmt.Sprintf("room %s not found", roomName))
	}

	return resp
}

// OnRooms - player is quering rooms
func (s *Server) OnRooms(d *json.RawMessage, p *core.Player) *core.Result {
	return core.NewSuccessResult(s.rooms)
}

// Handlers returns GameServer handlers
func (s *Server) Handlers() []*core.CommandHandler {
	return []*core.CommandHandler{
		{
			Name:   "join",
			Handle: s.OnJoin,
		},
		{
			Name:   "rooms",
			Handle: s.OnRooms,
		},
	}
}

// Listen starts the server
func (s *Server) Listen() {
	log.Printf("GameServer is started on %v\n", s.pattern)

	http.Handle(s.pattern, websocket.Handler(func(ws *websocket.Conn) {
		player := core.NewPlayer(ws, s.Handlers())

		s.clientsLock.Lock()
		defer s.clientsLock.Unlock()

		s.players = append(s.players, player)
	}))
}
