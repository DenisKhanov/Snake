//go:build windows

package main

import (
	_ "embed"
	"fmt"
	"github.com/DenisKhanov/Snake/game"
	"os"
)

//go:embed libmcfgthread-1.dll
var libmcfgthread []byte //need for run game on windows

//go:embed SDL2.dll
var sdl2 []byte //need for run game on windows

// main is the entry point of the program that performs the following steps:
// 1. Checks if the required DLL files (`libmcfgthread-1.dll` and `SDL2.dll`) exist in the current directory.
// 2. If any DLL is missing, it extracts the corresponding DLL file from embedded resources into the current directory.
// 3. If all DLLs are present (or successfully extracted), it runs the game using the `RunGame` function from the `game` package.
//
// The function checks for the existence of the necessary DLL files, and if any are missing,
// it calls `extractDLL` to extract the DLLs from the embedded byte slices (`libmcfgthread` and `sdl2`).
// If extraction fails, it prints an error message and exits the program with a non-zero status code.
//
// The `RunGame` function is called to start the game after ensuring the required DLLs are present.
func main() {
	limbFile := "libmcfgthread-1.dll"
	sdlFile := "SDL2.dll"
	// check, file is have
	if _, err := os.Stat(limbFile); err != nil {
		// extract embedded DLL in to the current directory
		if err = extractDLL("libmcfgthread-1.dll", libmcfgthread); err != nil {
			fmt.Println("Failed to extract DLL:", err)
			os.Exit(1)
		}
	}
	if _, err := os.Stat(sdlFile); err != nil {
		err = extractDLL("SDL2.dll", sdl2)
		if err != nil {
			fmt.Println("Failed to extract DLL:", err)
			os.Exit(1)
		}
	}
	game.RunGame()

}

// extractDLL saves the provided byte data to a file with the specified filename.
//
// This function attempts to write the given byte slice (`data`) to a file specified by `filename`
// with permissions set to 0644. If an error occurs during the file writing process, it returns
// an error with a descriptive message.
//
// Parameters:
//
//	filename (string): The name of the file to which the data will be written.
//	data ([]byte): The byte slice containing the data to be saved in the file.
//
// Returns:
//
//	error: If there is an error writing the data to the file, an error is returned; otherwise, nil.
func extractDLL(filename string, data []byte) error {
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing file %s: %w", filename, err)
	}
	return nil
}
