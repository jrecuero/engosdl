package engosdl

import (
	"fmt"
	"strconv"
)

// SceneEventData is the event structure used by the scene manager.
type SceneEventData struct {
	*Object
	scene IScene
	index int
}

// NewSceneEvent creates a new scene event instance.
func NewSceneEvent(scene IScene, index int) *Event {
	Logger.Trace().Str("scene-event", scene.GetName()).Str("index", strconv.Itoa(index)).Msg("new scene-event")
	return &Event{
		Object: NewObject("scene-event"),
		data: &SceneEventData{
			Object: NewObject(scene.GetName()),
			scene:  scene,
			index:  index,
		},
	}
}

// ISceneManager represents the interface for the scene handler.
type ISceneManager interface {
	IObject
	AddScene(IScene) bool
	DeleteScene(string) bool
	DoFrameEnd()
	DoFrameStart()
	DoInit()
	GetActiveScene() IScene
	GetScene(string) IScene
	GetSceneByName(string) IScene
	GetScenes() []IScene
	GetStandbyScene() IScene
	OnAfterUpdate()
	OnRender()
	OnEnable()
	OnStart()
	OnUpdate()
	RestartScene() bool
	SetActiveFirstScene() IScene
	SetActiveLastScene() IScene
	SetActiveNextScene() IScene
	SetActivePrevScene() IScene
	SetActiveScene(IScene) bool
	SwapBack() bool
	SwapFromSceneTo(IScene) bool
}

// ActiveScene represents the active scene.
type ActiveScene struct {
	scene IScene
	index int
}

// SceneManager is the default implementation for the scene handler interface.
type SceneManager struct {
	*Object
	scenes       []IScene
	activeScene  *ActiveScene
	standByScene *ActiveScene
	eventPoolID  string
}

var _ ISceneManager = (*SceneManager)(nil)

// NewSceneManager creates a new scene handler instance.
func NewSceneManager(name string) *SceneManager {
	Logger.Trace().Str("scene-manager", name).Msg("new scene handler")
	return &SceneManager{
		Object:       NewObject(name),
		scenes:       []IScene{},
		activeScene:  &ActiveScene{},
		standByScene: &ActiveScene{},
	}
}

// AddScene adds a new scene to the scene handler
func (h *SceneManager) AddScene(scene IScene) bool {
	Logger.Trace().Str("scene-manager", h.GetName()).Msg("AddScene")
	if scene := h.GetScene(scene.GetID()); scene != nil {
		return false
	}
	h.scenes = append(h.scenes, scene)
	return true
}

// DeleteScene deletes the scene given by the name.
func (h *SceneManager) DeleteScene(name string) bool {
	Logger.Trace().Str("scene-manager", h.GetName()).Msg("DeleteScene")
	if scene, i := h.getScene(name); scene != nil {
		h.scenes = append(h.scenes[:i], h.scenes[i+1:]...)
		return true
	}
	return false
}

// DoFrameEnd calls all methods to run at the end of a tick frame.
func (h *SceneManager) DoFrameEnd() {
	if activeScene := h.GetActiveScene(); activeScene != nil {
		activeScene.DoFrameEnd()
	}
	// Read the event pool for any scene change.
	if pool := GetEventManager().GetPool(h.eventPoolID); pool != nil {
		if event, _ := pool.Pop(); event != nil {
			data := event.GetData().(*SceneEventData)
			scene := data.scene
			index := data.index
			h.setActiveScene(scene, index)
		}
	}
}

// DoFrameStart calls all methods to run at the start of a tick frame.
func (h *SceneManager) DoFrameStart() {
	if activeScene := h.GetActiveScene(); activeScene != nil {
		activeScene.DoFrameStart()
	}
}

// DoInit initializes all scene manager resources.
func (h *SceneManager) DoInit() {
	Logger.Trace().Str("scene-manager", h.GetName()).Msg("DoInit")
}

// GetActiveScene returns the scene handler active scene at that time.
func (h *SceneManager) GetActiveScene() IScene {
	return h.activeScene.scene
}

// getScene returns the scene and index by the given scene name.
func (h *SceneManager) getScene(id string) (IScene, int) {
	for i, scene := range h.GetScenes() {
		if scene.GetID() == id {
			return scene, i
		}
	}
	return nil, -1
}

// getSceneByName returns the scene and index by the given scene name.
func (h *SceneManager) getSceneByName(name string) (IScene, int) {
	for i, scene := range h.GetScenes() {
		if scene.GetName() == name {
			return scene, i
		}
	}
	return nil, -1
}

// GetScene returns the scene for the given name.
func (h *SceneManager) GetScene(id string) IScene {
	if scene, _ := h.getScene(id); scene != nil {
		return scene
	}
	return nil
}

// GetSceneByName returns the scene for the given name.
func (h *SceneManager) GetSceneByName(name string) IScene {
	if scene, _ := h.getSceneByName(name); scene != nil {
		return scene
	}
	return nil
}

// GetScenes returns all scenes in the scene handler.
func (h *SceneManager) GetScenes() []IScene {
	return h.scenes
}

// GetStandbyScene returns the standby scene in the scene manager.
func (h *SceneManager) GetStandbyScene() IScene {
	return h.standByScene.scene
}

// OnAfterUpdate calls all scene OnAfterUpdate, which should run after DoUpdate
// runs and before DoRender.
func (h *SceneManager) OnAfterUpdate() {
	if activeScene := h.GetActiveScene(); activeScene != nil {
		activeScene.OnAfterUpdate()
	}
}

