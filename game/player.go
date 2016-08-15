package game

import "sync"

// IPlayerConnection describes abstract game player
type IPlayerConnection interface {
	GetName() string
	Send(data interface{})
	Receive(data interface{})
}

// Player holds game player data
type Player struct {
	id         int
	name       string
	currTurn   PlayerTurn
	conn       IPlayerConnection
	turnCh     chan *PlayerTurn
	isActive   bool
	activeLock *sync.Mutex
}

// NewPlayer creates new player
func NewPlayer(turnCh chan *PlayerTurn, conn IPlayerConnection, id int) *Player {
	p := Player{
		turnCh: turnCh,
		conn:   conn,
		id:     id,
	}

	return &p
}

// StartTurn should be called at the start of turn
func (p *Player) StartTurn() {
	p.activeLock.Lock()
	defer p.activeLock.Unlock()

	p.isActive = true
}

// EndTurn should be callled at the end of turn
func (p *Player) EndTurn() {
	p.activeLock.Lock()
	defer p.activeLock.Unlock()

	p.isActive = false
}

// Handle handles data from socket
func (p *Player) Handle(data interface{}) {
	switch data.(type) {
	case PlayerInfo:
		p.name = data.(PlayerInfo).name
	case PlayerTurn:
		p.activeLock.Lock()
		defer p.activeLock.Unlock()
		if p.isActive {
			currTurn := data.(PlayerTurn)
			currTurn.PlayerId = p.id
			p.turnCh <- &currTurn
			p.isActive = false
		}
	}
}
