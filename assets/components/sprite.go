package components

import (
	"fmt"

	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// Sprite represents a component that can display a sprite texture.
type Sprite struct {
	*engosdl.Component
	filename string
	width    int32
	height   int32
	centered bool
	renderer *sdl.Renderer
	texture  *sdl.Texture
}

// NewSprite creates a new sprite instance.
func NewSprite(name string, filename string, renderer *sdl.Renderer, centered bool) *Sprite {
	engosdl.Logger.Trace().Str("component", "sprite").Str("sprite", name).Msg("new sprite")
	return &Sprite{
		Component: engosdl.NewComponent(name),
		filename:  filename,
		renderer:  renderer,
		centered:  centered,
	}
}

// onCollision checks when there is a collision with other entity.
func (c *Sprite) onCollision(params ...interface{}) bool {
	collisionEntityOne := params[0].(*engosdl.Entity)
	collisionEntityTwo := params[1].(*engosdl.Entity)
	if c.GetEntity().GetID() == collisionEntityOne.GetID() || c.GetEntity().GetID() == collisionEntityTwo.GetID() {
		fmt.Printf("%s sprite onCollision %s with %s\n", c.GetEntity().GetName(), collisionEntityOne.GetName(), collisionEntityTwo.GetName())
		if collisionEntityOne.GetDieOnCollision() {
			engosdl.GetEngine().DestroyEntity(collisionEntityOne)
		}
		if collisionEntityTwo.GetDieOnCollision() {
			engosdl.GetEngine().DestroyEntity(collisionEntityTwo)
		}
	}
	return true
}

// OnStart is called first time the component is enabled.
func (c *Sprite) OnStart() {
	engosdl.Logger.Trace().Str("component", "sprite").Str("sprite", c.GetName()).Msg("OnStart")
	c.textureFromBMP()
	delegate := engosdl.GetEngine().GetEventHandler().GetDelegateHandler().GetCollisionDelegate()
	engosdl.GetEngine().GetEventHandler().GetDelegateHandler().RegisterToDelegate(delegate, c.onCollision)
}

// OnDraw is called for every draw tick.
func (c *Sprite) OnDraw() {
	// engosdl.Logger.Trace().Str("sprite", spr.GetName()).Msg("OnDraw")
	x := int32(c.GetEntity().GetTransform().GetPosition().X)
	y := int32(c.GetEntity().GetTransform().GetPosition().Y)
	width := c.width * int32(c.GetEntity().GetTransform().GetScale().X)
	height := c.height * int32(c.GetEntity().GetTransform().GetScale().Y)
	var displayAt *sdl.Rect
	if c.centered {
		displayAt = &sdl.Rect{X: x - c.width/2, Y: y - c.height/2, W: width, H: height}
	} else {
		displayAt = &sdl.Rect{X: x, Y: y, W: width, H: height}
	}
	c.renderer.CopyEx(c.texture,
		&sdl.Rect{X: 0, Y: 0, W: c.width, H: c.height},
		displayAt,
		// &sdl.Rect{X: x, Y: y, W: width, H: height},
		0,
		&sdl.Point{},
		sdl.FLIP_NONE)
}

// textureFromBMP creates a texture from a BMP image file.
func (c *Sprite) textureFromBMP() {
	img, err := sdl.LoadBMP(c.filename)
	if err != nil {
		engosdl.Logger.Error().Err(err)
		panic(err)
	}
	defer img.Free()
	c.texture, err = c.renderer.CreateTextureFromSurface(img)
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
	// c.center.X = c.GetParent().GetTransform().GetPosition().X + float64(spr.width)/2
	// c.center.Y = c.GetParent().GetTransform().GetPosition().Y + float64(spr.height)/2
}
