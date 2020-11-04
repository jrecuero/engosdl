package engosdl

import "github.com/veandco/go-sdl2/sdl"

// IComponent represents the interface for any component to be added to any
// Entity
type IComponent interface {
	IObject
	DoCycleEnd()
	DoCycleStart()
	DoLoad(IComponent)
	DoUnLoad()
	GetActive() bool
	GetDelegates() []IDelegate
	GetEntity() IEntity
	OnCollision(IEntity)
	OnAwake()
	OnDraw()
	OnEnable()
	OnStart()
	OnUpdate()
	SetActive(bool)
	SetEntity(IEntity)
}

// ICollisionBox represets the interface for any collider collision box.
type ICollisionBox interface {
	GetRect() *sdl.Rect
	GetCenter() *Vector
	GetRadius() float64
}

// ICollider represents a special kind of component that implement collisions.
type ICollider interface {
	IComponent
	GetCollisionBox() ICollisionBox
}

// ISprite represents the interface for any sprite component.
type ISprite interface {
	IComponent
	GetFilename() string
	SetDestroyOnOutOfBounds(bool)
}

// Component represents the default IComponent implementation.
type Component struct {
	*Object
	entity IEntity
	active bool
	loaded bool
}

var _ IComponent = (*Component)(nil)

// NewComponent creates a new component instance.
func NewComponent(name string) *Component {
	Logger.Trace().Str("component", name).Msg("new component")
	return &Component{
		Object: NewObject(name),
		entity: nil,
		active: true,
		loaded: false,
	}
}

// DoCycleEnd calls all methods to run at the end of a tick cycle.
func (c *Component) DoCycleEnd() {
}

// DoCycleStart calls all methods to run at the start of a tick cycle.
func (c *Component) DoCycleStart() {
}

// DoLoad is called when component is loaded by the entity.
func (c *Component) DoLoad(component IComponent) {
	c.loaded = true
	// fmt.Printf("load: %#v\n", reflect.TypeOf(component).String())
	component.OnStart()
}

// DoUnLoad is called when component is unloaded by the entity.
func (c *Component) DoUnLoad() {
	c.loaded = false
}

// GetActive returns if component is active or not
func (c *Component) GetActive() bool {
	return c.active
}

// GetDelegates returns all delegates registered to the component.
func (c *Component) GetDelegates() []IDelegate {
	return nil
}

// GetEntity return the component entity parent.
func (c *Component) GetEntity() IEntity {
	return c.entity
}

// OnCollision is called when entity collides with other entity.
func (c *Component) OnCollision(entity IEntity) {
}

// OnAwake should create all component resources that don't have any dependency
// with any other component or entity.
func (c *Component) OnAwake() {
	Logger.Trace().Str("component", c.name).Msg("OnAwake")
}

// OnDraw is called for every draw tick.
func (c *Component) OnDraw() {
	Logger.Trace().Str("component", c.name).Msg("OnDraw")
}

// OnEnable is called every time the component is enabled.
func (c *Component) OnEnable() {
	Logger.Trace().Str("component", c.name).Msg("OnEnable")
}

// OnStart is called first time the component is enabled.
func (c *Component) OnStart() {
	Logger.Trace().Str("component", c.name).Msg("OnStart")
}

// OnUpdate is called for every update tick.
func (c *Component) OnUpdate() {
	// Logger.Trace().Str("component", c.name).Msg("OnUpdate")
}

// SetActive sets component active attribute
func (c *Component) SetActive(active bool) {
	c.active = active
}

// SetEntity sets component new entity instance.
func (c *Component) SetEntity(entity IEntity) {
	c.entity = entity
}
