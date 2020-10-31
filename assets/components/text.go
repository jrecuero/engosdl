package components

import (
	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Text represents a component that can display some text.
type Text struct {
	*engosdl.Component
	fontFile string
	font     *ttf.Font
	fontSize int
	color    sdl.Color
	message  string
	renderer *sdl.Renderer
	texture  *sdl.Texture
	width    int32
	height   int32
}

// NewText create a new text instance.
func NewText(name string, fontFile string, fontSize int, color sdl.Color, message string, renderer *sdl.Renderer) *Text {
	engosdl.Logger.Trace().Str("component", "text").Str("text", name).Msg("new text")
	return &Text{
		Component: engosdl.NewComponent(name),
		fontFile:  fontFile,
		fontSize:  fontSize,
		color:     color,
		message:   message,
		renderer:  renderer,
	}
}

// OnStart is called first time the component is enabled.
func (c *Text) OnStart() {
	engosdl.Logger.Trace().Str("component", "text").Str("text", c.GetName()).Msg("OnStart")
	c.texttureFromTTF()
}

// OnDraw is called for every draw tick.
func (c *Text) OnDraw() {
	x := int32(c.GetEntity().GetTransform().GetPosition().X)
	y := int32(c.GetEntity().GetTransform().GetPosition().Y)
	width := c.width * int32(c.GetEntity().GetTransform().GetScale().X)
	height := c.height * int32(c.GetEntity().GetTransform().GetScale().Y)
	c.renderer.CopyEx(c.texture,
		&sdl.Rect{X: 0, Y: 0, W: c.width, H: c.height},
		&sdl.Rect{X: x, Y: y, W: width, H: height},
		0,
		&sdl.Point{},
		sdl.FLIP_NONE)
}

// texttureFromTTF creates a textture from a ttf file.
func (c *Text) texttureFromTTF() {
	var err error
	c.font, err = ttf.OpenFont(c.fontFile, c.fontSize)
	if err != nil {
		engosdl.Logger.Error().Err(err)
		panic(err)
	}
	surface, err := c.font.RenderUTF8Solid(c.message, c.color)
	if err != nil {
		engosdl.Logger.Error().Err(err)
		panic(err)
	}
	defer surface.Free()

	c.texture, err = c.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		engosdl.Logger.Error().Err(err)
		panic(err)
	}
	_, _, c.width, c.height, err = c.texture.Query()
	if err != nil {
		engosdl.Logger.Error().Err(err)
		panic(err)
	}
	c.GetEntity().GetTransform().SetDim(engosdl.NewVector(float64(c.width), float64(c.height)))
}
