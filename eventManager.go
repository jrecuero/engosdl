package engosdl

import "fmt"

// IEvent represents any event to be used in the pool event handler.
type IEvent interface {
	IObject
	GetData() IObject
	SetData(IObject)
}

// Event is the default implementation for IEvent interface.
type Event struct {
	*Object
	data IObject
}

var _ IEvent = (*Event)(nil)

// NewEvent creates a new event instance.
func NewEvent(name string, data IObject) *Event {
	Logger.Trace().Str("event", name).Str("data", data.GetName()).Msg("new event")
	return &Event{
		Object: NewObject(name),
		data:   data,
	}
}

// GetData returns event data.
func (e *Event) GetData() IObject {
	return e.data
}

// SetData sets event data.
func (e *Event) SetData(data IObject) {
	e.data = data
}

// IEventPool represents any event pool.
type IEventPool interface {
	IObject
	Add(IEvent) error
	Flush() error
	Next() (IEvent, error)
	Pool() []IEvent
	Pop() (IEvent, error)
}

// EventPool is the default implementation for IEventPool interface.
type EventPool struct {
	*Object
	pool []IEvent
}

var _ IEventPool = (*EventPool)(nil)

// NewEventPool creates a new event pool instance.
func NewEventPool(name string) *EventPool {
	Logger.Trace().Str("event-pool", name).Msg("new event-pool")
	return &EventPool{
		Object: NewObject(name),
		pool:   []IEvent{},
	}
}

// Add adds a new event to the event pool.
func (ep *EventPool) Add(event IEvent) error {
	Logger.Trace().Str("event-pool", ep.GetName()).Str("event", event.GetName()).Msg("Add")
	ep.pool = append(ep.pool, event)
	return nil
}

// Flush removes all events from the event pool.
func (ep *EventPool) Flush() error {
	Logger.Trace().Str("event-pool", ep.GetName()).Msg("Flush")
	ep.pool = []IEvent{}
	return nil
}

// Next returns the first element in the event pool. Not removed.
func (ep *EventPool) Next() (IEvent, error) {
	Logger.Trace().Str("event-pool", ep.GetName()).Msg("Next")
	if len(ep.pool) != 0 {
		return ep.pool[0], nil
	}
	return nil, fmt.Errorf("pool %s is empty", ep.GetName())
}

// Pool returns the whole pool
func (ep *EventPool) Pool() []IEvent {
	return ep.pool
}

// Pop returns the first element in the event pool. Remove event.
func (ep *EventPool) Pop() (IEvent, error) {
	Logger.Trace().Str("event-pool", ep.GetName()).Msg("Pop")
	var event IEvent
	if len(ep.pool) != 0 {
		event, ep.pool = ep.pool[0], ep.pool[1:]
		return event, nil
	}
	return nil, fmt.Errorf("pool %s is empty", ep.GetName())
}

// IEventManager represents the interface for the  event handler.
type IEventManager interface {
	IObject
	CreatePool(string) (string, error)
	DeletePool(string) error
	DoInit()
	GetIDForName(string) (string, error)
	GetPool(string) IEventPool
	GetPools() map[string]IEventPool
	OnStart()
}

// EventManager is the default implementation fort the event handler
// interface.
type EventManager struct {
	*Object
	pools    map[string]IEventPool
	poolsMap map[string]string
}

var _ IEventManager = (*EventManager)(nil)

// NewEventManager creates a new event handler instance.
func NewEventManager(name string) *EventManager {
	return &EventManager{
		Object:   NewObject(name),
		pools:    make(map[string]IEventPool),
		poolsMap: make(map[string]string),
	}
}

// CreatePool creates a new event pool in the event manager.
func (h EventManager) CreatePool(name string) (string, error) {
	Logger.Trace().Str("event-manager", h.GetName()).Str("pool", name).Msg("create pool")
	for keyName, key := range h.poolsMap {
		if keyName == name {
			return key, fmt.Errorf("pool with name %s already in event manager", name)
		}
	}
	id := nextIder()
	h.poolsMap[name] = id
	h.pools[id] = NewEventPool(name)
	return id, nil
}

// DeletePool deletes the event pool for the given identification.
func (h *EventManager) DeletePool(id string) error {
	Logger.Trace().Str("event-manager", h.GetName()).Str("pool", id).Msg("delete pool")
	if _, ok := h.pools[id]; ok {
		delete(h.pools, id)
		return nil
	}
	return fmt.Errorf("pool with id %s not found in event manager", id)
}

// DoInit initializes all event manager resources.
func (h *EventManager) DoInit() {
	Logger.Trace().Str("event-manager", h.GetName()).Msg("DoInit")
}

// GetIDForName returns the event pool identification for the given name.
func (h *EventManager) GetIDForName(name string) (string, error) {
	if id, ok := h.poolsMap[name]; ok {
		return id, nil
	}
	return "", fmt.Errorf("pool with name %s not found in event manager", name)
}

// GetPool returns the event pool for the given identification.
func (h *EventManager) GetPool(id string) IEventPool {
	if pool, ok := h.pools[id]; ok {
		return pool
	}
	return nil
}

// GetPools returns all event pools in the event manager.
func (h *EventManager) GetPools() map[string]IEventPool {
	return h.pools
}

// OnStart calls OnStart for all event handlers.
func (h *EventManager) OnStart() {
	Logger.Trace().Str("event-manager", h.GetName()).Msg("OnStart")
}
