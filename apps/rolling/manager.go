package main

import (
	"fmt"

	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
	"github.com/veandco/go-sdl2/sdl"
)

// GameManager is the application game manager.
type GameManager struct {
	*engosdl.GameManager
}

var _ engosdl.IGameManager = (*GameManager)(nil)

// NewGameManager created a new game manager instance.
func NewGameManager(name string) *GameManager {
	engosdl.Logger.Trace().Str("game-manager", name).Msg("new game-manager")
	return &GameManager{
		GameManager: engosdl.NewGameManager(name),
	}
}

// CreateAssets creates all application assets and resources. it is called
// before game engine starts in order to create all required assets and
// resources.
func (h *GameManager) CreateAssets() {
	scenePlay := engosdl.NewScene("play-scene", "play")
	scenePlay.SetSceneCode(h.createScenePlay())
	engosdl.GetEngine().AddScene(scenePlay)
}

func (h *GameManager) createScenePlay() func(*engosdl.Engine, engosdl.IScene) bool {
	return func(engine *engosdl.Engine, scene engosdl.IScene) bool {
		player := engosdl.NewEntity("player")
		player.SetTag("player")
		mouse := components.NewMouse("player/mouse", true)
		player.AddComponent(mouse)

		console := engosdl.NewEntity("console")
		console.GetTransform().SetPositionXY(10, 200)
		consoleText := components.NewText("console/text", "fonts/fira2.ttf", 12, sdl.Color{}, " ")
		console.AddComponent(consoleText)
		message := ""

		lookButton := engosdl.NewEntity("player-look-button")
		lookButton.GetTransform().SetPositionXY(10, 50)
		lookButton.AddComponent(components.NewBox("look/box", &engosdl.Rect{}, sdl.Color{R: 255}, true))
		lookButton.GetComponent(&components.Box{}).SetCustomOnUpdate(func(c engosdl.IComponent) {
			entity := c.GetEntity()
			x, y, _ := sdl.GetMouseState()
			box := entity.GetComponent(&components.Box{}).(*components.Box)
			if entity.IsInside(engosdl.NewVector(float64(x), float64(y))) {
				box.Color = sdl.Color{G: 255}
			} else {
				box.Color = sdl.Color{R: 255}
			}

		})
		lookText := components.NewText("player/look-text", "fonts/fira.ttf", 32, sdl.Color{B: 255}, "LOOK")
		lookText.AddDelegateToRegister(nil, player, &components.Mouse{}, func(c engosdl.IComponent) func(params ...interface{}) bool {
			return func(params ...interface{}) bool {
				mousePos := engosdl.NewVector(float64(params[0].(int32)), float64(params[1].(int32)))
				if c.GetEntity().IsInside(mousePos) {
					message += "You clicked inside LOOK button\n"
					consoleText.SetMessage(message)
					fmt.Println("you clicked inside LOOK button")
				} else {
					fmt.Println("you clicked outside LOOK button")
				}
				return true
			}
		}(lookText))
		lookButton.AddComponent(lookText)
		player.AddChild(lookButton)

		moveButton := engosdl.NewEntity("player-move-button")
		moveButton.GetTransform().SetPositionXY(100, 50)
		moveButton.AddComponent(components.NewBox("move/box", &engosdl.Rect{}, sdl.Color{B: 255}, true))
		moveText := components.NewText("player/move-text", "fonts/fira.ttf", 32, sdl.Color{R: 255}, "MOVE")
		moveText.AddDelegateToRegister(nil, player, &components.Mouse{}, func(c engosdl.IComponent) func(params ...interface{}) bool {
			return func(params ...interface{}) bool {
				mousePos := engosdl.NewVector(float64(params[0].(int32)), float64(params[1].(int32)))
				if c.GetEntity().IsInside(mousePos) {
					message += "You clicked inside MOVE button\n"
					consoleText.SetMessage(message)
					fmt.Println("you clicked inside MOVE button")
				} else {
					fmt.Println("you clicked outside MOVE button")
				}
				return true
			}
		}(moveText))
		moveButton.AddComponent(moveText)
		player.AddChild(moveButton)

		scene.AddEntity(player)
		scene.AddEntity(console)
		return true
	}
}

// DoFrameEnd is called at the end of every engine frame.
func (h *GameManager) DoFrameEnd() {
}

// DoFrameStart is called at the start of the game frame.
func (h *GameManager) DoFrameStart() {
}

// DoInit initializes internal game manager resources.
func (h *GameManager) DoInit() {
}
