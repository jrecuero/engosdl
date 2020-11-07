package components

import (
	"fmt"

	"github.com/jrecuero/engosdl"
)

// OutOfBounds represents a component that check if entity is inside
// window bounds
type OutOfBounds struct {
	*engosdl.Component
	leftCorner bool
}

// NewOutOfBounds creates a new out of bounds instance
func NewOutOfBounds(name string, leftCorner bool) *OutOfBounds {
	return &OutOfBounds{
		Component:  engosdl.NewComponent(name),
		leftCorner: leftCorner,
	}
}

// OnAwake should create all component resources that don't have any dependency
// with any other component or entity.
func (c *OutOfBounds) OnAwake() {
	engosdl.Logger.Trace().Str("component", "out-of-bounds").Str("out-of-bounds", c.GetName()).Msg("OnAwake")
	// Creates new delegate "out-of-bounds"
	c.SetDelegate(engosdl.GetEngine().GetEventHandler().GetDelegateHandler().CreateDelegate(c, "out-of-bounds"))
}

// OnStart is called the first time component is loaded.
func (c *OutOfBounds) OnStart() {
	engosdl.Logger.Trace().Str("component", "out-of-bounds").Str("out-of-bounds", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}

// OnUpdate is called for every update tick.
func (c *OutOfBounds) OnUpdate() {
	W := engosdl.GetEngine().GetWidth()
	H := engosdl.GetEngine().GetHeight()
	// x := c.GetEntity().GetTransform().GetPosition().X
	// y := c.GetEntity().GetTransform().GetPosition().Y
	// rect := c.GetEntity().GetTransform().GetRect()
	// x := float64(rect.X)
	// y := float64(rect.Y)
	// w := float64(rect.W)
	// h := float64(rect.H)
	x, y, w, h := c.GetEntity().GetTransform().GetRectExt()
	var testX, testY bool
	if c.leftCorner {
		testX = x < 0 || int32(x+w) > W
		testY = y < 0 || int32(y+h) > H
	} else {
		testX = (x+w) < 0 || int32(x) > W
		testY = (y+h) < 0 || int32(y) > H
	}
	if testX {
		fmt.Printf("[OutOfBounds] %s out of bounds %f\n", c.GetEntity().GetName(), x)
		// c.GetEntity().GetScene().DeleteEntity(c.GetEntity())
		engosdl.GetEventHandler().GetDelegateHandler().TriggerDelegate(c.GetDelegate(), true, c.GetEntity())
		// engosdl.GetEngine().DestroyEntity(c.GetEntity())

	}
	if testY {
		fmt.Printf("[OutOfBounds] %s out of bounds %f\n", c.GetEntity().GetName(), y)
		// c.GetEntity().GetScene().DeleteEntity(c.GetEntity())
		engosdl.GetEventHandler().GetDelegateHandler().TriggerDelegate(c.GetDelegate(), true, c.GetEntity())
		// engosdl.GetEngine().DestroyEntity(c.GetEntity())
	}
}
