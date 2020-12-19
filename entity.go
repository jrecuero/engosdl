package engosdl

import (
	"fmt"
	"reflect"
)

// ComponentToMarshal identifies component information being marshaled to be
// saved in JSON format.
type ComponentToMarshal struct {
	ComponentName string     `json:"component-type"`
	Component     IComponent `json:"component-data"`
}

// EntityToMarshal identifies entity information being marshaled to be saved
// in JSON format.
type EntityToMarshal struct {
	Entity     IEntity               `json:"entity-data"`
	Components []*ComponentToMarshal `json:"components"`
}

// ComponentToUnmarshal identifies component information from a JSON file
// required to build a new component instance.
type ComponentToUnmarshal struct {
	ComponentName string                 `json:"component-type"`
	Component     map[string]interface{} `json:"component-data"`
}

// EntityToUnmarshal identifies entity information from a JSON file required
// to build a new entity instance.
type EntityToUnmarshal struct {
	Entity     interface{}             `json:"entity-data"`
	Components []*ComponentToUnmarshal `json:"components"`
}

// IEntity represents the interface for any entity. Any object in the
// game has to implement this interface.
type IEntity interface {
	IObject
	AddChild(IEntity) bool
	AddComponent(IComponent) IEntity
	AddComponentExt(IComponent, IEntity) IEntity
	DeleteChild(string) bool
	DeleteChildByName(string) bool
	DoDestroy()
	DoDump() *EntityToMarshal
	DoFrameEnd()
	DoFrameStart()
	DoLoad()
	DoUnLoad()
	GetActive() bool
	GetCache(string) (interface{}, error)
	GetChild(string) IEntity
	GetChildByName(string) IEntity
	GetChildren() []IEntity
	GetComponent(IComponent) IComponent
	GetComponents() []IComponent
	GetDelegateForComponent(IComponent) IDelegate
	GetDieOnCollision() bool
	GetDieOnOutOfBounds() bool
	GetLayer() int
	GetParent() IEntity
	GetRenderable() bool
	GetScene() IScene
	GetTag() string
	GetTransform() ITransform
	IsInside(*Vector) bool
	OnRender()
	OnEnable()
	OnStart()
	OnUpdate()
	RemoveComponent(IComponent) bool
	RemoveComponents() bool
	SetActive(bool) IEntity
	SetCache(string, interface{})
	SetCustomOnUpdate(func(IEntity))
	SetDieOnCollision(bool) IEntity
	SetDieOnOutOfBounds(bool) IEntity
	SetLayer(int) IEntity
	SetParent(IEntity) IEntity
	SetRenderable(bool)
	SetScene(IScene) IEntity
	SetTag(string) IEntity
	Unmarshal(*EntityToUnmarshal)
}

// Entity is the default implementation for IEntity.
type Entity struct {
	*Object
	active             bool
	Renderable         bool   `json:"renderable"`
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
	DieOnOutOfBounds   bool `json:"die-on-out-of-bounds"`
	customOnUpdate     func(IEntity)
	cache              map[string]interface{}
}

var _ IEntity = (*Entity)(nil)

// NewEntity creates a new entity instance.
func NewEntity(name string) *Entity {
	Logger.Trace().Str("entity", name).Msg("new entity")
	return &Entity{
		Object:             NewObject(name),
		active:             true,
		Renderable:         true,
		Layer:              LayerMiddle,
		Tag:                "",
		parent:             nil,
		children:           []IEntity{},
		scene:              nil,
		Transform:          NewTransform(),
		components:         []IComponent{},
		loadedComponents:   []IComponent{},
		unloadedComponents: []IComponent{},
		DieOnCollision:     false,
		DieOnOutOfBounds:   false,
		customOnUpdate:     nil,
		cache:              make(map[string]interface{}),
	}
}

// AddChild adds a new child to entity children.
func (entity *Entity) AddChild(child IEntity) bool {
	entity.children = append(entity.children, child)
	child.SetParent(entity)
	// Child entity inherits layer from parent.
	child.SetLayer(entity.GetLayer())
	return true
}

// AddComponent adds a new component to the entity.
func (entity *Entity) AddComponent(component IComponent) IEntity {
	// Logger.Trace().Str("entity", entity.GetName()).
	// 	Str("component", component.GetName()).
	// 	Str("type", reflect.TypeOf(component).String()).
	// 	Msg("add component")
	// for _, comp := range entity.GetComponents() {
	// 	if reflect.TypeOf(comp) == reflect.TypeOf(component) {
	// 		err := fmt.Errorf("component type %s already exist", reflect.TypeOf(component))
	// 		Logger.Error().Err(err).Msg("AddComponent error")
	// 		panic(err)
	// 	}
	// }
	// component.SetEntity(entity)
	// // component.OnAwake()
	// entity.components = append(entity.components, component)
	// entity.unloadedComponents = append(entity.unloadedComponents, component)
	// return entity
	return entity.AddComponentExt(component, entity)
}

// AddComponentExt adds a new component to the entity, but it provides the
// entity as an additional parameter. Required for custom entities with
// methods not defined in IEntity interface.
func (entity *Entity) AddComponentExt(component IComponent, extEntity IEntity) IEntity {
	Logger.Trace().Str("entity", extEntity.GetName()).
		Str("component", component.GetName()).
		Str("type", reflect.TypeOf(component).String()).
		Msg("add component")
	for _, comp := range extEntity.GetComponents() {
		if reflect.TypeOf(comp) == reflect.TypeOf(component) {
			err := fmt.Errorf("component type %s already exist", reflect.TypeOf(component))
			Logger.Error().Err(err).Msg("AddComponent error")
			panic(err)
		}
	}
	component.SetEntity(extEntity)
	entity.components = append(entity.components, component)
	entity.unloadedComponents = append(entity.unloadedComponents, component)
	return extEntity
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
	for _, component := range entity.GetComponents() {
		component.DoDestroy()
	}
	entity.components = []IComponent{}
	entity.loadedComponents = []IComponent{}
	entity.unloadedComponents = []IComponent{}

	// for _, component := range entity.GetComponents() {
	// 	if !component.GetRemoveOnDestroy() {
	// 		entity.unloadedComponents = append(entity.unloadedComponents, component)
	// 	}
	// }
}

