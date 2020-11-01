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

// NewShootBullet creates a instance of shoot bullet
func NewShootBullet(name string) *ShootBullet {
	engosdl.Logger.Trace().Str("component", "shootbullet").Str("shootbullet", name).Msg("new shoot-bullet")
	shootBullet := &ShootBullet{
		Component: engosdl.NewComponent(name),
	}

	return shootBullet
}

// OnStart is called first time the component is enabled.
func (c *ShootBullet) OnStart() {
	for _, component := range c.GetEntity().GetComponents() {
		for _, delegate := range component.GetDelegates() {
			if delegate.GetEventName() == "shoot" {
				// shootBullet.delegate = delegate
				engosdl.GetEngine().GetEventHandler().GetDelegateHandler().RegisterToDelegate(delegate, c.shootBulletSignature)
			}
		}
	}
}

func (c *ShootBullet) shootBulletSignature(...interface{}) bool {
	x, y := c.GetEntity().GetTransform().GetPosition().Get()
	w, h := c.GetEntity().GetTransform().GetDim().Get()
	c.counter++
	bullet := engosdl.NewEntity("bullet" + strconv.Itoa(c.counter))
	bulletSprite := NewSprite("bullet-sprite", "images/player_bullet.bmp", engosdl.GetEngine().GetRenderer(), true)
	bulletMoveTo := NewMoveTo("bullet-move-to", engosdl.NewVector(0, -5))
	bulletOutOfBounds := NewOutOfBounds("bullet-out-of-bounds")
	bulletCollider2D := NewCollider2D("bullet-collider-2D")
	bullet.AddComponent(bulletSprite)
	bullet.AddComponent(bulletMoveTo)
	bullet.AddComponent(bulletOutOfBounds)
	bullet.AddComponent(bulletCollider2D)
	bullet.GetTransform().SetPosition(engosdl.NewVector((x + w/2), (y + h/2)))
	bullet.SetDieOnCollision(true)
	c.GetEntity().GetScene().AddEntity(bullet)
	return true
}
