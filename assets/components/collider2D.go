package components

import (
	"math"

	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

type collisionBox struct {
	rect   *sdl.Rect
	center *engosdl.Vector
	radius float64
}

// GetRect returns collision box rectangle.
func (c *collisionBox) GetRect() *sdl.Rect {
	return c.rect
}

// GetCenter returns collision box center point as a vector.
func (c *collisionBox) GetCenter() *engosdl.Vector {
	return c.center
}

// GetRadius returns collision box radius.
func (c *collisionBox) GetRadius() float64 {
	return c.radius
}

// Collider2D represents a component that check for 2D collisions.
type Collider2D struct {
	*engosdl.Component
	collisionBox *collisionBox
}

// NewCollider2D create a new collider-2D instance.
func NewCollider2D(name string) *Collider2D {
	return &Collider2D{
		Component:    engosdl.NewComponent(name),
		collisionBox: &collisionBox{
			// center: engosdl.NewVector(0, 0),
			// radius: 0,
		},
	}
}

// GetCollisionBox returns the collision box for the parent entity.
func (c *Collider2D) GetCollisionBox() engosdl.ICollisionBox {
	// x, y := c.GetEntity().GetTransform().GetPosition().Get()
	// w, h := c.GetEntity().GetTransform().GetDim().Get()
	// rect := c.GetEntity().GetTransform().GetRect()
	// x := float64(rect.X)
	// y := float64(rect.Y)
	// w := float64(rect.W)
	// h := float64(rect.H)
	x, y, w, h := c.GetEntity().GetTransform().GetRectExt()
	c.collisionBox.center = engosdl.NewVector(x+(w/2), y+(h/2))
	// Set collision box radius as 75% of the minimum radius.
	c.collisionBox.radius = (math.Min(w, h) / 2) * 0.75
	c.collisionBox.rect = &sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}
	return c.collisionBox
}

// OnUpdate is called for every update tick.
func (c *Collider2D) OnUpdate() {
}
