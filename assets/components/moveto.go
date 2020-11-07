package components

import "github.com/jrecuero/engosdl"

// MoveTo represents a component that moves a entity.
type MoveTo struct {
	*engosdl.Component
	speed *engosdl.Vector
}

// NewMoveTo creates a new move-to instance.
func NewMoveTo(name string, speed *engosdl.Vector) *MoveTo {
	engosdl.Logger.Trace().Str("component", "move-to").Str("move-to", name).Msg("new move-to")
	return &MoveTo{
		Component: engosdl.NewComponent(name),
		speed:     speed,
	}
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
	if c.CanRegisterTo(engosdl.OutOfBoundsName) {
		if component := c.GetEntity().GetComponent(&OutOfBounds{}); component != nil {
			if outOfBoundsComponent, ok := component.(*OutOfBounds); ok {
				if delegate := outOfBoundsComponent.GetDelegate(); delegate != nil {
					// engosdl.GetEngine().GetEventHandler().GetDelegateHandler().RegisterToDelegate(delegate, c.onOutOfBounds)
					c.AddDelegateToRegister(delegate, nil, nil, c.onOutOfBounds)
				}
			}
		}
	}
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
