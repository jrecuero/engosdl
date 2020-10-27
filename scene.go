package engosdl

// IScene represents the interface for any game scene
type IScene interface {
	IObject
	AddGameObject(IGameObject) bool
	GetGameObject(name string) IGameObject
	OnAwake()
	OnEnable()
	OnStart()
	OnUpdate()
	OnDraw()
}

// Scene is the default implementation for IScene interface.
type Scene struct {
	*Object
	gameObjects []IGameObject
}

var _ IScene = (*Scene)(nil)

// NewScene creates a new scene instance
func NewScene(name string) *Scene {
	Logger.Trace().Str("scene", name).Msg("new scene")
	return &Scene{
		Object:      NewObject(name),
		gameObjects: []IGameObject{},
	}
}

// AddGameObject adds a new game object to the scene.
func (scene *Scene) AddGameObject(gobj IGameObject) bool {
	Logger.Trace().Str("scene", scene.name).Str("gameobject", gobj.GetName()).Msg("add game object")
	scene.gameObjects = append(scene.gameObjects, gobj)
	gobj.SetScene(scene)
	return true
}

// GetGameObject returns a game object for the given name.
func (scene *Scene) GetGameObject(name string) IGameObject {
	for _, gobj := range scene.gameObjects {
		if gobj.GetName() == name {
			return gobj
		}
	}
	return nil
}

// OnAwake calls all game object OnAwake methods.
func (scene *Scene) OnAwake() {
	for _, gobj := range scene.gameObjects {
		gobj.OnAwake()
	}
}

// OnEnable calls all game object OnEnable methods.
func (scene *Scene) OnEnable() {
	for _, gobj := range scene.gameObjects {
		gobj.OnEnable()
	}
}

// OnStart calls all game objects OnStart methods.
func (scene *Scene) OnStart() {
	for _, gobj := range scene.gameObjects {
		gobj.OnStart()
	}
}

// OnUpdate calls all game objects OnUpdate methods.
func (scene *Scene) OnUpdate() {
	for _, gobj := range scene.gameObjects {
		gobj.OnUpdate()
	}
}

// OnDraw calls all game objects OnDraw methods.
func (scene *Scene) OnDraw() {
	for _, gobj := range scene.gameObjects {
		gobj.OnDraw()
	}
}
