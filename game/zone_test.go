package game

import "testing"

func Test_ZoneSpawn(t *testing.T) {
	x := Zone{maxEnergy: 10, currEnergy: 0.0, energyGenerationRate: 11.0}
	x.Spawn()
	if x.currEnergy != 10 {
		t.Error("curr energy should rise and not overflow")
	}
}
