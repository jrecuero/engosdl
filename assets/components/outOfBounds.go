package components

import (
	"fmt"

	"github.com/jrecuero/engosdl"
)

// OutOfBounds represents a component that check if game object is inside
// window bounds
type OutOfBounds struct {
	*engosdl.Component
}

// OnUpdate is called for every update tick.
func (c *OutOfBounds) OnUpdate() {
	w := engosdl.GetEngine().GetWidth()
	h := engosdl.GetEngine().GetHeight()
	x := c.GetGameObject().GetTransform().GetPosition().X
	y := c.GetGameObject().GetTransform().GetPosition().Y
	if (x+c.GetGameObject().GetTransform().GetDim().X) < 0 || int32(x) > w {
		fmt.Printf("%s out of bounds %f\n", c.GetGameObject().GetName(), x)
		c.GetGameObject().GetScene().DeleteGameObject(c.GetGameObject())
	}
	if (y+c.GetGameObject().GetTransform().GetDim().Y) < 0 || int32(y) > h {
		fmt.Printf("%s out of bounds %f\n", c.GetGameObject().GetName(), y)
		c.GetGameObject().GetScene().DeleteGameObject(c.GetGameObject())
	}
}

// NewOutOfBounds creates a new out of bounds instance
func NewOutOfBounds(name string, gobj *engosdl.GameObject) *OutOfBounds {
	return &OutOfBounds{
		Component: engosdl.NewComponent(name, gobj),
	}
}
