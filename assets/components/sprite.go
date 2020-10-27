package components

import (
	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// Sprite represents a component that can display a sprite texture.
type Sprite struct {
	*engosdl.Component
	texture *sdl.Texture
}

// NewSprite creates a new sprite instance.
func NewSprite(name string) *Sprite {
	engosdl.Logger.Trace().Str("sprite", name).Msg("new sprite")
	return &Sprite{
		Component: engosdl.NewComponent(name),
	}
}
