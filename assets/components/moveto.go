package components

import "github.com/jrecuero/engosdl"

// MoveTo represents a component that moves a entity.
type MoveTo struct {
	*engosdl.Component
	speed *engosdl.Vector
}

// OnUpdate is called for every update tick.
func (m *MoveTo) OnUpdate() {
	position := m.GetEntity().GetTransform().GetPosition()
	position.X += m.speed.X
	position.Y += m.speed.Y
}

// NewMoveTo creates a new moveto instance.
func NewMoveTo(name string, entity *engosdl.Entity, speed *engosdl.Vector) *MoveTo {
	engosdl.Logger.Trace().Str("component", "moveto").Str("moveto", name).Msg("new move-to")
	return &MoveTo{
		Component: engosdl.NewComponent(name, entity),
		speed:     speed,
	}
}
