package components

import (
	"fmt"
	"reflect"

	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// ComponentNameKeyboard is the name to refer keyboard component
var ComponentNameKeyboard string = reflect.TypeOf(&Keyboard{}).String()

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor(ComponentNameKeyboard, CreateKeyboard)
	}
}

// Keyboard represents a component that can take keyboard input
type Keyboard struct {
	*engosdl.Component
	keys map[int]bool
}

// NewKeyboard creates a new keyboard instance.
// It creates delegate "on-keyboard".
// It registers to "on-out-of-bounds" delegate.
func NewKeyboard(name string) *Keyboard {
	engosdl.Logger.Trace().Str("component", "keyboard").Str("keyboard", name).Msg("new keyboard")
	result := &Keyboard{
		Component: engosdl.NewComponent(name),
		keys:      make(map[int]bool),
	}
	return result
}

// CreateKeyboard implements keyboard constructor used by component manager.
func CreateKeyboard(params ...interface{}) engosdl.IComponent {
	return NewKeyboard(params[0].(string))
}

// DefaultAddDelegateToRegister will proceed to add default delegate to
// register for the component.
func (c *Keyboard) DefaultAddDelegateToRegister() {
}

// OnAwake should create all component resources that don't have any dependency
// with any other component or entity.
func (c *Keyboard) OnAwake() {
	engosdl.Logger.Trace().Str("component", "keyboard").Str("keyboard", c.GetName()).Msg("OnAwake")
	// Create new delegate "on-keyboard"
	name := fmt.Sprintf("on-keyboard/%s", c.GetName())
	c.SetDelegate(engosdl.GetDelegateManager().CreateDelegate(c, name))
	c.Component.OnAwake()
}

// OnStart is called first time the component is enabled.
func (c *Keyboard) OnStart() {
	engosdl.Logger.Trace().Str("component", "move-to").Str("move-to", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}

// OnUpdate is called for every update tick.
func (c *Keyboard) OnUpdate() {
	keys := sdl.GetKeyboardState()
	if keys[sdl.SCANCODE_LEFT] == 1 {
		engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), true, sdl.SCANCODE_LEFT)
	}
	if keys[sdl.SCANCODE_RIGHT] == 1 {
		engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), true, sdl.SCANCODE_RIGHT)
	}
	if keys[sdl.SCANCODE_UP] == 1 {
		engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), true, sdl.SCANCODE_UP)
	}
	if keys[sdl.SCANCODE_DOWN] == 1 {
		engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), true, sdl.SCANCODE_DOWN)
	}
	if keys[sdl.SCANCODE_RETURN] == 1 {
		if _, ok := c.keys[sdl.SCANCODE_RETURN]; !ok {
			c.keys[sdl.SCANCODE_RETURN] = true
		}
		// engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), false, sdl.SCANCODE_RETURN)
	}
	if keys[sdl.SCANCODE_N] == 1 {
		engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), false, sdl.SCANCODE_N)
	}
	if keys[sdl.SCANCODE_P] == 1 {
		engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), false, sdl.SCANCODE_P)
	}
	if keys[sdl.SCANCODE_SPACE] == 1 {
		engosdl.Logger.Trace().Str("component", "keyboard").Str("keyboard", c.GetName()).Msg("space key pressed")
	}

	if keys[sdl.SCANCODE_RETURN] == 0 && c.keys[sdl.SCANCODE_RETURN] {
		engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), false, sdl.SCANCODE_RETURN)
		c.keys[sdl.SCANCODE_RETURN] = false
	}

}