// OnRender calls all scene OnRender methods.
func (h *SceneManager) OnRender() {
	if activeScene := h.GetActiveScene(); activeScene != nil {
		activeScene.OnRender()
	}
}

// OnEnable calls all scene OnEnable methods.
func (h *SceneManager) OnEnable() {
	if activeScene := h.GetActiveScene(); activeScene != nil {
		activeScene.OnEnable()
	}
}

// OnStart calls all scene OnStart methods.
func (h *SceneManager) OnStart() {
	Logger.Trace().Str("scene-manager", h.GetName()).Msg("OnStart")
	var err error
	// Create event pool in event manager.
	if h.eventPoolID, err = GetEventManager().CreatePool("scene-manager-pool"); err != nil {
		Logger.Error().Err(err)
		panic(err)
	}
}

// OnUpdate calls all scene OnUpdate methods.
func (h *SceneManager) OnUpdate() {
	if activeScene := h.GetActiveScene(); activeScene != nil {
		activeScene.OnUpdate()
	}
}

// RestartScene restart the active scene.
func (h *SceneManager) RestartScene() bool {
	if scene := h.GetActiveScene(); scene != nil {
		return h.SetActiveScene(scene)
	}
	return false
}

// SetActiveFirstScene sets the first scene as the active one.
func (h *SceneManager) SetActiveFirstScene() IScene {
	if len(h.GetScenes()) > 0 {

		scene := h.GetScenes()[0]
		h.setActiveScene(scene, 0)
		return scene
	}
	return nil
}

// SetActiveLastScene sets the last scene as the active one.
func (h *SceneManager) SetActiveLastScene() IScene {
	length := len(h.GetScenes())
	if length > 0 {
		scene := h.GetScenes()[length-1]
		h.setActiveScene(scene, length-1)
		return scene
	}
	return nil
}

// SetActiveNextScene sets the next scene as the active one.
func (h *SceneManager) SetActiveNextScene() IScene {
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
func (h *SceneManager) SetActivePrevScene() IScene {
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
func (h *SceneManager) setActiveScene(scene IScene, index int) {
	fmt.Println("Audit Before UnLoading")
	fmt.Println("----------------------")
	GetDelegateManager().AuditDelegates()
	GetDelegateManager().AuditRegisters()
	if h.activeScene.scene != nil {
		h.activeScene.scene.AuditEntities()
		// h.activeScene.scene.DoUnLoad()
		h.activeScene.scene.DoDestroy()
	}
	fmt.Println("Audit After UnLoading")
	fmt.Println("---------------------")
	GetDelegateManager().AuditDelegates()
	GetDelegateManager().AuditRegisters()
	h.activeScene.scene = scene
	h.activeScene.index = index
	h.activeScene.scene.GetSceneCode()(GetEngine(), h.activeScene.scene)
	h.activeScene.scene.DoLoad()
	h.activeScene.scene.OnStart()
	fmt.Println("Audit After Loading")
	fmt.Println("-------------------")
	GetDelegateManager().AuditDelegates()
	GetDelegateManager().AuditRegisters()
	h.activeScene.scene.AuditEntities()
	// h.activeScene.scene.DoDump()
}

// SetActiveScene sets the given scene as the active scene.
func (h *SceneManager) SetActiveScene(scene IScene) bool {
	for index, scn := range h.GetScenes() {
		if scn == scene {
			// h.setActiveScene(scene, index)
			if pool := GetEventManager().GetPool(h.eventPoolID); pool != nil {
				pool.Add(NewSceneEvent(scene, index))
			}
			return true
		}
	}
	return false
}

// SwapFromSceneTo swaps from the active scene to a new one. The former active
// scene is moved to standby.
func (h *SceneManager) SwapFromSceneTo(newScene IScene) bool {
	if scene, index := h.getScene(newScene.GetID()); scene != nil {
		fmt.Println("Audit Before Swap From")
		fmt.Println("----------------------")
		GetDelegateManager().AuditDelegates()
		GetDelegateManager().AuditRegisters()
		h.activeScene.scene.AuditEntities()

		h.activeScene.scene.DoSwapFrom()
		h.standByScene.scene = h.activeScene.scene
		h.standByScene.index = h.activeScene.index
		h.activeScene.scene = scene
		h.activeScene.index = index
		h.activeScene.scene.GetSceneCode()(GetEngine(), h.activeScene.scene)
		h.activeScene.scene.DoLoad()
		h.activeScene.scene.OnStart()

		fmt.Println("Audit After Swap From")
		fmt.Println("---------------------")
		GetDelegateManager().AuditDelegates()
		GetDelegateManager().AuditRegisters()
		h.activeScene.scene.AuditEntities()

		return true
	}
	return false
}

// SwapBack swaps back to the standby scene.
func (h *SceneManager) SwapBack() bool {
	if h.standByScene.scene == nil {
		return false
	}
	fmt.Println("Audit Before Swap Back")
	fmt.Println("----------------------")
	GetDelegateManager().AuditDelegates()
	GetDelegateManager().AuditRegisters()
	h.activeScene.scene.AuditEntities()

	h.activeScene.scene.DoUnLoad()
	h.activeScene.scene = h.standByScene.scene
	h.activeScene.index = h.standByScene.index
	h.activeScene.scene.DoSwapBack()
	h.activeScene.scene.OnStart()

	fmt.Println("Audit After Swap Back")
	fmt.Println("---------------------")
	GetDelegateManager().AuditDelegates()
	GetDelegateManager().AuditRegisters()
	h.activeScene.scene.AuditEntities()

	return true
}
