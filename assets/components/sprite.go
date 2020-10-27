package components

import (
	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// Sprite represents a component that can display a sprite texture.
type Sprite struct {
	*engosdl.Component
	filename string
	// center   *engosdl.Vector
	width    int32
	height   int32
	renderer *sdl.Renderer
	texture  *sdl.Texture
}

func (spr *Sprite) textureFromBMP() {
	img, err := sdl.LoadBMP(spr.filename)
	if err != nil {
		engosdl.Logger.Error().Err(err)
		panic(err)
	}
	defer img.Free()
	spr.texture, err = spr.renderer.CreateTextureFromSurface(img)
	if err != nil {
		engosdl.Logger.Error().Err(err)
		panic(err)
	}
	_, _, spr.width, spr.height, err = spr.texture.Query()
	if err != nil {
		engosdl.Logger.Error().Err(err)
		panic(err)
	}
	// spr.center.X = spr.GetParent().GetTransform().GetPosition().X + float64(spr.width)/2
	// spr.center.Y = spr.GetParent().GetTransform().GetPosition().Y + float64(spr.height)/2
}

// OnStart is called first time the component is enabled.
func (spr *Sprite) OnStart() {
	engosdl.Logger.Trace().Str("component", "sprite").Str("sprite", spr.GetName()).Msg("OnStart")
	spr.textureFromBMP()
}

// OnDraw is called for every draw tick.
func (spr *Sprite) OnDraw() {
	// engosdl.Logger.Trace().Str("sprite", spr.GetName()).Msg("OnDraw")
	x := int32(spr.GetGameObject().GetTransform().GetPosition().X)
	y := int32(spr.GetGameObject().GetTransform().GetPosition().Y)
	width := spr.width * int32(spr.GetGameObject().GetTransform().GetScale().X)
	height := spr.height * int32(spr.GetGameObject().GetTransform().GetScale().Y)
	spr.renderer.CopyEx(spr.texture,
		&sdl.Rect{X: 0, Y: 0, W: spr.width, H: spr.height},
		&sdl.Rect{X: x, Y: y, W: width, H: height},
		0,
		&sdl.Point{},
		sdl.FLIP_NONE)
}

// NewSprite creates a new sprite instance.
func NewSprite(name string, gobj *engosdl.GameObject, filename string, renderer *sdl.Renderer) *Sprite {
	engosdl.Logger.Trace().Str("component", "sprite").Str("sprite", name).Msg("new sprite")
	return &Sprite{
		Component: engosdl.NewComponent(name, gobj),
		filename:  filename,
		renderer:  renderer,
		// center:    engosdl.NewVector(0, 0),
	}
}
