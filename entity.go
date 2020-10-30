package engosdl

import (
	"fmt"
	"reflect"
)

// IEntity represents the interface for any entity. Any object in the
// game has to implement this interface.
type IEntity interface {
	IObject
	AddChild(IEntity) bool
	AddComponent(IComponent) IEntity
	DeleteChild(string) bool
	GetActive() bool
	GetChild(string) IEntity
	GetChildren() []IEntity
	GetComponent(interface{}) IComponent
	GetComponents() []IComponent
	GetLayer() int32
	GetParent() IEntity
	GetScene() IScene
	GetTag() string
	GetTransform() ITransform
	Load()
	OnAwake()
	OnCycleEnd()
	OnCycleStart()
	OnDraw()
	OnEnable()
	OnStart()
	OnUpdate()
	SetActive(bool) IEntity
	SetLayer(int32) IEntity
	SetParent(IEntity) IEntity
	SetScene(IScene) IEntity
	SetTag(string) IEntity
	Unload()
}

// Entity is the default implementation for IEntity.
type Entity struct {
	*Object
	active             bool
	layer              int32
	tag                string
	parent             IEntity
	children           []IEntity
	scene              IScene
	transform          ITransform
	components         []IComponent
	loadedComponents   []IComponent
	unloadedComponents []IComponent
	loaded             bool
}

var _ IEntity = (*Entity)(nil)

// NewEntity creates a new entity instance.
func NewEntity(name string) *Entity {
	Logger.Trace().Str("entity", name).Msg("new entity")
	return &Entity{
		Object:             NewObject(name),
		active:             true,
		layer:              0,
		tag:                "",
		parent:             nil,
		children:           []IEntity{},
		scene:              nil,
		transform:          NewTransform(),
		components:         []IComponent{},
		loadedComponents:   []IComponent{},
		unloadedComponents: []IComponent{},
		loaded:             false,
	}
}

// AddChild adds a new child to entity children.
func (entity *Entity) AddChild(child IEntity) bool {
	entity.children = append(entity.children, child)
	return true
}

// AddComponent adds a new component to the entity.
func (entity *Entity) AddComponent(component IComponent) IEntity {
	Logger.Trace().Str("entity", entity.name).
		Str("component", component.GetName()).
		Str("type", reflect.TypeOf(component).String()).
		Msg("add component")
	for _, comp := range entity.GetComponents() {
		if reflect.TypeOf(comp) == reflect.TypeOf(component) {
			err := fmt.Errorf("component type %s already exist", reflect.TypeOf(component))
			Logger.Error().Err(err)
			panic(err)
		}
	}
	entity.components = append(entity.components, component)
	entity.unloadedComponents = append(entity.unloadedComponents, component)
	return entity
}

// DeleteChild removes a child from entity children.
func (entity *Entity) DeleteChild(name string) bool {
	if child, i := entity.getChild(name); child != nil {
		entity.children = append(entity.children[:i], entity.children[i+1:]...)
		return true
	}
	return false
}

// GetActive returns if the entity is active (enable) or not (disable).
func (entity *Entity) GetActive() bool {
	return entity.active
}

// getChild returns child and index by child name from entity children.
func (entity *Entity) getChild(name string) (IEntity, int) {
	for i, child := range entity.GetChildren() {
		if child.GetName() == name {
			return child, i
		}
	}
	return nil, -1
}

// GetChild returns a child by name from entity children.
func (entity *Entity) GetChild(name string) IEntity {
	if child, _ := entity.getChild(name); child != nil {
		return child
	}
	return nil
}

// GetChildren returns entity children.
func (entity *Entity) GetChildren() []IEntity {
	return entity.children
}

// GetComponent returns the given component from the entity.
func (entity *Entity) GetComponent(k interface{}) IComponent {
	return nil
}

// GetComponents returns all entity components.
func (entity *Entity) GetComponents() []IComponent {
	return entity.components
}

// GetLayer returns the  layer where the entity has been placed.
func (entity *Entity) GetLayer() int32 {
	return entity.layer
}

// GetParent returns entity parent.
func (entity *Entity) GetParent() IEntity {
	return entity.parent
}

// GetScene returns the scene where the entity has been placed.
func (entity *Entity) GetScene() IScene {
	return entity.scene
}

// GetTag returns the entity tag.
func (entity *Entity) GetTag() string {
	return entity.tag
}

// GetTransform returns the entity transform.
func (entity *Entity) GetTransform() ITransform {
	return entity.transform
}

// loadUnloadedComponents proceeds to load any unloaded component.
func (entity *Entity) loadUnloadedComponents() {
	unloaded := []IComponent{}
	for _, component := range entity.unloadedComponents {
		if component.GetActive() {
			component.Load()
			entity.loadedComponents = append(entity.loadedComponents, component)
		} else {
			unloaded = append(unloaded, component)
		}
	}
	entity.unloadedComponents = unloaded
}

// Load is called when object is loaded by the scene.
func (entity *Entity) Load() {
	entity.loaded = true
	entity.OnStart()
	entity.loadUnloadedComponents()
}

// OnAwake calls all component OnAwake methods.
func (entity *Entity) OnAwake() {
	for _, component := range entity.GetComponents() {
		component.OnAwake()
	}
}

// OnCycleEnd calls all methods to run at the end of a tick cycle.
func (entity *Entity) OnCycleEnd() {
}

// OnCycleStart calls all methods to run at the start of a tick cycle.
func (entity *Entity) OnCycleStart() {
	entity.loadUnloadedComponents()
	for _, component := range entity.loadedComponents {
		component.OnCycleStart()
	}
}

// OnDraw calls all component OnDraw methods.
func (entity *Entity) OnDraw() {
	for _, component := range entity.loadedComponents {
		component.OnDraw()
	}
}

// OnEnable calls all component OnEnable methods.
func (entity *Entity) OnEnable() {
	for _, component := range entity.GetComponents() {
		component.OnEnable()
	}
}

// OnStart calls all component OnStart methods.
func (entity *Entity) OnStart() {
	for _, component := range entity.GetComponents() {
		component.OnStart()
	}
}

// OnUpdate calls all component OnUpdate methods.
func (entity *Entity) OnUpdate() {
	for _, component := range entity.loadedComponents {
		component.OnUpdate()
	}
}

// SetActive sets if the entity is active (enable) or not (disable).
func (entity *Entity) SetActive(active bool) IEntity {
	entity.active = active
	return entity
}

// SetLayer sets the entity layer where it will be placed.
func (entity *Entity) SetLayer(layer int32) IEntity {
	entity.layer = layer
	return entity
}

// SetParent sets entity parent.
func (entity *Entity) SetParent(parent IEntity) IEntity {
	entity.parent = parent
	return entity
}

// SetScene sets the scene where the entity will be placed.
func (entity *Entity) SetScene(scene IScene) IEntity {
	entity.scene = scene
	return entity
}

// SetTag sets the entity tag.
func (entity *Entity) SetTag(tag string) IEntity {
	entity.tag = tag
	return entity
}

// Unload is called when object is unloaded by the scene.
func (entity *Entity) Unload() {
	entity.loaded = false
	for _, component := range entity.loadedComponents {
		component.Unload()
	}
	entity.loadedComponents = []IComponent{}
	entity.unloadedComponents = []IComponent{}
}
