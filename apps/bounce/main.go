package main

import (
	"fmt"

	"github.com/jrecuero/engosdl"
)

func main() {
	fmt.Println("bounce game")
	if engine := engosdl.NewEngine("bounce", 800, 400, NewGameManager("rolling-game-manager")); engine != nil {
		engine.RunEngine(nil)
	}
}
