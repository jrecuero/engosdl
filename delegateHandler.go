package engosdl

// IDelegate represents any delegate to be used in the delegate event handler.
type IDelegate interface {
	IObject
}

// TDelegateSignature represents the callback for any method to be registered
// to a delegate.
type TDelegateSignature func(...interface{}) bool

// IDelegateHandler represents the interface for the delefate event handler.
type IDelegateHandler interface {
	IObject
	CreateDelegate(IObject, string) IDelegate
	RegisterToDelegate(IDelegate, TDelegateSignature) (string, bool)
	UnregisterFromDelegate(string) bool
	TriggerDelegate(IDelegate, ...interface{})
}

// Delegate is the default implementation for delegate interface.
type Delegate struct {
	*Object
	obj    IObject
	evName string
}

// GetObject returns delegate object.
func (d *Delegate) GetObject() IObject {
	return d.obj
}

// GetEventName returns delegate evet name.
func (d *Delegate) GetEventName() string {
	return d.evName
}

// NewDelegate creates a new delegate instance.
func NewDelegate(name string, obj IObject, evName string) *Delegate {
	Logger.Trace().Str("delegate", name).Str("evName", evName).Msg("create new delegate")
	return &Delegate{
		Object: NewObject(name),
		obj:    obj,
		evName: evName,
	}
}

// Register is the structure used to store any delegate registration
type Register struct {
	*Object
	Delegate  IDelegate
	Signature TDelegateSignature
}

// NewRegister creates a new register instance.
func NewRegister(name string, delegate IDelegate, signature TDelegateSignature) *Register {
	Logger.Trace().Str("register", name).Msg("create new register")
	return &Register{
		Object:    NewObject(name),
		Delegate:  delegate,
		Signature: signature,
	}
}

// DelegateHandler is the default implementation for event handler interface.
type DelegateHandler struct {
	*Object
	delegates []IDelegate
	registers []*Register
}

// CreateDelegate creates a new delefate in the delegate handler
func (h *DelegateHandler) CreateDelegate(obj IObject, evName string) IDelegate {
	delegate := NewDelegate(obj.GetName()+"/"+evName, obj, evName)
	h.delegates = append(h.delegates, delegate)
	return delegate
}

// RegisterToDelegate registers a method to a delegate.
func (h *DelegateHandler) RegisterToDelegate(delegate IDelegate, signature TDelegateSignature) (string, bool) {
	register := NewRegister("", delegate, signature)
	h.registers = append(h.registers, register)
	return register.GetID(), true
}

// TriggerDelegate calls all signatures registered to a given delegate.
func (h *DelegateHandler) TriggerDelegate(delegate IDelegate, params ...interface{}) {
	for _, register := range h.registers {
		if register.Delegate.GetName() == delegate.GetName() {
			register.Signature(params...)
		}
	}
}

// UnregisterFromDelegate unresgisters the given register from the delegate.
func (h *DelegateHandler) UnregisterFromDelegate(string) bool {
	return true
}

// NewDelegateHandler creates a new delegate handler instance.
func NewDelegateHandler(name string) *DelegateHandler {
	Logger.Trace().Str("delegate-handler", name).Msg("create new delegate handler")
	return &DelegateHandler{
		Object:    NewObject(name),
		delegates: []IDelegate{},
		registers: []*Register{},
	}
}
