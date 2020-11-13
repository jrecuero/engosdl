package main

import (
	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
	"github.com/veandco/go-sdl2/sdl"
)

// GameManager is the application game manager.
type GameManager struct {
	*engosdl.GameManager
}

var _ engosdl.IGameManager = (*GameManager)(nil)

// NewGameManager creates a new game manager instance.
func NewGameManager(name string) *GameManager {
	engosdl.Logger.Trace().Str("game-manager", name).Msg("new game-manager")
	return &GameManager{
		GameManager: engosdl.NewGameManager(name),
	}
}

// CreateAssets is called before game engine starts in order to create all
// required assets.
func (h *GameManager) CreateAssets() {
	createAssets(engosdl.GetEngine())
}

// DoFrameEnd is called at the end of the engine frame.
func (h *GameManager) DoFrameEnd() {
}

// DoFrameStart is called at the start of engine frame.
func (h *GameManager) DoFrameStart() {
	// Count number of enemies at the end of a frame. Takes some functionality
	// from the enemy controller.
	activeScene := engosdl.GetSceneManager().GetActiveScene()
	if activeScene.GetTag() == "battle" {
		enemies := activeScene.GetEntitiesByTag("enemy")
		// fmt.Printf("there are %d enemies\n", len(enemies))
		if len(enemies) == 0 {
			winner := engosdl.NewEntity("winner")
			winnerKeyboard := components.NewKeyboard("winner-keyboard")
			winnerKeyboard.DefaultAddDelegateToRegister()
			winnerText := components.NewText("winner-text", "fonts/lato.ttf", 24, sdl.Color{R: 0, G: 0, B: 255}, "You Won...type any key", engosdl.GetRenderer())
			winnerText.DefaultAddDelegateToRegister()
			winnerText.AddDelegateToRegister(nil, nil, &components.Keyboard{}, func(params ...interface{}) bool {
				key := params[0].(int)
				if key == sdl.SCANCODE_RETURN {
					engosdl.GetEngine().GetSceneManager().SetActiveFirstScene()
				}
				return true
			})
			winner.AddComponent(winnerKeyboard)
			winner.AddComponent(winnerText)
			winner.SetTag("winner")
			activeScene.AddEntity(winner)
			// // engosdl.GetEngine().GetSceneManager().SetActiveFirstScene()
			// winner := activeScene.GetEntityByName("winner")
			// winner.SetActive(true)
		}
	}
}
