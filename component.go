package engosdl

import (
	"encoding/json"
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
	DoDump(IComponent)
	DoFrameEnd()
	DoFrameStart() bool
	DoLoad(IComponent) bool
	DoUnLoad()
	GetActive() bool
	GetDelegate() IDelegate
	GetEntity() IEntity
	GetRemoveOnDestroy() bool
	OnAwake()
	OnEnable()
	OnRender()
	OnStart()
	OnUpdate()
	RemoveDelegateToRegister(IDelegate, IEntity, IComponent) error
	SetActive(bool)
	SetDelegate(IDelegate)
	SetEntity(IEntity)
	SetRemoveOnDestroy(bool)
	Unmarshal(map[string]interface{})
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
	entity          IEntity
	active          bool
	delegate        IDelegate
	registers       []IRegister
	removeOnDestroy bool
}

var _ IComponent = (*Component)(nil)

// NewComponent creates a new component instance.
func NewComponent(name string) *Component {
	Logger.Trace().Str("component", name).Msg("new component")
	return &Component{
		Object:          NewObject(name),
		entity:          nil,
		active:          true,
		delegate:        nil,
		registers:       []IRegister{},
		removeOnDestroy: true,
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
		GetDelegateManager().DeregisterFromDelegate(register.GetRegisterID())
		register.SetDelegate(nil)
	}
	// Delete delegate being created.
	if c.GetDelegate() != nil {
		GetDelegateManager().DeleteDelegate(c.GetDelegate())
	}
	c.registers = []IRegister{}
	c.delegate = nil
	c.SetLoaded(false)
	c.SetStarted(false)
}

// DoDump dumps component in JSON format.
func (c *Component) DoDump(component IComponent) {
	if result, err := json.Marshal(component); err == nil {
		fmt.Printf("%s\n", result)
	}
}

// DoFrameEnd calls all methods to run at the end of a tick frame.
func (c *Component) DoFrameEnd() {
}

// DoFrameStart calls all methods to run at the start of a tick frame.
func (c *Component) DoFrameStart() bool {
	return c.GetStarted()
}

// DoLoad is called when component is loaded by the entity.
func (c *Component) DoLoad(component IComponent) bool {
	Logger.Trace().Str("component", c.GetName()).Msg("DoLoad")
	// fmt.Printf("load: %#v\n", reflect.TypeOf(component).String())
	// component.OnStart()
	return c.GetLoaded()
}

// DoUnLoad is called when component is unloaded by the entity. If this method
// is overwritten in any child class, Component  method has to be called,
// because it is in charge to deregister all delegates and registers.
func (c *Component) DoUnLoad() {
	Logger.Trace().Str("component", c.GetName()).Msg("DoUnLoad")
	// Deregister all register entries from delegate handler
	for _, register := range c.registers {
		GetDelegateManager().DeregisterFromDelegate(register.GetRegisterID())
		// Register is created based on delegates for those belonging to the
		// DelegateManager. Those are unchanged and delegate does not have to
		// be cleared.
		if register.GetDelegate() != nil &&
			register.GetDelegate().GetID() != GetDelegateManager().GetCollisionDelegate().GetID() &&
			register.GetDelegate().GetID() != GetDelegateManager().GetDestroyDelegate().GetID() &&
			register.GetDelegate().GetID() != GetDelegateManager().GetLoadDelegate().GetID() {
			register.SetDelegate(nil)

		}
	}
	// Delete delegate being created.
	if c.GetDelegate() != nil {
		GetDelegateManager().DeleteDelegate(c.GetDelegate())
	}
	c.SetLoaded(false)
	c.SetStarted(false)
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

// GetRemoveOnDestroy returns component remove on destroy. If this attribute
// is true, the component will be removed from the entity when scene is
// fully unloaded.
func (c *Component) GetRemoveOnDestroy() bool {
	return c.removeOnDestroy
}

// OnAwake should create all component resources that don't have any dependency
// with any other component or entity.
func (c *Component) OnAwake() {
	Logger.Trace().Str("component", c.GetName()).Msg("OnAwake")
	c.SetLoaded(true)
}

// OnEnable is called every time the component is enabled.
func (c *Component) OnEnable() {
	Logger.Trace().Str("component", c.GetName()).Msg("OnEnable")
}

// OnRender is called for every render tick.
func (c *Component) OnRender() {
	// Logger.Trace().Str("component", c.GetName()).Msg("OnRender")
}

// OnStart is called first time the component is enabled.
func (c *Component) OnStart() {
	Logger.Trace().Str("component", c.GetName()).Msg("OnStart")
	if !c.GetStarted() {
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
				if registerID, ok := GetDelegateManager().RegisterToDelegate(c, delegate, register.GetSignature()); ok {
					register.SetRegisterID(registerID)
					continue
				}
			}
			Logger.Error().Err(fmt.Errorf("register for component %s failed", c.GetName())).Str("component", c.GetName()).Msg("registration error")
			// panic("Failure at register " + register.GetName())
		}
		c.SetStarted(true)
	}
}

// OnUpdate is called for every update tick.
func (c *Component) OnUpdate() {
	// Logger.Trace().Str("component", c.GetName()).Msg("OnUpdate")
}

// RemoveDelegateToRegister removes a register for the given delegate in
// delegate handler. Register is being removed from the component too.
func (c *Component) RemoveDelegateToRegister(delegate IDelegate, entity IEntity, component IComponent) error {
	Logger.Trace().Str("component", c.GetName()).Msg("RemoveDelegateToRegister")
	index := -1
	for i, register := range c.registers {
		if delegate != nil && delegate == register.GetDelegate() {
			index = i
			break
		}
		if entity == nil {
			entity = c.GetEntity()
		}
		if _component := entity.GetComponent(component); _component != nil {
			if _delegate := _component.GetDelegate(); _delegate != nil {
				if _delegate == register.GetDelegate() {
					index = i
					break
				}
			}
		}
	}
	if index != -1 {
		GetDelegateManager().DeregisterFromDelegate(c.registers[index].GetRegisterID())
		c.registers = append(c.registers[:index], c.registers[index+1:]...)
		return nil
	}
	return fmt.Errorf("register not found for given input")
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

// SetRemoveOnDestroy sets component remove on destroy. If this attribute
// is true, the component will be removed from the entity when scene is
// fully unloaded.
func (c *Component) SetRemoveOnDestroy(remove bool) {
	c.removeOnDestroy = remove
}

// Unmarshal takes a ComponentToMarshal instance and  creates a new entity
// instance.
func (c *Component) Unmarshal(data map[string]interface{}) {
	c.SetName(data["name"].(string))
}
