package main

import (
	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/exp/rand"
)

// GameManager is the flier application game manager.
type GameManager struct {
	*engosdl.GameManager
	player *engosdl.Entity
	score  *engosdl.Entity
}

var _ engosdl.IGameManager = (*GameManager)(nil)

// NewGameManager created a new game manager instance.
func NewGameManager(name string) *GameManager {
	engosdl.Logger.Trace().Str("game-manager", name).Msg("new game-manager")
	return &GameManager{
		GameManager: engosdl.NewGameManager(name),
	}
}

// CreateAssets creates all flier assets and resources. it is called before
// game engine starts in order to create all required assets and resources.
func (h *GameManager) CreateAssets() {
	playScene := engosdl.NewScene("flier-play-scene-1", "play")
	playScene.SetSceneCode(h.createScenePlay())
	engosdl.GetEngine().AddScene(playScene)
}

// createBullet creates every bullet entity being shoot by the player.
func (h *GameManager) createBullet() engosdl.IEntity {
	bullet := engosdl.NewEntity("bullet")
	bullet.SetTag("bullet")
	bullet.SetLayer(engosdl.LayerBottom)
	box := components.NewBox("bullet-box", &sdl.Rect{W: 10, H: 10}, sdl.Color{R: 255}, true)
	box.DefaultAddDelegateToRegister()
	body := components.NewBody("bullet-body", true)
	move := components.NewMoveTo("bullet-move", engosdl.NewVector(10, 0))
	bullet.AddComponent(box)
	bullet.AddComponent(body)
	bullet.AddComponent(move)
	return bullet
}

// createCoin creates every coin in the scene.
func (h *GameManager) createCoin() engosdl.IEntity {
	coin := engosdl.NewEntity("coin")
	coin.GetTransform().SetPositionXY(900.0, 200.0)
	coin.SetTag("coin")
	coin.SetDieOnOutOfBounds(true)
	coin.AddComponent(components.NewBox("coin/box", &sdl.Rect{W: 32, H: 32}, sdl.Color{B: 255}, true))
	coin.AddComponent(components.NewBody("coin/body", true))
	coin.AddComponent(components.NewMoveTo("coin/move-to", engosdl.NewVector(-1, 0)))
	coin.GetComponent(&components.Box{}).AddDelegateToRegister(nil, nil, &components.Body{}, func(params ...interface{}) bool {
		entity := params[0].(engosdl.IEntity)
		if outAt := params[1].(int); outAt == engosdl.Left {
			if entity.GetID() == coin.GetID() {
				engosdl.GetEngine().DestroyEntity(coin)
			}
		}
		return true
	})
	return coin
}

