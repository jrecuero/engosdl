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
	x := c.GetEntity().GetTransform().GetPosition().X
	y := c.GetEntity().GetTransform().GetPosition().Y
	c.counter++
	bullet := engosdl.NewEntity("bullet" + strconv.Itoa(c.counter))
	bullet.GetTransform().SetPosition(engosdl.NewVector(x, y))
	bullet.SetDieOnCollision(true)
	bulletSprite := NewSprite("bullet-sprite", "images/player_bullet.bmp", engosdl.GetEngine().GetRenderer())
	bulletMoveTo := NewMoveTo("bullet-moveto", engosdl.NewVector(0, -5))
	bulletOutOfBounds := NewOutOfBounds("bullet-out-of-bounds")
	bulletCollider2D := NewCollider2D("bullet-collider-2D")
	bullet.AddComponent(bulletSprite)
	bullet.AddComponent(bulletMoveTo)
	bullet.AddComponent(bulletOutOfBounds)
	bullet.AddComponent(bulletCollider2D)
	c.GetEntity().GetScene().AddEntity(bullet)
	// bullet.OnStart()
	return true
}
