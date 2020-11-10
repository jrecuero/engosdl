package engosdl

// IEvent represents any event to be used in the pool event handler.
type IEvent interface {
	IObject
	Add(IEvent) bool
	Next() IEvent
	Pop() IEvent
}

// IEventHandler represents the interface for the  event handler.
type IEventHandler interface {
	IObject
	// AddEvent(IEvent) bool
	// GetPool() []IEvent
	// NextEventInPool() IEvent
	OnStart()
	// PopEventInPool() IEvent
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
		Object: NewObject(name),
	}
}

// OnStart calls OnStart for all event handlers.
func (h *EventHandler) OnStart() {
	Logger.Trace().Str("event-handler", h.GetName()).Msg("OnStart")
}
