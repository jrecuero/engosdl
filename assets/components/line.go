package components

import (
	"reflect"

	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// ComponentNameLine is the name to refer line it component.
var ComponentNameLine string = reflect.TypeOf(&Line{}).String()

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor(ComponentNameLine, CreateLine)
	}
}

// Line represents a component.
type Line struct {
	*engosdl.Component
	renderer *sdl.Renderer
	endPoint *engosdl.Vector
	color    sdl.Color
}

// var _ engosdl.ILine = (*Line)(nil)

// NewLine creates a new Line instance.
func NewLine(name string, endPoint *engosdl.Vector, color sdl.Color) *Line {
	engosdl.Logger.Trace().Str("component", "Line").Str("Line", name).Msg("new Line")
	return &Line{
		Component: engosdl.NewComponent(name),
		endPoint:  endPoint,
		color:     color,
	}
}

// CreateLine implements Line constructor used by component
// manager.
func CreateLine(params ...interface{}) engosdl.IComponent {
	if len(params) == 3 {
		return NewLine(params[0].(string), params[1].(*engosdl.Vector), params[2].(sdl.Color))
	}
	return NewLine("", nil, sdl.Color{})
}

// DefaultAddDelegateToRegister will proceed to add default delegates to
// register to the component.
func (c *Line) DefaultAddDelegateToRegister() {
	// c.AddDelegateToRegister(<DELEGATE>, nil, <OTHER-Line>, <SIGNATURE>)
}

// DoDestroy should destroy all component resources. This is called when
// component is removed from the scene and resources are not anymore
// required.
func (c *Line) DoDestroy() {
	engosdl.Logger.Trace().Str("component", "Line").Str("Line", c.GetName()).Msg("DoDestroy")
	c.Component.DoDestroy()
}

// DoUnLoad is called when component is unloaded from scene.
func (c *Line) DoUnLoad() {
	engosdl.Logger.Trace().Str("component", "Line").Str("Line", c.GetName()).Msg("DoUnLoad")
	c.Component.DoUnLoad()
}

// OnAwake is called when component is first loaded into the scene and all
// component resources have to be created. No resources dependent with other
// components or entities can be created at this point.
func (c *Line) OnAwake() {
	engosdl.Logger.Trace().Str("component", "Line").Str("Line", c.GetName()).Msg("OnAwake")
	c.GetEntity().GetTransform().SetDim(engosdl.NewVector(float64(c.endPoint.X), float64(c.endPoint.Y)))
	c.Component.OnAwake()
}

// OnRender is called every engine frame when component has to be rendered.
func (c *Line) OnRender() {
	x, y, w, h := c.GetEntity().GetTransform().GetRectExt()
	c.renderer.SetDrawColor(c.color.R, c.color.G, c.color.B, c.color.A)
	c.renderer.DrawLine(int32(x), int32(y), int32(w), int32(h))
	c.renderer.DrawLine(int32(x), int32(y), int32(w+1), int32(h+1))
	c.renderer.DrawLine(int32(x), int32(y), int32(w+2), int32(h+2))
	c.renderer.DrawLine(int32(x), int32(y), int32(w+3), int32(h+3))
}

// OnStart is called at the end of the component being loaded by the scene.
// Any component resource dependent from other entities or components has
// to be created at this point.
func (c *Line) OnStart() {
	engosdl.Logger.Trace().Str("component", "Line").Str("Line", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}

// OnUpdate is called every engine frame in order to update any component
// resource.
func (c *Line) OnUpdate() {
	c.Component.OnUpdate()
}

// Unmarshal takes information from a ComponentToUnmarshal instance and
// creates a new component instance.
func (c *Line) Unmarshal(data map[string]interface{}) {
	c.Component.Unmarshal(data)
}
