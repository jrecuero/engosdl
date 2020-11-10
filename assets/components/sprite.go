package components

import (
	"fmt"

	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// Sprite represents a component that can display a sprite texture.
type Sprite struct {
	*engosdl.Component
	filename             string
	width                int32
	height               int32
	renderer             *sdl.Renderer
	texture              *sdl.Texture
	destroyOnOutOfBounds bool
	camera               *sdl.Rect
}

var _ engosdl.ISprite = (*Sprite)(nil)

// NewSprite creates a new sprite instance.
// It registers to on-collision and on-out-of-bounds delegate.
func NewSprite(name string, filename string, renderer *sdl.Renderer) *Sprite {
	engosdl.Logger.Trace().Str("component", "sprite").Str("sprite", name).Msg("new sprite")
	result := &Sprite{
		Component:            engosdl.NewComponent(name),
		filename:             filename,
		renderer:             renderer,
		destroyOnOutOfBounds: true,
		camera:               nil,
	}
	result.AddDelegateToRegister(engosdl.GetEngine().GetEventHandler().GetDelegateHandler().GetCollisionDelegate(), nil, nil, result.onCollision)
	result.AddDelegateToRegister(nil, nil, &OutOfBounds{}, result.onOutOfBounds)
	return result
}

// DoUnLoad is called when component is unloaded, so all resources have
// to be released.
func (c *Sprite) DoUnLoad() {
	c.texture.Destroy()
}

// GetCamera returns the camera used to display the sprite
func (c *Sprite) GetCamera() *sdl.Rect {
	return c.camera
}

// GetFilename returns filename used for the sprite.
func (c *Sprite) GetFilename() []string {
	return []string{c.filename}
}

// OnAwake should create all component resources that don't have any dependency
// with any other component or entity.
func (c *Sprite) OnAwake() {
	engosdl.Logger.Trace().Str("component", "sprite").Str("sprite", c.GetName()).Msg("OnAwake")
	c.textureFromBMP()
	c.GetEntity().GetTransform().SetDim(engosdl.NewVector(float64(c.width), float64(c.height)))
}

// OnDraw is called for every draw tick.
func (c *Sprite) OnDraw() {
	// engosdl.Logger.Trace().Str("sprite", spr.GetName()).Msg("OnDraw")
	x := int32(c.GetEntity().GetTransform().GetPosition().X)
	y := int32(c.GetEntity().GetTransform().GetPosition().Y)
	width := int32(float64(c.width) * c.GetEntity().GetTransform().GetScale().X)
	height := int32(float64(c.height) * c.GetEntity().GetTransform().GetScale().Y)
	var displayFrom *sdl.Rect
	var displayAt *sdl.Rect
	displayFrom = &sdl.Rect{X: 0, Y: 0, W: c.width, H: c.width}
	displayAt = &sdl.Rect{X: x, Y: y, W: width, H: height}

	c.renderer.CopyEx(c.texture,
		displayFrom,
		displayAt,
		0,
		&sdl.Point{},
		sdl.FLIP_NONE)
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

// onOutOfBounds checks if the entity has gone out of bounds.
func (c *Sprite) onOutOfBounds(params ...interface{}) bool {
	if c.destroyOnOutOfBounds {
		entity := params[0].(engosdl.IEntity)
		if entity.GetID() == c.GetEntity().GetID() {
			engosdl.GetEngine().DestroyEntity(c.GetEntity())
		}
	}
	return true
}

// OnStart is called first time the component is enabled.
func (c *Sprite) OnStart() {
	// Register to: "on-collision" and "out-of-bounds"
	engosdl.Logger.Trace().Str("component", "sprite").Str("sprite", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}

// SetCamera sets the camera used to display the sprite.
func (c *Sprite) SetCamera(camera *sdl.Rect) {
	c.camera = camera
}

// SetDestroyOnOutOfBounds sets internal attribute used to destroy sprite when
// it is out of bounds or no.
func (c *Sprite) SetDestroyOnOutOfBounds(destroy bool) {
	c.destroyOnOutOfBounds = destroy
}

// textureFromBMP creates a texture from a BMP image file.
func (c *Sprite) textureFromBMP() {
	img, err := sdl.LoadBMP(c.filename)
	if err != nil {
		engosdl.Logger.Error().Err(err).Msg("LoadBMP error")
		panic(err)
	}
	defer img.Free()
	c.texture, err = c.renderer.CreateTextureFromSurface(img)
	if err != nil {
		engosdl.Logger.Error().Err(err).Msg("CreateTextureFromSurface error")
		panic(err)
	}
	_, _, c.width, c.height, err = c.texture.Query()
	if err != nil {
		engosdl.Logger.Error().Err(err).Msg("Query error")
		panic(err)
	}
}
