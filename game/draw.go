// Package game contains the core functionality for the Snake game, including game logic, rendering, geometry handling, and snake behavior.
package game

import (
	"fmt"
	"log"
	"math"
)

// drawWorld renders the background of the game area.
//
// This method fills a rectangular region representing the game world with a specific color.
func (g *Game) drawWorld() {
	g.cv.BeginPath()
	g.cv.SetFillStyle("#78909C")
	g.cv.FillRect(g.gameAreaSP.X, g.gameAreaSP.Y, g.gameAreaEP.X-15, g.gameAreaEP.Y-15)
	g.cv.Stroke()
}

// drawGridGameArea renders a grid within the game area.
//
// This method draws evenly spaced vertical and horizontal lines to create a grid.
func (g *Game) drawGridGameArea() {
	g.cv.BeginPath()
	g.cv.SetStrokeStyle("#5D4037")
	g.cv.SetLineWidth(0.5)
	for i := 0; i < 20+1; i++ {
		g.cv.MoveTo(g.gameAreaSP.X+float64(i)*g.cellH, g.gameAreaSP.Y)
		g.cv.LineTo(g.gameAreaSP.X+float64(i)*g.cellH, g.gameAreaEP.Y)
		g.cv.MoveTo(g.gameAreaSP.X, g.gameAreaSP.Y+float64(i)*g.cellW)
		g.cv.LineTo(g.gameAreaEP.X, g.gameAreaSP.Y+float64(i)*g.cellW)
	}
	g.cv.Stroke()
}

// drawSnakeHead renders the snake's head on the game canvas at the specified position.
//
// The snake's head is drawn as an ellipse with eyes, nostrils, and a tongue to create a more detailed visual representation.
//
// Parameters:
// - x (float64): The x-coordinate of the snake's head position.
// - y (float64): The y-coordinate of the snake's head position.
// - side (float64): The size of the square cell that the snake's head fits into, used to calculate proportions for the head and its features.
func (g *Game) drawSnakeHead(x, y, side float64) {
	//Draw snake head's main ellipse
	centerX := x + side/2
	centerY := y + side/2
	radiusX := side / 2
	radiusY := side * 0.6 / 2

	g.cv.SetFillStyle("#039BE5")
	g.cv.BeginPath()
	g.cv.Ellipse(centerX, centerY, radiusX, radiusY, 0, 0, 2*math.Pi, false)
	g.cv.Fill()

	// Draw eyes
	eyeRadius := side * 0.1
	eyeOffsetX := side * 0.2
	eyeOffsetY := side * 0.2

	g.cv.SetFillStyle("#ffffff")
	g.cv.BeginPath()
	g.cv.Arc(centerX-eyeOffsetX, centerY-eyeOffsetY, eyeRadius, 0, 2*math.Pi, false) // Левый глаз
	g.cv.Arc(centerX+eyeOffsetX, centerY-eyeOffsetY, eyeRadius, 0, 2*math.Pi, false) // Правый глаз
	g.cv.Fill()

	g.cv.SetFillStyle("#000000")
	g.cv.BeginPath()
	g.cv.Arc(centerX-eyeOffsetX, centerY-eyeOffsetY, eyeRadius*0.4, 0, 2*math.Pi, false) // Левый зрачок
	g.cv.Arc(centerX+eyeOffsetX, centerY-eyeOffsetY, eyeRadius*0.4, 0, 2*math.Pi, false) // Правый зрачок
	g.cv.Fill()

	// Draw nostrils
	nostrilRadius := side * 0.03
	nostrilOffsetX := side * 0.1
	nostrilOffsetY := side - 38

	g.cv.SetFillStyle("#000000")
	g.cv.BeginPath()
	g.cv.Arc(centerX-nostrilOffsetX, centerY-nostrilOffsetY, nostrilRadius, 0, 2*math.Pi, false) // Left nostril
	g.cv.Arc(centerX+nostrilOffsetX, centerY-nostrilOffsetY, nostrilRadius, 0, 2*math.Pi, false) // Right nostril
	g.cv.Fill()

	// Draw tongue
	tongueWidth := side * 0.05
	tongueLength := side * 0.5

	g.cv.SetFillStyle("#ff0000")
	g.cv.BeginPath()
	g.cv.MoveTo(centerX, centerY+radiusY*0.8)
	g.cv.LineTo(centerX-tongueWidth, centerY+radiusY+tongueLength/2)
	g.cv.LineTo(centerX, centerY+radiusY+tongueLength)
	g.cv.LineTo(centerX+tongueWidth, centerY+radiusY+tongueLength/2)
	g.cv.ClosePath()
	g.cv.Fill()
}

