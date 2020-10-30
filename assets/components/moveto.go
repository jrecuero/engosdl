package components

import "github.com/jrecuero/engosdl"

// MoveTo represents a component that moves a entity.
type MoveTo struct {
	*engosdl.Component
	speed *engosdl.Vector
}

// NewMoveTo creates a new moveto instance.
func NewMoveTo(name string, speed *engosdl.Vector) *MoveTo {
	engosdl.Logger.Trace().Str("component", "moveto").Str("moveto", name).Msg("new move-to")
	return &MoveTo{
		Component: engosdl.NewComponent(name),
		speed:     speed,
	}
}

// OnUpdate is called for every update tick.
func (c *MoveTo) OnUpdate() {
	position := c.GetEntity().GetTransform().GetPosition()
	position.X += c.speed.X
	position.Y += c.speed.Y
}
