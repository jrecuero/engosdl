package entities

import (
	"fmt"

	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
)

// Body2D represents an skeleton entity with all required components to
// create a 2D body.
type Body2D struct {
	*engosdl.Entity
	Filenames       []string
	NumberOfSprites int
	Format          int
	LeftCorner      bool
	MoveTo          *engosdl.Vector
}

// NewBody2D creates a new body 2D entity.
func NewBody2D(name string, filenames []string, numberOfSprites int, format int, leftCorner bool, moveTo *engosdl.Vector) *Body2D {
	engosdl.Logger.Trace().Str("entity", "body-2d").Str("body-2d", name).Msg("new body-2d")
	return &Body2D{
		Entity:          engosdl.NewEntity(name),
		Filenames:       filenames,
		NumberOfSprites: numberOfSprites,
		Format:          format,
		LeftCorner:      leftCorner,
		MoveTo:          moveTo,
	}
}

// DoLoad is called when object is loaded by the scene.
func (entity *Body2D) DoLoad() {
	engosdl.Logger.Trace().Str("entity", "body-2d").Str("body-2d", entity.GetName()).Msg("DoLoad")
	sprite := components.NewSprite(fmt.Sprintf("body-2d-%s/sprite", entity.GetName()),
		entity.Filenames,
		entity.NumberOfSprites,
		entity.Format)
	sprite.AddDelegateToRegister(nil, nil, &components.OutOfBounds{}, sprite.DefaultOnOutOfBounds)
	outOfBounds := components.NewOutOfBounds(fmt.Sprintf("body-2d-%s/out-of-bounds", entity.GetName()), entity.LeftCorner)
	outOfBounds.DefaultAddDelegateToRegister()
	moveTo := components.NewMoveTo(fmt.Sprintf("body-2d-%s/move-to", entity.GetName()), entity.MoveTo)
	moveTo.DefaultAddDelegateToRegister()
	collider2D := components.NewCollider2D(fmt.Sprintf("body-2d-%s/collider-2d", entity.GetName()))
	entity.AddComponent(sprite)
	entity.AddComponent(outOfBounds)
	entity.AddComponent(moveTo)
	entity.AddComponent(collider2D)
	entity.Entity.DoLoad()
}
