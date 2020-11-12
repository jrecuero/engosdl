package components

import (
	"fmt"

	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterComponent(&Sprite{})
	}
}

// Sprite represents a component that can display multiple
// sprites, which can be animated.
type Sprite struct {
	*engosdl.Component
	Filenames            []string `json:"filenames"`
	width                int32
	height               int32
	renderer             *sdl.Renderer
	textures             []*sdl.Texture
	DestroyOnOutOfBounds bool `json:"destroy-on-out-of-bounds"`
	camera               *sdl.Rect
	fileImageIndex       int
	SpriteTotal          int `json:"sprite-total"`
	spriteIndex          int
	resources            []engosdl.IResource
}

var _ engosdl.ISprite = (*Sprite)(nil)

// NewSprite creates a new sprite instance.
// It resgiters to "on-collision" and "on-out-of-bounds" delegates.
func NewSprite(name string, filenames []string, numberOfSprites int, renderer *sdl.Renderer) *Sprite {
	engosdl.Logger.Trace().Str("component", "multi-sprite").Str("multi-sprite", name).Msg("new multi-sprite")
	result := &Sprite{
		Component:            engosdl.NewComponent(name),
		Filenames:            filenames,
		renderer:             renderer,
		textures:             []*sdl.Texture{},
		DestroyOnOutOfBounds: true,
		camera:               nil,
		fileImageIndex:       0,
		SpriteTotal:          numberOfSprites,
		spriteIndex:          0,
		resources:            []engosdl.IResource{},
	}
	return result
}

// DefaultAddDelegateToRegister will proceed to add default delegate to
// register for the component.
func (c *Sprite) DefaultAddDelegateToRegister() {
	c.AddDelegateToRegister(engosdl.GetDelegateHandler().GetCollisionDelegate(), nil, nil, c.onCollision)
	c.AddDelegateToRegister(nil, nil, &OutOfBounds{}, c.onOutOfBounds)
}

// DoDestroy calls all methods to clean up sprite.
func (c *Sprite) DoDestroy() {
	for _, texture := range c.textures {
		texture.Destroy()
	}
	c.textures = []*sdl.Texture{}
	c.resources = []engosdl.IResource{}
	c.Component.DoDestroy()
}

// DoUnLoad is called when component is unloaded, so all resources have
// to be released.
func (c *Sprite) DoUnLoad() {
	// for _, texture := range c.textures {
	// 	texture.Destroy()
	// }
	// c.textures = []*sdl.Texture{}
	// c.resources = []engosdl.IResource{}
	c.Component.DoUnLoad()
}

// GetCamera returns the camera used to display the sprite
func (c *Sprite) GetCamera() *sdl.Rect {
	return c.camera
}

// GetFileImageIndex returns sprite sheet file image index currently used.
func (c *Sprite) GetFileImageIndex() int {
	return c.fileImageIndex
}

// GetFilename returns filenames used for the sprite.
func (c *Sprite) GetFilename() []string {
	return c.Filenames
}

// GetSpriteIndex returns sprite sheet sprite index currently used.
func (c *Sprite) GetSpriteIndex() int {
	return c.spriteIndex
}

// LoadSprite loads the sprite from the filename.
func (c *Sprite) LoadSprite() {
	engosdl.Logger.Trace().Str("component", "multi-sprite").Str("sprite", c.GetName()).Msg("LoadSprite")
	c.loadTexturesFromBMP()
	// TODO: assuming SpriteSheet is horizontal.
	c.GetEntity().GetTransform().SetDim(engosdl.NewVector(float64(c.width/int32(c.SpriteTotal)), float64(c.height)))
}

// loadTexturesFromBMP creates textures for every BMP image file.
func (c *Sprite) loadTexturesFromBMP() {
	for _, filename := range c.Filenames {
		if len(c.resources) == 0 && len(c.textures) == 0 {
			var err error
			resource := engosdl.GetResourceHandler().CreateResource(c.GetName(), filename)
			texture := resource.GetTextureFromSurface()
			_, _, c.width, c.height, err = texture.Query()
			if err != nil {
				engosdl.Logger.Error().Err(err).Msg("Query error")
				panic(err)
			}
			c.resources = append(c.resources, resource)
			c.textures = append(c.textures, texture)
		}
	}
}

// NextFileImage increases by one file image index.
func (c *Sprite) NextFileImage() int {
	c.fileImageIndex = (c.fileImageIndex + 1) % len(c.Filenames)
	return c.fileImageIndex
}

// NextSprite increases by one the sprite index.
func (c *Sprite) NextSprite() int {
	c.spriteIndex = (c.spriteIndex + 1) % c.SpriteTotal
	return c.spriteIndex
}

// OnAwake should create all component resources that don't have any dependency
// with any other component or entity.
func (c *Sprite) OnAwake() {
	engosdl.Logger.Trace().Str("component", "multi-sprite").Str("sprite", c.GetName()).Msg("OnAwake")
	c.loadTexturesFromBMP()
	// TODO: assuming SpriteSheet is horizontal.
	c.GetEntity().GetTransform().SetDim(engosdl.NewVector(float64(c.width/int32(c.SpriteTotal)), float64(c.height)))
	c.Component.OnAwake()
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
	if c.DestroyOnOutOfBounds {
		entity := params[0].(engosdl.IEntity)
		if entity.GetID() == c.GetEntity().GetID() {
			engosdl.GetEngine().DestroyEntity(c.GetEntity())
		}
	}
	return true
}

// OnRender is called for every render tick.
func (c *Sprite) OnRender() {
	// engosdl.Logger.Trace().Str("sprite", spr.GetName()).Msg("OnRender")
	x, y, width, height := c.GetEntity().GetTransform().GetRectExt()
	var displayFrom *sdl.Rect
	var displayAt *sdl.Rect
	spriteX := (c.spriteIndex * int(c.width)) / int(c.SpriteTotal)

	displayFrom = &sdl.Rect{X: int32(spriteX), Y: 0, W: c.width / int32(c.SpriteTotal), H: c.height}
	displayAt = &sdl.Rect{X: int32(x), Y: int32(y), W: int32(width), H: int32(height)}

	c.renderer.CopyEx(c.textures[c.fileImageIndex],
		displayFrom,
		displayAt,
		0,
		&sdl.Point{},
		sdl.FLIP_NONE)
}

// OnStart is called first time the component is enabled.
func (c *Sprite) OnStart() {
	// Register to: "on-collision" and "out-of-bounds"
	engosdl.Logger.Trace().Str("component", "sprite").Str("sprite", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}

//PreviousFileImage decreases by one file image index.
func (c *Sprite) PreviousFileImage() int {
	c.fileImageIndex = (c.fileImageIndex - 1) % len(c.Filenames)
	return c.fileImageIndex
}

// PreviousSprite decreases by one sprite index.
func (c *Sprite) PreviousSprite() int {
	c.spriteIndex = (c.spriteIndex - 1) % c.SpriteTotal
	return c.spriteIndex
}

// SetCamera sets the camera used to display the sprite.
func (c *Sprite) SetCamera(camera *sdl.Rect) {
	c.camera = camera
}

// SetDestroyOnOutOfBounds sets internal attribute used to destroy sprite when
// it is out of bounds or no.
func (c *Sprite) SetDestroyOnOutOfBounds(destroy bool) {
	c.DestroyOnOutOfBounds = destroy
}
