package game

type starbase struct {
	pos   point
	owner int
}

type starbaseInfo struct {
	Pos   point `json:"pos"`
	Owner int   `json:"owner"`
}

func newStarbase(pos point, owner int) *starbase {
	s := starbase{
		pos:   pos,
		owner: owner,
	}
	return &s
}

func (s *starbase) toInfo() *starbaseInfo {
	i := starbaseInfo{
		Pos:   s.pos,
		Owner: s.owner,
	}
	return &i
}
