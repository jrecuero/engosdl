package engosdl

import (
	"encoding/json"
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
	DeleteChildByName(string) bool
	DoDestroy()
	DoDump()
	DoFrameEnd()
	DoFrameStart()
	DoLoad()
	DoUnLoad()
	GetActive() bool
	GetChild(string) IEntity
	GetChildByName(string) IEntity
	GetChildren() []IEntity
	GetComponent(IComponent) IComponent
	GetComponents() []IComponent
	GetDelegateForComponent(IComponent) IDelegate
	GetDieOnCollision() bool
	GetLayer() int
	GetParent() IEntity
	GetScene() IScene
	GetTag() string
	GetTransform() ITransform
	OnRender()
	OnEnable()
	OnStart()
	OnUpdate()
	RemoveComponent(IComponent) bool
	RemoveComponents() bool
	SetActive(bool) IEntity
	SetDieOnCollision(bool) IEntity
	SetLayer(int) IEntity
	SetParent(IEntity) IEntity
	SetScene(IScene) IEntity
	SetTag(string) IEntity
}

// Entity is the default implementation for IEntity.
type Entity struct {
	*Object
	active             bool
	Layer              int    `json:"layer"`
	Tag                string `json:"tag"`
	parent             IEntity
	children           []IEntity
	scene              IScene
	Transform          ITransform `json:"transform"`
	components         []IComponent
	loadedComponents   []IComponent
	unloadedComponents []IComponent
	DieOnCollision     bool `json:"die-on-collision"`
}

var _ IEntity = (*Entity)(nil)

// NewEntity creates a new entity instance.
func NewEntity(name string) *Entity {
	Logger.Trace().Str("entity", name).Msg("new entity")
	return &Entity{
		Object:             NewObject(name),
		active:             true,
		Layer:              LayerMiddle,
		Tag:                "",
		parent:             nil,
		children:           []IEntity{},
		scene:              nil,
		Transform:          NewTransform(),
		components:         []IComponent{},
		loadedComponents:   []IComponent{},
		unloadedComponents: []IComponent{},
	}
}

// AddChild adds a new child to entity children.
func (entity *Entity) AddChild(child IEntity) bool {
	entity.children = append(entity.children, child)
	return true
}

// AddComponent adds a new component to the entity.
func (entity *Entity) AddComponent(component IComponent) IEntity {
	Logger.Trace().Str("entity", entity.GetName()).
		Str("component", component.GetName()).
		Str("type", reflect.TypeOf(component).String()).
		Msg("add component")
	for _, comp := range entity.GetComponents() {
		if reflect.TypeOf(comp) == reflect.TypeOf(component) {
			err := fmt.Errorf("component type %s already exist", reflect.TypeOf(component))
			Logger.Error().Err(err).Msg("AddComponent error")
			panic(err)
		}
	}
	component.SetEntity(entity)
	// component.OnAwake()
	entity.components = append(entity.components, component)
	entity.unloadedComponents = append(entity.unloadedComponents, component)
	return entity
}

// DeleteChild removes a child from entity children using child ID.
func (entity *Entity) DeleteChild(id string) bool {
	if child, i := entity.getChild(id); child != nil {
		entity.children = append(entity.children[:i], entity.children[i+1:]...)
		return true
	}
	return false
}

// DeleteChildByName removes a child from entity children using child name
func (entity *Entity) DeleteChildByName(name string) bool {
	if child, i := entity.getChildByName(name); child != nil {
		entity.children = append(entity.children[:i], entity.children[i+1:]...)
		return true
	}
	return false
}

// DoDestroy calls all methods to clean up entity.
func (entity *Entity) DoDestroy() {
	Logger.Trace().Str("entity", entity.GetName()).Msg("DoDestroy")
	entity.SetLoaded(false)
	for _, component := range entity.loadedComponents {
		component.DoDestroy()
	}
	entity.unloadedComponents = []IComponent{}
	for _, component := range entity.GetComponents() {
		if !component.GetRemoveOnDestroy() {
			entity.unloadedComponents = append(entity.unloadedComponents, component)
		}
	}
	entity.components = entity.unloadedComponents
	entity.loadedComponents = []IComponent{}
}

// DoDump dumps entity in JSON format.
func (entity *Entity) DoDump() {
	if result, err := json.Marshal(entity); err == nil {
		fmt.Printf("%s\n", result)
		for i, component := range entity.GetComponents() {
			fmt.Printf("%d %s\n", i, reflect.TypeOf(component))
			component.DoDump(component)
		}
	}
}

// DoFrameEnd calls all methods to run at the end of a tick frame.
func (entity *Entity) DoFrameEnd() {
}

// DoFrameStart calls all methods to run at the start of a tick frame.
func (entity *Entity) DoFrameStart() {
	entity.loadUnloadedComponents()
	for _, component := range entity.loadedComponents {
		if started := component.DoFrameStart(); !started {
			component.OnStart()
		}
	}
}

// DoLoad is called when object is loaded by the scene.
func (entity *Entity) DoLoad() {
	Logger.Trace().Str("entity", entity.GetName()).Msg("DoLoad")
	entity.SetLoaded(true)
	// entity.OnStart()
	entity.loadUnloadedComponents()
}

