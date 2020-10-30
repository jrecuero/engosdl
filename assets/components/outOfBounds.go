package components

import (
	"fmt"

	"github.com/jrecuero/engosdl"
)

// OutOfBounds represents a component that check if entity is inside
// window bounds
type OutOfBounds struct {
	*engosdl.Component
}

// OnUpdate is called for every update tick.
func (c *OutOfBounds) OnUpdate() {
	w := engosdl.GetEngine().GetWidth()
	h := engosdl.GetEngine().GetHeight()
	x := c.GetEntity().GetTransform().GetPosition().X
	y := c.GetEntity().GetTransform().GetPosition().Y
	if (x+c.GetEntity().GetTransform().GetDim().X) < 0 || int32(x) > w {
		fmt.Printf("%s out of bounds %f\n", c.GetEntity().GetName(), x)
		// c.GetEntity().GetScene().DeleteEntity(c.GetEntity())
		engosdl.GetEngine().DestroyEntity(c.GetEntity())
	}
	if (y+c.GetEntity().GetTransform().GetDim().Y) < 0 || int32(y) > h {
		fmt.Printf("%s out of bounds %f\n", c.GetEntity().GetName(), y)
		// c.GetEntity().GetScene().DeleteEntity(c.GetEntity())
		engosdl.GetEngine().DestroyEntity(c.GetEntity())
	}
}

// NewOutOfBounds creates a new out of bounds instance
func NewOutOfBounds(name string, entity *engosdl.Entity) *OutOfBounds {
	return &OutOfBounds{
		Component: engosdl.NewComponent(name, entity),
	}
}
