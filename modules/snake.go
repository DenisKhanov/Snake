package modules

import (
	"slices"
)

type Snake struct {
	Parts []Point
	Size  int
}

func NewSnake() *Snake {
	return &Snake{}
}

func (s *Snake) Len() int {
	return len(s.Parts)
}

func (s *Snake) Add(point Point) {
	s.Parts = append([]Point{point}, s.Parts...)
}
func (s *Snake) IsSnake(point Point) bool {
	return slices.Contains(s.Parts, point)
}

func (s *Snake) CutIfSnake(point Point) bool {
	i := 0
	for ; i < len(s.Parts); i++ {
		if s.Parts[i] == point {
			s.Parts = s.Parts[0:i]
			return true
		}
	}
	return false
}

func (s *Snake) Head() Point {
	if len(s.Parts) == 0 {
		return Point{-1, -1}
	}
	return s.Parts[0]
}
func (s *Snake) Tail() Point {
	if len(s.Parts) == 0 {
		return Point{-1, -1}
	}
	return s.Parts[len(s.Parts)-1]
}
func (s *Snake) Reset() {
	x, y, length := 1, 1, 3 //snake position and length
	for i := length - 1; i >= 0; i-- {
		s.Parts = append(s.Parts, Point{float64(x + i), float64(y)})
		s.Size++
	}

}

func (s *Snake) Move(directional Dir) {
	lastPoint := s.Parts[0]
	s.Parts[0] = directional.Exec(s.Parts[0])
	for i := range s.Parts[1:] {
		s.Parts[i+1], lastPoint = lastPoint, s.Parts[i+1]
	}
}
