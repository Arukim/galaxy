package game

// IPlayer describes abstract game player
type IPlayerConnection interface {
	GetName() string
	Send(data interface{})
	Receive(data interface{})
}

type Player struct {
	name string
	conn IPlayerConnection
}

func (p *Player) Handle(data interface{}) {
	switch data.(type) {
	case PlayerInfo:
		p.name = data.(PlayerInfo).name
	}
}
