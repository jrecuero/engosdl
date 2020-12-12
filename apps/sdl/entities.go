package main

import (
	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
	"github.com/veandco/go-sdl2/sdl"
)

// createEntityBackground creates the background
func createEntityBackground(engine *engosdl.Engine, filename string) engosdl.IEntity {
	bg := engosdl.NewEntity("background")
	bg.SetLayer(engosdl.LayerBackground)
	bgSprite := components.NewScrollSprite("bg-sprite", filename, engosdl.FormatBMP)
	bgSprite.DefaultAddDelegateToRegister()
	bgSprite.SetScroll(engosdl.NewVector(0, -1))
	// bgSprite.SetCamera(&engosdl.Rect{X: 0, Y: 0, W: 400, H: 800})
	bgMoveTo := components.NewMoveTo("bg-move", engosdl.NewVector(0, -5))
	bgMoveTo.DefaultAddDelegateToRegister()
	bg.AddComponent(bgSprite)
	bg.AddComponent(bgMoveTo)
	return bg
}

// createWinner creates text at the end of battle scene.
func createWinner(engine *engosdl.Engine) engosdl.IEntity {
	winner := engosdl.NewEntity("winner")
	winnerKeyboard := components.NewKeyboard("winner-keyboard", map[int]bool{sdl.SCANCODE_RETURN: false})
	winnerKeyboard.DefaultAddDelegateToRegister()
	winnerText := components.NewText("winner-text", "fonts/lato.ttf", 24, sdl.Color{R: 0, G: 0, B: 255}, "You Won..\ntype any key")
	winnerText.DefaultAddDelegateToRegister()
	winnerText.AddDelegateToRegister(nil, nil, &components.Keyboard{}, func(params ...interface{}) bool {
		key := params[0].(int)
		if key == sdl.SCANCODE_RETURN {
			if engosdl.GetSceneManager().GetActiveScene().GetName() == "play-scene-1" {
				engosdl.GetSceneManager().SetActiveNextScene()
			} else if engosdl.GetSceneManager().GetActiveScene().GetName() == "play-scene-2" {
				engosdl.GetSceneManager().SetActiveFirstScene()
			}
		}
		return true
	})
	winner.AddComponent(winnerKeyboard)
	winner.AddComponent(winnerText)
	winner.SetTag("winner")
	return winner
}
