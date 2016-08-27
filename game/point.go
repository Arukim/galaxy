package game

import "math"

type point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Sum two points into new point
func (a *point) Add(b point) point {
	return point{X: a.X + b.X, Y: a.Y + b.Y}
}

func (a *point) GetDistance(b point) int {
	return int(math.Sqrt(float64((a.X-b.X)*(a.X-b.X) + (a.Y-b.Y)*(a.Y-b.Y))))
}

func wrapBetween(value, a, b int) int {
	if value < a {
		value = b - (a - value)
	} else if value >= b {
		value = value - b + a
	}
	return value
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
