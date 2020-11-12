package main

import (
	"github.com/jrecuero/engosdl"
)

// GameManager is the application game manager.
type GameManager struct {
	*engosdl.Object
}

var _ engosdl.IGameManager = (*GameManager)(nil)

// NewGameManager creates a new game manager instance.
func NewGameManager(name string) *GameManager {
	engosdl.Logger.Trace().Str("game-manager", name).Msg("new game-manager")
	return &GameManager{
		Object: engosdl.NewObject(name),
	}
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
			// engosdl.GetEngine().GetSceneManager().SetActiveFirstScene()
			winner := activeScene.GetEntityByName("winner")
			winner.SetActive(true)
		}
	}
}
