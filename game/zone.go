package game

import "math/rand"

// Zone describe one map point
type zone struct {
	maxEnergy            int
	currEnergy           float64
	energyGenerationRate float64
}

// Spawn calculates new turn energy level
func (z *zone) Spawn() {
	z.currEnergy += z.energyGenerationRate
	if int(z.currEnergy) > z.maxEnergy {
		z.currEnergy = float64(z.maxEnergy)
	}
}

// Init Zone
func (z *zone) Init() {
	z.maxEnergy = rand.Intn(100)
	z.currEnergy = rand.Float64() * float64(z.maxEnergy)
	z.energyGenerationRate = rand.Float64() * float64(z.maxEnergy) / 5
}

// info about 1 map zone, is sent to player
type zoneInfo struct {
	Energy     int              `json:"energy"`
	Pos        *point           `json:"pos"`
	Spaceships []*spaceshipInfo `json:"spaceships,omitempty"`
}
