package components

import (
	"math"
	"reflect"

	"github.com/jrecuero/engosdl"
)

// ComponentNameCollider2D is the name to refer collider-2 component.
var ComponentNameCollider2D string = reflect.TypeOf(&Collider2D{}).String()

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor(ComponentNameCollider2D, CreateCollider2D)
	}
}

// CollisionBox identify the model to be used for collisions.
type CollisionBox struct {
	rect   *engosdl.Rect
	center *engosdl.Vector
	radius float64
}

// GetRect returns collision box rectangle.
func (c *CollisionBox) GetRect() *engosdl.Rect {
	return c.rect
}

// GetCenter returns collision box center point as a vector.
func (c *CollisionBox) GetCenter() *engosdl.Vector {
	return c.center
}

// GetRadius returns collision box radius.
func (c *CollisionBox) GetRadius() float64 {
	return c.radius
}

// Collider2D represents a component that check for 2D collisions.
type Collider2D struct {
	*engosdl.Component
	collisionBox *CollisionBox
}

// NewCollider2D create a new collider-2D instance.
func NewCollider2D(name string) *Collider2D {
	return &Collider2D{
		Component:    engosdl.NewComponent(name),
		collisionBox: &CollisionBox{},
	}
}

// CreateCollider2D implements collider-2D constructor used by component
// manager
func CreateCollider2D(params ...interface{}) engosdl.IComponent {
	if len(params) == 1 {
		return NewCollider2D(params[0].(string))
	}
	return NewCollider2D("")
}

// GetCollisionBox returns the collision box for the parent entity.
func (c *Collider2D) GetCollisionBox() engosdl.ICollisionBox {
	x, y, w, h := c.GetEntity().GetTransform().GetRectExt()
	c.collisionBox.center = engosdl.NewVector(x+(w/2), y+(h/2))
	// Set collision box radius as 75% of the minimum radius.
	c.collisionBox.radius = (math.Min(w, h) / 2) * 0.75
	c.collisionBox.rect = &engosdl.Rect{X: x, Y: y, W: w, H: h}
	return c.collisionBox
}

// OnUpdate is called for every update tick.
func (c *Collider2D) OnUpdate() {
}
