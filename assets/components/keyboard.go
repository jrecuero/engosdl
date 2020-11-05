package components

import (
	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// Keyboard represents a component that can take keyboard input
type Keyboard struct {
	*engosdl.Component
	speed    *engosdl.Vector
	position *engosdl.Vector
}

// NewKeyboard creates a new keyboard instance.
func NewKeyboard(name string, speed *engosdl.Vector) *Keyboard {
	engosdl.Logger.Trace().Str("component", "keyboard").Str("keyboard", name).Msg("new keyboard")
	return &Keyboard{
		Component: engosdl.NewComponent(name),
		speed:     speed,
		position:  engosdl.NewVector(0, 0),
	}
}

// onOutOfBounds checks if the entity has gone out of bounds.
func (c *Keyboard) onOutOfBounds(params ...interface{}) bool {
	position := c.GetEntity().GetTransform().GetPosition()
	position.X = c.position.X
	position.Y = c.position.Y
	return true
}

// OnStart is called first time the component is enabled.
func (c *Keyboard) OnStart() {
	engosdl.Logger.Trace().Str("component", "move-to").Str("move-to", c.GetName()).Msg("OnStart")
	if component := c.GetEntity().GetComponent(&OutOfBounds{}); component != nil {
		if outOfBoundsComponent, ok := component.(*OutOfBounds); ok {
			if delegate := outOfBoundsComponent.GetDelegate(); delegate != nil {
				// engosdl.GetEngine().GetEventHandler().GetDelegateHandler().RegisterToDelegate(delegate, c.onOutOfBounds)
				c.AddDelegateToRegister(delegate, nil, nil, c.onOutOfBounds)
			}
		}
	}
	c.Component.OnStart()
}

// OnUpdate is called for every update tick.
func (c *Keyboard) OnUpdate() {
	keys := sdl.GetKeyboardState()
	position := c.GetEntity().GetTransform().GetPosition()
	c.position.X = position.X
	c.position.Y = position.Y
	if keys[sdl.SCANCODE_LEFT] == 1 {
		position.X -= c.speed.X
	}
	if keys[sdl.SCANCODE_RIGHT] == 1 {
		position.X += c.speed.X
	}
	if keys[sdl.SCANCODE_UP] == 1 {
		position.Y -= c.speed.Y
	}
	if keys[sdl.SCANCODE_DOWN] == 1 {
		position.Y += c.speed.Y
	}
	if keys[sdl.SCANCODE_SPACE] == 1 {
		engosdl.Logger.Trace().Str("component", "keyboard").Str("keyboard", c.GetName()).Msg("space key pressed")
	}
	// if keys[sdl.SCANCODE_TAB] == 1 {
	// 	scale := k.GetEntity().GetTransform().GetScale()
	// 	scale.X = 1
	// }
}
