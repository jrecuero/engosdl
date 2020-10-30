package engosdl

// IEvent represents any event to be used in the pool event handler.
type IEvent interface {
	IObject
}

// IEventHandler represents the interface for the event handler.
type IEventHandler interface {
	IObject
	GetDelegateHandler() IDelegateHandler
	GetPoolHandler() IPoolHandler
	OnStart()
}

// EventHandler is the default implementation fort the event handler
// interface.
type EventHandler struct {
	*Object
	delegateHandler IDelegateHandler
}

// NewEventHandler creates a new event handler instance.
func NewEventHandler(name string) *EventHandler {
	return &EventHandler{
		delegateHandler: NewDelegateHandler("delegate-handler"),
	}
}

// GetDelegateHandler returns the delegate event handler.
func (h *EventHandler) GetDelegateHandler() IDelegateHandler {
	return h.delegateHandler
}

// GetPoolHandler returns the pool event handler.
func (h *EventHandler) GetPoolHandler() IPoolHandler {
	return nil
}

// OnStart calls OnStart for all event handlers.
func (h *EventHandler) OnStart() {
	h.GetDelegateHandler().OnStart()
	// h.GetPoolHandler().OnStart()
}