// DoDump dumps entity in JSON format.
func (entity *Entity) DoDump() *EntityToMarshal {
	toDump := &EntityToMarshal{
		Entity:     entity,
		Components: []*ComponentToMarshal{},
	}
	// if result, err := json.Marshal(entity); err == nil {
	// fmt.Printf("%s\n", result)
	for _, component := range entity.GetComponents() {
		// fmt.Printf("%d %s\n", i, reflect.TypeOf(component))
		toDump.Components = append(toDump.Components, &ComponentToMarshal{
			ComponentName: reflect.TypeOf(component).String(),
			Component:     component,
		})
	}
	// }
	// result, err := json.MarshalIndent(toDump, "", "    ")
	// if err != nil {
	// 	Logger.Error().Err(err)
	// 	panic(err)
	// }
	// fmt.Printf("%s\n", result)
	// if err := ioutil.WriteFile("entities.json", result, 0644); err != nil {
	// 	Logger.Error().Err(err)
	// 	panic(err)
	// }
	return toDump
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

// GetCache retrieves entity cache data for the given key.
func (entity *Entity) GetCache(key string) (interface{}, error) {
	if data, ok := entity.cache[key]; ok {
		return data, nil
	}
	return nil, fmt.Errorf("no cache data found for key: %s", key)
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

// GetDieOnOutOfBounds returns if the entity should be destroyed when it
// is out of bounds.
func (entity *Entity) GetDieOnOutOfBounds() bool {
	return entity.DieOnOutOfBounds
}

// GetLayer returns the  layer where the entity has been placed.
func (entity *Entity) GetLayer() int {
	return entity.Layer
}

// GetParent returns entity parent.
func (entity *Entity) GetParent() IEntity {
	return entity.parent
}

// GetRenderable return entity renderable attribute. Entity is not being
// rendered if this is false.
func (entity *Entity) GetRenderable() bool {
	return entity.Renderable
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

// IsInside returns if the given position is inside the entity rectangle.
func (entity *Entity) IsInside(pos *Vector) bool {
	rect := entity.GetTransform().GetRect()
	return pos.InRect(rect)
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
	if entity.GetRenderable() {
		for _, component := range entity.loadedComponents {
			if component.GetActive() {
				component.OnRender()
			}
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
	if entity.customOnUpdate != nil {
		entity.customOnUpdate(entity)
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

// SetCache sets a new cache entry for the given key and data.
func (entity *Entity) SetCache(key string, data interface{}) {
	entity.cache[key] = data
}

// SetCustomOnUpdate sets a custom method to be called OnUpdate.
func (entity *Entity) SetCustomOnUpdate(customCall func(IEntity)) {
	entity.customOnUpdate = customCall
}

// SetDieOnCollision sets if the entity should be destroyed in any collision.
func (entity *Entity) SetDieOnCollision(die bool) IEntity {
	entity.DieOnCollision = die
	return entity
}

// SetDieOnOutOfBounds sets if the entity should be destroyed if it is out
// of bounds.
func (entity *Entity) SetDieOnOutOfBounds(die bool) IEntity {
	entity.DieOnOutOfBounds = die
	return entity
}

// SetEnabled sets the entity to be enabled.
func (entity *Entity) SetEnabled(enabled bool) {
	for _, component := range entity.components {
		component.SetEnabled(enabled)
	}
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

// SetRenderable sets entity renderable attribute. Entity is not being
// rendered if this attribute is false.
func (entity *Entity) SetRenderable(renderable bool) {
	entity.Renderable = renderable
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

// Unmarshal takes a EntityToMarshal instance and  creates a new entity
// instance.
func (entity *Entity) Unmarshal(instance *EntityToUnmarshal) {
	obj := instance.Entity.(map[string]interface{})
	entity.SetName(obj["name"].(string))
	entity.SetTag(obj["tag"].(string))
	entity.SetLayer(int(obj["layer"].(float64)))
	entity.SetRenderable(obj["renderable"].(bool))
	entity.SetDieOnCollision(obj["die-on-collision"].(bool))
	entity.SetDieOnOutOfBounds(obj["die-on-out-of-bounds"].(bool))
	position := obj["transform"].(map[string]interface{})["position"].(map[string]interface{})
	scale := obj["transform"].(map[string]interface{})["scale"].(map[string]interface{})
	dimension := obj["transform"].(map[string]interface{})["dimension"].(map[string]interface{})
	rotation := obj["transform"].(map[string]interface{})["rotation"]
	entity.GetTransform().SetPosition(NewVector(position["X"].(float64), position["Y"].(float64)))
	entity.GetTransform().SetScale(NewVector(scale["X"].(float64), scale["Y"].(float64)))
	entity.GetTransform().SetDim(NewVector(dimension["X"].(float64), dimension["Y"].(float64)))
	entity.GetTransform().SetRotation(rotation.(float64))
	for _, comp := range instance.Components {
		constructor := GetComponentManager().Constructors[comp.ComponentName]
		component := constructor()
		component.Unmarshal(comp.Component)
		entity.AddComponent(component)
	}
}
