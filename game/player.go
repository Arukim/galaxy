package game

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/arukim/galaxy/core"
)

// Player holds game player data
type player struct {
	id         int
	Name       string
	currTurn   int
	client     *core.Client
	isActive   bool
	activeLock *sync.Mutex
	galaxy     *galaxy

	turnStartCh chan bool
	turnEndCh   chan bool

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
	Turn       int             `json:"turn"`
	PlayerID   int             `json:"-"`
	Spaceships []spaceshipTurn `json:"spaceships"`
}

type spaceshipTurn struct {
	ID     int    `json:"id"`
	Action string `json:"action"`
	Pos    point  `json:"pos"`
}

// NewPlayer creates new player
func newPlayer(turnCh chan bool, turnEnd chan bool, c *core.Client, g *galaxy, id int) *player {
	p := player{
		turnStartCh: turnCh,
		turnEndCh:   turnEnd,
		client:      c,
		id:          id,
		Name:        c.Name,
		Spaceships:  []*spaceship{},
		galaxy:      g,
		activeLock:  &sync.Mutex{},
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
	p.Spaceships = append(p.Spaceships, newSpaceship(pos, p.id, 0))
}

func (p *player) onTurn(d *json.RawMessage, c *core.Client) *core.Result {
	var resp *core.Result
	var turn playerTurn
	json.Unmarshal(*d, &turn)

	p.activeLock.Lock()
	defer p.activeLock.Unlock()
	if p.isActive && p.currTurn == turn.Turn {
		turn.PlayerID = p.id
		p.turnStartCh <- true
		p.isActive = false
		go p.processTurn(&turn)
		resp = core.NewSuccessResult("turn accepted")
	} else {
		resp = core.NewErrorResult("turn not accepted")
	}

	return resp
}

func (p *player) processTurn(t *playerTurn) {

	log.Printf("%v", t.Spaceships[0])
shipLoop:
	for _, shipTurn := range t.Spaceships {
		// TODO check income data
		ship := p.Spaceships[shipTurn.ID]
		if ship == nil {
			continue shipLoop
		}

		if shipTurn.Pos.GetDistance(point{}) <= ship.level {
			ship.pos = p.galaxy.wrapPoint(ship.pos.Add(shipTurn.Pos))
			log.Printf("new pos of %v is %v", ship.id, ship.pos)
		}
		//if shipTurn.Pos
		shipTurn.ID++
	}
	p.turnEndCh <- true
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
