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
	currTurn   playerTurn
	conn       IPlayerConnection
	turnCh     chan *playerTurn
	isActive   bool
	activeLock *sync.Mutex
}

// NewPlayer creates new player
func NewPlayer(turnCh chan *playerTurn, conn IPlayerConnection, id int) *Player {
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
