package main

import (
	"fmt"
	"math/rand"
	"reflect"

	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
	"github.com/veandco/go-sdl2/sdl"
)

// ComponentNameAlive is the name to refer Alive component.
var ComponentNameAlive string = reflect.TypeOf(&Alive{}).String()

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor(ComponentNameAlive, CreateAlive)
	}
}

// Alive represents a component.
type Alive struct {
	*engosdl.Component
	Tick      int `json:"tick"`
	counter   int
	Row       int
	Col       int
	Size      int32
	Color     sdl.Color
	firstCell bool
}

// var _ engosdl.IAlive = (*Alive)(nil)

// NewAlive create a new Alive instance.
func NewAlive(name string, tick int, row int, col int, size int32, color sdl.Color) *Alive {
	engosdl.Logger.Trace().Str("component", "Alive").Str("Alive", name).Msg("new Alive")
	return &Alive{
		Component: engosdl.NewComponent(name),
		Row:       row,
		Col:       col,
		Tick:      tick,
		Size:      size,
		Color:     color,
		counter:   0,
		firstCell: true,
	}
}

// CreateAlive implements Alive constructor used by component
// manager.
func CreateAlive(params ...interface{}) engosdl.IComponent {
	if len(params) == 6 {
		return NewAlive(params[0].(string), params[1].(int), params[2].(int), params[3].(int), params[4].(int32), params[5].(sdl.Color))
	}
	return NewAlive("", 0, 0, 0, 0, sdl.Color{})
}

// DefaultAddDelegateToRegister will proceed to add default delegates to
// register to the component.
func (c *Alive) DefaultAddDelegateToRegister() {
	// c.AddDelegateToRegister(<DELEGATE>, nil, <OTHER-COMPONENT>, <SIGNATURE>)
}

// DoDestroy should destroy all component resources. This is called when
// component is removed from the scene and resources are not anymore
// required.
func (c *Alive) DoDestroy() {
	engosdl.Logger.Trace().Str("component", "Alive").Str("Alive", c.GetName()).Msg("DoDestroy")
	c.Component.DoDestroy()
}

// DoUnLoad is called when component is unloaded from scene.
func (c *Alive) DoUnLoad() {
	engosdl.Logger.Trace().Str("component", "Alive").Str("Alive", c.GetName()).Msg("DoUnLoad")
	c.Component.DoUnLoad()
}

// OnAwake is called when component is first loaded into the scene and all
// component resources have to be created. No resources dependent with other
// components or entities can be created at this point.
func (c *Alive) OnAwake() {
	engosdl.Logger.Trace().Str("component", "Alive").Str("Alive", c.GetName()).Msg("OnAwake")
	c.Component.OnAwake()
}

// OnRender is called every engine frame when component has to be rendered.
func (c *Alive) OnRender() {
}

// OnStart is called at the end of the component being loaded by the scene.
// Any component resource dependent from other entities or components has
// to be created at this point.
func (c *Alive) OnStart() {
	engosdl.Logger.Trace().Str("component", "Alive").Str("Alive", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}

// OnUpdate is called every engine frame in order to update any component
// resource.
func (c *Alive) OnUpdate() {
	var row, col int
	if c.counter == c.Tick {
		scene := c.GetEntity().GetScene()
		if controller := scene.GetEntityByName("controller"); controller != nil {
			if board := controller.GetComponent(&Board{}); board != nil {
				if c.firstCell {
					row, col = c.Row, c.Col
					c.firstCell = false
				} else {
					var colValue int
					values := []int{-1, 0, 1}
					values2 := []int{-1, 1}
					iRow := rand.Intn(len(values))
					rowValue := values[iRow]
					if rowValue == 0 {
						iCol := rand.Intn(len(values2))
						colValue = values2[iCol]
					} else {
						iCol := rand.Intn(len(values))
						colValue = values[iCol]
					}
					row = c.Row + rowValue
					col = c.Col + colValue
				}
				// if board.(*Board).IsCellFree(row, col) {
				if board.(*Board).UseCell(row, col) {
					child := NewPixel(fmt.Sprintf("%s-child", c.GetEntity().GetName()))
					box := components.NewBox("child/box", &engosdl.Rect{W: float64(c.Size), H: float64(c.Size)}, c.Color, true)
					child.AddComponent(box)
					c.GetEntity().AddChild(child)
					scene.AddEntity(child)
					board.(*Board).AddEntityAt(child, row, col)
					pos, _ := board.(*Board).GetPositionFromCell(row, col)
					child.GetTransform().SetPosition(pos)
					c.Row = row
					c.Col = col
					fmt.Printf("new pixel at: %d, %d\n", row, col)
					c.counter = 0
				}
			}
		}
	} else {
		c.counter++
	}
	c.Component.OnUpdate()
}

// Unmarshal takes information from a ComponentToUnmarshal instance and
// creates a new component instance.
func (c *Alive) Unmarshal(data map[string]interface{}) {
	c.Component.Unmarshal(data)
	c.Tick = int(data["tick"].(float64))
}
