package engosdl

// IEvent represents any event to be used in the pool event handler.
type IEvent interface {
	IObject
	Add(IEvent) bool
	Next() IEvent
	Pop() IEvent
}

// IEventManager represents the interface for the  event handler.
type IEventManager interface {
	IObject
	// AddEvent(IEvent) bool
	// GetPool() []IEvent
	// NextEventInPool() IEvent
	OnStart()
	// PopEventInPool() IEvent
}

// EventManager is the default implementation fort the event handler
// interface.
type EventManager struct {
	*Object
	delegateManager IDelegateManager
}

// NewEventManager creates a new event handler instance.
func NewEventManager(name string) *EventManager {
	return &EventManager{
		Object: NewObject(name),
	}
}

// OnStart calls OnStart for all event handlers.
func (h *EventManager) OnStart() {
	Logger.Trace().Str("event-handler", h.GetName()).Msg("OnStart")
}
