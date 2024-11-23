![Snake logo](https://github.com/github.com/DenisKhanov/Snake/game/assets/SnakeGO.png)


# Snake Game in Go V1.0

Welcome to the Snake Game written in Go! This repository contains the implementation of a classic Snake game using the Go programming language and SDL2 for graphics rendering. The game features a snake that grows as it eats food, avoids obstacles, and plays with simple controls.

![GIF](https://github.com/github.com/DenisKhanov/Snake/game/assets/snake.gif)

## Features

- **Classic Snake Gameplay**: Control the snake with the arrow keys to eat food and grow longer.
- **Game Instructions**: Easy-to-read game instructions displayed on the screen.
- **High Scores**: Tracks your score and displays it in real-time.
- **Game Over Mechanism**: The game ends when the snake collides with the walls.
- **Graphics**: Custom fonts and colorful visuals to enhance the user experience.

## Prerequisites

Before you begin, ensure you have the following dependencies installed:

- **Go**: A Go runtime environment to compile and run the game.
- **SDL2**: Used for graphical rendering. The `SDL2.dll` and `libmcfgthread-1.dll` files are embedded in the project for Windows users, and will automatically be extracted when running the game.

## Installation

1. **Clone the repository**:
    ```bash
    git clone https://github.com/yourusername/snake-game.git
    cd snake-game
    ```

2. **Install Go dependencies**:
   If you're using Go modules, simply run:
    ```bash
    go mod tidy
    ```

3. **Run the game**:
   To run the game, execute the `main.go` file:
    ```bash
    go run main.go
    ```

## Download the executable file

### Run on Linux (ubuntu/debian)
1. You can download [`SnakeGO`](https://github.com/github.com/DenisKhanov/Snake/SnakeGO)file for Linux.
2. You need install SDL2 library
    ```bash
    sudo apt install libsdl2-dev
    ```
   and install OpenGL
    ```bash
    sudo apt install mesa-common-dev libgl1-mesa-dev
    ```
3. You can run `SnakeGO`

### Run on Windows

1. You need download [`SnakeGO.exe`](https://github.com/github.com/DenisKhanov/Snake/SnakeGO.exe) file for Windows.
2. [`SDL2.dll`](https://github.com/github.com/DenisKhanov/Snake/cmd/SDL2.dll) and [`libmcfgthread-1.dll`](https://github.com/github.com/DenisKhanov/Snake/cmd/libmcfgthread-1.dll) files should be created automatically when...
   files should be created automatically when running on windows.
3. If this did not happen, then you can download these files by clicking on them and place them in the directory next to the executable file.


## How to Play

- Use the **arrow keys ← ↑ → ↓** to control the direction of the snake.
- **Eat food** to grow the snake.
- The game ends if the snake collides with the boundaries of the game area.
- Track your **score** and how many food items you've eaten on the right side of the screen.
- **Restart the game** after it ends by pressing **ENTER** key.

## Key Functions and Features

### `Game` Struct
This is the core struct of the game which holds all the data and state information related to the game.
```go
type Game struct {
    cv         *canvas.Canvas
    wnd        *sdlcanvas.Window
    param      *GameParam
    snake      *Snake
    food       Point
    fonts      Fonts
    gameAreaSP Point
    gameAreaEP Point
    cellW      float64
    cellH      float64
    side       float64
    score      int
    ateFood    int
    gameOver   bool
    needMove   bool
    needUpdateInfo bool
}
```

### Movement Directions
The game uses the `Dir` type to handle the snake's movement. Directions are encoded using constants:
```go
const (
    up = iota
    right
    down
    left
)
```

### Game Controls
The user can control the snake using the arrow keys:
- **Up Arrow**: Move Up
- **Right Arrow**: Move Right
- **Down Arrow**: Move Down
- **Left Arrow**: Move Left

### Snake Logic
The `Snake` struct holds the segments of the snake and manages its movement. The snake's body is represented as a slice of `Point` structs, and its movement is handled by the `Exec` method:
```go
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
```

### Game Logic
- **Direction Check**: The snake cannot reverse direction. This is managed by the `CheckParallel` method:
```go
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
```

### Fonts and Rendering
Custom fonts are loaded for the game interface from ./game/assets:
```go
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
```

### Game Parameters
The `GameParam` struct defines the window and game area sizes:
```go
type GameParam struct {
    windowW int
    windowH int
    gameW   float64
    gameH   float64
    speed   int
}
```

## Contributing

Feel free to fork this repository, open an issue, or create a pull request to contribute to this project. If you have any suggestions or improvements, I’d love to hear from you!


---

Enjoy playing the Snake game! 🐍😄