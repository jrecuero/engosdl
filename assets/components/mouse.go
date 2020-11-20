package components

import (
	"fmt"
	"reflect"

	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// ComponentNameMouse is the name to refer mouse component.
var ComponentNameMouse string = reflect.TypeOf(&Mouse{}).String()

const (
	// ButtonLeft is the state when left buttom is click.
	ButtonLeft uint32 = 1 << (sdl.BUTTON_LEFT - 1)

	// ButtonMiddle is the state when middle buttom is click.
	ButtonMiddle uint32 = 1 << (sdl.BUTTON_MIDDLE - 1)

	// ButtonRight is the state when right buttom is click.
	ButtonRight uint32 = 1 << (sdl.BUTTON_RIGHT - 1)

	// ButtonX1 is the state when x1 buttom is click
	ButtonX1 uint32 = 1 << (sdl.BUTTON_X1 - 1)

	// ButtonX2 is the state when x2 buttom is click.
	ButtonX2 uint32 = 1 << (sdl.BUTTON_X2 - 1)
)

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor(ComponentNameMouse, CreateMouse)
	}
}

// Mouse represents a component that can take mouse input.
type Mouse struct {
	*engosdl.Component
	buttons map[uint32]bool
	OnClick bool `json:"on-click"`
}

// NewMouse creates a new mouse instance.
func NewMouse(name string, onClick bool) *Mouse {
	engosdl.Logger.Trace().Str("component", "mouse").Str("mouse", name).Msg("new mouse")
	result := &Mouse{
		Component: engosdl.NewComponent(name),
		buttons:   make(map[uint32]bool),
		OnClick:   onClick,
	}
	return result
}

// CreateMouse implements mouse constructor used by component manager.
func CreateMouse(params ...interface{}) engosdl.IComponent {
	if len(params) == 2 {
		return NewMouse(params[0].(string), params[1].(bool))
	}
	return NewMouse("", true)
}

// OnAwake should create al component resources that don't have any dependency
// with any other component or entity.
func (c *Mouse) OnAwake() {
	engosdl.Logger.Trace().Str("component", "mouse").Str("mouse", c.GetName()).Msg("OnAwake")
	c.Component.OnAwake()
}

// OnUpdate is called for every update frame.
func (c Mouse) OnUpdate() {
	x, y, state := sdl.GetMouseState()
	switch state {
	case ButtonLeft:
		c.buttons[ButtonLeft] = true
		if !c.OnClick {
			fmt.Printf("left mouse click at (%d, %d) : %d\n", x, y, state)
		}
		break
	case ButtonMiddle:
		c.buttons[ButtonMiddle] = true
		if !c.OnClick {
			fmt.Printf("middle mouse click at (%d, %d) : %d\n", x, y, state)
		}
		break
	case ButtonRight:
		c.buttons[ButtonRight] = true
		if !c.OnClick {
			fmt.Printf("right mouse click at (%d, %d) : %d\n", x, y, state)
		}
		break
	case ButtonX1:
		c.buttons[ButtonX1] = true
		if !c.OnClick {
			fmt.Printf("x1 mouse click at (%d, %d) : %d\n", x, y, state)
		}
		break
	case ButtonX2:
		c.buttons[ButtonX2] = true
		if !c.OnClick {
			fmt.Printf("x1 mouse click at (%d, %d) : %d\n", x, y, state)
		}
		break
	default:
		if c.OnClick {
			if c.buttons[ButtonLeft] {
				c.buttons[ButtonLeft] = false
				fmt.Printf("left mouse click at (%d, %d) : %d\n", x, y, state)
			}
			if c.buttons[ButtonMiddle] {
				c.buttons[ButtonMiddle] = false
				fmt.Printf("middle mouse click at (%d, %d) : %d\n", x, y, state)
			}
			if c.buttons[ButtonRight] {
				c.buttons[ButtonRight] = false
				fmt.Printf("right mouse click at (%d, %d) : %d\n", x, y, state)
			}
			if c.buttons[ButtonX1] {
				c.buttons[ButtonX1] = false
				fmt.Printf("x1 mouse click at (%d, %d) : %d\n", x, y, state)
			}
			if c.buttons[ButtonX2] {
				c.buttons[ButtonX2] = false
				fmt.Printf("x2 mouse click at (%d, %d) : %d\n", x, y, state)
			}
		}
		break
	}
}

// OnStart is called first time component is enabled.
func (c *Mouse) OnStart() {
	engosdl.Logger.Trace().Str("component", "mouse").Str("mouse", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}

// Unmarshal takes a ComponentToMarshal instance and  creates a new entity
// instance.
func (c *Mouse) Unmarshal(data map[string]interface{}) {
	c.Component.Unmarshal(data)
	c.OnClick = data["on-click"].(bool)
}
