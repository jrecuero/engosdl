package engosdl

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
	OnStart()
	PopEventInPool() IEvent
}
