package main

import (
	"fmt"

	_ "net/http/pprof"

	"github.com/jrecuero/engosdl"
)

func main() {
	fmt.Println("engosdl app")
	if engine := engosdl.NewEngine("engosdl app", 400, 600, NewGameManager("app-game-manager")); engine != nil {
		// engine.DoInit()
		// createAssets(engine)
		engine.RunEngine(nil)
	}
}
