package components

import (
	"reflect"

	"github.com/jrecuero/engosdl"
)

// ShooterSignatureT is the type used to pass the function called to create
// a new bullet.
type ShooterSignatureT func() engosdl.IEntity

// ComponentNameShooter is the name to refer shooter component.
var ComponentNameShooter string = reflect.TypeOf(&Shooter{}).String()

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor(ComponentNameShooter, CreateShooter)
	}
}

// Shooter represents a component that shoot any given bullet.
type Shooter struct {
	*engosdl.Component
	// delegate engosdl.IDelegate
	counter   int
	newBullet ShooterSignatureT
}

// NewShooter creates a instance of shooter.
// It registers to "on-shoot" delegate.
func NewShooter(name string, newBullet ShooterSignatureT) *Shooter {
	engosdl.Logger.Trace().Str("component", "shoot-bullet").Str("shoot-bullet", name).Msg("new shoot-bullet")
	result := &Shooter{
		Component: engosdl.NewComponent(name),
		newBullet: newBullet,
	}
	return result
}

// CreateShooter implements shoot bullet constructor used by component
// manager.
func CreateShooter(params ...interface{}) engosdl.IComponent {
	if len(params) == 2 {
		return NewShooter(params[0].(string), params[1].(ShooterSignatureT))
	}
	return NewShooter("", nil)
}

// DefaultAddDelegateToRegister will proceed to add default delegate to
// register for the component.
func (c *Shooter) DefaultAddDelegateToRegister() {
	c.AddDelegateToRegister(nil, nil, &KeyShooter{}, c.ShooterSignature)
}

// OnStart is called first time the component is enabled.
func (c *Shooter) OnStart() {
	engosdl.Logger.Trace().Str("component", "shoot-bullet").Str("shoot-bullet", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}

// ShooterSignature trigger a bullet.
func (c *Shooter) ShooterSignature(...interface{}) bool {
	x, y := c.GetEntity().GetTransform().GetPosition().Get()
	w, h := c.GetEntity().GetTransform().GetDim().Get()
	c.counter++
	// bullet := engosdl.NewEntity("bullet" + strconv.Itoa(c.counter))
	// bullet.SetTag("bullet")
	// bullet.SetParent(c.GetEntity())
	// bulletSprite := NewSprite("bullet-sprite", []string{"images/player_bullet.bmp"}, 1, engosdl.FormatBMP)
	// bulletSprite.DefaultAddDelegateToRegister()
	// bulletMoveTo := NewMoveTo("bullet-move-to", c.Speed)
	// bulletMoveTo.DefaultAddDelegateToRegister()
	// bulletOutOfBounds := NewOutOfBounds("bullet-out-of-bounds", false)
	// bulletOutOfBounds.DefaultAddDelegateToRegister()
	// bulletCollider2D := NewCollider2D("bullet-collider-2D")
	// bulletCollider2D.DefaultAddDelegateToRegister()
	// bullet.SetLayer(engosdl.LayerBottom)
	// bullet.AddComponent(bulletSprite)
	// bullet.AddComponent(bulletMoveTo)
	// bullet.AddComponent(bulletOutOfBounds)
	// bullet.AddComponent(bulletCollider2D)
	// bulletSprite.LoadSprite()
	bullet := c.newBullet()
	bulletW, bulletH := bullet.GetTransform().GetDim().Get()
	bullet.GetTransform().SetPosition(engosdl.NewVector((x + w/2 - bulletW/2), (y + h/2 - bulletH/2)))
	bullet.SetDieOnCollision(false)
	c.GetEntity().GetScene().AddEntity(bullet)
	return true
}

// Unmarshal takes a ComponentToMarshal instance and  creates a new entity
// instance.
func (c *Shooter) Unmarshal(data map[string]interface{}) {
	c.Component.Unmarshal(data)
}
