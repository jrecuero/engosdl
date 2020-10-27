package components

import (
	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// Keyboard represents a component that can take keyboard input
type Keyboard struct {
	*engosdl.Component
}

// OnUpdate is called for every update tick.
func (k *Keyboard) OnUpdate() {
	keys := sdl.GetKeyboardState()
	if keys[sdl.SCANCODE_LEFT] == 1 {
	}
	if keys[sdl.SCANCODE_RIGHT] == 1 {
	}
	if keys[sdl.SCANCODE_UP] == 1 {
		position := k.GetGameObject().GetTransform().GetPosition()
		position.Y -= 10
	}
	if keys[sdl.SCANCODE_DOWN] == 1 {
		position := k.GetGameObject().GetTransform().GetPosition()
		position.Y += 10
	}
	if keys[sdl.SCANCODE_SPACE] == 1 {
		scale := k.GetGameObject().GetTransform().GetScale()
		scale.X = 2
	}
	if keys[sdl.SCANCODE_TAB] == 1 {
		scale := k.GetGameObject().GetTransform().GetScale()
		scale.X = 1
	}
}

// NewKeyboard creates a new keyboard instance.
func NewKeyboard(name string, gobj *engosdl.GameObject) *Keyboard {
	engosdl.Logger.Trace().Str("component", "keyboard").Str("keyboard", name).Msg("new keyboard")
	return &Keyboard{
		Component: engosdl.NewComponent(name, gobj),
	}
}
