package components

import (
	"reflect"
	"strconv"

	"github.com/jrecuero/engosdl"
)

// ComponentNameShootBullet is the name to refer shoot bullet component.
var ComponentNameShootBullet string = reflect.TypeOf(&ShootBullet{}).String()

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor(ComponentNameShootBullet, CreateShootBullet)
	}
}

// ShootBullet represents a component that shoot a bullet.
type ShootBullet struct {
	*engosdl.Component
	// delegate engosdl.IDelegate
	Speed   *engosdl.Vector
	counter int
}

// NewShootBullet creates a instance of shoot bullet.
// It registers to "on-shoot" delegate.
func NewShootBullet(name string, speed *engosdl.Vector) *ShootBullet {
	engosdl.Logger.Trace().Str("component", "shoot-bullet").Str("shoot-bullet", name).Msg("new shoot-bullet")
	result := &ShootBullet{
		Component: engosdl.NewComponent(name),
		Speed:     speed,
	}
	return result
}

// CreateShootBullet implements shoot bullet constructor used bu component
// manager.
func CreateShootBullet(params ...interface{}) engosdl.IComponent {
	if len(params) == 2 {
		return NewShootBullet(params[0].(string), params[1].(*engosdl.Vector))
	}
	return NewShootBullet("", engosdl.NewVector(0, 0))
}

// DefaultAddDelegateToRegister will proceed to add default delegate to
// register for the component.
func (c *ShootBullet) DefaultAddDelegateToRegister() {
	c.AddDelegateToRegister(nil, nil, &KeyShooter{}, c.ShootBulletSignature)
}

// OnStart is called first time the component is enabled.
func (c *ShootBullet) OnStart() {
	engosdl.Logger.Trace().Str("component", "shoot-bullet").Str("shoot-bullet", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}

// ShootBulletSignature trigger a bullet.
func (c *ShootBullet) ShootBulletSignature(...interface{}) bool {
	x, y := c.GetEntity().GetTransform().GetPosition().Get()
	w, h := c.GetEntity().GetTransform().GetDim().Get()
	c.counter++
	bullet := engosdl.NewEntity("bullet" + strconv.Itoa(c.counter))
	bullet.SetTag("bullet")
	bullet.SetParent(c.GetEntity())
	bulletSprite := NewSprite("bullet-sprite", []string{"images/player_bullet.bmp"}, 1, engosdl.FormatBMP)
	bulletSprite.DefaultAddDelegateToRegister()
	bulletMoveTo := NewMoveTo("bullet-move-to", c.Speed)
	bulletMoveTo.DefaultAddDelegateToRegister()
	bulletOutOfBounds := NewOutOfBounds("bullet-out-of-bounds", false)
	bulletOutOfBounds.DefaultAddDelegateToRegister()
	bulletCollider2D := NewCollider2D("bullet-collider-2D")
	bulletCollider2D.DefaultAddDelegateToRegister()
	bullet.SetLayer(engosdl.LayerBottom)
	bullet.AddComponent(bulletSprite)
	bullet.AddComponent(bulletMoveTo)
	bullet.AddComponent(bulletOutOfBounds)
	bullet.AddComponent(bulletCollider2D)
	bulletSprite.LoadSprite()
	bulletW, bulletH := bullet.GetTransform().GetDim().Get()
	bullet.GetTransform().SetPosition(engosdl.NewVector((x + w/2 - bulletW/2), (y + h/2 - bulletH/2)))
	bullet.SetDieOnCollision(false)
	c.GetEntity().GetScene().AddEntity(bullet)
	return true
}

// Unmarshal takes a ComponentToMarshal instance and  creates a new entity
// instance.
func (c *ShootBullet) Unmarshal(data map[string]interface{}) {
	c.Component.Unmarshal(data)
	speed := data["speed"].(map[string]interface{})
	c.Speed = engosdl.NewVector(speed["X"].(float64), speed["Y"].(float64))
}
