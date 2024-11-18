package main

import "github.com/DenisKhanov/Snake/modules"

func main() {
	snake := modules.NewSnake()
	snake.Reset()
	game := modules.NewGame()
	game.SetSnake(snake)
	game.Run()
}
