package components

import (
	"github.com/jrecuero/engosdl"
)

// ShootBullet represents a component that shoot a bullet.
type ShootBullet struct {
	*engosdl.Component
	delegate engosdl.IDelegate
}

func (c *ShootBullet) shootBulletSignature(...interface{}) bool {
	x := c.GetGameObject().GetTransform().GetPosition().X
	y := c.GetGameObject().GetTransform().GetPosition().Y
	bullet := engosdl.NewGameObject("bullet")
	bullet.GetTransform().SetPosition(engosdl.NewVector(x, y))
	bulletSprite := NewSprite("bullet-sprite", bullet, "images/player_bullet.bmp", engosdl.GetEngine().GetRenderer())
	bulletMoveTo := NewMoveTo("bullet-moveto", bullet, engosdl.NewVector(0, -1))
	bulletOutOfBounds := NewOutOfBounds("bullet-out-of-bounds", bullet)
	bullet.AddComponent(bulletSprite)
	bullet.AddComponent(bulletMoveTo)
	bullet.AddComponent(bulletOutOfBounds)
	c.GetGameObject().GetScene().AddGameObject(bullet)
	bullet.OnStart()
	return true
}

// NewShootBullet creates a instance of shoot bullet
func NewShootBullet(name string, gobj *engosdl.GameObject) *ShootBullet {
	engosdl.Logger.Trace().Str("component", "shootbullet").Str("shootbullet", name).Msg("new shoot-bullet")
	shootBullet := &ShootBullet{
		Component: engosdl.NewComponent(name, gobj),
	}
	for _, component := range gobj.GetComponents() {
		for _, delegate := range component.GetDelegates() {
			if delegate.GetEventName() == "shoot" {
				shootBullet.delegate = delegate
				engosdl.GetEngine().GetEventHandler().GetDelegateHandler().RegisterToDelegate(delegate, shootBullet.shootBulletSignature)
			}
		}
	}
	return shootBullet
}
