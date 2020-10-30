package components

import (
	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// Keyboard represents a component that can take keyboard input
type Keyboard struct {
	*engosdl.Component
	speed *engosdl.Vector
}

// NewKeyboard creates a new keyboard instance.
func NewKeyboard(name string, speed *engosdl.Vector) *Keyboard {
	engosdl.Logger.Trace().Str("component", "keyboard").Str("keyboard", name).Msg("new keyboard")
	return &Keyboard{
		Component: engosdl.NewComponent(name),
		speed:     speed,
	}
}

// OnUpdate is called for every update tick.
func (c *Keyboard) OnUpdate() {
	keys := sdl.GetKeyboardState()
	position := c.GetEntity().GetTransform().GetPosition()
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
