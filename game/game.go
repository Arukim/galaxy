package game

import (
	"sync"
	"time"

	"github.com/arukim/galaxy/helpers"
)

// Game contains game information
type Game struct {
	Map     *Map
	Players []*Player
	Mutex   *sync.Mutex

	maxPlayers, maxTurns int
	mapWidth, mapHeight  int
	currTurn             int

	playersCount      int
	turnCh            chan *PlayerTurn
	turnProccessCount chan bool

	turnTimeout time.Duration
}

// NewGame const for game
func NewGame(maxPlayers int, maxTurns int, turnTimeout time.Duration) *Game {
	g := new(Game)
	g.maxPlayers = maxPlayers
	g.maxTurns = maxTurns
	g.turnTimeout = turnTimeout
	g.turnCh = make(chan *PlayerTurn, maxPlayers)
	g.turnProccessCount = make(chan bool, maxPlayers)

	g.Players = make([]*Player, maxPlayers)
	g.Map = NewMap(10, 10)

	g.Mutex = &sync.Mutex{}

	return g
}

// AddPlayer adds new player into the game
func (g *Game) AddPlayer(conn IPlayerConnection) {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	g.Players[g.playersCount] = NewPlayer(g.turnCh, conn, g.playersCount)
	g.playersCount++
}

// MakeTurn Proccess one game turn
func (g *Game) MakeTurn() {
	tInfo := TurnInfo{turn: g.currTurn}

	// send turn info to players
	for _, p := range g.Players {
		p.StartTurn()
		go p.conn.Send(&tInfo)
	}
	// wait until all make turn or timeout
	var timeout = helpers.NewTimeout(g.turnTimeout)
	receivedTurns := 0

mainReceiveLoop:
	// or all players will send turn or timeout should work
	for receivedTurns < g.playersCount {
		select {
		case t := <-g.turnCh:
			go g.makePlayerTurn(t)
			receivedTurns++
		case <-timeout.Alarm:
			break mainReceiveLoop
		}
	}

	for _, p := range g.Players {
		p.EndTurn()
	}

	// now when all players are locked, we can clean turnCh channel
cleanUpLoop:
	for {
		select {
		case t := <-g.turnCh:
			go g.makePlayerTurn(t)
			receivedTurns++
		default:
			break cleanUpLoop
		}
	}

	// now wait until all player turns are made
	for receivedTurns > 0 {
		select {
		case <-g.turnProccessCount:
			receivedTurns--
		}
	}

	// all player turns are applied, calculate turn end

	g.currTurn++
}

func (g *Game) makePlayerTurn(t *PlayerTurn) {
	g.turnProccessCount <- true
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
