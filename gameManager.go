package engosdl

// IGameManager represents the game manager interface.
type IGameManager interface {
	IObject
	CreateAssets()
	DoFrameEnd()
	DoFrameStart()
	DoInit()
	OnAfterUpdate()
	OnRender()
	OnStart()
	OnUpdate()
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

// CreateAssets is called before game engine starts in order to create all
// required assets.
func (h *GameManager) CreateAssets() {
	Logger.Trace().Str("game-manager", h.GetName()).Msg("create assets")
}

// DoInit initializes game manager resources. It is called after all engine
// resources have been initialized.
func (h *GameManager) DoInit() {
	Logger.Trace().Str("game-manager", h.GetName()).Msg("DoInit")
}

// DoFrameEnd is called at the end of the frame.
func (h *GameManager) DoFrameEnd() {
	Logger.Trace().Str("game-manager", h.GetName()).Msg("DoFrameEnd")
}

// DoFrameStart is called  at the start of the frame.
func (h *GameManager) DoFrameStart() {
	Logger.Trace().Str("game-manager", h.GetName()).Msg("DoFrameStart")
}

// OnAfterUpdate is called after all updates have been running.
func (h *GameManager) OnAfterUpdate() {
}

// OnRender runs after all renders have been running.
func (h *GameManager) OnRender() {
}

// OnStart starts the game manager.
func (h *GameManager) OnStart() {
	Logger.Trace().Str("game-manager", h.GetName()).Msg("OnStart")
}

// OnUpdate runs before any updates have been running.
func (h *GameManager) OnUpdate() {
}
