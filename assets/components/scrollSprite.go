package components

import (
	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// ScrollSprite represents an sprite that scroll across the display.
type ScrollSprite struct {
	*Sprite
	scroll *engosdl.Vector
}

var _ engosdl.ISprite = (*ScrollSprite)(nil)

// NewScrollSprite creates a new sprite instance.
func NewScrollSprite(name string, filename string, renderer *sdl.Renderer) *ScrollSprite {
	engosdl.Logger.Trace().Str("component", "sprite").Str("sprite", name).Msg("new sprite")
	return &ScrollSprite{
		// NewSprite is not called because it registers to onOutOfBounds and
		// onCollision, and by default scroll sprite does not want to register
		// to any of then.
		Sprite: &Sprite{
			Component:            engosdl.NewComponent(name),
			filenames:            []string{filename},
			renderer:             renderer,
			destroyOnOutOfBounds: true,
			camera:               nil,
			spriteTotal:          1,
		},
		scroll: engosdl.NewVector(0, -1),
	}
}

// OnRender is called for every render tick.
func (c *ScrollSprite) OnRender() {
	// engosdl.Logger.Trace().Str("sprite", spr.GetName()).Msg("OnRender")
	x := int32(c.GetEntity().GetTransform().GetPosition().X)
	y := int32(c.GetEntity().GetTransform().GetPosition().Y)
	width := c.width * int32(c.GetEntity().GetTransform().GetScale().X)
	height := c.height * int32(c.GetEntity().GetTransform().GetScale().Y)
	W, H, _ := engosdl.GetRenderer().GetOutputSize()
	if c.scroll.Y == -1 {
		y = y % height
	} else if c.scroll.X == -1 {
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
	if c.scroll.Y == -1 && (y+height) < H {
		c.renderer.CopyEx(c.textures[0],
			&sdl.Rect{X: 0, Y: 0, W: width, H: height},
			&sdl.Rect{X: x, Y: y + height, W: width, H: height},
			0,
			&sdl.Point{},
			sdl.FLIP_NONE)
	} else if c.scroll.X == -1 && (x+width) < W {
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
	c.scroll = scroll
}
