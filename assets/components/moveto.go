package components

import "github.com/jrecuero/engosdl"

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterComponent(&MoveTo{})
	}
}

// MoveTo represents a component that moves a entity.
type MoveTo struct {
	*engosdl.Component
	Speed *engosdl.Vector `json:"speed"`
}

// NewMoveTo creates a new move-to instance.
// It registers to "on-out-of-bounds" delegate.
func NewMoveTo(name string, speed *engosdl.Vector) *MoveTo {
	engosdl.Logger.Trace().Str("component", "move-to").Str("move-to", name).Msg("new move-to")
	result := &MoveTo{
		Component: engosdl.NewComponent(name),
		Speed:     speed,
	}
	return result
}

// DefaultAddDelegateToRegister will proceed to add default delegate to
// register for the component.
func (c *MoveTo) DefaultAddDelegateToRegister() {
	c.AddDelegateToRegister(nil, nil, &OutOfBounds{}, c.onOutOfBounds)
}

// GetSpeed returns component speed.
func (c *MoveTo) GetSpeed() *engosdl.Vector {
	return c.Speed
}

// onOutOfBounds checks if the entity has gone out of bounds.
func (c *MoveTo) onOutOfBounds(params ...interface{}) bool {
	position := c.GetEntity().GetTransform().GetPosition()
	position.X -= c.Speed.X
	position.Y -= c.Speed.Y
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
	position.X += c.Speed.X
	position.Y += c.Speed.Y
}

// SetSpeed sets movement speed.
func (c *MoveTo) SetSpeed(speed *engosdl.Vector) {
	c.Speed = speed
}
