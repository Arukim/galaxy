package game

import (
	"sync"

	"github.com/arukim/galaxy/core"
)

// Player holds game player data
type player struct {
	id         int
	Name       string
	currTurn   playerTurn
	client     *core.Client
	turnCh     chan *playerTurn
	isActive   bool
	activeLock *sync.Mutex
	galaxy     *galaxy

	Spaceships []*spaceship
}

type turnInfo struct {
	Turn       int              `json:"turn"`
	Spaceships []*spaceshipInfo `json:"spaceships"`
}

// NewPlayer creates new player
func newPlayer(turnCh chan *playerTurn, c *core.Client, g *galaxy, id int) *player {
	p := player{
		turnCh:     turnCh,
		client:     c,
		id:         id,
		Name:       c.Name,
		Spaceships: []*spaceship{},
		galaxy:     g,
	}

	return &p
}

func (p *player) sendTurnInfo(turn int) {
	ti := turnInfo{
		Turn:       turn,
		Spaceships: []*spaceshipInfo{},
	}

	for _, s := range p.Spaceships {
		ti.Spaceships = append(ti.Spaceships, p.galaxy.getSpaceshipInfo(s))
	}

	p.client.Router.Send("turnInfo", &ti)
}

func (p *player) handlers() []*core.CommandHandler {
	return []*core.CommandHandler{}
}

func (p *player) init() {
	p.Spaceships = append(p.Spaceships, p.galaxy.spawnSpaceship(p.id))
}

// Send something to client
func (p *player) send(cmd string, data interface{}) {
	p.client.Router.Send(cmd, data)
}

// Connect player to current game
func (p *player) connect() {
	p.client.Router.SetHandlers(p.handlers())
}

// Disconnect player to game server
func (p *player) disconnect() {
	p.client.Router.SetHandlers(p.client.DefaultHandlers)
}

// StartTurn should be called at the start of turn
func (p *player) startTurn() {
	p.activeLock.Lock()
	defer p.activeLock.Unlock()

	p.isActive = true
}

// EndTurn should be callled at the end of turn
func (p *player) endTurn() {
	p.activeLock.Lock()
	defer p.activeLock.Unlock()

	p.isActive = false
}

// Handle handles data from socket
func (p *player) handle(data interface{}) {
	switch data.(type) {
	case playerTurn:
		p.activeLock.Lock()
		defer p.activeLock.Unlock()
		if p.isActive {
			currTurn := data.(playerTurn)
			currTurn.PlayerID = p.id
			p.turnCh <- &currTurn
			p.isActive = false
		}
	}
}
