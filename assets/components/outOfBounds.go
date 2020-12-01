package components

import (
	"reflect"

	"github.com/jrecuero/engosdl"
)

// ComponentNameOutOfBounds is the name to refer out of bounds component.
var ComponentNameOutOfBounds string = reflect.TypeOf(&OutOfBounds{}).String()

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor(ComponentNameOutOfBounds, CreateOutOfBounds)
	}
}

// OutOfBounds represents a component that check if entity is inside
// window bounds
type OutOfBounds struct {
	*engosdl.Component
	LeftCorner bool `json:"left-corner"`
}

// NewOutOfBounds creates a new out of bounds instance
// It creates delegate "on-out-of-bounds"
func NewOutOfBounds(name string, leftCorner bool) *OutOfBounds {
	return &OutOfBounds{
		Component:  engosdl.NewComponent(name),
		LeftCorner: leftCorner,
	}
}

// CreateOutOfBounds implements out of bounds constructor used by component
// manager.
func CreateOutOfBounds(params ...interface{}) engosdl.IComponent {
	if len(params) == 2 {
		return NewOutOfBounds(params[0].(string), params[1].(bool))
	}
	return NewOutOfBounds("", true)
}

// OnAwake should create all component resources that don't have any dependency
// with any other component or entity.
func (c *OutOfBounds) OnAwake() {
	engosdl.Logger.Trace().Str("component", "out-of-bounds").Str("out-of-bounds", c.GetName()).Msg("OnAwake")
	// Creates new delegate "out-of-bounds"
	c.SetDelegate(engosdl.GetDelegateManager().CreateDelegate(c, "on-out-of-bounds"))
	c.Component.OnAwake()
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
	x, y, w, h := c.GetEntity().GetTransform().GetRectExt()
	var testLeft, testRight, testUp, testDown bool
	if c.LeftCorner {
		// testX = x < 0 || int32(x+w) > W
		// testY = y < 0 || int32(y+h) > H
		testLeft = x < 0
		testRight = int32(x+w) > W
		testUp = y < 0
		testDown = int32(y+h) > H
	} else {
		testLeft = (x + w) < 0
		testRight = int32(x) > W
		testUp = (y + h) < 0
		testDown = int32(y) > H
	}
	if testLeft {
		// fmt.Printf("[OutOfBounds] %s out of bounds %f\n", c.GetEntity().GetName(), x)
		engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), true, c.GetEntity(), engosdl.Left)
	}
	if testRight {
		// fmt.Printf("[OutOfBounds] %s out of bounds %f\n", c.GetEntity().GetName(), x)
		engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), true, c.GetEntity(), engosdl.Right)
	}
	if testUp {
		// fmt.Printf("[OutOfBounds] %s out of bounds %f\n", c.GetEntity().GetName(), y)
		engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), true, c.GetEntity(), engosdl.Up)
	}
	if testDown {
		// fmt.Printf("[OutOfBounds] %s out of bounds %f\n", c.GetEntity().GetName(), y)
		engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), true, c.GetEntity(), engosdl.Down)
	}
}

// Unmarshal takes a ComponentToMarshal instance and  creates a new entity
// instance.
func (c *OutOfBounds) Unmarshal(data map[string]interface{}) {
	c.Component.Unmarshal(data)
	c.LeftCorner = data["left-corner"].(bool)
}
