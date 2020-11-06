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
		Sprite: NewSprite(name, filename, renderer),
		scroll: engosdl.NewVector(0, -1),
	}
}

// OnDraw is called for every draw tick.
func (c *ScrollSprite) OnDraw() {
	// engosdl.Logger.Trace().Str("sprite", spr.GetName()).Msg("OnDraw")
	x := int32(c.GetEntity().GetTransform().GetPosition().X)
	y := int32(c.GetEntity().GetTransform().GetPosition().Y)
	width := c.width * int32(c.GetEntity().GetTransform().GetScale().X)
	height := c.height * int32(c.GetEntity().GetTransform().GetScale().Y)
	// var displayFrom *sdl.Rect
	// var displayAt *sdl.Rect
	W, H, _ := engosdl.GetEngine().GetRenderer().GetOutputSize()
	if c.scroll.Y == -1 {
		y = y % height
	} else if c.scroll.X == -1 {
		x = x % width
	}
	displayFrom := &sdl.Rect{X: 0, Y: 0, W: width, H: height}
	displayAt := &sdl.Rect{X: x, Y: y, W: width, H: height}
	c.renderer.CopyEx(c.texture,
		displayFrom,
		displayAt,
		0,
		&sdl.Point{},
		sdl.FLIP_NONE)
	if c.scroll.Y == -1 && (y+height) < H {
		c.renderer.CopyEx(c.texture,
			&sdl.Rect{X: 0, Y: 0, W: width, H: height},
			&sdl.Rect{X: x, Y: y + height, W: width, H: height},
			0,
			&sdl.Point{},
			sdl.FLIP_NONE)
	} else if c.scroll.X == -1 && (x+width) < W {
		c.renderer.CopyEx(c.texture,
			&sdl.Rect{X: 0, Y: 0, W: width, H: height},
			&sdl.Rect{X: x + width, Y: y, W: width, H: height},
			0,
			&sdl.Point{},
			sdl.FLIP_NONE)
	}
}

// SetScroll sets sprite image scroll.
func (c *ScrollSprite) SetScroll(scroll *engosdl.Vector) {
	c.scroll = scroll
}
