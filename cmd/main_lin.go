//go:build linux

package main

import (
	"github.com/DenisKhanov/Snake/game"
)

// main is the entry point of the program that performs the following steps:
//
// The `RunGame` function is called to start the game.
func main() {
	game.RunGame()
}
