// Package game contains the core functionality for the Snake game, including game logic, rendering, geometry handling, and snake behavior.
package game

import (
	_ "embed"
	"fmt"
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"
)

//go:embed  assets/samuraiterrapingradital.ttf
var samuraiFont []byte

//go:embed assets/Dejavusansmono.ttf
var dejavuFont []byte

//go:embed assets/Righteous-Regular.ttf
var righteousFont []byte

//go:embed assets/SnakeGO.png
var backgroundImage []byte

const (
	cellsCount = 20
	startSpeed = 300
)

// Fonts holds the font styles used in the game for different text stile.
type Fonts struct {
	main   *canvas.Font
	middle *canvas.Font
	small  *canvas.Font
}

// GameParam holds the configuration parameters for the game window and game area.
// It includes the dimensions of the window and game area, as well as the speed of the game.
type GameParam struct {
	windowW int
	windowH int
	gameW   float64
	gameH   float64
	speed   int
}

// NewGameParam creates and returns a new instance of GameParam with default values.
// These values include the window size, game area size, and the initial speed of the game.
// The returned GameParam is used to configure the game environment when creating a new game.
func NewGameParam() *GameParam {
	return &GameParam{
		windowW: 1030,
		windowH: 730,
		gameW:   700.0,
		gameH:   700.0,
		speed:   startSpeed,
	}
}

// Game represents the state and behavior of the Snake game. It holds the
// game configuration, game area properties, and manages the snake, food,
// score, and game state.
type Game struct {
	cv  *canvas.Canvas
	wnd *sdlcanvas.Window

	param *GameParam
	snake *Snake
	food  Point
	fonts Fonts

	gameAreaSP Point
	gameAreaEP Point
	cellW      float64
	cellH      float64
	side       float64

	score          int
	ateFood        int
	gameOver       bool
	needMove       bool
	needUpdateInfo bool
}

// NewGame creates a new instance of the Game struct.
// It initializes the game window and canvas with specified window size
// and other game parameters, such as the game area dimensions and cell sizes.
//
// The function creates the window with a title and calculates the width and height
// of each cell in the grid based on the game area dimensions and a predefined constant
// `cellsCount` (which determines the number of cells in the grid).
// If the window creation fails, the function will panic.
func NewGame(param *GameParam) *Game {
	wnd, cv, err := sdlcanvas.CreateWindow(param.windowW, param.windowH, "Welcome to the Snake game written in Golang")
	if err != nil {
		panic(err)
	}

	cellW := param.gameW / cellsCount
	cellH := param.gameH / cellsCount
	return &Game{
		cv:         cv,
		wnd:        wnd,
		param:      param,
		gameAreaSP: Point{15, 15},
		gameAreaEP: Point{15 + param.gameW, 15 + param.gameH},
		cellW:      cellW,
		cellH:      cellH,
		side:       math.Min(cellW-1*2, cellH-1*2),
		gameOver:   false,
	}
}

// initFonts initializes the fonts used in the game.
// It loads three different font files for different text styles
// and assigns them to the game's `fonts` field.
//
// The function will panic if any font fails to load.
func (g *Game) initFonts() {
	mainFont, err := g.cv.LoadFont(samuraiFont)
	if err != nil {
		panic(err)
	}
	instructionFont, err := g.cv.LoadFont(dejavuFont)
	if err != nil {
		panic(err)
	}
	easyFont, err := g.cv.LoadFont(righteousFont)
	if err != nil {
		panic(err)
	}

	fonts := Fonts{
		main:   mainFont,
		middle: instructionFont,
		small:  easyFont,
	}
	g.fonts = fonts
}

// setSnake sets the provided snake instance to the game object.
// It assigns the passed *Snake object to the `g.snake` field,
// allowing the game to track and update the snake's state.
func (g *Game) setSnake(snake *Snake) {
	g.snake = snake
}

// run starts the main game loop for the Snake game.
// It initializes the game logic handling, food generation, and rendering loop.
func (g *Game) run() {
	go g.handleGameLogic()
	g.foodGeneration()
	g.renderLoop()
}

