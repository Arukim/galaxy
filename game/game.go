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

	maxPlayers, maxTurns int
	mapWidth, mapHeight  int
	currTurn             int

	playersCount int

	turnTimeout time.Duration
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
func (g *Game) AddPlayer(conn IPlayerConnection) {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	p := Player{conn: conn}
	g.Players[g.playersCount] = &p
	g.playersCount++
}

// MakeTurn Proccess one game turn
func (g *Game) MakeTurn() {
	tInfo := TurnInfo{turn: g.currTurn}
	// send turn info to players
	for _, p := range g.Players {
		go p.conn.Send(&tInfo)
	}
	// wait until all make turn or timeout
	// calculate turn

	g.currTurn++
}

// Start func to start the game
func (g *Game) Start() {
	g.Map.Init()
	players := make([]string, g.playersCount)
	for i, p := range g.Players {
		players[i] = p.name
	}

	gInfo := GameInfo{
		maxTurns:  g.maxTurns,
		mapWidth:  g.mapWidth,
		mapHeight: g.mapHeight,
		players:   players,
	}

	for _, p := range g.Players {
		go p.conn.Send(gInfo)
	}

	go func() {
		if g.currTurn < g.maxTurns {
			g.MakeTurn()
		} else {
			// Find
		}
	}()
}
