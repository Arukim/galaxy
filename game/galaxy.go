package game

import (
	"log"
	"math/rand"
)

// Map is storing game map
type galaxy struct {
	zones  []*zone
	width  int
	height int
}

// Get to get Row/Col of map
func (g *galaxy) get(r, c int) *zone {
	if r*g.width+c >= g.width*g.height {
		return nil
	}
	return g.zones[r*g.width+c]
}

// NewGalaxy creates new instance of Map
func newGalaxy(width, height int) *galaxy {
	g := new(galaxy)

	g.width = width
	g.height = height
	g.zones = make([]*zone, g.width*g.height)

	return g
}

// Init to initialize Map
func (g *galaxy) init() {
	for i := range g.zones {
		g.zones[i] = &zone{}
		g.zones[i].Init()
	}
	log.Println("map generated")
}

// Spawn the energy
func (g *galaxy) spawn() {
	for _, z := range g.zones {
		z.Spawn()
	}
}

// spawn ship at start location
func (g *galaxy) spawnSpaceship(owner int) *spaceship {
	pos := point{
		X: rand.Intn(g.width),
		Y: rand.Intn(g.height),
	}
	return newSpaceship(pos, owner)
}

// Get info about one spaceship (surround view and etc)
func (g *galaxy) getSpaceshipInfo(s *spaceship) *spaceshipInfo {
	res := s.spaceshipInfo(true)
	viewZone := getAllPointsInCircle(s.Radar)

	for _, v := range viewZone {
		p := s.Pos.Add(v)
		zone := g.get(p.X, p.Y)

		if zone != nil {
			zInfo := &zoneInfo{
				Pos:    v,
				Energy: int(zone.currEnergy),
			}
			if v.X == 0 && v.Y == 0 {
				zInfo.Spaceships = []*spaceshipInfo{
					s.spaceshipInfo(false),
				}
			}
			res.View = append(res.View, zInfo)
		}
	}

	return res
}
