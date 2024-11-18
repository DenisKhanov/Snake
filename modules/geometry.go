package modules

type Point struct {
	X, Y float64
}

const (
	up = iota
	right
	down
	left
)

type Dir int

func (d Dir) Exec(point Point) Point {
	switch d {
	case up:
		return Point{point.X, point.Y + 1}
	case down:
		return Point{point.X, point.Y - 1}
	case left:
		return Point{point.X - 1, point.Y}
	case right:
		return Point{point.X + 1, point.Y}
	default:
		return point
	}
}

func (d Dir) CheckParallel(d2 Dir) bool {
	switch d {
	case up:
		return d2 == down
	case right:
		return d2 == left
	case down:
		return d2 == up
	case left:
		return d2 == right
	default:
		return false
	}
}