// createScenePlay creates the main scene to play.
func (h *GameManager) createScenePlay() func(engine *engosdl.Engine, scene engosdl.IScene) bool {
	return func(engine *engosdl.Engine, scene engosdl.IScene) bool {
		h.player.GetTransform().SetPosition(engosdl.NewVector(100, 100))
		// playerSprite := components.NewSprite("player-sprite", []string{"images/plane.png"}, 1, engosdl.FormatPNG)
		playerSprite := components.NewBox("player-box", &sdl.Rect{X: 0, Y: 0, W: 64, H: 64}, sdl.Color{R: 255}, true)
		// playerSprite.DefaultAddDelegateToRegister()
		playerSprite.AddDelegateToRegister(nil, nil, &components.OutOfBounds{}, playerSprite.DefaultOnOutOfBounds)
		// playerSprite.AddDelegateToRegister(engosdl.GetDelegateManager().GetCollisionDelegate(), nil, nil, func(params ...interface{}) bool {
		// 	c := playerSprite
		// 	collisionEntityOne := params[0].(*engosdl.Entity)
		// 	collisionEntityTwo := params[1].(*engosdl.Entity)
		// 	if c.GetEntity().GetID() == collisionEntityOne.GetID() || c.GetEntity().GetID() == collisionEntityTwo.GetID() {
		// 		if collisionEntityOne.GetTag() == "wall" || collisionEntityTwo.GetTag() == "wall" {
		// 			engosdl.GetEngine().DestroyEntity(c.GetEntity())
		// 		}
		// 	}
		// 	return true
		// })
		playerKeyboard := components.NewKeyboard("player-keyboard", components.KeyboardStandardMoveAndShoot)
		playerKeyboard.DefaultAddDelegateToRegister()
		playerMoveIt := components.NewMoveIt("player-move-it", engosdl.NewVector(5, 5))
		playerMoveIt.DefaultAddDelegateToRegister()
		playerMoveIt.AddDelegateToRegister(engosdl.GetDelegateManager().GetCollisionDelegate(), nil, nil, func(params ...interface{}) bool {
			c := playerMoveIt
			collisionEntityOne := params[0].(*engosdl.Entity)
			collisionEntityTwo := params[1].(*engosdl.Entity)
			if c.GetEntity().GetID() == collisionEntityOne.GetID() || c.GetEntity().GetID() == collisionEntityTwo.GetID() {
				if collisionEntityOne.GetTag() == "wall" || collisionEntityTwo.GetTag() == "wall" {
					x, y := c.GetEntity().GetTransform().GetPosition().Get()
					c.GetEntity().GetTransform().SetPosition(engosdl.NewVector(x-c.LastMove.X, y-c.LastMove.Y))
				}
			}
			return true
		})
		playerCollider2D := components.NewCollider2D("player-collider")
		playerShooter := components.NewShooter("player-shooter", h.createBullet)
		playerShooter.AddDelegateToRegister(nil, nil, &components.Keyboard{}, playerShooter.ShooterSignature)
		h.player.AddComponent(playerSprite)
		h.player.AddComponent(playerKeyboard)
		h.player.AddComponent(playerMoveIt)
		h.player.AddComponent(playerCollider2D)
		h.player.AddComponent(playerShooter)

		waller := engosdl.NewEntity("waller")
		waller.AddComponent(components.NewTimer("waller-timer", 100))
		wallerCaller := engosdl.NewComponent("waller-caller")
		wallerCaller.AddDelegateToRegister(nil, nil, &components.Timer{}, func(params ...interface{}) bool {
			scene.AddEntity(h.createWall())
			return true
		})
		waller.AddComponent(wallerCaller)

		coiner := engosdl.NewEntity("coiner")
		coiner.AddComponent(components.NewTimer("coiner-timer", 200))
		coinerCaller := engosdl.NewComponent("cointer-caller")
		coinerCaller.AddDelegateToRegister(nil, nil, &components.Timer{}, func(params ...interface{}) bool {
			scene.AddEntity(h.createCoin())
			return true
		})
		coiner.AddComponent(coinerCaller)

		scene.AddEntity(h.player)
		scene.AddEntity(waller)
		scene.AddEntity(coiner)

		scene.SetCollisionMode(engosdl.ModeBox)
		return true
	}
}

// createWall creates every wall in the scene.
func (h *GameManager) createWall() engosdl.IEntity {
	length := rand.Intn(5)
	up := rand.Intn(2)
	pos := 0
	if up != 0 {
		pos = 400 - length*64
	}
	wall := engosdl.NewEntity("wall2")
	wall.GetTransform().SetPosition(engosdl.NewVector(900, float64(pos)))
	wall.GetTransform().SetScaleXY(1, float64(length))
	wall.SetTag("wall")
	wall.SetDieOnOutOfBounds(true)
	wall.AddComponent(components.NewBox("wall2/box", &sdl.Rect{W: 64, H: 64}, sdl.Color{G: 255, A: 255}, true))
	wall.AddComponent(components.NewCollider2D("wall2/collider-2d"))
	wall.AddComponent(components.NewMoveTo("wall2/move-to", engosdl.NewVector(-2, 0)))
	wall.AddComponent(components.NewOutOfBounds("wall2/out-of-bounds", false))
	wall.GetComponent(&components.Box{}).AddDelegateToRegister(nil, nil, &components.OutOfBounds{}, func(params ...interface{}) bool {
		entity := params[0].(engosdl.IEntity)
		outAt := params[1].(int)
		if outAt == engosdl.Left {
			if entity.GetID() == wall.GetID() {
				engosdl.GetEngine().DestroyEntity(wall)
			}
		}
		return true
	})
	return wall
}

// DoFrameEnd is called at the end of every engine frame.
func (h *GameManager) DoFrameEnd() {
}

// DoFrameStart is called at the start of the game frame.
func (h *GameManager) DoFrameStart() {
}

// DoInit initializes internal game manager resources.
func (h *GameManager) DoInit() {
	h.player = engosdl.NewEntity("player")
	h.player.SetTag("player")

	h.score = engosdl.NewEntity("score")
}
