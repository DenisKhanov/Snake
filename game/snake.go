// Package game contains the core functionality for the Snake game, including game logic, rendering, geometry handling, and snake behavior.
package game

import (
	"slices"
)

// Snake represents the game snake.
// Fields:
// - Direction: snake direction for to go next step
// - Parts: an array of points that define the positions of the snake's segments on the game field.
// - Size: the current size of the snake (number of segments).
type Snake struct {
	Direction Dir
	Parts     []Point
	Size      int
}

// NewSnake creates and returns a new instance of the Snake struct.
//
// This function initializes the Snake object without any predefined state.

func NewSnake() *Snake {
	return &Snake{}
}

// Len returns the current length of the snake.
//
// This method calculates the length by counting the number of parts
// in the snake's body (s.Parts).
//
// Returns:
//
//	int - The total number of parts in the snake.
func (s *Snake) Len() int {
	return len(s.Parts)
}

// Add inserts a new point at the head of the snake.
//
// This method extends the snake by adding a new part at the beginning
// of the `s.Parts` slice, representing the snake's head.
//
// Parameters:
//   - point (Point): The coordinates of the new part to be added.
func (s *Snake) Add(point Point) {
	s.Parts = append([]Point{point}, s.Parts...)
}

// IsSnake checks if a given point is part of the snake's body.
//
// This method determines whether the specified point is within the `s.Parts`
// slice, representing the snake's current body parts.
//
// Parameters:
//   - point (Point): The point to check for presence in the snake's body.
//
// Returns:
//   - bool: `true` if the point is part of the snake, otherwise `false`.
func (s *Snake) IsSnake(point Point) bool {
	return slices.Contains(s.Parts, point)
}

// CutIfSnake checks if a given point is part of the snake's body
// and, if so, cuts the snake at that point.
//
// This method iterates through the snake's body (`s.Parts`) to find the specified point.
// If the point is found, the snake's body is truncated up to that point,
// effectively removing all parts after it.
//
// Parameters:
//   - point (Point): The point to check and cut the snake at.
//
// Returns:
//   - bool: `true` if the point is part of the snake and the body was cut, otherwise `false`.
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

// Head retrieves the current position of the snake's head.
//
// If the snake has no parts (i.e., it has not been initialized or is empty),
// this method returns a default invalid position (-1, -1).
//
// Returns:
//   - Point: The coordinates of the snake's head or (-1, -1) if the snake is empty.
func (s *Snake) Head() Point {
	if len(s.Parts) == 0 {
		return Point{-1, -1}
	}
	return s.Parts[0]
}

// Tail retrieves the current position of the snake's tail.
//
// If the snake has no parts (i.e., it has not been initialized or is empty),
// this method returns a default invalid position (-1, -1).
//
// Returns:
//   - Point: The coordinates of the snake's tail or (-1, -1) if the snake is empty.
func (s *Snake) Tail() Point {
	if len(s.Parts) == 0 {
		return Point{-1, -1}
	}
	return s.Parts[len(s.Parts)-1]
}

// Reset reinitialized the snake to its starting state.
//
// This method resets the snake's position and direction. It clears the existing snake parts
// and then sets the snake's head and body at a defined starting position. By default,
// the snake starts at position (1, 1) with a length of 3 and moves to the right.
//
// The snake's size is updated as the parts are added to the snake's body.
//
// Side Effects:
//   - Resets the snake's parts to a new slice of length 0.
//   - Sets the snake's direction to "right".
//   - Initializes the snake's body at a starting position with a default length of 3.
func (s *Snake) Reset() {
	s.Parts = []Point{}
	s.Direction = right
	x, y, length := 1, 1, 3 //snake position and length
	for i := length - 1; i >= 0; i-- {
		s.Parts = append(s.Parts, Point{float64(x + i), float64(y)})
		s.Size++
	}
}

// Move updates the snake's position based on the given direction.
//
// This method moves the snake by updating its head to the new position according to the
// provided direction. It also updates the positions of the body parts, shifting each part
// to the position of the previous one, effectively simulating movement.
//
// The movement is handled by modifying the snake's head using the Exec method of the
// provided direction, and then shifting the remaining body parts sequentially.
//
// Parameters:
//   - directional (Dir): The direction in which the snake should move. This can be one of
//     the constants up, down, left, or right.
func (s *Snake) Move(directional Dir) {
	lastPoint := s.Parts[0]
	s.Parts[0] = directional.Exec(s.Parts[0])
	for i := range s.Parts[1:] {
		s.Parts[i+1], lastPoint = lastPoint, s.Parts[i+1]
	}
}
