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
	DefaultAddDelegateToRegister()
	DefaultOnCollision(...interface{}) bool
	DefaultOnDestroy(...interface{}) bool
	DefaultOnLoad(...interface{}) bool
	DefaultOnOutOfBounds(...interface{}) bool
	DoDestroy()
	DoFrameEnd()
	DoFrameStart() bool
	DoLoad(IComponent) bool
	DoUnLoad()
	GetActive() bool
	GetDelegate() IDelegate
	GetEntity() IEntity
	OnAwake()
	OnEnable()
	OnRender()
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
	GetFileImageIndex() int
	GetFilename() []string
	LoadSprite()
	GetSpriteIndex() int
	NextFileImage() int
	NextSprite() int
	PreviousFileImage() int
	PreviousSprite() int
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
	started   bool
	delegate  IDelegate
	registers []IRegister
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
		started:   false,
		delegate:  nil,
		registers: []IRegister{},
	}
}

// AddDelegateToRegister adds a new delegate that component should register.
func (c *Component) AddDelegateToRegister(delegate IDelegate, entity IEntity, component IComponent, signature TDelegateSignature) IComponent {
	Logger.Trace().Str("component", c.GetName()).Msg("AddDelegateToRegister")
	register := NewRegister("new-register", c, entity, component, delegate, signature)
	c.registers = append(c.registers, register)
	return c
}

// DefaultAddDelegateToRegister will proceed to add default delegate to
// register for the component.
func (c *Component) DefaultAddDelegateToRegister() {
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

// DoDestroy calls all methods to clean up component.
func (c *Component) DoDestroy() {
	Logger.Trace().Str("component", c.GetName()).Msg("DoDestroy")
	// Deregister all register entries from delegate handler
	for _, register := range c.registers {
		GetDelegateHandler().DeregisterFromDelegate(register.GetRegisterID())
		register.SetDelegate(nil)
	}
	// Delete delegate being created.
	if c.GetDelegate() != nil {
		GetDelegateHandler().DeleteDelegate(c.GetDelegate())
	}
	c.registers = []IRegister{}
	c.delegate = nil
	c.loaded = false
	c.started = false
}

// DoFrameEnd calls all methods to run at the end of a tick frame.
func (c *Component) DoFrameEnd() {
}

// DoFrameStart calls all methods to run at the start of a tick frame.
func (c *Component) DoFrameStart() bool {
	return c.started
}

// DoLoad is called when component is loaded by the entity.
func (c *Component) DoLoad(component IComponent) bool {
	Logger.Trace().Str("component", c.GetName()).Msg("DoLoad")
	// fmt.Printf("load: %#v\n", reflect.TypeOf(component).String())
	// component.OnStart()
	return c.loaded
}

// DoUnLoad is called when component is unloaded by the entity. If this method
// is overwritten in any child class, Component  method has to be called,
// because it is in charge to deregister all delegates and registers.
func (c *Component) DoUnLoad() {
	Logger.Trace().Str("component", c.GetName()).Msg("DoUnLoad")
	// Deregister all register entries from delegate handler
	for _, register := range c.registers {
		GetDelegateHandler().DeregisterFromDelegate(register.GetRegisterID())
		// Register is created based on delegates for those belonging to the
		// DelegateHandler. Those are unchanged and delegate does not have to
		// be cleared.
		if register.GetDelegate() != nil &&
			register.GetDelegate().GetID() != GetDelegateHandler().GetCollisionDelegate().GetID() &&
			register.GetDelegate().GetID() != GetDelegateHandler().GetDestroyDelegate().GetID() &&
			register.GetDelegate().GetID() != GetDelegateHandler().GetLoadDelegate().GetID() {
			register.SetDelegate(nil)

		}
	}
	// Delete delegate being created.
	if c.GetDelegate() != nil {
		GetDelegateHandler().DeleteDelegate(c.GetDelegate())
	}
	c.loaded = false
	c.started = false
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
	c.loaded = true
}

// OnEnable is called every time the component is enabled.
func (c *Component) OnEnable() {
	Logger.Trace().Str("component", c.name).Msg("OnEnable")
}

// OnRender is called for every render tick.
func (c *Component) OnRender() {
	// Logger.Trace().Str("component", c.name).Msg("OnRender")
}

// OnStart is called first time the component is enabled.
func (c *Component) OnStart() {
	Logger.Trace().Str("component", c.name).Msg("OnStart")
	if !c.started {
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
				if registerID, ok := GetDelegateHandler().RegisterToDelegate(c, delegate, register.GetSignature()); ok {
					register.SetRegisterID(registerID)
					continue
				}
			}
			Logger.Error().Err(fmt.Errorf("register for component %s failed", c.GetName())).Str("component", c.GetName()).Msg("registration error")
			// panic("Failure at register " + register.GetName())
		}
		c.started = true
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
