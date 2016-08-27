package game

type spaceship struct {
	pos       point
	radar     int
	level     int
	owner     int
	energy    int
	maxEnergy int
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
		pos:       pos,
		radar:     1,
		level:     1,
		owner:     owner,
		energy:    1,
		maxEnergy: 100,
	}
}

// create spaceshipInfo
func (s *spaceship) toInfo(isFull bool) *spaceshipInfo {
	i := &spaceshipInfo{
		Level:     s.level,
		Owner:     s.owner,
		MaxEnergy: s.maxEnergy,
	}

	if isFull {
		i.Radar = s.radar
		i.Energy = s.energy
	}

	return i
}

func (s *spaceship) collect(z *zone) {
	canAbsorb := s.maxEnergy - s.energy
	available := int(z.currEnergy)

	if available >= canAbsorb {
		s.energy = s.maxEnergy
		z.currEnergy -= float64(canAbsorb)
	} else {
		s.energy += available
		z.currEnergy -= float64(available)
	}
}

func (s *spaceship) getScore() int {
	return s.energy
}
