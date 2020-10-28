package components

import (
	"time"

	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// KeyShooter represents a component that can create a game object when
// a key is being pressed.
type KeyShooter struct {
	*engosdl.Component
	key       int
	cooldown  time.Duration
	lastshoot time.Time
	delegate  engosdl.IDelegate
}

// GetDelegates returns all delegates registered to the component.
func (k *KeyShooter) GetDelegates() []engosdl.IDelegate {
	return []engosdl.IDelegate{k.delegate}
}

// OnUpdate is called for every update tick.
func (k *KeyShooter) OnUpdate() {
	keys := sdl.GetKeyboardState()
	if keys[sdl.SCANCODE_SPACE] == 1 {
		engosdl.Logger.Trace().Str("component", "keyshooter").Str("keyshooter", k.GetName()).Msg("space key pressed")
		if time.Since(k.lastshoot) >= k.cooldown {

			// x := k.GetGameObject().GetTransform().GetPosition().X
			// y := k.GetGameObject().GetTransform().GetPosition().Y
			// bullet := engosdl.NewGameObject("bullet")
			// bullet.GetTransform().SetPosition(engosdl.NewVector(x, y))
			// bulletSprite := NewSprite("bullet-sprite", bullet, "images/player_bullet.bmp", engosdl.GetEngine().GetRenderer())
			// bulletMoveTo := NewMoveTo("bullet-moveto", bullet, engosdl.NewVector(0, -1))
			// bullet.AddComponent(bulletSprite)
			// bullet.AddComponent(bulletMoveTo)
			// k.GetGameObject().GetScene().AddGameObject(bullet)
			// bullet.OnStart()

			engosdl.GetEngine().GetEventHandler().GetDelegateHandler().TriggerDelegate(k.delegate)

			k.lastshoot = time.Now()

		}
	}
}

// NewKeyShooter creates a new keyshooter instance
func NewKeyShooter(name string, gobj *engosdl.GameObject, key int) *KeyShooter {
	engosdl.Logger.Trace().Str("component", "keyshooter").Str("keyshooter", name).Msg("new key-shooter")
	keyShooter := &KeyShooter{
		Component: engosdl.NewComponent(name, gobj),
		key:       key,
		cooldown:  500 * time.Millisecond,
	}
	keyShooter.delegate = engosdl.GetEngine().GetEventHandler().GetDelegateHandler().CreateDelegate(keyShooter, "shoot")
	return keyShooter
}
