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

// KeyboardStandardMove contains all directional keys.
var KeyboardStandardMove map[int]bool = map[int]bool{
	sdl.SCANCODE_LEFT:  true,
	sdl.SCANCODE_RIGHT: true,
	sdl.SCANCODE_UP:    true,
	sdl.SCANCODE_DOWN:  true,
}

// KeyboardStandardMoveAndShoot contains all directional keys and space.
var KeyboardStandardMoveAndShoot map[int]bool = map[int]bool{
	sdl.SCANCODE_LEFT:  true,
	sdl.SCANCODE_RIGHT: true,
	sdl.SCANCODE_UP:    true,
	sdl.SCANCODE_DOWN:  true,
	sdl.SCANCODE_SPACE: true,
}

// Keyboard represents a component that can take keyboard input
type Keyboard struct {
	*engosdl.Component
	keys      map[int]bool
	KeysToMap map[int]bool `json:"keys-to-map"`
}

// NewKeyboard creates a new keyboard instance.
// It creates delegate "on-keyboard".
// It registers to "on-out-of-bounds" delegate.
func NewKeyboard(name string, keysToMap map[int]bool) *Keyboard {
	engosdl.Logger.Trace().Str("component", "keyboard").Str("keyboard", name).Msg("new keyboard")
	result := &Keyboard{
		Component: engosdl.NewComponent(name),
		keys:      make(map[int]bool),
		KeysToMap: keysToMap,
	}
	return result
}

// CreateKeyboard implements keyboard constructor used by component manager.
// It creates delegate "on-keyboard".
// It registers to "on-out-of-bounds" delegate.
func CreateKeyboard(params ...interface{}) engosdl.IComponent {
	if len(params) == 2 {
		return NewKeyboard(params[0].(string), params[1].(map[int]bool))
	}
	return NewKeyboard("", make(map[int]bool))
}

// DefaultAddDelegateToRegister will proceed to add default delegate to
// register for the component.
func (c *Keyboard) DefaultAddDelegateToRegister() {
}

// OnAwake should create all component resources that don't have any dependency
// with any other component or entity.
// It creates delegate "on-keyboard".
func (c *Keyboard) OnAwake() {
	engosdl.Logger.Trace().Str("component", "keyboard").Str("keyboard", c.GetName()).Msg("OnAwake")
	// Create new delegate "on-keyboard"
	name := fmt.Sprintf("on-keyboard/%s", c.GetName())
	c.SetDelegate(engosdl.GetDelegateManager().CreateDelegate(c, name))
	c.Component.OnAwake()
}

// OnStart is called first time the component is enabled.
func (c *Keyboard) OnStart() {
	engosdl.Logger.Trace().Str("component", "keyboard").Str("keyboard", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}

// OnUpdate is called for every update tick.
func (c *Keyboard) OnUpdate() {
	keys := sdl.GetKeyboardState()
	for key, trigger := range c.KeysToMap {
		if trigger && keys[key] == 1 {
			engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), true, key)
		} else if !trigger && keys[key] == 1 {
			if _, ok := c.keys[key]; !ok {
				c.keys[key] = true
			}
		} else if !trigger && keys[key] == 0 && c.keys[key] {
			engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), false, key)
			c.keys[key] = false
		}
	}
}
