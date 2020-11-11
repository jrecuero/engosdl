package components

import (
	"strconv"

	"github.com/jrecuero/engosdl"
)

// ShootBullet represents a component that shoot a bullet.
type ShootBullet struct {
	*engosdl.Component
	// delegate engosdl.IDelegate
	counter int
}

// NewShootBullet creates a instance of shoot bullet.
// It registers to "on-shoot" delegate.
func NewShootBullet(name string) *ShootBullet {
	engosdl.Logger.Trace().Str("component", "shoot-bullet").Str("shoot-bullet", name).Msg("new shoot-bullet")
	result := &ShootBullet{
		Component: engosdl.NewComponent(name),
	}
	return result
}

// DefaultAddDelegateToRegister will proceed to add default delegate to
// register for the component.
func (c *ShootBullet) DefaultAddDelegateToRegister() {
	c.AddDelegateToRegister(nil, nil, &KeyShooter{}, c.shootBulletSignature)
}

// OnStart is called first time the component is enabled.
func (c *ShootBullet) OnStart() {
	engosdl.Logger.Trace().Str("component", "shoot-bullet").Str("shoot-bullet", c.GetName()).Msg("OnStart")
	// if component := c.GetEntity().GetComponent(&KeyShooter{}); component != nil {
	// 	if delegate := component.GetDelegate(); delegate != nil {
	// 		c.AddDelegateToRegister(delegate, nil, nil, c.shootBulletSignature)
	// 	}
	// }
	c.Component.OnStart()
}

func (c *ShootBullet) shootBulletSignature(...interface{}) bool {
	x, y := c.GetEntity().GetTransform().GetPosition().Get()
	w, h := c.GetEntity().GetTransform().GetDim().Get()
	c.counter++
	bullet := engosdl.NewEntity("bullet" + strconv.Itoa(c.counter))
	bullet.SetTag("bullet")
	// bulletSprite := NewSprite("bullet-sprite", "images/player_bullet.bmp", engosdl.GetRenderer())
	bulletSprite := NewSprite("bullet-sprite", []string{"images/player_bullet.bmp"}, 1, engosdl.GetRenderer())
	bulletSprite.DefaultAddDelegateToRegister()
	bulletMoveTo := NewMoveTo("bullet-move-to", engosdl.NewVector(0, -5))
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
	bullet.SetDieOnCollision(true)
	c.GetEntity().GetScene().AddEntity(bullet)
	return true
}