// drawSnake renders the snake on the game canvas.
//
// The snake is drawn part by part, with the first part being the head and the rest of the body alternating between two different colors for visual distinction.
func (g *Game) drawSnake() {
	g.cv.BeginPath()
	for i, point := range g.snake.Parts {
		switch {
		case i == 0: //draw head
			g.drawSnakeHead(g.gameAreaSP.X+point.X*g.cellW+1, g.gameAreaSP.Y+point.Y*g.cellH+1, g.side)
		case i%2 == 0:
			g.cv.SetFillStyle("#00BCD4")
			g.cv.FillRect(
				g.gameAreaSP.X+point.X*g.cellW+1,
				g.gameAreaSP.Y+point.Y*g.cellH+1,
				g.cellW-1*2,
				g.cellH-1*2,
			)
		default:
			g.cv.SetFillStyle("#4DD0E1")
			g.cv.FillRect(
				g.gameAreaSP.X+point.X*g.cellW+1,
				g.gameAreaSP.Y+point.Y*g.cellH+1,
				g.cellW-1*2,
				g.cellH-1*2,
			)
		}
	}
	g.cv.Stroke()
}

// drawApple renders an apple on the game canvas at the specified position.
//
// The apple consists of three parts: a circular body, a leaf, and a stalk.
// The apple is drawn inscribed within a square cell, with its size determined by `sizeCell`.
//
// Parameters:
// - x (float64): The x-coordinate of the apple's position.
// - y (float64): The y-coordinate of the apple's position.
// - sizeCell (float64): The size of the cell the apple fits into (used to calculate radius and proportions).
func (g *Game) drawApple(x, y, sizeCell float64) {
	// Draw main an apple circle inscribed in a square
	radius := sizeCell / 2
	centerX := x + radius
	centerY := y + radius

	g.cv.SetFillStyle("#7CB342")
	g.cv.BeginPath()
	g.cv.Arc(centerX, centerY, radius, 0, 2*math.Pi, false)
	g.cv.Fill()

	// Draw an apple leaf
	g.cv.SetFillStyle("#1B5E20")
	g.cv.BeginPath()
	g.cv.MoveTo(centerX-5, centerY-radius*0.1)
	g.cv.BezierCurveTo(
		centerX-radius*0.8, centerY-radius*1.2,
		centerX+radius*0.6, centerY-radius*1.2,
		centerX+radius*0.2, centerY-radius*0.8,
	)
	g.cv.ClosePath()
	g.cv.Fill()

	// Draw an apple stalk
	stemWidth := sizeCell * 0.1
	stemHeight := sizeCell * 0.2
	g.cv.SetFillStyle("#8B4513")
	g.cv.FillRect(centerX-stemWidth/2, centerY-radius, stemWidth, -stemHeight)
	g.cv.Stroke()
}

// drawGameInfo displays the current game statistics on the screen.
//
// This method shows the current score, the number of food items eaten, the current speed of the snake, and the FPS.
func (g *Game) drawGameInfo() {
	g.cv.SetFillStyle("#4CAF50")
	g.cv.BeginPath()
	g.cv.SetFont(g.fonts.main, 25)

	//draw score
	text := fmt.Sprintf("Your score: %d", g.score)
	g.cv.FillText(text, g.param.gameW+50, 50)

	// food
	text = fmt.Sprintf("You ate food: %d", g.ateFood)
	g.cv.FillText(text, g.param.gameW+50, 85)

	// speed
	text = fmt.Sprintf("Your speed: %d", startSpeed-g.param.speed+5)
	g.cv.FillText(text, g.param.gameW+50, 120)

	g.cv.Stroke()
}

