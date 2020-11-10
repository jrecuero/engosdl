package components

import "github.com/jrecuero/engosdl"

// MoveTo represents a component that moves a entity.
type MoveTo struct {
	*engosdl.Component
	speed *engosdl.Vector
}

// NewMoveTo creates a new move-to instance.
// It registers to on-out-of-bounds delegate.
func NewMoveTo(name string, speed *engosdl.Vector) *MoveTo {
	engosdl.Logger.Trace().Str("component", "move-to").Str("move-to", name).Msg("new move-to")
	result := &MoveTo{
		Component: engosdl.NewComponent(name),
		speed:     speed,
	}
	result.AddDelegateToRegister(nil, nil, &OutOfBounds{}, result.onOutOfBounds)
	return result
}

// GetSpeed returns component speed.
func (c *MoveTo) GetSpeed() *engosdl.Vector {
	return c.speed
}

// onOutOfBounds checks if the entity has gone out of bounds.
func (c *MoveTo) onOutOfBounds(params ...interface{}) bool {
	position := c.GetEntity().GetTransform().GetPosition()
	position.X -= c.speed.X
	position.Y -= c.speed.Y
	return true
}

// OnStart is called first time the component is enabled.
func (c *MoveTo) OnStart() {
	engosdl.Logger.Trace().Str("component", "move-to").Str("move-to", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}

// OnUpdate is called for every update tick.
func (c *MoveTo) OnUpdate() {
	position := c.GetEntity().GetTransform().GetPosition()
	position.X += c.speed.X
	position.Y += c.speed.Y
}

// SetSpeed sets movement speed.
func (c *MoveTo) SetSpeed(speed *engosdl.Vector) {
	c.speed = speed
}
