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
	players        []*core.Player
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

	r.players = []*core.Player{}
	r.lock = &sync.Mutex{}

	return &r
}

// AddPlayer into game room
func (r *Room) AddPlayer(p *core.Player) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.players = append(r.players, p)
	log.Println(r.players)
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

		g.AddPlayers(r.players[:r.PlayersPerGame])

		r.players = r.players[r.PlayersPerGame:]
		r.PlayersCount -= r.PlayersPerGame

		go g.Start()
	}
}
