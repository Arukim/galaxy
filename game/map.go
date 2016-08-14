package game

import "log"

// Map is storing game map
type Map struct {
	Zones  []Zone
	Width  int
	Height int
}

// Get to get Row/Col of map
func (m *Map) Get(r, c int) *Zone {
	if r*m.Width+c >= m.Width*m.Height {
		return nil
	}
	return &m.Zones[r*m.Width+c]
}

// NewMap creates new instance of Map
func NewMap(width, height int) *Map {
	m := new(Map)

	m.Width = width
	m.Height = height
	m.Zones = make([]Zone, m.Width*m.Height)

	return m
}

// Init to initialize Map
func (m *Map) Init() {
	for i := range m.Zones {
		zone := &m.Zones[i]

		zone.Init()

		log.Printf("zone %v", *zone)
	}
}

// Spawn the energy
func (m *Map) Spawn() {
	for i := range m.Zones {
		z := &m.Zones[i]

		z.Spawn()
	}
}
