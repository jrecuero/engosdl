package engosdl

// IGameManager represents the game manager interface.
type IGameManager interface {
	IObject
	DoFrameEnd()
	DoFrameStart()
}

// GameManager is the default implementation for game manager interface.
type GameManager struct {
	*Object
}

// NewGameManager creates a new game manager instance.
func NewGameManager(name string) *GameManager {
	Logger.Trace().Str("game-manager", name).Msg("new game-manager")
	result := &GameManager{
		Object: NewObject(name),
	}
	return result
}

// DoFrameEnd is called at the end of the frame.
func (h *GameManager) DoFrameEnd() {
}

// DoFrameStart is called  at the start of the frame.
func (h *GameManager) DoFrameStart() {
}
