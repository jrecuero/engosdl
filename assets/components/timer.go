package components

import (
	"fmt"
	"reflect"

	"github.com/jrecuero/engosdl"
)

// ComponentNameTimer is the name to refer timer component.
var ComponentNameTimer string = reflect.TypeOf(&Timer{}).String()

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor(ComponentNameTimer, CreateTimer)
	}
}

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

// CreateTimer implements timer constructor used by component manager.
func CreateTimer(params ...interface{}) engosdl.IComponent {
	if len(params) == 2 {
		return NewTimer(params[0].(string), params[1].(int))
	}
	return NewTimer("", 0)
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

// Unmarshal takes a ComponentToMarshal instance and  creates a new entity
// instance.
func (t *Timer) Unmarshal(data map[string]interface{}) {
	t.Component.Unmarshal(data)
	t.Tick = int(data["key"].(float64))
}