// drawInstructions renders the game instructions on the canvas.
//
// This method displays the basic controls for the game, including how to move the snake, how to grow the snake, and how to shorten it if it eats its own tail.
func (g *Game) drawInstructions() {
	g.cv.BeginPath()
	g.cv.SetFillStyle("#FFEE58")
	g.cv.SetFont(g.fonts.main, 20)
	text := fmt.Sprint("Game Instructions:")
	g.cv.FillText(text, g.param.gameW+50, 215)
	g.cv.Stroke()

	g.cv.BeginPath()
	g.cv.SetFillStyle("#CFD8DC")
	g.cv.SetFont(g.fonts.middle, 15)
	text = fmt.Sprint("Use keys ← ↑ → ↓ to move snake")
	g.cv.FillText(text, g.param.gameW+30, 245)

	text = fmt.Sprint("Raise     to grow +++")
	g.cv.FillText(text, g.param.gameW+30, 275)

	text = fmt.Sprint("If you eat your tail, ")
	g.cv.FillText(text, g.param.gameW+30, 305)
	text = fmt.Sprint(" the snake will shorten---")
	g.cv.FillText(text, g.param.gameW+70, 325)
	g.cv.Stroke()

	g.drawApple(g.param.gameW+90, 265, g.side*0.6)
}

// drawAboutCreator displays information about the game's creator on the screen.
//
// This method renders a brief description of the game and credits the creator.
// The text is displayed at the specified coordinates.
func (g *Game) drawAboutCreator(x, y float64) {
	g.cv.BeginPath()
	g.cv.SetFillStyle("#00897B")
	g.cv.SetFont(g.fonts.small, 15)
	text := fmt.Sprint("This game  was created in the Golang")
	g.cv.FillText(text, x, y)
	text = fmt.Sprint("by Denis Khanov")
	g.cv.FillText(text, x, y+20)
	g.cv.Stroke()
}

// drawFPS displays information about FPS
func (g *Game) drawFPS() {
	g.cv.BeginPath()
	g.cv.SetFillStyle("#FFEE58")
	g.cv.SetFont(g.fonts.small, 15)
	text := fmt.Sprintf("FPS: %.1f", g.wnd.FPS())
	g.cv.FillText(text, 5, 14)
	g.cv.Stroke()
}

// drawContacts displays contact information and clickable links for the game's repository and the creator's Telegram profile.
//
// This method shows the game's repository URL and the creator's Telegram handle as clickable links.
func (g *Game) drawContacts() {
	g.cv.BeginPath()
	g.cv.SetFillStyle("#00897B")
	g.cv.SetFont(g.fonts.small, 15)
	text := fmt.Sprint("Game's repo:")
	g.cv.FillText(text, g.param.gameW+130, g.param.gameH-10)
	text = fmt.Sprint("Telegram:")
	g.cv.FillText(text, g.param.gameW+130, g.param.gameH+10)

	g.cv.SetFillStyle("#1A237E")
	text = fmt.Sprint("@DenKhan")
	g.cv.FillText(text, g.param.gameW+200, g.param.gameH+10)
	text = fmt.Sprint("@GitHub")
	g.cv.FillText(text, g.param.gameW+225, g.param.gameH-10)

	onTheLinc := func(x, x1, x2 float64, y, y1, y2 float64) bool {
		return x >= x1 && x <= x2 && y <= y1 && y >= y2
	}

	g.wnd.MouseUp = func(button, x, y int) {
		if button == 1 && onTheLinc(float64(x), g.param.gameW+200, g.param.gameW+300,
			float64(y), g.param.gameH+10, g.param.gameH-5) {
			if err := openURL("https://t.me/DenKhan"); err != nil {
				log.Println(err)
			}
		} else if button == 1 && onTheLinc(float64(x), g.param.gameW+225, g.param.gameW+300,
			float64(y), g.param.gameH-10, g.param.gameH-20) {
			if err := openURL("https://github.com/DenisKhanov/Snake"); err != nil {
				log.Println(err)
			}
		}
	}
	g.cv.Stroke()
}

// drawGameOver displays the "Game Over" message and instructions on the screen.
//
// This method renders a prominent "Game Over" text and provides instructions to restart or exit the game.
// The text is displayed at the specified coordinates.
//
// Parameters:
// - x, y (float64): The starting position for rendering the "Game Over" text.
func (g *Game) drawGameOver(x, y float64) {
	g.cv.BeginPath()
	g.cv.SetFillStyle("#C2185B")
	g.cv.SetFont(g.fonts.main, 60)
	text := fmt.Sprintf("Game over")
	g.cv.FillText(text, x, y)
	g.cv.Stroke()

	g.cv.BeginPath()
	g.cv.SetFillStyle("#1B5E20")
	g.cv.SetFont(g.fonts.small, 15)
	text = fmt.Sprintf("Press 'ENTER' for start new game")
	g.cv.FillText(text, x-60, y+40)
	text = fmt.Sprintf("Press 'ESC' for close game")
	g.cv.FillText(text, x+225, y+40)
	g.cv.Stroke()

}
