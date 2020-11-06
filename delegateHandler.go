package engosdl

const (
	delegateHandlerName   = "delegate-handler"
	collisionDelegate     = "on-collision"
	collisionDelegateName = delegateHandlerName + "/" + collisionDelegate
	destroyDelegate       = "on-destroy"
	destroyDelegateName   = delegateHandlerName + "/" + destroyDelegate
	loadDelegate          = "on-load"
	loadDelegateName      = delegateHandlerName + "/" + loadDelegate
)

// IDelegate represents any delegate to be used in the delegate event handler.
type IDelegate interface {
	IObject
	GetObject() IObject
	GetEventName() string
}

// TDelegateSignature represents the callback for any method to be registered
// to a delegate.
type TDelegateSignature func(...interface{}) bool

// IRegister represents all information required to register to a delegate.
type IRegister interface {
	GetComponent() IComponent
	GetDelegate() IDelegate
	GetEntity() IEntity
	GetRegisterID() string
	GetSignature() TDelegateSignature
	SetComponent(IComponent) IRegister
	SetDelegate(IDelegate) IRegister
	SetEntity(IEntity) IRegister
	SetRegisterID(string) IRegister
	SetSignature(TDelegateSignature) IRegister
}

// IDelegateHandler represents the interface for the delefate event handler.
type IDelegateHandler interface {
	IObject
	CreateDelegate(IObject, string) IDelegate
	DeleteDelegate(IDelegate) bool
	DeregisterFromDelegate(string) bool
	GetCollisionDelegate() IDelegate
	GetDestroyDelegate() IDelegate
	GetLoadDelegate() IDelegate
	OnStart()
	RegisterToDelegate(IDelegate, TDelegateSignature) (string, bool)
	TriggerDelegate(IDelegate, ...interface{})
}

// Delegate is the default implementation for delegate interface.
type Delegate struct {
	*Object
	obj    IObject
	evName string
}

// var _ IDelegate = (*Delegate)(nil)

// GetObject returns delegate object.
func (d *Delegate) GetObject() IObject {
	return d.obj
}

// GetEventName returns delegate event name.
func (d *Delegate) GetEventName() string {
	return d.evName
}

// NewDelegate creates a new delegate instance.
func NewDelegate(name string, obj IObject, evName string) *Delegate {
	Logger.Trace().Str("delegate", name).Str("evName", evName).Msg("new delegate")
	return &Delegate{
		Object: NewObject(name),
		obj:    obj,
		evName: evName,
	}
}

// Register is the structure used to store any delegate registration
type Register struct {
	*Object
	entity     IEntity
	component  IComponent
	delegate   IDelegate
	signature  TDelegateSignature
	registerID string
}

var _ IRegister = (*Register)(nil)

// NewRegister creates a new register instance.
func NewRegister(name string, entity IEntity, component IComponent, delegate IDelegate, signature TDelegateSignature) *Register {
	Logger.Trace().Str("register", name).Msg("create new register")
	return &Register{
		Object:    NewObject(name),
		entity:    entity,
		component: component,
		delegate:  delegate,
		signature: signature,
	}
}

// GetComponent returns register component.
func (r *Register) GetComponent() IComponent {
	return r.component
}

// GetDelegate returns register delegate.
func (r *Register) GetDelegate() IDelegate {
	return r.delegate
}

// GetEntity returns register entity.
func (r *Register) GetEntity() IEntity {
	return r.entity
}

// GetRegisterID returns the registerID.
func (r *Register) GetRegisterID() string {
	return r.registerID
}

// GetSignature returns the register signature.
func (r *Register) GetSignature() TDelegateSignature {
	return r.signature
}

// SetComponent sets the register component.
func (r *Register) SetComponent(component IComponent) IRegister {
	r.component = component
	return r
}

// SetDelegate sets the register delegate.
func (r *Register) SetDelegate(delegate IDelegate) IRegister {
	r.delegate = delegate
	return r
}

// SetEntity sets the register entity.
func (r *Register) SetEntity(entity IEntity) IRegister {
	r.entity = entity
	return r
}

// SetRegisterID sets the registerID.
func (r *Register) SetRegisterID(id string) IRegister {
	r.registerID = id
	return r
}

