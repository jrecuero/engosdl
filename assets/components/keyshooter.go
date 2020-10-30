package components

import (
	"time"

	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// KeyShooter represents a component that can create a entity when
// a key is being pressed.
type KeyShooter struct {
	*engosdl.Component
	key       int
	cooldown  time.Duration
	lastshoot time.Time
	delegate  engosdl.IDelegate
}

// NewKeyShooter creates a new keyshooter instance
func NewKeyShooter(name string, key int) *KeyShooter {
	engosdl.Logger.Trace().Str("component", "keyshooter").Str("keyshooter", name).Msg("new key-shooter")
	keyShooter := &KeyShooter{
		Component: engosdl.NewComponent(name),
		key:       key,
		cooldown:  500 * time.Millisecond,
	}
	return keyShooter
}

// GetDelegates returns all delegates registered to the component.
func (c *KeyShooter) GetDelegates() []engosdl.IDelegate {
	return []engosdl.IDelegate{c.delegate}
}

// OnStart is called first time the component is enabled.
func (c *KeyShooter) OnStart() {
	c.delegate = engosdl.GetEngine().GetEventHandler().GetDelegateHandler().CreateDelegate(c, "shoot")
}

// OnUpdate is called for every update tick.
func (c *KeyShooter) OnUpdate() {
	keys := sdl.GetKeyboardState()
	if keys[sdl.SCANCODE_SPACE] == 1 {
		engosdl.Logger.Trace().Str("component", "keyshooter").Str("keyshooter", c.GetName()).Msg("space key pressed")
		if time.Since(c.lastshoot) >= c.cooldown {
			engosdl.GetEngine().GetEventHandler().GetDelegateHandler().TriggerDelegate(c.delegate)
			c.lastshoot = time.Now()
		}
	}
}
