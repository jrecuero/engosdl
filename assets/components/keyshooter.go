package components

import (
	"reflect"
	"time"

	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// ComponentNameKeyShooter is the name to refer key shooter component.
var ComponentNameKeyShooter string = reflect.TypeOf(&KeyShooter{}).String()

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor(ComponentNameKeyShooter, CreateKeyShooter)
	}
}

// KeyShooter represents a component that can create a entity when
// a key is being pressed.
type KeyShooter struct {
	*engosdl.Component
	Key       int           `json:"key"`
	Cooldown  time.Duration `json:"cooldown"`
	lastshoot time.Time
}

// NewKeyShooter creates a new keyshooter instance
// It creates delegate "on-shoot".
func NewKeyShooter(name string, key int) *KeyShooter {
	engosdl.Logger.Trace().Str("component", "key-shooter").Str("key-shooter", name).Msg("new key-shooter")
	keyShooter := &KeyShooter{
		Component: engosdl.NewComponent(name),
		Key:       key,
		Cooldown:  500 * time.Millisecond,
	}
	return keyShooter
}

// CreateKeyShooter implements key shooter constructor used by component
// manager.
func CreateKeyShooter(params ...interface{}) engosdl.IComponent {
	return NewKeyShooter(params[0].(string), params[0].(int))
}

// OnAwake should create all component resources that don't have any dependency
// with any other component or entity.
func (c *KeyShooter) OnAwake() {
	engosdl.Logger.Trace().Str("component", "key-shooter").Str("key-shooter", c.GetName()).Msg("OnAwake")
	// Create new delegate "shoot"
	c.SetDelegate(engosdl.GetDelegateManager().CreateDelegate(c, "on-shoot"))
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
		if time.Since(c.lastshoot) >= c.Cooldown {
			engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), true)
			c.lastshoot = time.Now()
		}
	}
}
