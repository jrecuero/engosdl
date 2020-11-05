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
	w := engosdl.GetEngine().GetWidth()
	h := engosdl.GetEngine().GetHeight()
	x := c.GetEntity().GetTransform().GetPosition().X
	y := c.GetEntity().GetTransform().GetPosition().Y
	var testX, testY bool
	if c.leftCorner {
		testX = x < 0 || int32(x+c.GetEntity().GetTransform().GetDim().X) > w
		testY = y < 0 || int32(y+c.GetEntity().GetTransform().GetDim().Y) > h
	} else {
		testX = (x+c.GetEntity().GetTransform().GetDim().X) < 0 || int32(x) > w
		testY = (y+c.GetEntity().GetTransform().GetDim().Y) < 0 || int32(y) > h
	}
	if testX {
		fmt.Printf("%s out of bounds %f\n", c.GetEntity().GetName(), x)
		// c.GetEntity().GetScene().DeleteEntity(c.GetEntity())
		engosdl.GetEventHandler().GetDelegateHandler().TriggerDelegate(c.GetDelegate(), c.GetEntity())
		// engosdl.GetEngine().DestroyEntity(c.GetEntity())

	}
	if testY {
		fmt.Printf("%s out of bounds %f\n", c.GetEntity().GetName(), y)
		// c.GetEntity().GetScene().DeleteEntity(c.GetEntity())
		engosdl.GetEventHandler().GetDelegateHandler().TriggerDelegate(c.GetDelegate(), c.GetEntity())
		// engosdl.GetEngine().DestroyEntity(c.GetEntity())
	}
}
