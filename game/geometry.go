// Package game contains the core functionality for the Snake game, including game logic, rendering, geometry handling, and snake behavior.
package game

// Point represents a 2D coordinate with X and Y values.
// This struct is commonly used to represent positions
// of game elements (e.g., snake, food) in a 2D space.
type Point struct {
	X, Y float64
}

// IsCorner checks whether a given Point is located at one of the four corners.
func (p Point) IsCorner() bool {
	return p.X == 0 && p.Y == 0 || p.X == 0 && p.Y == cellsCount-1 ||
		p.X == cellsCount-1 && p.Y == 0 || p.X == cellsCount-1 && p.Y == cellsCount-1
}

// IsEdge checks whether a given Point is located at one of the four edge.
func (p Point) IsEdge() bool {
	return p.X == 0 || p.Y == 0 || p.X == cellsCount-1 || p.Y == cellsCount-1
}

// Direction constants for snake movement.
const (
	up = iota
	right
	down
	left
)

type Dir int

// Exec moves the point based on the given Direction (up, down, left, or right).
// It modifies the X or Y coordinate of the point depending on the Direction.
// - `up`: Increases the Y coordinate by 1 (moves the point upwards).
// - `down`: Decreases the Y coordinate by 1 (moves the point downwards).
// - `left`: Decreases the X coordinate by 1 (moves the point leftward).
// - `right`: Increases the X coordinate by 1 (moves the point rightward).
// If an invalid Direction is provided, the point remains unchanged.
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

// FromKey returns the corresponding Direction based on the key code passed as an argument.
// The key codes correspond to the arrow keys on the keyboard:
// - 80: Left arrow key → Returns "left" Direction.
// - 82: Up arrow key → Returns "down" Direction (Note: this seems reversed in your code, should probably be "up").
// - 79: Right arrow key → Returns "right" Direction.
// - 81: Down arrow key → Returns "up" Direction (Note: this also seems reversed, should probably be "down").
// If the key code does not match any of the above, it returns "right" as the default Direction.
func (d Dir) FromKey(ceyKode int) Dir {
	switch ceyKode {
	case 80: //left
		return left
	case 82: //up
		return down
	case 79: //right
		return right
	case 81: //down
		return up
	default:
		return right
	}
}

// CheckParallel checks if the new Direction is opposite (parallel) to the current Direction.
// This method helps to prevent the snake from reversing Direction (which would result in it colliding with itself).
//
// The method compares the current Direction (`d`) with the new Direction (`newDir`) and returns:
// - `true` if the new Direction is directly opposite (i.e., the snake would collide with itself if it moved that way).
// - `false` otherwise.
func (d Dir) CheckParallel(newDir Dir) bool {
	switch d {
	case up:
		return newDir == down
	case right:
		return newDir == left
	case down:
		return newDir == up
	case left:
		return newDir == right
	default:
		return false
	}
}
