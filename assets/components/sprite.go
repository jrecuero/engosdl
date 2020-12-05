package components

import (
	"fmt"
	"reflect"

	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// ComponentNameSprite is the name to refer sprite component
var ComponentNameSprite string = reflect.TypeOf(&Sprite{}).String()

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor(ComponentNameSprite, CreateSprite)
	}
}

// Sprite represents a component that can display multiple
// sprites, which can be animated.
type Sprite struct {
	*engosdl.Component
	Filenames      []string `json:"filenames"`
	width          int32
	height         int32
	renderer       *sdl.Renderer
	textures       []*sdl.Texture
	camera         *sdl.Rect
	fileImageIndex int
	SpriteTotal    int `json:"sprite-total"`
	spriteIndex    int
	resources      []engosdl.IResource
	Format         int `json:"format"`
}

var _ engosdl.ISprite = (*Sprite)(nil)

// NewSprite creates a new sprite instance.
// It register to "collision" delegate.
// It register to "out-of-bounds" delegate.
func NewSprite(name string, filenames []string, numberOfSprites int, format int) *Sprite {
	engosdl.Logger.Trace().Str("component", "sprite").Str("sprite", name).Msg("new sprite")
	result := &Sprite{
		Component:      engosdl.NewComponent(name),
		Filenames:      filenames,
		renderer:       engosdl.GetRenderer(),
		textures:       []*sdl.Texture{},
		camera:         nil,
		fileImageIndex: 0,
		SpriteTotal:    numberOfSprites,
		spriteIndex:    0,
		resources:      []engosdl.IResource{},
		Format:         format,
	}
	return result
}

// CreateSprite implements sprite constructor used by component manager.
// It register to "collision" delegate.
// It register to "out-of-bounds" delegate.
func CreateSprite(params ...interface{}) engosdl.IComponent {
	if len(params) == 4 {
		return NewSprite(params[0].(string), params[1].([]string), params[2].(int), params[3].(int))
	}
	return NewSprite("", []string{}, 1, engosdl.FormatBMP)
}

// DefaultAddDelegateToRegister will proceed to add default delegate to
// register for the component.
// It register to "collision" delegate.
// It register to "out-of-bounds" delegate.
func (c *Sprite) DefaultAddDelegateToRegister() {
	c.AddDelegateToRegister(engosdl.GetDelegateManager().GetCollisionDelegate(), nil, nil, c.DefaultOnCollision)
	c.AddDelegateToRegister(nil, nil, &OutOfBounds{}, c.DefaultOnOutOfBounds)
}

// DefaultOnCollision checks when there is a collision with other entity.
func (c *Sprite) DefaultOnCollision(params ...interface{}) bool {
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

// DefaultOnOutOfBounds checks if the entity has gone out of bounds.
func (c *Sprite) DefaultOnOutOfBounds(params ...interface{}) bool {
	if c.GetEntity().GetDieOnOutOfBounds() {
		entity := params[0].(engosdl.IEntity)
		if entity.GetID() == c.GetEntity().GetID() {
			engosdl.GetEngine().DestroyEntity(c.GetEntity())
		}
	}
	return true
}

// DoDestroy calls all methods to clean up sprite.
func (c *Sprite) DoDestroy() {
	engosdl.Logger.Trace().Str("component", "sprite").Str("sprite", c.GetName()).Msg("DoDestroy")
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
	engosdl.Logger.Trace().Str("component", "sprite").Str("sprite", c.GetName()).Msg("DoUnLoad")
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
	engosdl.Logger.Trace().Str("component", "sprite").Str("sprite", c.GetName()).Msg("LoadSprite")
	c.loadTextures()
	// TODO: assuming SpriteSheet is horizontal.
	c.GetEntity().GetTransform().SetDim(engosdl.NewVector(float64(c.width/int32(c.SpriteTotal)), float64(c.height)))
}

// loadTextures creates textures for every image file.
func (c *Sprite) loadTextures() {
	for _, filename := range c.Filenames {
		if len(c.resources) == 0 && len(c.textures) == 0 {
			var err error
			resource := engosdl.GetResourceManager().CreateResource(c.GetName(), filename, c.Format)
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
	engosdl.Logger.Trace().Str("component", "sprite").Str("sprite", c.GetName()).Msg("OnAwake")
	c.loadTextures()
	// TODO: assuming SpriteSheet is horizontal.
	c.GetEntity().GetTransform().SetDim(engosdl.NewVector(float64(c.width/int32(c.SpriteTotal)), float64(c.height)))
	c.Component.OnAwake()
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

// Unmarshal takes information from a ComponentToUnmarshal instance and
//  creates a new component instance.
func (c *Sprite) Unmarshal(data map[string]interface{}) {
	c.Component.Unmarshal(data)
	// c.SetName(data["name"].(string))
	c.Filenames = []string{}
	for _, filename := range data["filenames"].([]interface{}) {
		c.Filenames = append(c.Filenames, filename.(string))
	}
	c.SpriteTotal = int(data["sprite-total"].(float64))
	c.Format = int(data["format"].(float64))
}
