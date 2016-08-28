package game

import (
	"fmt"
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
	turnStartCh       chan bool
	turnEndCh         chan bool
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
	s.turnStartCh = make(chan bool, s.maxPlayers)
	s.turnEndCh = make(chan bool, s.maxPlayers)
	s.turnProccessCount = make(chan bool, s.maxPlayers)

	s.galaxy = newGalaxy(10, 10)

	return &s
}

// AddClients populate game with players
func (s *Server) AddClients(clients []*core.Client) {
	s.clientsCount = len(clients)
	s.Players = make([]*player, len(clients))
	for i, c := range clients {
		s.Players[i] = newPlayer(s.turnStartCh, s.turnEndCh, c, s.galaxy, i)
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
		case <-s.turnStartCh:
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
		case <-s.turnStartCh:
			receivedTurns++
		default:
			break cleanUpLoop
		}
	}

	// now wait until all player turns are made
	for receivedTurns > 0 {
		select {
		case <-s.turnEndCh:
			receivedTurns--
		}
	}

	// all player turns are applied, calculate turn end
	// lets eat energy
	for _, p := range s.Players {
		p.collectEnergy()
	}
	// now is galaxy turn
	s.galaxy.spawn()

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

			winner := s.findWinner()

			gameInfo := gameResult{Result: fmt.Sprintf("player %v wins", winner.Name)}

			s.broadcast("gameResult", &gameInfo)

			for _, p := range s.Players {
				p.disconnect()
			}

			core.GameStats.AddPlayedGame()
			return
		}
	}
}

func (s *Server) findWinner() *player {
	var winner *player
	maxScore := 0

	for _, p := range s.Players {
		pScore := p.getScore()
		if winner == nil || pScore > maxScore {
			winner = p
			maxScore = pScore
		}
	}

	return winner
}
