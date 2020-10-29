package engosdl

// IScene represents the interface for any game scene
type IScene interface {
	IObject
	AddGameObject(IGameObject) bool
	DeleteGameObject(IGameObject) bool
	GetGameObject(name string) IGameObject
	GetGameObjects() []IGameObject
	Load()
	OnAfterUpdate()
	OnAwake()
	OnCycleEnd()
	OnCycleStart()
	OnDraw()
	OnEnable()
	OnStart()
	OnUpdate()
	Unload()
}

// Scene is the default implementation for IScene interface.
type Scene struct {
	*Object
	gameObjects         []IGameObject
	toDeleteGameObjects []IGameObject
	loadedGameObjects   []IGameObject
	unloadedGameObjects []IGameObject
	loaded              bool
}

var _ IScene = (*Scene)(nil)

// NewScene creates a new scene instance
func NewScene(name string) *Scene {
	Logger.Trace().Str("scene", name).Msg("new scene")
	return &Scene{
		Object:              NewObject(name),
		gameObjects:         []IGameObject{},
		toDeleteGameObjects: []IGameObject{},
		loadedGameObjects:   []IGameObject{},
		unloadedGameObjects: []IGameObject{},
		loaded:              false,
	}
}

// AddGameObject adds a new game object to the scene.
func (scene *Scene) AddGameObject(gobj IGameObject) bool {
	Logger.Trace().Str("scene", scene.GetName()).Str("gameobject", gobj.GetName()).Msg("add game object")
	scene.gameObjects = append(scene.gameObjects, gobj)
	scene.unloadedGameObjects = append(scene.unloadedGameObjects, gobj)
	gobj.SetScene(scene)
	return true
}

// DeleteGameObject deletes a game object from the scene.
func (scene *Scene) DeleteGameObject(gobj IGameObject) bool {
	Logger.Trace().Str("scene", scene.GetName()).Str("gameobject", gobj.GetName()).Msg("delete game object")
	for _, travObj := range scene.gameObjects {
		if travObj == gobj {
			// GameObject to be deleted in OnAfterUpdate method.
			// scene.gameObjects = append(scene.gameObjects[:i], scene.gameObjects[i+1:]...)
			scene.toDeleteGameObjects = append(scene.toDeleteGameObjects, gobj)
			return true
		}
	}
	return false
}

// getGameObject returns game object and index for the given name.
func (scene *Scene) getGameObject(name string) (IGameObject, int) {
	for i, gobj := range scene.gameObjects {
		if gobj.GetName() == name {
			return gobj, i
		}
	}
	return nil, -1
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

// GetGameObjects returns all game objects in the scene.
func (scene *Scene) GetGameObjects() []IGameObject {
	return scene.gameObjects
}

// getIndexInLoadedGameObject return the index for the given game object in
// loadedGameObject array.
func (scene Scene) getIndexInLoadedGameObject(gobj IGameObject) (int, bool) {
	for i, obj := range scene.loadedGameObjects {
		if obj.GetName() == gobj.GetName() {
			return i, true
		}
	}
	return -1, false
}

// getIndexInUnloadedGameObject return the index for the given game object in
// unloadedGameObject array.
func (scene Scene) getIndexInUnloadedGameObject(gobj IGameObject) (int, bool) {
	for i, obj := range scene.unloadedGameObjects {
		if obj.GetName() == gobj.GetName() {
			return i, true
		}
	}
	return -1, false
}

// Load is call when scene is loaded in the scene handler.
func (scene *Scene) Load() {
	scene.loaded = true
	scene.loadUnloadedGameObjects()
}

// loadUnloadedGameObjects proceeds to load any unloaded game object
func (scene *Scene) loadUnloadedGameObjects() {
	unloaded := []IGameObject{}
	for _, gobj := range scene.unloadedGameObjects {
		if gobj.GetActive() {
			gobj.Load()
			scene.loadedGameObjects = append(scene.loadedGameObjects, gobj)
		} else {
			unloaded = append(unloaded, gobj)
		}
	}
	scene.unloadedGameObjects = unloaded
}

// OnAfterUpdate calls executed after all DoUpdates have been executed and
// before OnDraw.
func (scene *Scene) OnAfterUpdate() {
	// Delete all GameObjects being marked to be deleted
	if len(scene.toDeleteGameObjects) != 0 {
		for _, gobj := range scene.toDeleteGameObjects {
			if _, i := scene.getGameObject(gobj.GetName()); i != -1 {
				scene.gameObjects = append(scene.gameObjects[:i], scene.gameObjects[i+1:]...)
			}
			if index, ok := scene.getIndexInLoadedGameObject(gobj); ok {
				scene.loadedGameObjects = append(scene.loadedGameObjects[:index], scene.loadedGameObjects[index+1:]...)
			}
			if index, ok := scene.getIndexInUnloadedGameObject(gobj); ok {
				scene.unloadedGameObjects = append(scene.unloadedGameObjects[:index], scene.unloadedGameObjects[index+1:]...)
			}
		}
		scene.toDeleteGameObjects = []IGameObject{}
	}
}

// OnAwake calls all game object OnAwake methods.
func (scene *Scene) OnAwake() {
	for _, gobj := range scene.gameObjects {
		gobj.OnAwake()
	}
}

// OnCycleEnd calls all methods to run at the end of a tick cycle.
func (scene *Scene) OnCycleEnd() {
}

// OnCycleStart calls all methods to run at the start of a tick cycle.
func (scene *Scene) OnCycleStart() {
	scene.loadUnloadedGameObjects()
	for _, gobj := range scene.loadedGameObjects {
		gobj.OnCycleStart()
	}
}

// OnDraw calls all game objects OnDraw methods.
func (scene *Scene) OnDraw() {
	for _, gobj := range scene.loadedGameObjects {
		gobj.OnDraw()
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
	for _, gobj := range scene.loadedGameObjects {
		gobj.OnUpdate()
	}
}

// Unload is called when scene is unloaded from the scene handler.
func (scene *Scene) Unload() {
	scene.loaded = false
	for _, gobj := range scene.loadedGameObjects {
		gobj.Unload()
	}
	scene.loadedGameObjects = []IGameObject{}
	scene.unloadedGameObjects = []IGameObject{}
}
