package engosdl

// IScene represents the interface for any game scene
type IScene interface {
	AddGameObject(IGameObject) bool
	GetGameObject(interface{}) IGameObject
}

// NewScene creates a new scene instance.
func NewScene() IScene {
	return nil
}
