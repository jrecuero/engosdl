package main

import (
	"fmt"

	"github.com/jrecuero/engosdl"
)

func main() {
	fmt.Println("rol player game")
	if engine := engosdl.NewEngine("flier", 800, 400, NewGameManager("rolling-game-manager")); engine != nil {
		engine.RunEngine(nil)
	}
}
