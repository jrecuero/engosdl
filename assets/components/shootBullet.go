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
// It registers to on-key-shooter delegate.
func NewShootBullet(name string) *ShootBullet {
	engosdl.Logger.Trace().Str("component", "shoot-bullet").Str("shoot-bullet", name).Msg("new shoot-bullet")
	result := &ShootBullet{
		Component: engosdl.NewComponent(name),
	}
	result.AddDelegateToRegister(nil, nil, &KeyShooter{}, result.shootBulletSignature)
	return result
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
	// bulletSprite := NewSprite("bullet-sprite", "images/player_bullet.bmp", engosdl.GetEngine().GetRenderer())
	bulletSprite := NewMultiSprite("bullet-sprite", []string{"images/player_bullet.bmp"}, engosdl.GetEngine().GetRenderer())
	bulletMoveTo := NewMoveTo("bullet-move-to", engosdl.NewVector(0, -5))
	bulletOutOfBounds := NewOutOfBounds("bullet-out-of-bounds", false)
	bulletCollider2D := NewCollider2D("bullet-collider-2D")
	bullet.SetLayer(engosdl.LayerBottom)
	bullet.AddComponent(bulletSprite)
	bullet.AddComponent(bulletMoveTo)
	bullet.AddComponent(bulletOutOfBounds)
	bullet.AddComponent(bulletCollider2D)
	bulletW, bulletH := bullet.GetTransform().GetDim().Get()
	bullet.GetTransform().SetPosition(engosdl.NewVector((x + w/2 - bulletW/2), (y + h/2 - bulletH/2)))
	bullet.SetDieOnCollision(true)
	c.GetEntity().GetScene().AddEntity(bullet)
	return true
}
