package components

import (
	"fmt"
	"math"
	"reflect"

	"github.com/jrecuero/engosdl"
)

type collisionBox struct {
	center *engosdl.Vector
	radius float64
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

// OnUpdate is called for every update tick.
func (c *Collider2D) OnUpdate() {
	transform := c.GetEntity().GetTransform()
	x := transform.GetPosition().X
	y := transform.GetPosition().Y
	w := transform.GetDim().X
	h := transform.GetDim().Y
	c.collisionBox.center = engosdl.NewVector(x, y)
	c.collisionBox.radius = math.Min(w, h) / 2
	// fmt.Printf("collision-box at (%f, %f) radius: %f\n", x, y, math.Min(w, h)/2)
	for _, entity := range c.GetEntity().GetScene().GetEntities() {
		if entity == c.GetEntity() {
			continue
		}
		for _, component := range entity.GetComponents() {
			if reflect.TypeOf(component) == reflect.TypeOf(c) {
				cx := entity.GetTransform().GetPosition().X
				cy := entity.GetTransform().GetPosition().Y
				cw := entity.GetTransform().GetDim().X
				ch := entity.GetTransform().GetDim().Y
				cradius := math.Min(cw, ch) / 2
				distance := math.Sqrt(math.Pow(x-cx, 2) + math.Pow(y-cy, 2))
				if distance < (c.collisionBox.radius + cradius) {
					fmt.Printf("%f check collision %s with %s\n", y, c.GetEntity().GetName(), entity.GetName())
					delegate := engosdl.GetEngine().GetEventHandler().GetDelegateHandler().GetCollisionDelegate()
					engosdl.GetEngine().GetEventHandler().GetDelegateHandler().TriggerDelegate(delegate, c.GetEntity(), entity)
				}
			}
		}
	}
}
