package components

import (
	"math"
	"reflect"

	"github.com/jrecuero/engosdl"
)

// ComponentNameBody is the name to refer body component.
var ComponentNameBody string = reflect.TypeOf(&Body{}).String()

// Body events
const (
	// BodyEvOutOfBounds is the event when body gets out of bounds.
	BodyEvOutOfBounds int = 1
)

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor(ComponentNameBody, CreateBody)
	}
}

// Body represents a component for any physical body component.
type Body struct {
	*engosdl.Component
	collisionBox  *CollisionBox
	AllBodyForOOB bool `json:"all-body-for-oob"`
}

// NewBody creates a new body instance.
// It creates delegate "body"
func NewBody(name string, origin bool) *Body {
	engosdl.Logger.Trace().Str("component", "body").Str("body", name).Msg("new body")
	return &Body{
		Component:     engosdl.NewComponent(name),
		AllBodyForOOB: origin,
		collisionBox:  &CollisionBox{},
	}
}

// CreateBody implements body constructor used by component manager.
// It creates delegate "body"
func CreateBody(params ...interface{}) engosdl.IComponent {
	if len(params) == 2 {
		return NewBody(params[0].(string), params[1].(bool))
	}
	return NewBody("", false)
}

// GetCollisionBox returns the collision box for the parent entity.
func (c *Body) GetCollisionBox() engosdl.ICollisionBox {
	x, y, w, h := c.GetEntity().GetTransform().GetRectExt()
	c.collisionBox.center = engosdl.NewVector(x+(w/2), y+(h/2))
	// Set collision box radius as 75% of the minimum radius.
	c.collisionBox.radius = (math.Min(w, h) / 2) * 0.75
	c.collisionBox.rect = &engosdl.Rect{X: x, Y: y, W: w, H: h}
	return c.collisionBox
}

// OnAwake should create all component resources that don't have any dependency
// with any other component or entity.
// It creates delegate "body"
func (c *Body) OnAwake() {
	engosdl.Logger.Trace().Str("component", "body").Str("body", c.GetName()).Msg("OnAwake")
	c.SetDelegate(engosdl.GetDelegateManager().CreateDelegate(c, "body"))
	c.Component.OnAwake()
}

// OnUpdate is called for every update tick.
func (c *Body) OnUpdate() {
	W := engosdl.GetEngine().GetWidth()
	H := engosdl.GetEngine().GetHeight()
	x, y, w, h := c.GetEntity().GetTransform().GetRectExt()
	var testLeft, testRight, testUp, testDown bool
	if c.AllBodyForOOB {
		testLeft = (x + w) < 0
		testRight = int32(x) > W
		testUp = (y + h) < 0
		testDown = int32(y) > H
	} else {
		testLeft = x < 0
		testRight = int32(x+w) > W
		testUp = y < 0
		testDown = int32(y+h) > H
	}
	if testLeft {
		engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), true, c.GetEntity(), BodyEvOutOfBounds, engosdl.Left)
	}
	if testRight {
		engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), true, c.GetEntity(), BodyEvOutOfBounds, engosdl.Right)
	}
	if testUp {
		engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), true, c.GetEntity(), BodyEvOutOfBounds, engosdl.Up)
	}
	if testDown {
		engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), true, c.GetEntity(), BodyEvOutOfBounds, engosdl.Down)
	}
}

// Unmarshal takes a ComponentToMarshal instance and  creates a new entity
// instance.
func (c *Body) Unmarshal(data map[string]interface{}) {
	c.Component.Unmarshal(data)
	c.AllBodyForOOB = data["all-body-for-oob"].(bool)
}
