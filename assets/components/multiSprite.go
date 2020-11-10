package components

import (
	"fmt"

	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// MultiSprite represents a component that can display multiple
// sprites, which can be animated.
type MultiSprite struct {
	*engosdl.Component
	filenames            []string
	width                int32
	height               int32
	renderer             *sdl.Renderer
	textures             []*sdl.Texture
	destroyOnOutOfBounds bool
	camera               *sdl.Rect
	index                int
}

var _ engosdl.ISprite = (*MultiSprite)(nil)

// NewMultiSprite creates a new multi sprite instance.
// It registers to on-collision and on-out-of-bounds delegates.
func NewMultiSprite(name string, filenames []string, renderer *sdl.Renderer) *MultiSprite {
	engosdl.Logger.Trace().Str("component", "multi-sprite").Str("multi-sprite", name).Msg("new multi-sprite")
	result := &MultiSprite{
		Component:            engosdl.NewComponent(name),
		filenames:            filenames,
		renderer:             renderer,
		textures:             []*sdl.Texture{},
		destroyOnOutOfBounds: true,
		camera:               nil,
		index:                0,
	}
	result.AddDelegateToRegister(engosdl.GetEngine().GetEventHandler().GetDelegateHandler().GetCollisionDelegate(), nil, nil, result.onCollision)
	result.AddDelegateToRegister(nil, nil, &OutOfBounds{}, result.onOutOfBounds)
	return result
}

// DoUnLoad is called when component is unloaded, so all resources have
// to be released.
func (c *MultiSprite) DoUnLoad() {
	for _, texture := range c.textures {
		texture.Destroy()
	}
}

// GetCamera returns the camera used to display the sprite
func (c *MultiSprite) GetCamera() *sdl.Rect {
	return c.camera
}

// GetFilename returns filenames used for the sprite.
func (c *MultiSprite) GetFilename() []string {
	return c.filenames
}

// OnAwake should create all component resources that don't have any dependency
// with any other component or entity.
func (c *MultiSprite) OnAwake() {
	engosdl.Logger.Trace().Str("component", "multi-sprite").Str("sprite", c.GetName()).Msg("OnAwake")
	c.loadTexturesFromBMP()
	c.GetEntity().GetTransform().SetDim(engosdl.NewVector(float64(c.width), float64(c.height)))
}

// OnDraw is called for every draw tick.
func (c *MultiSprite) OnDraw() {
	// engosdl.Logger.Trace().Str("sprite", spr.GetName()).Msg("OnDraw")
	x, y, width, height := c.GetEntity().GetTransform().GetRectExt()
	var displayFrom *sdl.Rect
	var displayAt *sdl.Rect

	displayFrom = &sdl.Rect{X: 0, Y: 0, W: c.width, H: c.height}
	displayAt = &sdl.Rect{X: int32(x), Y: int32(y), W: int32(width), H: int32(height)}

	c.renderer.CopyEx(c.textures[c.index],
		displayFrom,
		displayAt,
		0,
		&sdl.Point{},
		sdl.FLIP_NONE)
}

// onCollision checks when there is a collision with other entity.
func (c *MultiSprite) onCollision(params ...interface{}) bool {
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

// onOutOfBounds checks if the entity has gone out of bounds.
func (c *MultiSprite) onOutOfBounds(params ...interface{}) bool {
	if c.destroyOnOutOfBounds {
		entity := params[0].(engosdl.IEntity)
		if entity.GetID() == c.GetEntity().GetID() {
			engosdl.GetEngine().DestroyEntity(c.GetEntity())
		}
	}
	return true
}

// OnStart is called first time the component is enabled.
func (c *MultiSprite) OnStart() {
	// Register to: "on-collision" and "out-of-bounds"
	engosdl.Logger.Trace().Str("component", "sprite").Str("sprite", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}

// SetCamera sets the camera used to display the sprite.
func (c *MultiSprite) SetCamera(camera *sdl.Rect) {
	c.camera = camera
}

// SetDestroyOnOutOfBounds sets internal attribute used to destroy sprite when
// it is out of bounds or no.
func (c *MultiSprite) SetDestroyOnOutOfBounds(destroy bool) {
	c.destroyOnOutOfBounds = destroy
}

// loadTexturesFromBMP creates textures for every BMP image file.
func (c *MultiSprite) loadTexturesFromBMP() {
	for _, filename := range c.filenames {
		img, err := sdl.LoadBMP(filename)
		if err != nil {
			engosdl.Logger.Error().Err(err).Msg("LoadBMP error")
			panic(err)
		}
		defer img.Free()
		texture, err := c.renderer.CreateTextureFromSurface(img)
		if err != nil {
			engosdl.Logger.Error().Err(err).Msg("CreateTextureFromSurface error")
			panic(err)
		}
		_, _, c.width, c.height, err = texture.Query()
		if err != nil {
			engosdl.Logger.Error().Err(err).Msg("Query error")
			panic(err)
		}
		c.textures = append(c.textures, texture)
	}
}
