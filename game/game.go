package game

import (
	"sync"
	"time"
)

// Game contains game information
type Game struct {
	Map     *Map
	Players []*Player
	Mutex   *sync.Mutex

	maxPlayers int
	maxTurns   int
	mapWidth   int
	mapHeight  int

	playersCount int

	turnTimeout time.Duration
}

// Player is game player
type Player struct {
}

// NewGame const for game
func NewGame(maxPlayers int, maxTurns int, turnTimeout time.Duration) *Game {
	g := new(Game)
	g.maxPlayers = maxPlayers
	g.maxTurns = maxTurns
	g.turnTimeout = turnTimeout

	g.Players = make([]*Player, maxPlayers)
	g.Map = NewMap(10, 10)

	g.Mutex = &sync.Mutex{}

	return g
}

// AddPlayer adds new player into the game
func (g *Game) AddPlayer(p *Player) {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	g.Players[g.playersCount] = p
	g.playersCount++
}

// Start func to start the game
func (g *Game) Start() {
	g.Map.Init()
}
