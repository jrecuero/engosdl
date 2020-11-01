package components

import (
	"math"

	"github.com/jrecuero/engosdl"
)

type collisionBox struct {
	center *engosdl.Vector
	radius float64
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
		Component: engosdl.NewComponent(name),
		collisionBox: &collisionBox{
			center: engosdl.NewVector(0, 0),
			radius: 0,
		},
	}
}

// GetCollisionBox returns the collision box for the parent entity.
func (c *Collider2D) GetCollisionBox() engosdl.ICollisionBox {
	x, y := c.GetEntity().GetTransform().GetPosition().Get()
	w, h := c.GetEntity().GetTransform().GetDim().Get()
	c.collisionBox.center = engosdl.NewVector(x+w/2, y+h/y)
	c.collisionBox.radius = (math.Min(w, h) / 2) * 0.9
	return c.collisionBox
}

// OnUpdate is called for every update tick.
func (c *Collider2D) OnUpdate() {
	// x, y := c.GetEntity().GetTransform().GetPosition().Get()
	// w, h := c.GetEntity().GetTransform().GetDim().Get()
	// c.collisionBox.center = engosdl.NewVector(x, y)
	// c.collisionBox.radius = math.Min(w, h) / 2
	// // fmt.Printf("collision-box at (%f, %f) radius: %f\n", x, y, math.Min(w, h)/2)
	// for _, entity := range c.GetEntity().GetScene().GetEntities() {
	// 	if entity == c.GetEntity() {
	// 		continue
	// 	}
	// 	for _, component := range entity.GetComponents() {
	// 		if reflect.TypeOf(component) == reflect.TypeOf(c) {
	// 			cx, cy := entity.GetTransform().GetPosition().Get()
	// 			cw, ch := entity.GetTransform().GetDim().Get()
	// 			cradius := math.Min(cw, ch) / 2
	// 			distance := math.Sqrt(math.Pow(x-cx, 2) + math.Pow(y-cy, 2))
	// 			if distance < (c.collisionBox.radius + cradius) {
	// 				// fmt.Printf("%f check collision %s with %s\n", y, c.GetEntity().GetName(), entity.GetName())
	// 				// delegate := engosdl.GetEngine().GetEventHandler().GetDelegateHandler().GetCollisionDelegate()
	// 				// engosdl.GetEngine().GetEventHandler().GetDelegateHandler().TriggerDelegate(delegate, c.GetEntity(), entity)
	// 			}
	// 		}
	// 	}
	// }
}
