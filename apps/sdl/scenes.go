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
	titleText := components.NewText("title-text", "fonts/lato.ttf", 32, sdl.Color{R: 0, G: 0, B: 255}, "PLAY")
	titleText.DefaultAddDelegateToRegister()
	titleText.AddDelegateToRegister(nil, nil, &components.Keyboard{}, func(params ...interface{}) bool {
		key := params[0].(int)
		if key == sdl.SCANCODE_RETURN {
			engosdl.GetEngine().GetSceneManager().SetActiveNextScene()
		}
		return true
	})
	titleText.AddDelegateToRegister(nil, nil, &components.Mouse{}, func(params ...interface{}) bool {
		x := float64(params[0].(int32))
		y := float64(params[1].(int32))
		// button := params[2].(uint32)
		x1, y1, w1, h1 := title.GetTransform().GetRectExt()
		// fmt.Printf("title at (%f, %f) (%f, %f)\n", x1, y1, w1, h1)
		// fmt.Printf("mouse has clicked %d at (%f, %f)\n", button, x, y)
		if x >= x1 && x <= (x1+w1) && y >= y1 && y <= (y1+h1) {
			// fmt.Print("click inside text")
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
	titleMouse := components.NewMouse("title-mouse", true)

	title.AddComponent(titleText)
	title.AddComponent(titleKeyboard)
	title.AddComponent(titleOutOfBounds)
	title.AddComponent(titleMoveIt)
	title.AddComponent(titleMouse)

	scene.AddEntity(title)
	return true
}
