package modules

import (
	"fmt"
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
	"math/rand"
	"sync"
	"time"
)

const (
	GameW = 700.0
	GameH = 700.0
)

type Game struct {
	cv  *canvas.Canvas
	wnd *sdlcanvas.Window

	snake    *Snake
	food     Point
	score    int
	ateFood  int
	speed    int
	gameOver bool
	needMove bool
}

func NewGame() *Game {
	wnd, cv, err := sdlcanvas.CreateWindow(1000, 730, "Welcome in snake game")
	if err != nil {
		panic(err)
	}
	return &Game{
		cv:       cv,
		wnd:      wnd,
		speed:    300,
		gameOver: false,
	}
}
func (g *Game) SetSnake(snake *Snake) {
	g.snake = snake
}

func (g *Game) Run() {
	go g.SnakeMovement()
	g.FoodGeneration()
	g.RenderLoop()
}

func (g *Game) FoodGeneration() {
	minNum := 0
	maxNum := 20
	for {
		if !g.gameOver {
			randX := rand.Intn(maxNum-minNum) + minNum
			randY := rand.Intn(maxNum-minNum) + minNum
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
}

func (g *Game) SnakeMovement() {
	var snakeTimer = time.NewTimer(time.Duration(g.speed) * time.Millisecond)
	var snakeDirection Dir = right
	var snakeLock sync.Mutex

	resetTime := func() {
		snakeTimer.Reset(time.Duration(g.speed) * time.Millisecond)
	}
	resetTime()

	//scan keyboard
	g.wnd.KeyUp = func(code int, rn rune, name string) {
		if code < 79 && code > 82 || g.needMove {
			return
		}

		snakeLock.Lock()

		newDir := snakeDirection
		switch code {
		case 80: //left
			newDir = left
		case 82: //up
			newDir = down
		case 79: //right
			newDir = right
		case 81: //down
			newDir = up
		}

		if !snakeDirection.CheckParallel(newDir) {
			snakeDirection = newDir
		}

		snakeLock.Unlock()
	}
	//loop
	for {
		<-snakeTimer.C
		snakeLock.Lock()
		if !g.gameOver {
			newPos := snakeDirection.Exec(g.snake.Parts[0])
			if newPos.X < 0 || newPos.X > 19 || newPos.Y < 0 || newPos.Y > 19 {
				g.gameOver = true
			}

			if g.snake.CutIfSnake(newPos) {
				newSize := len(g.snake.Parts)
				g.score = g.score / g.snake.Size * newSize //correct score according new snake size
				g.snake.Size = newSize
			}

			//is food
			isFood := false
			if newPos == g.food {
				g.ateFood += 1
				g.snake.Size++
				switch { //algorithm calculate score according food position
				case newPos.X == 0 && newPos.Y == 0 || newPos.X == 0 && newPos.Y == 19 ||
					newPos.X == 19 && newPos.Y == 0 || newPos.X == 19 && newPos.Y == 19:
					g.score += 1000 / g.speed * 4
				case newPos.X == 0 || newPos.Y == 0 || newPos.X == 19 || newPos.Y == 19:
					g.score += 1000 / g.speed * 2
				default:
					g.score += 1000 / g.speed
				}
				g.FoodGeneration()
				g.snake.Add(newPos)
				g.FoodGeneration()
				g.speed -= 5
				isFood = true
			}

			if !isFood && !g.gameOver {
				g.snake.Move(snakeDirection)
				g.needMove = false
			}
		}
		snakeLock.Unlock()
		resetTime()
	}

}
func (g *Game) RenderLoop() {
	gameAreaSP := Point{15, 15}
	gameAreaEP := Point{15 + GameW, 15 + GameH}
	cellW := GameW / 20
	cellH := GameH / 20

	mainFont, err := g.cv.LoadFont("./asset/samuraiterrapingradital.ttf")
	if err != nil {
		panic(err)
	}
	secondFont, err := g.cv.LoadFont("./asset/deadpack.ttf")
	if err != nil {
		panic(err)
	}

	g.wnd.MainLoop(func() {
		//clear world
		g.cv.ClearRect(0, 0, 1000, 750)
		//draw world
		g.cv.BeginPath()
		g.cv.SetFillStyle("#78909C")
		g.cv.FillRect(gameAreaSP.X, gameAreaSP.Y, gameAreaEP.X-15, gameAreaEP.Y-15)
		g.cv.Stroke()

		//draw grid
		g.cv.BeginPath()
		g.cv.SetStrokeStyle("#5D4037")
		g.cv.SetLineWidth(0.5)
		for i := 0; i < 20+1; i++ {
			g.cv.MoveTo(gameAreaSP.X+float64(i)*cellH, gameAreaSP.Y)
			g.cv.LineTo(gameAreaSP.X+float64(i)*cellH, gameAreaEP.Y)
			g.cv.MoveTo(gameAreaSP.X, gameAreaSP.Y+float64(i)*cellW)
			g.cv.LineTo(gameAreaEP.X, gameAreaSP.Y+float64(i)*cellW)
		}
		g.cv.Stroke()

		//draw snake
		g.cv.BeginPath()
		// Color for snake's head
		if len(g.snake.Parts) > 0 {
			g.cv.SetFillStyle("#BF360C")
			head := g.snake.Parts[0]
			g.cv.FillRect(
				gameAreaSP.X+head.X*cellW+1,
				gameAreaSP.Y+head.Y*cellH+1,
				cellW-1*2,
				cellH-1*2,
			)
		}
		//color for snake's body
		g.cv.SetFillStyle("#FF5722")
		for _, point := range g.snake.Parts[1:] {
			g.cv.FillRect(
				gameAreaSP.X+point.X*cellW+1,
				gameAreaSP.Y+point.Y*cellH+1,
				cellW-1*2,
				cellH-1*2,
			)
		}
		g.cv.Stroke()

		//draw food
		g.cv.BeginPath()
		g.cv.SetFillStyle("#4CAF50")
		g.cv.FillRect(
			gameAreaSP.X+g.food.X*cellW+1,
			gameAreaSP.Y+g.food.Y*cellH+1,
			cellW-1*2,
			cellH-1*2,
		)
		g.cv.Stroke()

		//draw score
		g.cv.BeginPath()
		g.cv.SetFont(mainFont, 25)
		text := fmt.Sprintf("Your score: %d", g.score)
		g.cv.FillText(text, GameW+50, 50)

		// food
		g.cv.BeginPath()
		g.cv.SetFont(mainFont, 25)
		text = fmt.Sprintf("You ate food: %d", g.ateFood)
		g.cv.FillText(text, GameW+50, 85)

		// speed
		g.cv.BeginPath()
		g.cv.SetFont(mainFont, 25)
		text = fmt.Sprintf("Your speed: %d", 300-g.speed)
		g.cv.FillText(text, GameW+50, 120)

		// about creator
		g.cv.BeginPath()
		g.cv.SetFillStyle("#00897B")
		g.cv.SetFont(secondFont, 10)
		text = fmt.Sprint("This game created by Denis Khanov")
		g.cv.FillText(text, GameW+40, GameH-30)
		text = fmt.Sprint("Telegram    DenKhan")
		g.cv.FillText(text, GameW+150, GameH)

		// game over
		if g.gameOver {
			g.cv.BeginPath()
			g.cv.SetFillStyle("#C2185B")
			g.cv.SetFont(mainFont, 60)
			text = fmt.Sprintf("Game over")
			g.cv.FillText(text, GameW/2-160, GameH/2)
		}
		g.cv.Stroke()
	})
}
