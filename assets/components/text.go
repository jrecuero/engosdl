package components

import (
	"strconv"

	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterComponent(&Text{})
	}
}

// Text represents a component that can display some text.
type Text struct {
	*engosdl.Component
	FontFile string `json:"font-filename"`
	font     engosdl.IFont
	FontSize int       `json:"font-size"`
	Color    sdl.Color `json:"color"`
	Message  string    `json:"message"`
	renderer *sdl.Renderer
	texture  *sdl.Texture
	width    int32
	height   int32
}

var _ engosdl.IText = (*Text)(nil)

// NewText create a new text instance.
func NewText(name string, fontFile string, fontSize int, color sdl.Color, message string, renderer *sdl.Renderer) *Text {
	engosdl.Logger.Trace().Str("component", "text").Str("text", name).Msg("new text")
	return &Text{
		Component: engosdl.NewComponent(name),
		FontFile:  fontFile,
		FontSize:  fontSize,
		Color:     color,
		Message:   message,
		renderer:  renderer,
	}
}

// loadTextureFromTTF creates a texture from a ttf file.
func (c *Text) loadTextureFromTTF() {
	var err error
	c.font = engosdl.GetFontHandler().CreateFont(c.GetName(), c.FontFile, c.FontSize)
	c.texture = c.font.GetTextureFromFont(c.Message, c.Color)
	_, _, c.width, c.height, err = c.texture.Query()
	if err != nil {
		engosdl.Logger.Error().Err(err).Msg("Query error")
		panic(err)
	}
	c.GetEntity().GetTransform().SetDim(engosdl.NewVector(float64(c.width), float64(c.height)))
}

// OnRender is called for every render tick.
func (c *Text) OnRender() {
	x, y, width, height := c.GetEntity().GetTransform().GetRectExt()
	c.renderer.CopyEx(c.texture,
		&sdl.Rect{X: 0, Y: 0, W: c.width, H: c.height},
		&sdl.Rect{X: int32(x), Y: int32(y), W: int32(width), H: int32(height)},
		0,
		&sdl.Point{},
		sdl.FLIP_NONE)
}

// OnStart is called first time the component is enabled.
func (c *Text) OnStart() {
	engosdl.Logger.Trace().Str("component", "text").Str("text", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
	c.loadTextureFromTTF()
}

// onUpdateStats updates text with entity stats changes.
func (c *Text) onUpdateStats(params ...interface{}) bool {
	life := params[0].(int)
	if life == 0 {
		c.SetActive(false)
	} else {
		c.Message = "Enemy Life: " + strconv.Itoa(life)
		c.loadTextureFromTTF()
	}
	return true
}

// SetColor sets text color.
func (c *Text) SetColor(color sdl.Color) engosdl.IText {
	c.Color = color
	return c
}

// SetFontFilename sets the filename with the font.
func (c *Text) SetFontFilename(filename string) engosdl.IText {
	c.FontFile = filename
	return c
}

// SetMessage sets the message to be displayed by the text component.
func (c *Text) SetMessage(message string) engosdl.IText {
	c.Message = message
	c.loadTextureFromTTF()
	return c
}
