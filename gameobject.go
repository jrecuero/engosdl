package engosdl

import (
	"fmt"
	"reflect"
)

// IGameObject represents the interface for any game object. Any object in the
// game has to implement this interface.
type IGameObject interface {
	IObject
	GetActive() bool
	SetActive(bool) IGameObject
	GetLayer() int32
	SetLayer(int32) IGameObject
	GetScene() IScene
	SetScene(IScene) IGameObject
	GetTag() string
	SetTag(string) IGameObject
	GetTransform() ITransform
	AddComponent(IComponent) IGameObject
	GetComponent(interface{}) IComponent
	GetComponents() []IComponent
	OnAwake()
	OnEnable()
	OnStart()
	OnUpdate()
	OnDraw()
}

// GameObject is the default implementation for IGameObject.
type GameObject struct {
	*Object
	active     bool
	layer      int32
	tag        string
	scene      IScene
	transform  ITransform
	components []IComponent
}

var _ IGameObject = (*GameObject)(nil)

// GetActive returns if the game object is active (enable) or not (disable).
func (gobj *GameObject) GetActive() bool {
	return gobj.active
}

// SetActive sets if the game object is active (enable) or not (disable).
func (gobj *GameObject) SetActive(active bool) IGameObject {
	gobj.active = active
	return gobj
}

// GetLayer returns the  layer where the game object has been placed.
func (gobj *GameObject) GetLayer() int32 {
	return gobj.layer
}

// SetLayer sets the game object layer where it will be placed.
func (gobj *GameObject) SetLayer(layer int32) IGameObject {
	gobj.layer = layer
	return gobj
}

// GetScene returns the scene where the game object has been placed.
func (gobj *GameObject) GetScene() IScene {
	return gobj.scene
}

// SetScene sets the scene where the game object will be placed.
func (gobj *GameObject) SetScene(scene IScene) IGameObject {
	gobj.scene = scene
	return gobj
}

// GetTag returns the game object tag.
func (gobj *GameObject) GetTag() string {
	return gobj.tag
}

// SetTag sets the game object tag.
func (gobj *GameObject) SetTag(tag string) IGameObject {
	gobj.tag = tag
	return gobj
}

// GetTransform returns the game object transform.
func (gobj *GameObject) GetTransform() ITransform {
	return gobj.transform
}

// AddComponent adds a new component to the game object.
func (gobj *GameObject) AddComponent(component IComponent) IGameObject {
	Logger.Trace().Str("gameobject", gobj.name).Str("component", component.GetName()).Msg("add component")
	gobj.components = append(gobj.components, component)
	return gobj
}

// GetComponent returns the given component from the game object.
func (gobj *GameObject) GetComponent(k interface{}) IComponent {
	return nil
}

// GetComponents returns all game object components.
func (gobj *GameObject) GetComponents() []IComponent {
	return gobj.components
}

// OnAwake calls all component OnAwake methods.
func (gobj *GameObject) OnAwake() {
	for _, component := range gobj.GetComponents() {
		component.OnAwake()
	}
}

// OnEnable calls all component OnEnable methods.
func (gobj *GameObject) OnEnable() {
	for _, component := range gobj.GetComponents() {
		component.OnEnable()
	}
}

// OnStart calls all component OnStart methods.
func (gobj *GameObject) OnStart() {
	for _, component := range gobj.GetComponents() {
		fmt.Println(reflect.TypeOf(component))
		component.OnStart()
	}
}

// OnUpdate calls all component OnUpdate methods.
func (gobj *GameObject) OnUpdate() {
	for _, component := range gobj.GetComponents() {
		component.OnUpdate()
	}
}

// OnDraw calls all component OnDraw methods.
func (gobj *GameObject) OnDraw() {
	for _, component := range gobj.GetComponents() {
		component.OnDraw()
	}
}

// NewGameObject creates a new game object instance.
func NewGameObject(name string) *GameObject {
	Logger.Trace().Str("gameobject", name).Msg("new game object")
	return &GameObject{
		Object:     NewObject(name),
		active:     true,
		layer:      0,
		tag:        "",
		scene:      nil,
		transform:  NewTransform(),
		components: []IComponent{},
	}
}
