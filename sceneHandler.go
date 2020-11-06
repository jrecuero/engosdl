package engosdl

// ISceneHandler represents the interface for the scene handler.
type ISceneHandler interface {
	IObject
	AddScene(IScene) bool
	DeleteScene(string) bool
	DoCycleEnd()
	DoCycleStart()
	GetActiveScene() IScene
	GetScene(string) IScene
	GetSceneByName(string) IScene
	GetScenes() []IScene
	SetActiveFirstScene() IScene
	SetActiveLastScene() IScene
	SetActiveNextScene() IScene
	SetActivePrevScene() IScene
	SetActiveScene(IScene) bool
	OnAfterUpdate()
	OnDraw()
	OnEnable()
	OnStart()
	OnUpdate()
}

// ActiveScene represents the active scene.
type ActiveScene struct {
	scene IScene
	index int
}

// SceneHandler is the default implementation for the scene handler interface.
type SceneHandler struct {
	*Object
	scenes      []IScene
	activeScene *ActiveScene
}

var _ ISceneHandler = (*SceneHandler)(nil)

// NewSceneHandler creates a new scene handler instance.
func NewSceneHandler(name string) *SceneHandler {
	Logger.Trace().Str("scene-handler", name).Msg("new scene handler")
	return &SceneHandler{
		Object:      NewObject(name),
		scenes:      []IScene{},
		activeScene: &ActiveScene{},
	}
}

// AddScene adds a new scene to the scene handler
func (h *SceneHandler) AddScene(scene IScene) bool {
	Logger.Trace().Str("scene-handler", h.GetName()).Msg("AddScene")
	if scene := h.GetScene(scene.GetID()); scene != nil {
		return false
	}
	h.scenes = append(h.scenes, scene)
	return true
}

// DeleteScene deletes the scene given by the name.
func (h *SceneHandler) DeleteScene(name string) bool {
	Logger.Trace().Str("scene-handler", h.GetName()).Msg("DeleteScene")
	if scene, i := h.getScene(name); scene != nil {
		h.scenes = append(h.scenes[:i], h.scenes[i+1:]...)
		return true
	}
	return false
}

// DoCycleEnd calls all methods to run at the end of a tick cycle.
func (h *SceneHandler) DoCycleEnd() {
	if activeScene := h.GetActiveScene(); activeScene != nil {
		activeScene.DoCycleEnd()
	}
}

// DoCycleStart calls all methods to run at the start of a tick cycle.
func (h *SceneHandler) DoCycleStart() {
	if activeScene := h.GetActiveScene(); activeScene != nil {
		activeScene.DoCycleStart()
	}
}

// GetActiveScene returns the scene handler active scene at that time.
func (h *SceneHandler) GetActiveScene() IScene {
	return h.activeScene.scene
}

// getScene returns the scene and index by the given scene name.
func (h *SceneHandler) getScene(id string) (IScene, int) {
	for i, scene := range h.GetScenes() {
		if scene.GetID() == id {
			return scene, i
		}
	}
	return nil, -1
}

// getSceneByName returns the scene and index by the given scene name.
func (h *SceneHandler) getSceneByName(name string) (IScene, int) {
	for i, scene := range h.GetScenes() {
		if scene.GetName() == name {
			return scene, i
		}
	}
	return nil, -1
}

// GetScene returns the scene for the given name.
func (h *SceneHandler) GetScene(id string) IScene {
	if scene, _ := h.getScene(id); scene != nil {
		return scene
	}
	return nil
}

// GetSceneByName returns the scene for the given name.
func (h *SceneHandler) GetSceneByName(name string) IScene {
	if scene, _ := h.getSceneByName(name); scene != nil {
		return scene
	}
	return nil
}

// GetScenes returns all scenes in the scene handler.
func (h *SceneHandler) GetScenes() []IScene {
	return h.scenes
}

// OnAfterUpdate calls all scene OnAfterUpdate, which should run after DoUpdate
// runs and before DoDraw.
func (h *SceneHandler) OnAfterUpdate() {
	if activeScene := h.GetActiveScene(); activeScene != nil {
		activeScene.OnAfterUpdate()
	}
}

// OnDraw calls all scene OnDraw methods.
func (h *SceneHandler) OnDraw() {
	if activeScene := h.GetActiveScene(); activeScene != nil {
		activeScene.OnDraw()
	}
}

// OnEnable calls all scene OnEnable methods.
func (h *SceneHandler) OnEnable() {
	if activeScene := h.GetActiveScene(); activeScene != nil {
		activeScene.OnEnable()
	}
}

// OnStart calls all scene OnStart methods.
func (h *SceneHandler) OnStart() {
	Logger.Trace().Str("scene-handler", h.GetName()).Msg("OnStart")
	if activeScene := h.GetActiveScene(); activeScene != nil {
		activeScene.OnStart()
	}
}

// OnUpdate calls all scene OnUpdate methods.
func (h *SceneHandler) OnUpdate() {
	if activeScene := h.GetActiveScene(); activeScene != nil {
		activeScene.OnUpdate()
	}
}

// SetActiveFirstScene sets the first scene as the active one.
func (h *SceneHandler) SetActiveFirstScene() IScene {
	if len(h.GetScenes()) > 0 {

		scene := h.GetScenes()[0]
		h.setActiveScene(scene, 0)
		return scene
	}
	return nil
}

// SetActiveLastScene sets the last scene as the active one.
func (h *SceneHandler) SetActiveLastScene() IScene {
	length := len(h.GetScenes())
	if length > 0 {
		scene := h.GetScenes()[length-1]
		h.setActiveScene(scene, length-1)
		return scene
	}
	return nil
}

// SetActiveNextScene sets the next scene as the active one.
func (h *SceneHandler) SetActiveNextScene() IScene {
	length := len(h.GetScenes())
	if length > 0 && h.activeScene.scene != nil && h.activeScene.index < length-1 {
		index := h.activeScene.index + 1
		scene := h.GetScenes()[index]
		h.setActiveScene(scene, index)
		return scene
	}
	return nil
}

// SetActivePrevScene set the previous scene as the active one.
func (h *SceneHandler) SetActivePrevScene() IScene {
	length := len(h.GetScenes())
	if length > 0 && h.activeScene.scene != nil && h.activeScene.index > 0 {
		index := h.activeScene.index - 1
		scene := h.GetScenes()[index]
		h.setActiveScene(scene, index)
		return scene
	}
	return nil
}

// setActiveScene set the given scene and index and active one. It proceeds
// to unload previous scene active and load new one.
func (h *SceneHandler) setActiveScene(scene IScene, index int) {
	if h.activeScene.scene != nil {
		h.activeScene.scene.DoUnLoad()
	}
	h.activeScene.scene = scene
	h.activeScene.index = index
	h.activeScene.scene.DoLoad()
}

// SetActiveScene sets the given scene as the active scene.
func (h *SceneHandler) SetActiveScene(scene IScene) bool {
	for i, scn := range h.GetScenes() {
		if scn == scene {
			h.setActiveScene(scene, i)
			return true
		}
	}
	return false
}