// handleGameLogic manages the core game loop, including snake movement, collision detection,
// food consumption, and scoring updates. It uses a timer to control the snake's speed
// and processes game logic in each iteration.
//
// The method performs the following tasks:
// - Processes player input to update the snake's Direction.
// - Checks for collisions with walls or the snake's own body, setting the gameOver flag if necessary.
// - Updates the snake's size and score if it eats food.
// - Adjusts the game's speed dynamically based on the snake's progress.
// - Resets the timer at the end of each loop iteration to maintain consistent movement intervals.
//
// This method runs continuously until the game is over or the application is exited.
func (g *Game) handleGameLogic() {
	var snakeTimer = time.NewTimer(time.Millisecond * time.Duration(g.param.speed))
	//keyboard scan
	g.processInput()
	//loop
	for {
		<-snakeTimer.C
		newPos := g.snake.Direction.Exec(g.snake.Parts[0])
		if g.collidesWithWall(newPos) {
			g.gameOver = true
		}
		//we cut off the snake if there is a new position on its body
		if g.snake.CutIfSnake(newPos) {
			newSize := len(g.snake.Parts)
			g.score = g.score / g.snake.Size * newSize //correct score according new snake size
			g.snake.Size = newSize
			g.needUpdateInfo = true
		}

		//snakes move and eat food
		if newPos == g.food {
			g.snake.Add(newPos)
			g.foodGeneration()
			g.ateFood += 1
			g.snake.Size++
			g.param.speed -= 5
			g.score += g.calculateScore(newPos)
			g.needUpdateInfo = true
		} else if !g.gameOver {
			g.snake.Move(g.snake.Direction)
			g.needMove = true
		}
		snakeTimer.Reset(time.Millisecond * time.Duration(g.param.speed))
	}
}

// foodGeneration generates a new food position on the grid.
//
// It randomly selects coordinates within the grid (cellsCount) and ensures
// the position does not overlap with the snake's body. The new position is
// stored in g.food.
func (g *Game) foodGeneration() {
	for {
		randX := rand.Intn(cellsCount)
		randY := rand.Intn(cellsCount)
		newPoint := Point{float64(randX), float64(randY)}
		check := true
		if g.snake.IsSnake(newPoint) {
			check = false
		}
		if check {
			g.food = newPoint
			return
		}
	}
}

// calculateScore calculates the score based on the position of the food consumed by the snake.
// The score is determined by the proximity of the food to the edges or corners of the game field,
// with higher rewards for food closer to the corners and edges.
//
// Parameters:
// - pos (Point): The position of the food that was consumed.
//
// Returns:
// - int: The calculated score based on the food's position and the current game speed.
//
// Scoring logic:
// - Food in the corners of the game field yields the highest score (multiplied by 4).
// - Food on the edges but not in the corners yields a moderate score (multiplied by 2).
// - Food elsewhere yields the base score (no multiplier).
func (g *Game) calculateScore(pos Point) int {
	switch {
	case pos.IsCorner():
		return 1000 / g.param.speed * 4
	case pos.IsEdge():
		return 1000 / g.param.speed * 2
	default:
		return 1000 / g.param.speed
	}
}

// collidesWithWall checks if the given position causes a collision with the game field boundaries.
//
// Parameters:
// - newPos (Point): The position to check for a boundary collision.
//
// Returns:
// - bool: True if the position is outside the game field boundaries, otherwise false.
//
// The method verifies if the X or Y coordinates of the position are less than 0
// or exceed the maximum number of cells in the game field (`cellsCount`).
func (g *Game) collidesWithWall(newPos Point) bool {
	return newPos.X < 0 || newPos.X >= cellsCount || newPos.Y < 0 || newPos.Y >= cellsCount
}

