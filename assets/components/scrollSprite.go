package components

import (
	"reflect"

	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// ComponentNameScrollSprite is the name to refer scroll sprite component.
var ComponentNameScrollSprite string = reflect.TypeOf(&ScrollSprite{}).String()

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor(ComponentNameScrollSprite, CreateScrollSprite)
	}
}

// ScrollSprite represents an sprite that scroll across the display.
type ScrollSprite struct {
	*Sprite
	Scroll *engosdl.Vector `json:"scroll"`
}

var _ engosdl.ISprite = (*ScrollSprite)(nil)

// NewScrollSprite creates a new sprite instance.
func NewScrollSprite(name string, filename string, format int) *ScrollSprite {
	engosdl.Logger.Trace().Str("component", "sprite").Str("sprite", name).Msg("new sprite")
	return &ScrollSprite{
		// NewSprite is not called because it registers to onOutOfBounds and
		// onCollision, and by default scroll sprite does not want to register
		// to any of then.
		Sprite: &Sprite{
			Component:            engosdl.NewComponent(name),
			Filenames:            []string{filename},
			renderer:             engosdl.GetRenderer(),
			DestroyOnOutOfBounds: true,
			camera:               nil,
			SpriteTotal:          1,
			Format:               format,
		},
		Scroll: engosdl.NewVector(0, -1),
	}
}

// CreateScrollSprite implements scroll sprite constructor used by component
// manager.
func CreateScrollSprite(params ...interface{}) engosdl.IComponent {
	if len(params) == 3 {
		return NewScrollSprite(params[0].(string), params[1].(string), params[2].(int))
	}
	return NewScrollSprite("", "", engosdl.FormatBMP)
}

// OnRender is called for every render tick.
func (c *ScrollSprite) OnRender() {
	// engosdl.Logger.Trace().Str("sprite", spr.GetName()).Msg("OnRender")
	x := int32(c.GetEntity().GetTransform().GetPosition().X)
	y := int32(c.GetEntity().GetTransform().GetPosition().Y)
	width := c.width * int32(c.GetEntity().GetTransform().GetScale().X)
	height := c.height * int32(c.GetEntity().GetTransform().GetScale().Y)
	W, H, _ := engosdl.GetRenderer().GetOutputSize()
	if c.Scroll.Y == -1 {
		y = y % height
	} else if c.Scroll.X == -1 {
		x = x % width
	}
	displayFrom := &sdl.Rect{X: 0, Y: 0, W: width, H: height}
	displayAt := &sdl.Rect{X: x, Y: y, W: width, H: height}
	c.renderer.CopyEx(c.textures[0],
		displayFrom,
		displayAt,
		0,
		&sdl.Point{},
		sdl.FLIP_NONE)
	if c.Scroll.Y == -1 && (y+height) < H {
		c.renderer.CopyEx(c.textures[0],
			&sdl.Rect{X: 0, Y: 0, W: width, H: height},
			&sdl.Rect{X: x, Y: y + height, W: width, H: height},
			0,
			&sdl.Point{},
			sdl.FLIP_NONE)
	} else if c.Scroll.X == -1 && (x+width) < W {
		c.renderer.CopyEx(c.textures[0],
			&sdl.Rect{X: 0, Y: 0, W: width, H: height},
			&sdl.Rect{X: x + width, Y: y, W: width, H: height},
			0,
			&sdl.Point{},
			sdl.FLIP_NONE)
	}
}

// OnStart is called first time the component is enabled.
func (c *ScrollSprite) OnStart() {
	// Register to: "on-collision" and "out-of-bounds"
	engosdl.Logger.Trace().Str("component", "scroll-sprite").Str("scroll-sprite", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}

// SetScroll sets sprite image scroll.
func (c *ScrollSprite) SetScroll(scroll *engosdl.Vector) {
	c.Scroll = scroll
}

// Unmarshal takes a ComponentToMarshal instance and  creates a new entity
// instance.
func (c *ScrollSprite) Unmarshal(data map[string]interface{}) {
	c.Component.Unmarshal(data)
	// c.SetName(data["name"].(string))
	c.Filenames = []string{}
	for _, filename := range data["filenames"].([]interface{}) {
		c.Filenames = append(c.Filenames, filename.(string))
	}
	c.SpriteTotal = int(data["sprite-total"].(float64))
	scroll := data["scroll"].(map[string]interface{})
	c.SetScroll(engosdl.NewVector(scroll["X"].(float64), scroll["Y"].(float64)))
	c.DestroyOnOutOfBounds = data["destroy-on-out-of-bounds"].(bool)
	c.Format = int(data["format"].(float64))
}
