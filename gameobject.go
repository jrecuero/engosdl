package engosdl

import (
	"fmt"
	"reflect"
)

// IGameObject represents the interface for any game object. Any object in the
// game has to implement this interface.
type IGameObject interface {
	IObject
	AddChild(IGameObject) bool
	AddComponent(IComponent) IGameObject
	DeleteChild(string) bool
	GetActive() bool
	GetChild(string) IGameObject
	GetChildren() []IGameObject
	GetComponent(interface{}) IComponent
	GetComponents() []IComponent
	GetLayer() int32
	GetParent() IGameObject
	GetScene() IScene
	GetTag() string
	GetTransform() ITransform
	OnAwake()
	OnDraw()
	OnEnable()
	OnStart()
	OnUpdate()
	SetActive(bool) IGameObject
	SetLayer(int32) IGameObject
	SetParent(IGameObject) IGameObject
	SetScene(IScene) IGameObject
	SetTag(string) IGameObject
}

// GameObject is the default implementation for IGameObject.
type GameObject struct {
	*Object
	active     bool
	layer      int32
	tag        string
	parent     IGameObject
	children   []IGameObject
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

// GetParent returns game object parent.
func (gobj *GameObject) GetParent() IGameObject {
	return gobj.parent
}

// SetParent sets game object parent.
func (gobj *GameObject) SetParent(parent IGameObject) IGameObject {
	gobj.parent = parent
	return gobj
}

// GetChildren returns game object children.
func (gobj *GameObject) GetChildren() []IGameObject {
	return gobj.children
}

// getChild returns child and index by child name from game object children.
func (gobj *GameObject) getChild(name string) (IGameObject, int) {
	for i, child := range gobj.GetChildren() {
		if child.GetName() == name {
			return child, i
		}
	}
	return nil, -1
}

// GetChild returns a child by name from game object children.
func (gobj *GameObject) GetChild(name string) IGameObject {
	if child, _ := gobj.getChild(name); child != nil {
		return child
	}
	return nil
}

// AddChild adds a new child to game object children.
func (gobj *GameObject) AddChild(child IGameObject) bool {
	gobj.children = append(gobj.children, child)
	return true
}

// DeleteChild removes a child from game object children.
func (gobj *GameObject) DeleteChild(name string) bool {
	if child, i := gobj.getChild(name); child != nil {
		gobj.children = append(gobj.children[:i], gobj.children[i+1:]...)
		return true
	}
	return false
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
		parent:     nil,
		children:   []IGameObject{},
		scene:      nil,
		transform:  NewTransform(),
		components: []IComponent{},
	}
}
