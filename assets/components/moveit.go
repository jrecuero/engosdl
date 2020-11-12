package components

import (
	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterComponent(&MoveIt{})
	}
}

// MoveIt represents a component that can take move-it input
type MoveIt struct {
	*engosdl.Component
	Speed    *engosdl.Vector `json:"speed"`
	position *engosdl.Vector
}

// NewMoveIt creates a new move-it instance.
// It registers to "on-out-of-bounds" delegate.
func NewMoveIt(name string, speed *engosdl.Vector) *MoveIt {
	engosdl.Logger.Trace().Str("component", "move-it").Str("move-it", name).Msg("new move-it")
	result := &MoveIt{
		Component: engosdl.NewComponent(name),
		Speed:     speed,
		position:  engosdl.NewVector(0, 0),
	}
	return result
}

// DefaultAddDelegateToRegister will proceed to add default delegate to
// register for the component.
func (c *MoveIt) DefaultAddDelegateToRegister() {
	c.AddDelegateToRegister(nil, nil, &OutOfBounds{}, c.DefaultOnOutOfBounds)
	c.AddDelegateToRegister(nil, nil, &Keyboard{}, c.onKeyboard)
}

// DefaultOnOutOfBounds checks if the entity has gone out of bounds.
func (c *MoveIt) DefaultOnOutOfBounds(params ...interface{}) bool {
	c.GetEntity().GetTransform().SetPosition(engosdl.NewVector(c.position.X, c.position.Y))
	return true
}

func (c *MoveIt) onKeyboard(params ...interface{}) bool {
	position := c.GetEntity().GetTransform().GetPosition()
	c.position.X = position.X
	c.position.Y = position.Y
	key := params[0].(int)
	switch key {
	case sdl.SCANCODE_LEFT:
		position.X -= c.Speed.X
		break
	case sdl.SCANCODE_RIGHT:
		position.X += c.Speed.X
		break
	case sdl.SCANCODE_UP:
		position.Y -= c.Speed.Y
		break
	case sdl.SCANCODE_DOWN:
		position.Y += c.Speed.Y
		break
	}
	return true
}

// OnStart is called first time the component is enabled.
func (c *MoveIt) OnStart() {
	engosdl.Logger.Trace().Str("component", "move-to").Str("move-to", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}

// OnUpdate is called for every update tick.
func (c *MoveIt) OnUpdate() {
}
