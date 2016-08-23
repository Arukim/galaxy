package game

import (
	"log"
	"time"

	"github.com/arukim/galaxy/core"
	"github.com/arukim/galaxy/helpers"
)

// Server (game) contains game information
type Server struct {
	Map     *Map
	Players []*core.Player

	//settings
	maxPlayers, maxTurns int
	mapWidth, mapHeight  int
	turnTimeout          time.Duration

	currTurn          int
	playersCount      int
	turnCh            chan *playerTurn
	turnProccessCount chan bool
}

// ServerSettings are Initial server settings
type ServerSettings struct {
	MaxPlayers, MaxTurns int
	MapWidth, MapHeight  int
	TurnTimeout          time.Duration
}

// NewServer const for game
func NewServer(settings *ServerSettings) *Server {
	s := Server{}

	s.maxPlayers = settings.MaxPlayers
	s.maxTurns = settings.MaxTurns
	s.turnTimeout = settings.TurnTimeout
	s.mapWidth = settings.MapWidth
	s.mapHeight = settings.MapHeight
	s.turnCh = make(chan *playerTurn, s.maxPlayers)
	s.turnProccessCount = make(chan bool, s.maxPlayers)

	s.Map = NewMap(10, 10)

	return &s
}

// AddPlayers populate game with players
func (s *Server) AddPlayers(players []*core.Player) {
	s.Players = players
	s.playersCount = len(players)
}

// MakeTurn Proccess one game turn
func (s *Server) MakeTurn() {
	log.Printf("Turn %d has started\n", s.currTurn)

	tInfo := turnInfo{Turn: s.currTurn}

	// send turn info to players
	s.broadcast("turnInfo", &tInfo)

	// wait until all make turn or timeout
	var timeout = helpers.NewTimeout(s.turnTimeout)
	receivedTurns := 0

mainReceiveLoop:
	// or all players will send turn or timeout should work
	for receivedTurns < s.playersCount {
		select {
		case t := <-s.turnCh:
			go s.makePlayerTurn(t)
			receivedTurns++
		case <-timeout.Alarm:
			break mainReceiveLoop
		}
	}
	/*
		for _, p := range g.Players {
					p.EndTurn()
		}*/

	// now when all players are locked, we can clean turnCh channel
cleanUpLoop:
	for {
		select {
		case t := <-s.turnCh:
			go s.makePlayerTurn(t)
			receivedTurns++
		default:
			break cleanUpLoop
		}
	}

	// now wait until all player turns are made
	for receivedTurns > 0 {
		select {
		case <-s.turnProccessCount:
			receivedTurns--
		}
	}

	// all player turns are applied, calculate turn end

	s.currTurn++
}

func (s *Server) makePlayerTurn(t *playerTurn) {
	s.turnProccessCount <- true
}

func (s *Server) broadcast(com string, data interface{}) {
	for _, p := range s.Players {
		go p.Router.Send(com, data)
	}
}

// Start func to start the game
func (s *Server) Start() {
	s.Map.Init()
	players := make([]string, s.playersCount)
	for i, p := range s.Players {
		players[i] = p.Name
		p.Router.SetHandlers([]*core.CommandHandler{})
	}

	gInfo := gameInfo{
		MaxTurns:  s.maxTurns,
		MapWidth:  s.mapWidth,
		MapHeight: s.mapHeight,
		Players:   players,
	}

	s.broadcast("gameInfo", &gInfo)

	for true {
		if s.currTurn < s.maxTurns {
			s.MakeTurn()
		} else {
			log.Println("game has ended")
			gameInfo := gameResult{Winner: "no one"}
			s.broadcast("gameResult", &gameInfo)

			for _, p := range s.Players {
				p.Router.SetHandlers(p.DefaultHandlers)
			}
			return
		}
	}
}
