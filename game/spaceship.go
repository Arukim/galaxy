package game

type spaceship struct {
	Pos   point
	Radar int
	Level int
	Owner int
}

// spaceshipInfo is sent to player
type spaceshipInfo struct {
	Level int         `json:"level"`
	Owner int         `json:"owner"`
	Radar int         `json:"radar"`
	View  []*zoneInfo `json:"view,omitempty"`
}

func newSpaceship(pos point, owner int) *spaceship {
	return &spaceship{
		Pos:   pos,
		Radar: 1,
		Level: 1,
		Owner: owner,
	}
}

// create spaceshipInfo
func (s *spaceship) spaceshipInfo(isFull bool) *spaceshipInfo {
	i := &spaceshipInfo{
		Level: s.Level,
		Owner: s.Owner,
	}

	if isFull {
		i.Radar = s.Radar
	}

	return i
}