// processInput handles keyboard input during the game.
//
// This method assigns a function to the `KeyUp` event of the game window.
//
// This method dynamically updates the behavior of the game in response to player input.
func (g *Game) processInput() {
	g.wnd.KeyUp = func(code int, rn rune, name string) {
		//game over keys
		if g.gameOver {
			switch name {
			case "Enter":
				g.restartGame()
				g.gameOver = false
				return
			case "Escape":
				sdl.Quit()
				os.Exit(1)
			}
		}
		//Direction's keys  ← ↑ → ↓
		if 79 <= code && code <= 82 && g.needMove {
			newDir := g.snake.Direction.FromKey(code)
			if !g.snake.Direction.CheckParallel(newDir) {
				g.snake.Direction = newDir
				g.needMove = false
			}
		}
	}
}

// renderLoop manages the rendering process and continuously updates the game window.
//
// This method uses the `MainLoop` function to handle the rendering cycle, drawing the game's visual elements on each frame.
//
// This loop ensures that the game visuals are consistently updated based on the game's current state.
func (g *Game) renderLoop() {
	logo, err := g.cv.LoadImage(backgroundImage)
	if err != nil {
		log.Println(err)
	}
	g.drawGameInfo()
	//draw game instructions for the player
	g.drawInstructions()
	// draw creator information
	g.drawAboutCreator(g.param.gameW+20, g.param.gameH-50)
	//draw contact details
	g.drawContacts()
	//draw logo
	g.cv.DrawImage(logo, g.param.gameW+40, g.param.gameH-350, 250, 250)

	//start loop
	g.wnd.MainLoop(func() {
		//clear game world
		g.cv.ClearRect(0, 0, g.param.gameW, g.param.gameH+30) // update game area
		//draw world
		g.drawWorld()
		//draw grid within the game area
		g.drawGridGameArea()

		g.drawFPS()
		//draw snake
		g.drawSnake()
		//draw food
		g.drawApple(g.gameAreaSP.X+g.food.X*g.cellW+1, g.gameAreaSP.Y+g.food.Y*g.cellH+1, g.side)
		// draw "Game Over" screen, if the game has ended
		if g.gameOver {
			g.drawGameOver(g.param.gameW/2-160, g.param.gameH/2)
		}
		// this is an optimization to avoid drawing relatively static information every frame
		if g.needUpdateInfo {
			//clear game world
			g.cv.ClearRect(750, 0, 280, 200) //update only GameInfo area
			//draw game information, such as score and speed
			g.drawGameInfo()
			g.needUpdateInfo = false
		}
	})
}

// restartGame resets the game state to its initial values, effectively restarting the game.
//
// This method resets the snake's position and state, sets the score and food count to zero,
// restores the default game speed, and flags the game as not over.
func (g *Game) restartGame() {
	g.snake.Reset()
	g.score = 0
	g.ateFood = 0
	g.param.speed = 300
	g.gameOver = false
}

// openURL opens the specified URL in the default web browser based on the operating system.
//
// It determines the appropriate command to use for opening the URL based on the current
// operating system (Windows, macOS, or Linux) and executes it. If there is an error executing
// the command, it returns an error message.
//
// Parameters:
// - url (string): The URL to open in the default web browser.
//
// Returns:
// - error: An error if the URL could not be opened; otherwise, nil.
func openURL(url string) error {
	var cmd *exec.Cmd

	// Define the command depending on the operating system
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/C", "start", url) // For Windows
	case "darwin":
		cmd = exec.Command("open", url) // For macOS
	default:
		cmd = exec.Command("xdg-open", url) // For Linux
	}

	// run command and return error if it's having
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("error opening URL: %v", err)
	}
	return nil
}

// RunGame initializes and starts a new game of Snake.
// It creates a new Snake object, resets it, initializes game parameters, and runs the game.
//
// The function does the following:
// 1. Creates a new Snake instance using NewSnake() and resets it.
// 2. Initializes the game parameters with NewGameParam().
// 3. Creates a new game instance with NewGame(gameParam) and sets up the game environment.
// 4. Initializes fonts for rendering and sets the Snake for the game.
// 5. Starts the game loop with the run method.
func RunGame() {
	snake := NewSnake()
	snake.Reset()
	gameParam := NewGameParam()
	game := NewGame(gameParam)
	game.initFonts()
	game.setSnake(snake)
	game.run()
}
