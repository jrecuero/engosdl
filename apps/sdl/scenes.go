package main

import (
	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
	"github.com/veandco/go-sdl2/sdl"
)

// // createSceneStats creates stats scene
// func createSceneStats(engine *engosdl.Engine, scene engosdl.IScene) bool {
// 	message := engosdl.NewEntity("message")
// 	message.GetTransform().SetPosition(engosdl.NewVector(50, 100))
// 	messageText := components.NewText("message-text", "fonts/lato.ttf", 16, sdl.Color{R: 0, G: 255, B: 0}, "player stats", engine.GetRenderer())
// 	messageText.DefaultAddDelegateToRegister()
// 	messageText.AddDelegateToRegister(nil, nil, &components.Keyboard{}, func(params ...interface{}) bool {
// 		key := params[0].(int)
// 		// if key == sdl.SCANCODE_P {
// 		if key == sdl.SCANCODE_RETURN {
// 			engosdl.GetEngine().GetSceneManager().SwapBack()
// 		}
// 		return true
// 	})
// 	messageKeyboard := components.NewKeyboard("title-keyboard")
// 	messageKeyboard.DefaultAddDelegateToRegister()

// 	message.AddComponent(messageText)
// 	message.AddComponent(messageKeyboard)

// 	scene.AddEntity(message)
// 	return true
// }

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
