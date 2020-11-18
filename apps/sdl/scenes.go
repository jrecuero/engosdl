package main

import (
	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
	"github.com/veandco/go-sdl2/sdl"
)

// createSceneTitle creates title scene.
func createSceneTitle(engine *engosdl.Engine, scene engosdl.IScene) bool {
	title := engosdl.NewEntity("title")
	title.GetTransform().SetPosition(engosdl.NewVector(175, 250))
	titleText := components.NewText("title-text", "fonts/lato.ttf", 32, sdl.Color{R: 0, G: 0, B: 255}, "PLAY", engine.GetRenderer())
	titleText.DefaultAddDelegateToRegister()
	titleText.AddDelegateToRegister(nil, nil, &components.Keyboard{}, func(params ...interface{}) bool {
		key := params[0].(int)
		if key == sdl.SCANCODE_RETURN {
			engosdl.GetEngine().GetSceneManager().SetActiveNextScene()
		}
		return true
	})
	titleKeyboard := components.NewKeyboard("title-keyboard")
	titleKeyboard.DefaultAddDelegateToRegister()
	titleOutOfBounds := components.NewOutOfBounds("title-out-of-bounds", true)
	titleOutOfBounds.DefaultAddDelegateToRegister()
	titleMoveIt := components.NewMoveIt("title-move-it", engosdl.NewVector(5, 0))
	titleMoveIt.DefaultAddDelegateToRegister()

	title.AddComponent(titleText)
	title.AddComponent(titleKeyboard)
	title.AddComponent(titleOutOfBounds)
	title.AddComponent(titleMoveIt)

	scene.AddEntity(title)
	return true
}
