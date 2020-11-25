package components

import (
	"reflect"

	"github.com/jrecuero/engosdl"
)

// ComponentNameMove is the name to refer move to component.
var ComponentNameMove string = reflect.TypeOf(&Move{}).String()

// NextMoveT is the signature to be used in order to provide the next position
// for the entity.
type NextMoveT func(engosdl.IComponent) *engosdl.Vector

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor(ComponentNameMove, CreateMove)
	}
}

// Move represents a component that moves a entity.
type Move struct {
	*engosdl.Component
	nextMove NextMoveT
	lastMove *engosdl.Vector
}

// NewMove creates a new Move instance.
// It registers to "on-out-of-bounds" delegate.
func NewMove(name string, nextMove NextMoveT) *Move {
	engosdl.Logger.Trace().Str("component", "Move").Str("Move", name).Msg("new Move")
	result := &Move{
		Component: engosdl.NewComponent(name),
		nextMove:  nextMove,
		lastMove:  engosdl.NewVector(0, 0),
	}
	return result
}

// CreateMove implements move to constructor used by component manager.
func CreateMove(params ...interface{}) engosdl.IComponent {
	if len(params) == 1 {
		return NewMove(params[0].(string), nil)
	}
	return NewMove("", nil)
}

// DefaultAddDelegateToRegister will proceed to add default delegate to
// register for the component.
func (c *Move) DefaultAddDelegateToRegister() {
	c.AddDelegateToRegister(nil, nil, &OutOfBounds{}, c.DefaultOnOutOfBounds)
}

// DefaultOnOutOfBounds checks if the entity has gone out of bounds.
func (c *Move) DefaultOnOutOfBounds(params ...interface{}) bool {
	position := c.GetEntity().GetTransform().GetPosition()
	position.X -= c.lastMove.X
	position.Y -= c.lastMove.Y
	return true
}

// OnStart is called first time the component is enabled.
func (c *Move) OnStart() {
	engosdl.Logger.Trace().Str("component", "Move").Str("Move", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}

// OnUpdate is called for every update tick.
func (c *Move) OnUpdate() {
	position := c.GetEntity().GetTransform().GetPosition()
	c.lastMove = c.nextMove(c)
	position.X += c.lastMove.X
	position.Y += c.lastMove.Y
}

// Unmarshal takes a ComponentToMarshal instance and  creates a new entity
// instance.
func (c *Move) Unmarshal(data map[string]interface{}) {
	c.Component.Unmarshal(data)
}
