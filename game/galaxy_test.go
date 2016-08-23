package game

import "testing"

func Test_ZoneGetOutOfBounds(t *testing.T) {
	g := newGalaxy(1, 1)
	z := g.get(1, 2)
	if z != nil {
		t.Error("should return nil pointer cell")
	}

	z = g.get(0, 0)
	if z == nil {
		t.Error(("should return non nill cell"))
	}
}
