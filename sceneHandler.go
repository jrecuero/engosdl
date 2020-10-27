package engosdl

// ISceneHandler represents the interface for the scene handler.
type ISceneHandler interface {
	IObject
	AddScene(IScene) bool
	GetScene(string) IScene
	Scenes() []IScene
	DeleteScene(string) bool
	OnAwake()
	OnEnable()
	OnStart()
	OnUpdate()
	OnDraw()
}

// SceneHandler is the default implementation for the scene handler interface.
type SceneHandler struct {
	*Object
	scenes []IScene
}

var _ ISceneHandler = (*SceneHandler)(nil)

// NewSceneHandler creates a new scene handler instance.
func NewSceneHandler(name string) *SceneHandler {
	Logger.Trace().Str("scene-handler", name).Msg("new scene handler")
	return &SceneHandler{
		Object: NewObject(name),
		scenes: []IScene{},
	}
}

// Scenes returns all scenes in the scene handler.
func (h *SceneHandler) Scenes() []IScene {
	return h.scenes
}

// getScene returns the scene and index by the given scene name.
func (h *SceneHandler) getScene(name string) (IScene, int) {
	for i, scene := range h.Scenes() {
		if scene.GetName() == name {
			return scene, i
		}
	}
	return nil, -1
}

// GetScene returns the scene for the given name.
func (h *SceneHandler) GetScene(name string) IScene {
	if scene, _ := h.getScene(name); scene != nil {
		return scene
	}
	return nil
}

// AddScene adds a new scene to the scene handler
func (h *SceneHandler) AddScene(scene IScene) bool {
	if scene := h.GetScene(scene.GetName()); scene != nil {
		return false
	}
	h.scenes = append(h.scenes, scene)
	return true
}

// DeleteScene deletes the scene given by the name.
func (h *SceneHandler) DeleteScene(name string) bool {
	if scene, i := h.getScene(name); scene != nil {
		h.scenes = append(h.scenes[:i], h.scenes[i+1:]...)
		return true
	}
	return false
}

// OnAwake calls all scene OnAwake methods.
func (h *SceneHandler) OnAwake() {
	for _, scene := range h.Scenes() {
		scene.OnAwake()
	}
}

// OnEnable calls all scene OnEnable methods.
func (h *SceneHandler) OnEnable() {
	for _, scene := range h.Scenes() {
		scene.OnEnable()
	}
}

// OnStart calls all scene OnStart methods.
func (h *SceneHandler) OnStart() {
	for _, scene := range h.Scenes() {
		scene.OnStart()
	}
}

// OnUpdate calls all scene OnUpdate methods.
func (h *SceneHandler) OnUpdate() {
	for _, scene := range h.Scenes() {
		scene.OnUpdate()
	}
}

// OnDraw calls all scene OnDraw methods.
func (h *SceneHandler) OnDraw() {
	for _, scene := range h.Scenes() {
		scene.OnDraw()
	}
}