// SetSignature sets the register signature.
func (r *Register) SetSignature(signature TDelegateSignature) IRegister {
	r.signature = signature
	return r
}

// DelegateHandler is the default implementation for event handler interface.
type DelegateHandler struct {
	*Object
	delegates []IDelegate
	registers []*Register
	defaults  map[string]IDelegate
}

// NewDelegateHandler creates a new delegate handler instance.
func NewDelegateHandler(name string) *DelegateHandler {
	Logger.Trace().Str("delegate-handler", name).Msg("new delegate handler")
	return &DelegateHandler{
		Object:    NewObject(name),
		delegates: []IDelegate{},
		registers: []*Register{},
		defaults:  make(map[string]IDelegate),
	}
}

// CreateDelegate creates a new delefate in the delegate handler
func (h *DelegateHandler) CreateDelegate(obj IObject, evName string) IDelegate {
	Logger.Trace().Str("delegate-handler", h.GetName()).Msg("CreateDelegate")
	delegate := NewDelegate(obj.GetName()+"/"+evName, obj, evName)
	h.delegates = append(h.delegates, delegate)
	return delegate
}

// DeleteDelegate deletes the given delegate from delegate handler and
// all registers
func (h *DelegateHandler) DeleteDelegate(delegate IDelegate) bool {
	Logger.Trace().Str("delegate-handler", h.GetName()).Msg("DeleteDelegate")
	for i, delegat := range h.delegates {
		if delegat.GetID() == delegate.GetID() {
			h.delegates = append(h.delegates[:i], h.delegates[i+1:]...)
			for j, register := range h.registers {
				if register.GetDelegate().GetID() == delegate.GetID() {
					h.registers = append(h.registers[:j], h.registers[j+1:]...)
				}
			}
			return true
		}
	}
	return false
}

// DeregisterFromDelegate unregistered the given register from the delegate.
func (h *DelegateHandler) DeregisterFromDelegate(registerID string) bool {
	for i, register := range h.registers {
		if register.GetRegisterID() == registerID {
			Logger.Trace().Str("delegate-handler", h.GetName()).Str("delegate", register.GetDelegate().GetName()).Msg("deregister-to-delegate")
			h.registers = append(h.registers[:i], h.registers[i+1:]...)
			return true
		}
	}
	return false
}

// GetCollisionDelegate returns default delegate for collisions.
func (h *DelegateHandler) GetCollisionDelegate() IDelegate {
	return h.defaults[collisionDelegate]
}

// GetDestroyDelegate returns default delegate when entity is destroyed.
func (h *DelegateHandler) GetDestroyDelegate() IDelegate {
	return h.defaults[destroyDelegate]
}

// GetLoadDelegate returns default delegate when entity is loaded/created.
func (h *DelegateHandler) GetLoadDelegate() IDelegate {
	return h.defaults[loadDelegate]
}

// OnStart initializes all delegate handler structure.
func (h *DelegateHandler) OnStart() {
	Logger.Trace().Str("delegate-handler", h.GetName()).Msg("OnStart")
	h.defaults[collisionDelegate] = h.CreateDelegate(h, collisionDelegate)
	h.defaults[destroyDelegate] = h.CreateDelegate(h, destroyDelegate)
	h.defaults[loadDelegate] = h.CreateDelegate(h, loadDelegate)
}

// RegisterToDelegate registers a method to a delegate.
func (h *DelegateHandler) RegisterToDelegate(delegate IDelegate, signature TDelegateSignature) (string, bool) {
	Logger.Trace().Str("delegate-handler", h.GetName()).Str("delegate", delegate.GetName()).Msg("register-to-delegate")
	register := NewRegister("", nil, nil, delegate, signature)
	register.SetRegisterID(register.GetID())
	h.registers = append(h.registers, register)
	return register.GetRegisterID(), true
}

// TriggerDelegate calls all signatures registered to a given delegate.
func (h *DelegateHandler) TriggerDelegate(delegate IDelegate, params ...interface{}) {
	for _, register := range h.registers {
		if register.GetDelegate() != nil && register.GetDelegate().GetID() == delegate.GetID() {
			register.GetSignature()(params...)
		}
	}
}
