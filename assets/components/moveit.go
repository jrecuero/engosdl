package components

import (
	"reflect"

	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// ComponentNameMoveIt is the name to refer move it component.
var ComponentNameMoveIt string = reflect.TypeOf(&MoveIt{}).String()

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor(ComponentNameMoveIt, CreateMoveIt)
	}
}

// MoveIt represents a component that can take move-it input
type MoveIt struct {
	*engosdl.Component
	Speed    *engosdl.Vector `json:"speed"`
	LastMove *engosdl.Vector `json:"last-move"`
}

// NewMoveIt creates a new move-it instance.
// It registers to "on-out-of-bounds" delegate.
func NewMoveIt(name string, speed *engosdl.Vector) *MoveIt {
	engosdl.Logger.Trace().Str("component", "move-it").Str("move-it", name).Msg("new move-it")
	result := &MoveIt{
		Component: engosdl.NewComponent(name),
		Speed:     speed,
		LastMove:  engosdl.NewVector(0, 0),
	}
	return result
}

// CreateMoveIt implements move it constructor used by the component manager.
func CreateMoveIt(params ...interface{}) engosdl.IComponent {
	if len(params) == 2 {
		return NewMoveIt(params[0].(string), params[1].(*engosdl.Vector))
	}
	return NewMoveIt("", engosdl.NewVector(0, 0))
}

// DefaultAddDelegateToRegister will proceed to add default delegate to
// register for the component.
func (c *MoveIt) DefaultAddDelegateToRegister() {
	c.AddDelegateToRegister(nil, nil, &OutOfBounds{}, c.DefaultOnOutOfBounds)
	c.AddDelegateToRegister(nil, nil, &Keyboard{}, c.onKeyboard)
}

// DefaultOnOutOfBounds checks if the entity has gone out of bounds.
func (c *MoveIt) DefaultOnOutOfBounds(params ...interface{}) bool {
	x, y := c.GetEntity().GetTransform().GetPosition().Get()
	c.GetEntity().GetTransform().SetPosition(engosdl.NewVector(x-c.LastMove.X, y-c.LastMove.Y))
	return true
}

func (c *MoveIt) onKeyboard(params ...interface{}) bool {
	position := c.GetEntity().GetTransform().GetPosition()
	key := params[0].(int)
	switch key {
	case sdl.SCANCODE_LEFT:
		position.X -= c.Speed.X
		c.LastMove.X = -1 * c.Speed.X
		c.LastMove.Y = 0
		break
	case sdl.SCANCODE_RIGHT:
		position.X += c.Speed.X
		c.LastMove.X = c.Speed.X
		c.LastMove.Y = 0
		break
	case sdl.SCANCODE_UP:
		position.Y -= c.Speed.Y
		c.LastMove.Y = -1 * c.Speed.Y
		c.LastMove.X = 0
		break
	case sdl.SCANCODE_DOWN:
		position.Y += c.Speed.Y
		c.LastMove.Y = c.Speed.Y
		c.LastMove.X = 0
		break
	}
	return true
}

// OnStart is called first time the component is enabled.
func (c *MoveIt) OnStart() {
	engosdl.Logger.Trace().Str("component", "move-to").Str("move-to", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}

// OnUpdate is called for every update tick.
func (c *MoveIt) OnUpdate() {
}

// Unmarshal takes a ComponentToMarshal instance and  creates a new entity
// instance.
func (c *MoveIt) Unmarshal(data map[string]interface{}) {
	c.Component.Unmarshal(data)
	speed := data["speed"].(map[string]interface{})
	c.Speed = engosdl.NewVector(speed["X"].(float64), speed["Y"].(float64))
	lastMove := data["last-move"].(map[string]interface{})
	c.LastMove = engosdl.NewVector(lastMove["X"].(float64), lastMove["Y"].(float64))
}
