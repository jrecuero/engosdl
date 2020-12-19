package components

import (
	"reflect"

	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// ComponentNameButton is the name to refer text component.
var ComponentNameButton string = reflect.TypeOf(&Button{}).String()

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor(ComponentNameButton, CreateButton)
	}
}

// Button represents a component that can display some text.
type Button struct {
	*engosdl.Component
	FontFile    string `json:"font-filename"`
	font        engosdl.IFont
	FontSize    int       `json:"font-size"`
	Color       sdl.Color `json:"color"`
	Message     string    `json:"message"`
	renderer    *sdl.Renderer
	texture     *sdl.Texture
	width       int32
	height      int32
	border      *engosdl.Rect
	borderColor sdl.Color
	filled      bool
}

var _ engosdl.IButton = (*Button)(nil)

// NewButton create a new text instance.
func NewButton(name string, fontFile string, fontSize int, color sdl.Color, message string,
	border *engosdl.Rect, borderColor sdl.Color, filled bool) *Button {
	engosdl.Logger.Trace().Str("component", "text").Str("text", name).Msg("new text")
	return &Button{
		Component:   engosdl.NewComponent(name),
		FontFile:    fontFile,
		FontSize:    fontSize,
		Color:       color,
		Message:     message,
		renderer:    engosdl.GetRenderer(),
		border:      border,
		borderColor: borderColor,
		filled:      filled,
	}
}

// CreateButton implements text constructor used by component manager.
func CreateButton(params ...interface{}) engosdl.IComponent {
	if len(params) == 8 {
		return NewButton(params[0].(string), params[1].(string), params[2].(int), params[3].(sdl.Color), params[4].(string),
			params[5].(*engosdl.Rect), params[6].(sdl.Color), params[7].(bool))
	}
	return NewButton("", "", 0, sdl.Color{}, "", &engosdl.Rect{}, sdl.Color{}, false)
}

// loadTextureFromTTF creates a texture from a ttf file.
func (c *Button) loadTextureFromTTF() {
	var err error
	c.font = engosdl.GetFontManager().CreateFont(c.GetName(), c.FontFile, c.FontSize)
	c.texture = c.font.GetTextureFromFont(c.Message, c.Color)
	_, _, c.width, c.height, err = c.texture.Query()
	if err != nil {
		engosdl.Logger.Error().Err(err).Msg("Query error")
		panic(err)
	}
	c.GetEntity().GetTransform().SetDim(engosdl.NewVector(float64(c.width), float64(c.height)))
}

// OnRender is called for every render tick.
func (c *Button) OnRender() {
	x, y, w, h := c.GetEntity().GetTransform().GetRectExt()
	c.renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	c.renderer.SetDrawColor(c.borderColor.R, c.borderColor.G, c.borderColor.B, c.borderColor.A)
	if c.filled {
		c.renderer.FillRect(&sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)})
	} else {
		c.renderer.DrawRect(&sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)})
	}

	if c.GetDirty() {
		color := c.Color
		if c.GetEnabled() {
			color.A = 255
		} else {
			color.A = 100
		}
		c.SetColor(color)
		c.SetDirty(false)
	}
	c.renderer.CopyEx(c.texture,
		&sdl.Rect{X: 0, Y: 0, W: c.width, H: c.height},
		&sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)},
		0,
		&sdl.Point{},
		sdl.FLIP_NONE)
}

// OnStart is called first time the component is enabled.
func (c *Button) OnStart() {
	engosdl.Logger.Trace().Str("component", "text").Str("text", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
	c.loadTextureFromTTF()
}

// OnUpdate updates button component.
func (c *Button) OnUpdate() {
	entity := c.GetEntity()
	x, y, _ := sdl.GetMouseState()
	if entity.IsInside(engosdl.NewVector(float64(x), float64(y))) {
		// cursor := sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_HAND)
		// sdl.SetCursor(cursor)
		c.borderColor.A = 255
	} else {
		// cursor := sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_ARROW)
		// sdl.SetCursor(cursor)
		c.borderColor.A = 0
	}
	// c.Component.OnUpdate()
}

// SetColor sets text color.
func (c *Button) SetColor(color sdl.Color) engosdl.IButton {
	c.Color = color
	c.loadTextureFromTTF()
	return c
}

// SetFontFilename sets the filename with the font.
func (c *Button) SetFontFilename(filename string) engosdl.IButton {
	c.FontFile = filename
	return c
}

// SetFontSize sets the font size.
func (c *Button) SetFontSize(size int) engosdl.IButton {
	c.FontSize = size
	c.loadTextureFromTTF()
	return c
}

// SetMessage sets the message to be displayed by the text component.
func (c *Button) SetMessage(message string) engosdl.IButton {
	c.Message = message
	c.loadTextureFromTTF()
	return c
}

// Unmarshal takes a ComponentToMarshal instance and  creates a new entity
// instance.
func (c *Button) Unmarshal(data map[string]interface{}) {
	c.Component.Unmarshal(data)
	c.FontFile = data["font-filename"].(string)
	c.FontSize = int(data["font-size"].(float64))
	color := data["color"].(map[string]interface{})
	c.Color = sdl.Color{R: uint8(color["R"].(float64)),
		G: uint8(color["G"].(float64)),
		B: uint8(color["B"].(float64)),
		A: uint8(color["A"].(float64))}
	c.Message = data["message"].(string)
}
