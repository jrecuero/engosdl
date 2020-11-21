package main

import "github.com/jrecuero/engosdl"

// GameManager is the flier application game manager.
type GameManager struct {
	*engosdl.GameManager
	player *engosdl.Entity
	score  *engosdl.Entity
}

var _ engosdl.IGameManager = (*GameManager)(nil)

// NewGameManager created a new game manager instance.
func NewGameManager(name string) *GameManager {
	engosdl.Logger.Trace().Str("game-manager", name).Msg("new game-manager")
	return &GameManager{
		GameManager: engosdl.NewGameManager(name),
	}
}

// CreateAssets creates all flier assets and resources. it is called before
// game engine starts in order to create all required assets and resources.
func (h *GameManager) CreateAssets() {
}

// DoFrameEnd is called at the end of every engine frame.
func (h *GameManager) DoFrameEnd() {
}

// DoFrameStart is called at the start of the game frame.
func (h *GameManager) DoFrameStart() {
}

// DoInit initializes internal game manager resources.
func (h *GameManager) DoInit() {
	h.player = engosdl.NewEntity("player")
	h.player.SetTag("player")

	h.score = engosdl.NewEntity("score")
}
