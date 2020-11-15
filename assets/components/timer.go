package components

import (
	"fmt"

	"github.com/jrecuero/engosdl"
)

// ITimer represents the timer component. Timer should be base in engine
// frames.
type ITimer interface {
	engosdl.IComponent
	GetTick() int
	SetTick(int)
}

// Timer is the default implementation for the timer component interface.
type Timer struct {
	*engosdl.Component
	Tick        int
	tickCounter int
}

var _ ITimer = (*Timer)(nil)

// NewTimer creates a new timer instance.
func NewTimer(name string, tick int) *Timer {
	engosdl.Logger.Trace().Str("timer", name).Msg("new timer")
	return &Timer{
		Component:   engosdl.NewComponent(name),
		Tick:        tick,
		tickCounter: 0,
	}
}

// GetTick returns the timer tick. Tick is the number of engine frames before
// the timer has to be triggered.
func (t *Timer) GetTick() int {
	return t.Tick
}

// SetTick sets the timer tick. This is the number of engine frames before the
// timer has to be triggered.
func (t *Timer) SetTick(tick int) {
	t.Tick = tick
}

// OnAwake is called the first time the component is loaded in the scene,
// it should create any independent resource from other components or
// entities.
func (t *Timer) OnAwake() {
	engosdl.Logger.Trace().Str("timer", t.GetName()).Msg("OnAwake")
	name := fmt.Sprintf("on-timer/%s", t.GetName())
	t.SetDelegate(engosdl.GetDelegateManager().CreateDelegate(t, name))
}

// OnUpdate is called every engine frame in order to update the component.
func (t *Timer) OnUpdate() {
	if t.tickCounter == t.Tick {
		t.tickCounter = 0
		engosdl.GetDelegateManager().TriggerDelegate(t.GetDelegate(), false)
	} else {
		t.tickCounter++
	}
}
