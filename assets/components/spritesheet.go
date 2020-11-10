package components

import (
	"fmt"

	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// SpriteSheet represents a component that can display multiple
// sprites, which can be animated.
type SpriteSheet struct {
	*engosdl.Component
	filenames            []string
	width                int32
	height               int32
	renderer             *sdl.Renderer
	textures             []*sdl.Texture
	destroyOnOutOfBounds bool
	camera               *sdl.Rect
	fileImageIndex       int
	spriteTotal          int
	spriteIndex          int
}

var _ engosdl.ISprite = (*SpriteSheet)(nil)

// NewSpriteSheet creates a new multi sprite instance.
// It resgiters to on-collision and on-out-of-bounds delegates.
func NewSpriteSheet(name string, filenames []string, numberOfSprites int, renderer *sdl.Renderer) *SpriteSheet {
	engosdl.Logger.Trace().Str("component", "multi-sprite").Str("multi-sprite", name).Msg("new multi-sprite")
	result := &SpriteSheet{
		Component:            engosdl.NewComponent(name),
		filenames:            filenames,
		renderer:             renderer,
		textures:             []*sdl.Texture{},
		destroyOnOutOfBounds: true,
		camera:               nil,
		fileImageIndex:       0,
		spriteTotal:          numberOfSprites,
		spriteIndex:          0,
	}
	result.AddDelegateToRegister(engosdl.GetEngine().GetEventHandler().GetDelegateHandler().GetCollisionDelegate(), nil, nil, result.onCollision)
	result.AddDelegateToRegister(nil, nil, &OutOfBounds{}, result.onOutOfBounds)
	return result
}

// DoUnLoad is called when component is unloaded, so all resources have
// to be released.
func (c *SpriteSheet) DoUnLoad() {
	for _, texture := range c.textures {
		texture.Destroy()
	}
}

// GetCamera returns the camera used to display the sprite
func (c *SpriteSheet) GetCamera() *sdl.Rect {
	return c.camera
}

// GetFileImageIndex returns sprite sheet file image index currently used.
func (c *SpriteSheet) GetFileImageIndex() int {
	return c.fileImageIndex
}

// GetFilename returns filenames used for the sprite.
func (c *SpriteSheet) GetFilename() []string {
	return c.filenames
}

// GetSpriteIndex returns sprite sheet sprite index currently used.
func (c *SpriteSheet) GetSpriteIndex() int {
	return c.spriteIndex
}

// loadTexturesFromBMP creates textures for every BMP image file.
func (c *SpriteSheet) loadTexturesFromBMP() {
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

// NextFileImage increases by one file image index.
func (c *SpriteSheet) NextFileImage() int {
	c.fileImageIndex = (c.fileImageIndex + 1) % len(c.filenames)
	return c.fileImageIndex
}

// NextSprite increases by one the sprite index.
func (c *SpriteSheet) NextSprite() int {
	c.spriteIndex = (c.spriteIndex + 1) % c.spriteTotal
	return c.spriteIndex
}

// OnAwake should create all component resources that don't have any dependency
// with any other component or entity.
func (c *SpriteSheet) OnAwake() {
	engosdl.Logger.Trace().Str("component", "multi-sprite").Str("sprite", c.GetName()).Msg("OnAwake")
	c.loadTexturesFromBMP()
	// TODO: assuming SpriteSheet is horizontal.
	c.GetEntity().GetTransform().SetDim(engosdl.NewVector(float64(c.width/int32(c.spriteTotal)), float64(c.height)))
}

// onCollision checks when there is a collision with other entity.
func (c *SpriteSheet) onCollision(params ...interface{}) bool {
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

// OnDraw is called for every draw tick.
func (c *SpriteSheet) OnDraw() {
	// engosdl.Logger.Trace().Str("sprite", spr.GetName()).Msg("OnDraw")
	x, y, width, height := c.GetEntity().GetTransform().GetRectExt()
	var displayFrom *sdl.Rect
	var displayAt *sdl.Rect
	spriteX := (c.spriteIndex * int(c.width)) / int(c.spriteTotal)

	displayFrom = &sdl.Rect{X: int32(spriteX), Y: 0, W: c.width / int32(c.spriteTotal), H: c.height}
	displayAt = &sdl.Rect{X: int32(x), Y: int32(y), W: int32(width), H: int32(height)}

	c.renderer.CopyEx(c.textures[c.fileImageIndex],
		displayFrom,
		displayAt,
		0,
		&sdl.Point{},
		sdl.FLIP_NONE)
}

// onOutOfBounds checks if the entity has gone out of bounds.
func (c *SpriteSheet) onOutOfBounds(params ...interface{}) bool {
	if c.destroyOnOutOfBounds {
		entity := params[0].(engosdl.IEntity)
		if entity.GetID() == c.GetEntity().GetID() {
			engosdl.GetEngine().DestroyEntity(c.GetEntity())
		}
	}
	return true
}

// OnStart is called first time the component is enabled.
func (c *SpriteSheet) OnStart() {
	// Register to: "on-collision" and "out-of-bounds"
	engosdl.Logger.Trace().Str("component", "sprite").Str("sprite", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}

//PreviousFileImage decreases by one file image index.
func (c *SpriteSheet) PreviousFileImage() int {
	c.fileImageIndex = (c.fileImageIndex - 1) % len(c.filenames)
	return c.fileImageIndex
}

// PreviousSprite decreases by one sprite index.
func (c *SpriteSheet) PreviousSprite() int {
	c.spriteIndex = (c.spriteIndex - 1) % c.spriteTotal
	return c.spriteIndex
}

// SetCamera sets the camera used to display the sprite.
func (c *SpriteSheet) SetCamera(camera *sdl.Rect) {
	c.camera = camera
}

// SetDestroyOnOutOfBounds sets internal attribute used to destroy sprite when
// it is out of bounds or no.
func (c *SpriteSheet) SetDestroyOnOutOfBounds(destroy bool) {
	c.destroyOnOutOfBounds = destroy
}
