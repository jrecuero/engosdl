package components

import (
	"fmt"

	"github.com/jrecuero/engosdl"
)

// OutOfBounds represents a component that check if entity is inside
// window bounds
type OutOfBounds struct {
	*engosdl.Component
	delegate engosdl.IDelegate
}

// NewOutOfBounds creates a new out of bounds instance
func NewOutOfBounds(name string) *OutOfBounds {
	return &OutOfBounds{
		Component: engosdl.NewComponent(name),
	}
}

// GetDelegate returns the out of bounds delegate.
func (c *OutOfBounds) GetDelegate() engosdl.IDelegate {
	return c.delegate
}

// OnAwake should create all component resources that don't have any dependency
// with any other component or entity.
func (c *OutOfBounds) OnAwake() {
	engosdl.Logger.Trace().Str("component", "out-of-bounds").Str("out-of-bounds", c.GetName()).Msg("OnAwake")
	delegateHandler := engosdl.GetEventHandler().GetDelegateHandler()
	c.delegate = delegateHandler.CreateDelegate(c, "out-of-bounds")
}

// OnStart is called the first time component is loaded.
func (c *OutOfBounds) OnStart() {
	engosdl.Logger.Trace().Str("component", "out-of-bounds").Str("out-of-bounds", c.GetName()).Msg("OnStart")
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
		engosdl.GetEventHandler().GetDelegateHandler().TriggerDelegate(c.delegate, c.GetEntity())
		// engosdl.GetEngine().DestroyEntity(c.GetEntity())

	}
	if (y+c.GetEntity().GetTransform().GetDim().Y) < 0 || int32(y) > h {
		fmt.Printf("%s out of bounds %f\n", c.GetEntity().GetName(), y)
		// c.GetEntity().GetScene().DeleteEntity(c.GetEntity())
		engosdl.GetEventHandler().GetDelegateHandler().TriggerDelegate(c.delegate, c.GetEntity())
		// engosdl.GetEngine().DestroyEntity(c.GetEntity())
	}
}
