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
}

// NewKeyShooter creates a new keyshooter instance
// It creates delegate "on-shoot".
func NewKeyShooter(name string, key int) *KeyShooter {
	engosdl.Logger.Trace().Str("component", "key-shooter").Str("key-shooter", name).Msg("new key-shooter")
	keyShooter := &KeyShooter{
		Component: engosdl.NewComponent(name),
		key:       key,
		cooldown:  500 * time.Millisecond,
	}
	return keyShooter
}

// OnAwake should create all component resources that don't have any dependency
// with any other component or entity.
func (c *KeyShooter) OnAwake() {
	engosdl.Logger.Trace().Str("component", "key-shooter").Str("key-shooter", c.GetName()).Msg("OnAwake")
	// Create new delegate "shoot"
	c.SetDelegate(engosdl.GetDelegateHandler().CreateDelegate(c, "on-shoot"))
	c.Component.OnAwake()
}

// OnStart is called first time the component is enabled.
func (c *KeyShooter) OnStart() {
	engosdl.Logger.Trace().Str("component", "key-shooter").Str("key-shooter", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}

// OnUpdate is called for every update tick.
func (c *KeyShooter) OnUpdate() {
	keys := sdl.GetKeyboardState()
	if keys[sdl.SCANCODE_SPACE] == 1 {
		engosdl.Logger.Trace().Str("component", "key-shooter").Str("key-shooter", c.GetName()).Msg("space key pressed")
		if time.Since(c.lastshoot) >= c.cooldown {
			engosdl.GetDelegateHandler().TriggerDelegate(c.GetDelegate(), true)
			c.lastshoot = time.Now()
		}
	}
}
