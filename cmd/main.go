package main

import "Snake/modules"

func main() {
	snake := modules.NewSnake()
	snake.Reset()
	game := modules.NewGame()
	game.SetSnake(snake)
	game.Run()
}
