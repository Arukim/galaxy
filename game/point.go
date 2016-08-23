package game

type point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Sum two points into new point
func (a *point) Add(b *point) *point {
	return &point{X: a.X + b.X, Y: a.Y + b.Y}
}

// Generate list of all point in circle
func getAllPointsInCircle(r int) []*point {
	res := []*point{}
	for w := -r; w <= r; w++ {
		for h := -r; h <= r; h++ {
			d := w*w + h*h
			if d <= r {
				res = append(res, &point{w, h})
			}
		}
	}
	return res
}
