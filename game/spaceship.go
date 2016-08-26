package game

type spaceship struct {
	Pos       point
	Radar     int
	Level     int
	Owner     int
	Energy    int
	MaxEnergy int
}

// spaceshipInfo is sent to player
type spaceshipInfo struct {
	Level     int         `json:"level"`
	Owner     int         `json:"owner"`
	Radar     int         `json:"radar"`
	Energy    int         `json:"energy"`
	MaxEnergy int         `json:"maxEnergy"`
	View      []*zoneInfo `json:"view,omitempty"`
}

func newSpaceship(pos point, owner int) *spaceship {
	return &spaceship{
		Pos:       pos,
		Radar:     1,
		Level:     1,
		Owner:     owner,
		Energy:    1,
		MaxEnergy: 100,
	}
}

// create spaceshipInfo
func (s *spaceship) spaceshipInfo(isFull bool) *spaceshipInfo {
	i := &spaceshipInfo{
		Level:     s.Level,
		Owner:     s.Owner,
		MaxEnergy: s.MaxEnergy,
	}

	if isFull {
		i.Radar = s.Radar
		i.Energy = s.Energy
	}

	return i
}

func (s *spaceship) collect(z *zone) {
	canAbsorb := s.MaxEnergy - s.Energy
	available := int(z.currEnergy)

	if available >= canAbsorb {
		s.Energy = s.MaxEnergy
		z.currEnergy -= float64(canAbsorb)
	} else {
		s.Energy += available
		z.currEnergy -= float64(available)
	}
}

func (s *spaceship) getScore() int {
	return s.Energy
}
