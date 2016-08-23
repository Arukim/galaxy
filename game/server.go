package game

import (
	"log"
	"time"

	"github.com/arukim/galaxy/core"
	"github.com/arukim/galaxy/helpers"
)

// Server (game) contains game information
type Server struct {
	galaxy  *galaxy
	Players []*player

	//settings
	maxPlayers, maxTurns int
	mapWidth, mapHeight  int
	turnTimeout          time.Duration

	currTurn          int
	clientsCount      int
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

	s.galaxy = newGalaxy(10, 10)

	return &s
}

// AddClients populate game with players
func (s *Server) AddClients(clients []*core.Client) {
	s.clientsCount = len(clients)
	s.Players = make([]*player, len(clients))
	for i, c := range clients {
		s.Players[i] = newPlayer(s.turnCh, c, s.galaxy, i)
	}
}

// MakeTurn Proccess one game turn
func (s *Server) MakeTurn() {
	log.Printf("Turn %d has started\n", s.currTurn)

	// send turn info to players
	for _, p := range s.Players {
		go p.sendTurnInfo(s.currTurn)
	}

	// wait until all make turn or timeout
	var timeout = helpers.NewTimeout(s.turnTimeout)
	receivedTurns := 0

mainReceiveLoop:
	// or all players will send turn or timeout should work
	for receivedTurns < s.clientsCount {
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
		go p.send(com, data)
	}
}

// Start func to start the game
func (s *Server) Start() {
	s.galaxy.init()
	players := make([]string, s.clientsCount)
	for i, p := range s.Players {
		players[i] = p.Name
		p.connect()
	}

	gInfo := gameInfo{
		MaxTurns:  s.maxTurns,
		MapWidth:  s.mapWidth,
		MapHeight: s.mapHeight,
		Players:   players,
	}

	s.broadcast("gameInfo", &gInfo)

	// seed start spaceship
	for _, p := range s.Players {
		p.init()
	}

	for true {
		if s.currTurn < s.maxTurns {
			s.MakeTurn()
		} else {
			log.Println("game has ended")
			gameInfo := gameResult{Winner: "no one"}
			s.broadcast("gameResult", &gameInfo)

			for _, p := range s.Players {
				p.disconnect()
			}
			return
		}
	}
}
