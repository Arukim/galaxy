package login

import (
	"log"
	"sync"
	"time"

	"github.com/arukim/galaxy/core"
	"github.com/arukim/galaxy/game"
)

// Room each room runs one game
type Room struct {
	Name           string `json:"name"`
	PlayersPerGame int    `json:"playersPerGame"`
	PlayersCount   int    `json:"playersCount"`
	clients        []*core.Client
	lock           *sync.Mutex
}

// RoomSettings contains setting for room creation
type RoomSettings struct {
	Name           string
	PlayersPerGame int
}

// NewRoom creates new room
func NewRoom(s RoomSettings) *Room {
	r := Room{
		Name:           s.Name,
		PlayersPerGame: s.PlayersPerGame,
	}

	r.clients = []*core.Client{}
	r.lock = &sync.Mutex{}

	return &r
}

// AddPlayer into game room
func (r *Room) AddClient(p *core.Client) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.clients = append(r.clients, p)
	log.Println(r.clients)
	r.PlayersCount++

	if r.PlayersCount >= r.PlayersPerGame {
		// lets start game
		s := game.ServerSettings{
			MaxPlayers:  r.PlayersPerGame,
			MaxTurns:    10,
			TurnTimeout: 1 * time.Second,
			MapWidth:    100,
			MapHeight:   100,
		}

		g := game.NewServer(&s)

		g.AddClients(r.clients[:r.PlayersPerGame])

		r.clients = r.clients[r.PlayersPerGame:]
		r.PlayersCount -= r.PlayersPerGame

		go g.Start()
	}
}
