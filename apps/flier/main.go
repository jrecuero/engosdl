package main

import (
	"fmt"

	"github.com/jrecuero/engosdl"
)

func main() {
	fmt.Println("flier game")
	if engine := engosdl.NewEngine("flier", 800, 400, NewGameManager("flier-game-manager")); engine != nil {
		engine.RunEngine(nil)
	}
}
