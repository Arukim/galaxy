package game

import (
	"encoding/json"
	"sync"

	"github.com/arukim/galaxy/core"
)

// Player holds game player data
type player struct {
	id         int
	Name       string
	currTurn   int
	client     *core.Client
	turnCh     chan *playerTurn
	isActive   bool
	activeLock *sync.Mutex
	galaxy     *galaxy

	Spaceships []*spaceship
	starbases  []*starbase
}

type turnInfo struct {
	Turn       int              `json:"turn"`
	Spaceships []*spaceshipInfo `json:"spaceships"`
	Starbases  []*starbaseInfo  `json:"starbases"`
}

// Client turn
type playerTurn struct {
	Turn       int              `json:"turn"`
	PlayerID   int              `json:"-"`
	Spaceships []*spaceshipTurn `json:"spaceships"`
}

type spaceshipTurn struct {
	ID     int    `json:"id"`
	Action string `json:"action"`
	Pos    point  `json:"point"`
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
		activeLock: &sync.Mutex{},
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

	for _, s := range p.starbases {
		ti.Starbases = append(ti.Starbases, s.toInfo())
	}

	p.activeLock.Lock()
	defer p.activeLock.Unlock()

	p.isActive = true
	p.currTurn = turn
	p.client.Router.Send("turnInfo", &ti)
}

func (p *player) collectEnergy() {
	for _, s := range p.Spaceships {
		s.collect(p.galaxy.getPos(s.pos))
	}
}

func (p *player) getScore() int {
	var total = 0
	for _, s := range p.Spaceships {
		total += s.getScore()
	}
	return total
}

func (p *player) handlers() []*core.CommandHandler {
	return []*core.CommandHandler{
		{
			Name:   "makeTurn",
			Handle: p.onTurn,
		},
	}
}

func (p *player) init() {
	pos := p.galaxy.getStartLocation(p.id)
	p.starbases = append(p.starbases, newStarbase(pos, p.id))
	p.Spaceships = append(p.Spaceships, newSpaceship(pos, p.id))
}

func (p *player) onTurn(d *json.RawMessage, c *core.Client) *core.Result {
	var resp *core.Result
	var turn playerTurn
	json.Unmarshal(*d, &turn)

	p.activeLock.Lock()
	defer p.activeLock.Unlock()
	if p.isActive && p.currTurn == turn.Turn {
		turn.PlayerID = p.id
		p.turnCh <- &turn
		p.isActive = false
		resp = core.NewSuccessResult("turn accepted")
	} else {
		resp = core.NewErrorResult("turn not accepted")
	}

	return resp
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

// EndTurn should be callled at the end of turn
func (p *player) endTurn() {
	p.activeLock.Lock()
	defer p.activeLock.Unlock()

	p.isActive = false
}
