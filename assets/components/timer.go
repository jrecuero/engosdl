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

// Timer is the default implementation for the timer component interface.
type Timer struct {
	*engosdl.Component
	Tick         int `json:"tick"`
	Times        int `json:"times"`
	tickCounter  int
	timesCounter int
}

var _ engosdl.ITimer = (*Timer)(nil)

// NewTimer creates a new timer instance.
func NewTimer(name string, tick int, times int) *Timer {
	engosdl.Logger.Trace().Str("timer", name).Msg("new timer")
	return &Timer{
		Component:    engosdl.NewComponent(name),
		Tick:         tick,
		tickCounter:  0,
		Times:        times,
		timesCounter: 0,
	}
}

// CreateTimer implements timer constructor used by component manager.
func CreateTimer(params ...interface{}) engosdl.IComponent {
	if len(params) == 3 {
		return NewTimer(params[0].(string), params[1].(int), params[2].(int))
	}
	return NewTimer("", 0, 0)
}

// GetTick returns the timer tick. Tick is the number of engine frames before
// the timer has to be triggered.
func (t *Timer) GetTick() int {
	return t.Tick
}

// GetTimes returns the number of times timer has to be triggered.
func (t *Timer) GetTimes() int {
	return t.Times
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
	if t.Times == -1 || t.timesCounter < t.Times {
		if t.tickCounter == t.Tick {
			t.timesCounter++
			t.tickCounter = 0
			engosdl.GetDelegateManager().TriggerDelegate(t.GetDelegate(), false)
		} else {
			t.tickCounter++
		}
	}
}

// SetTick sets the timer tick. This is the number of engine frames before the
// timer has to be triggered.
func (t *Timer) SetTick(tick int) {
	t.Tick = tick
}

// SetTimes sets the number of times timer has to be triggered.
func (t *Timer) SetTimes(times int) {
	t.Times = times
}

// Unmarshal takes a ComponentToMarshal instance and  creates a new entity
// instance.
func (t *Timer) Unmarshal(data map[string]interface{}) {
	t.Component.Unmarshal(data)
	t.Tick = int(data["tick"].(float64))
	t.Times = int(data["times"].(float64))
}
