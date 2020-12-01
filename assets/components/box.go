package components

import (
	"reflect"

	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// ComponentNameBox is the name to refer box component.
var ComponentNameBox string = reflect.TypeOf(&Box{}).String()

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor(ComponentNameBox, CreateBox)
	}
}

// Box represents a component that display a rectangle.
type Box struct {
	*engosdl.Component
	renderer *sdl.Renderer
	box      *sdl.Rect
	color    sdl.Color
	filled   bool
}

// NewBox create a new box instance.
func NewBox(name string, box *sdl.Rect, color sdl.Color, filled bool) *Box {
	engosdl.Logger.Trace().Str("component", "box").Str("box", name).Msg("new box")
	return &Box{
		Component: engosdl.NewComponent(name),
		renderer:  engosdl.GetRenderer(),
		box:       box,
		color:     color,
		filled:    filled,
	}
}

// CreateBox implements box constructor used by component manager.
func CreateBox(params ...interface{}) engosdl.IComponent {
	if len(params) == 4 {
		return NewBox(params[0].(string), params[1].(*sdl.Rect), params[2].(sdl.Color), params[3].(bool))
	}
	return NewBox("", nil, sdl.Color{}, false)
}

// DefaultAddDelegateToRegister will proceed to add default delegate to
// register for the component.
// It register to "collision" delegate.
// It register to "out-of-bounds" delegate.
func (c *Box) DefaultAddDelegateToRegister() {
	// c.AddDelegateToRegister(engosdl.GetDelegateManager().GetCollisionDelegate(), nil, nil, c.DefaultOnCollision)
	c.AddDelegateToRegister(nil, nil, &OutOfBounds{}, c.DefaultOnOutOfBounds)
}

// // DefaultOnCollision checks when there is a collision with other entity.
// func (c *Box) DefaultOnCollision(params ...interface{}) bool {
// 	collisionEntityOne := params[0].(*engosdl.Entity)
// 	collisionEntityTwo := params[1].(*engosdl.Entity)
// 	if c.GetEntity().GetID() == collisionEntityOne.GetID() || c.GetEntity().GetID() == collisionEntityTwo.GetID() {
// 		fmt.Printf("%s sprite onCollision %s with %s\n", c.GetEntity().GetName(), collisionEntityOne.GetName(), collisionEntityTwo.GetName())
// 		if collisionEntityOne.GetDieOnCollision() {
// 			engosdl.GetEngine().DestroyEntity(collisionEntityOne)
// 		}
// 		if collisionEntityTwo.GetDieOnCollision() {
// 			engosdl.GetEngine().DestroyEntity(collisionEntityTwo)
// 		}
// 	}
// 	return true
// }

// DefaultOnOutOfBounds checks if the entity has gone out of bounds.
func (c *Box) DefaultOnOutOfBounds(params ...interface{}) bool {
	if c.GetEntity().GetDieOnOutOfBounds() {
		entity := params[0].(engosdl.IEntity)
		if entity.GetID() == c.GetEntity().GetID() {
			engosdl.GetEngine().DestroyEntity(c.GetEntity())
		}
	}
	return true
}

// OnAwake should create all component resources that don't have any dependency
// with any other component or entity.
func (c *Box) OnAwake() {
	engosdl.Logger.Trace().Str("component", "box").Str("box", c.GetName()).Msg("OnAwake")
	c.GetEntity().GetTransform().SetDim(engosdl.NewVector(float64(c.box.W), float64(c.box.H)))
	c.Component.OnAwake()
}

// OnRender is called every engine frame in order to render component.
func (c *Box) OnRender() {
	x, y, w, h := c.GetEntity().GetTransform().GetRectExt()
	c.renderer.SetDrawColor(c.color.R, c.color.G, c.color.B, c.color.A)
	if c.filled {
		c.renderer.FillRect(&sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)})
	} else {
		c.renderer.DrawRect(&sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)})
	}
}
