package login

import (
	"sync"

	"github.com/arukim/galaxy/core"
)

// Room each room runs one game
type Room struct {
	Name           string `json:"name"`
	PlayersPerGame int    `json:"playersPerGame"`
	PlayersCount   int
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

	return &r
}

// AddPlayer into game room
func (r *Room) AddPlayer(p *core.Player) {

}
