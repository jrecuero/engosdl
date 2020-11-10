package engosdl

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

// IComponent represents the interface for any component to be added to any
// Entity
type IComponent interface {
	IObject
	AddDelegateToRegister(IDelegate, IEntity, IComponent, TDelegateSignature) IComponent
	DefaultOnCollision(...interface{}) bool
	DefaultOnDestroy(...interface{}) bool
	DefaultOnLoad(...interface{}) bool
	DefaultOnOutOfBounds(...interface{}) bool
	DoCycleEnd()
	DoCycleStart()
	DoLoad(IComponent)
	DoUnLoad()
	GetActive() bool
	GetDelegate() IDelegate
	GetEntity() IEntity
	OnAwake()
	OnDraw()
	OnEnable()
	OnStart()
	OnUpdate()
	SetActive(bool)
	SetDelegate(IDelegate)
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
	GetCamera() *sdl.Rect
	GetFilename() []string
	SetCamera(*sdl.Rect)
	SetDestroyOnOutOfBounds(bool)
}

// IText represents the interface for any text component.
type IText interface {
	SetFontFilename(string) IText
	SetColor(sdl.Color) IText
	SetMessage(string) IText
}

// Component represents the default IComponent implementation.
type Component struct {
	*Object
	entity    IEntity
	active    bool
	loaded    bool
	delegate  IDelegate
	registers []*Register
}

var _ IComponent = (*Component)(nil)

// NewComponent creates a new component instance.
func NewComponent(name string) *Component {
	Logger.Trace().Str("component", name).Msg("new component")
	return &Component{
		Object:    NewObject(name),
		entity:    nil,
		active:    true,
		loaded:    false,
		delegate:  nil,
		registers: []*Register{},
	}
}

// AddDelegateToRegister adds a new delegate that component should register.
func (c *Component) AddDelegateToRegister(delegate IDelegate, entity IEntity, component IComponent, signature TDelegateSignature) IComponent {
	Logger.Trace().Str("component", c.GetName()).Msg("AddDelegateToRegister")
	register := NewRegister("new-register", entity, component, delegate, signature)
	c.registers = append(c.registers, register)
	return c
}

// DefaultOnCollision is the component default callback when on collision
// delegate is triggered.
func (c *Component) DefaultOnCollision(...interface{}) bool {
	return true
}

// DefaultOnDestroy is the component default callback when on destroy delegate
// is triggered.
func (c *Component) DefaultOnDestroy(...interface{}) bool {
	return true
}

// DefaultOnLoad is the component default callback when on load delegate is
// triggered.
func (c *Component) DefaultOnLoad(...interface{}) bool {
	return true
}

// DefaultOnOutOfBounds is the component default callback when on load delegate
// is triggered.
func (c *Component) DefaultOnOutOfBounds(...interface{}) bool {
	return true
}

// DoCycleEnd calls all methods to run at the end of a tick cycle.
func (c *Component) DoCycleEnd() {
}

// DoCycleStart calls all methods to run at the start of a tick cycle.
func (c *Component) DoCycleStart() {
}

// DoLoad is called when component is loaded by the entity.
func (c *Component) DoLoad(component IComponent) {
	Logger.Trace().Str("component", c.GetName()).Msg("DoLoad")
	c.loaded = true
	// fmt.Printf("load: %#v\n", reflect.TypeOf(component).String())
	component.OnStart()
}

// DoUnLoad is called when component is unloaded by the entity. If this method
// is overwritten in any child class, Component  method has to be called,
// because it is in charge to deregister all delegates and registers.
func (c *Component) DoUnLoad() {
	Logger.Trace().Str("component", c.GetName()).Msg("DoUnLoad")
	// Deregister all register entries from delegate handler
	for _, register := range c.registers {
		GetEngine().GetEventHandler().GetDelegateHandler().DeregisterFromDelegate(register.GetRegisterID())
	}
	// Delete delegate being created.
	if c.GetDelegate() != nil {
		GetEngine().GetEventHandler().GetDelegateHandler().DeleteDelegate(c.GetDelegate())
	}
	c.loaded = false
}

// GetActive returns if component is active or not
func (c *Component) GetActive() bool {
	return c.active
}

// GetDelegate returns delegates created by the component.
func (c *Component) GetDelegate() IDelegate {
	return c.delegate
}

// GetEntity return the component entity parent.
func (c *Component) GetEntity() IEntity {
	return c.entity
}

// OnAwake should create all component resources that don't have any dependency
// with any other component or entity.
func (c *Component) OnAwake() {
	Logger.Trace().Str("component", c.name).Msg("OnAwake")
}

// OnDraw is called for every draw tick.
func (c *Component) OnDraw() {
	// Logger.Trace().Str("component", c.name).Msg("OnDraw")
}

// OnEnable is called every time the component is enabled.
func (c *Component) OnEnable() {
	Logger.Trace().Str("component", c.name).Msg("OnEnable")
}

// OnStart is called first time the component is enabled.
func (c *Component) OnStart() {
	Logger.Trace().Str("component", c.name).Msg("OnStart")
	for _, register := range c.registers {
		delegate := register.GetDelegate()
		// Retrieve delegate from entity and component provided.
		if delegate == nil {
			entity := register.GetEntity()
			if entity == nil {
				entity = c.GetEntity()
			}
			if component := entity.GetComponent(register.GetComponent()); component != nil {
				if delegate = component.GetDelegate(); delegate != nil {
					register.SetDelegate(delegate)
				}
			}
		}
		if delegate != nil {
			if registerID, ok := GetEngine().GetEventHandler().GetDelegateHandler().RegisterToDelegate(delegate, register.GetSignature()); ok {
				register.SetRegisterID(registerID)
				continue
			}
		}
		Logger.Error().Err(fmt.Errorf("register for component %s failed", c.GetName())).Str("component", c.GetName()).Msg("registration error")
		// panic("Failure at register " + register.GetName())
	}
}

// OnUpdate is called for every update tick.
func (c *Component) OnUpdate() {
	// Logger.Trace().Str("component", c.name).Msg("OnUpdate")
}

// SetActive sets component active attribute
func (c *Component) SetActive(active bool) {
	c.active = active
}

// SetDelegate sets component delegate.
func (c *Component) SetDelegate(delegate IDelegate) {
	c.delegate = delegate
}

// SetEntity sets component new entity instance.
func (c *Component) SetEntity(entity IEntity) {
	c.entity = entity
}
