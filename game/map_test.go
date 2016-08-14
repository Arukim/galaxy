package game

import "testing"

func Test_ZoneGetOutOfBounds(t *testing.T) {
	m := NewMap(1, 1)
	z := m.Get(1, 2)
	if z != nil {
		t.Error("should return nil pointer cell")
	}

	z = m.Get(0, 0)
	if z == nil {
		t.Error(("should return non nill cell"))
	}
}
