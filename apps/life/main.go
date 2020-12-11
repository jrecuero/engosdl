package main

import (
	"fmt"

	"github.com/jrecuero/engosdl"
)

func main() {
	fmt.Println("life game")
	if engine := engosdl.NewEngine("life", 800, 800, NewGameManager("life-game-manager")); engine != nil {
		engine.RunEngine(nil)
	}
}
