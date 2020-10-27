package engosdl

// IEvent represents any event to be used in the pool event handler.
type IEvent interface {
	IObject
}

// IEventPool represents the pool where events are stores.
type IEventPool interface {
	IObject
	Add(IEvent) bool
	Next() IEvent
	Pop() IEvent
}

// IPoolHandler represents the interface for the pool event handler.
type IPoolHandler interface {
	IObject
	GetPool() IEventPool
	AddEvent(IEvent) bool
	NextEventInPool() IEvent
	PopEventInPool() IEvent
}

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
	CreateDelegate() IDelegate
	RegisterToDelegate(IDelegate, TDelegateSignature) (int64, bool)
	UnregisterFromDeleate(int64) bool
	TriggerDelegate(IDelegate, ...interface{})
}

// IEventHandler represents the interface for the event handler.
type IEventHandler interface {
	IObject
	GetPoolHandler() IPoolHandler
	GetDelegateHandler() IDelegateHandler
}

// EventHandler is the default implementation fort the event handler
// interface.
type EventHandler struct {
	*Object
}

// GetPoolHandler returns the pool event handler.
func (h *EventHandler) GetPoolHandler() IPoolHandler {
	return nil
}

// GetDelegateHandler returns the delegate event handler.
func (h *EventHandler) GetDelegateHandler() IDelegateHandler {
	return nil
}

// NewEventHandler creates a new event handler instance.
func NewEventHandler(name string) *EventHandler {
	return &EventHandler{}
}