// DoUnLoad is called when object is unloaded by the scene.
func (entity *Entity) DoUnLoad() {
	Logger.Trace().Str("entity", entity.GetName()).Msg("DoUnLoad")
	entity.SetLoaded(false)
	for _, component := range entity.loadedComponents {
		component.DoUnLoad()
	}
	entity.loadedComponents = []IComponent{}
	entity.unloadedComponents = []IComponent{}
	for _, component := range entity.GetComponents() {
		entity.unloadedComponents = append(entity.unloadedComponents, component)
	}
}

// GetActive returns if the entity is active (enable) or not (disable).
func (entity *Entity) GetActive() bool {
	return entity.active
}

// getChild returns child and index by child id from entity children.
func (entity *Entity) getChild(id string) (IEntity, int) {
	for i, child := range entity.GetChildren() {
		if child.GetID() == id {
			return child, i
		}
	}
	return nil, -1
}

// getChildByName returns child and index by child name from entity children.
func (entity *Entity) getChildByName(name string) (IEntity, int) {
	for i, child := range entity.GetChildren() {
		if child.GetName() == name {
			return child, i
		}
	}
	return nil, -1
}

// GetChild returns a child by id from entity children.
func (entity *Entity) GetChild(id string) IEntity {
	if child, _ := entity.getChild(id); child != nil {
		return child
	}
	return nil
}

// GetChildByName returns a child by name from entity children.
func (entity *Entity) GetChildByName(name string) IEntity {
	if child, _ := entity.getChildByName(name); child != nil {
		return child
	}
	return nil
}

// GetChildren returns entity children.
func (entity *Entity) GetChildren() []IEntity {
	return entity.children
}

// GetComponent returns the given component from the entity.
func (entity *Entity) GetComponent(typ IComponent) IComponent {
	for _, component := range entity.GetComponents() {
		if reflect.TypeOf(component) == reflect.TypeOf(typ) {
			return component
		}
	}
	return nil
}

// GetComponents returns all entity components.
func (entity *Entity) GetComponents() []IComponent {
	return entity.components
}

// GetDelegateForComponent returns the delegate for the given component.
func (entity *Entity) GetDelegateForComponent(typ IComponent) IDelegate {
	if component := entity.GetComponent(typ); component != nil {
		if delegate := component.GetDelegate(); delegate != nil {
			return delegate
		}
	}
	return nil
}

// GetDieOnCollision returns if the entity should be destroyed with any
// collision.
func (entity *Entity) GetDieOnCollision() bool {
	return entity.DieOnCollision
}

// GetLayer returns the  layer where the entity has been placed.
func (entity *Entity) GetLayer() int {
	return entity.Layer
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
	return entity.Tag
}

// GetTransform returns the entity transform.
func (entity *Entity) GetTransform() ITransform {
	return entity.Transform
}

// loadUnloadedComponents proceeds to load any unloaded component.
func (entity *Entity) loadUnloadedComponents() {
	unloaded := []IComponent{}
	for _, component := range entity.unloadedComponents {
		if component.GetActive() {
			if loaded := component.DoLoad(component); !loaded {
				component.OnAwake()
			}
			entity.loadedComponents = append(entity.loadedComponents, component)
		} else {
			unloaded = append(unloaded, component)
		}
	}
	entity.unloadedComponents = unloaded
}

// OnRender calls all component OnRender methods.
func (entity *Entity) OnRender() {
	for _, component := range entity.loadedComponents {
		if component.GetActive() {
			component.OnRender()
		}
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
	Logger.Trace().Str("entity", entity.GetName()).Msg("OnStart")
	for _, component := range entity.GetComponents() {
		if component.GetActive() {
			component.OnStart()
		}
	}
}

// OnUpdate calls all component OnUpdate methods.
func (entity *Entity) OnUpdate() {
	for _, component := range entity.loadedComponents {
		if component.GetActive() {
			component.OnUpdate()
		}
	}
}

// RemoveComponent removes the given component.
func (entity *Entity) RemoveComponent(component IComponent) bool {
	Logger.Trace().Str("entity", entity.GetName()).
		Str("component", component.GetName()).
		Str("type", reflect.TypeOf(component).String()).
		Msg("remove component")
	for i, comp := range entity.GetComponents() {
		if reflect.TypeOf(comp) == reflect.TypeOf(component) {
			comp.DoUnLoad()
			entity.components = append(entity.components[:i], entity.components[i+1:]...)
			return true
		}
	}
	return false
}

// RemoveComponents removes all components.
func (entity *Entity) RemoveComponents() bool {
	Logger.Trace().Str("entity", entity.GetName()).
		Msg("remove components")
	for _, comp := range entity.GetComponents() {
		comp.DoUnLoad()
	}
	entity.components = []IComponent{}
	return true
}

// SetActive sets if the entity is active (enable) or not (disable).
func (entity *Entity) SetActive(active bool) IEntity {
	entity.active = active
	return entity
}

// SetDieOnCollision sets if the entity should be destroyed in any collision.
func (entity *Entity) SetDieOnCollision(die bool) IEntity {
	entity.DieOnCollision = true
	return entity
}

// SetLayer sets the entity layer where it will be placed.
func (entity *Entity) SetLayer(layer int) IEntity {
	entity.Layer = layer
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
	entity.Tag = tag
	return entity
}
