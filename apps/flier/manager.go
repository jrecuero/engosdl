package main

import (
	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
	"github.com/veandco/go-sdl2/sdl"
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
		playerKeyboard := components.NewKeyboard("player-keyboard")
		playerKeyboard.DefaultAddDelegateToRegister()
		playerMoveIt := components.NewMoveIt("player-move-it", engosdl.NewVector(2, 2))
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
		h.player.AddComponent(playerSprite)
		h.player.AddComponent(playerKeyboard)
		h.player.AddComponent(playerMoveIt)
		h.player.AddComponent(playerCollider2D)

		// wall1 := engosdl.NewEntity("wall")
		// wall1.SetTag("wall")
		// wall1.GetTransform().SetPosition(engosdl.NewVector(600, 150))
		// // wall1.GetTransform().SetScale(engosdl.NewVector(1, 4))
		// wallSprite := components.NewSprite("wall/`-sprite", []string{"images/cube.bmp"}, 1, engosdl.FormatBMP)
		// wallSprite.AddDelegateToRegister(nil, nil, &components.OutOfBounds{}, wallSprite.DefaultOnOutOfBounds)
		// wallOutOfBounds := components.NewOutOfBounds("wall/1-out-of-bounds", false)
		// wallOutOfBounds.DefaultAddDelegateToRegister()
		// wallMoveTo := components.NewMove("wall/1-move-to", func() components.NextMoveT {
		// 	counter := 0
		// 	speed := 5.0
		// 	return func(c engosdl.IComponent) *engosdl.Vector {
		// 		position := engosdl.NewVector(0, 0)
		// 		mark := counter % 200
		// 		if mark%2 == 0 {
		// 		} else if mark < 50 {
		// 			position.X = -1 * speed
		// 		} else if mark < 100 {
		// 			position.Y = -1 * speed
		// 		} else if mark < 150 {
		// 			position.X = speed
		// 		} else {
		// 			position.Y = speed
		// 		}
		// 		counter++
		// 		return position
		// 	}
		// }())
		// wallMoveTo.DefaultAddDelegateToRegister()
		// wallCollider2D := components.NewCollider2D("wall/1-collider-2D")
		// wall1.AddComponent(wallSprite)
		// wall1.AddComponent(wallOutOfBounds)
		// wall1.AddComponent(wallMoveTo)
		// wall1.AddComponent(wallCollider2D)

		// wall2 := entities.NewBody2D("obstable/2",
		// 	[]string{"images/cube.bmp"},
		// 	1,
		// 	engosdl.FormatBMP,
		// 	false,
		// 	engosdl.NewVector(-2, 0))
		// wall2.SetTag("wall")
		// wall2.GetTransform().SetPosition(engosdl.NewVector(400, 200))
		// // wall.GetTransform().SetScale(engosdl.NewVector(4, 4))

		wall1 := engosdl.NewEntity("wall1")
		wall1.GetTransform().SetPosition(engosdl.NewVector(400, 200))
		wall1.SetTag("wall")
		wallBox := components.NewBox("wall1/box", &sdl.Rect{X: 0, Y: 0, W: 64, H: 64}, sdl.Color{B: 255}, true)
		wallCollider2D := components.NewCollider2D("wall1/collider-2d")
		wall1.AddComponent(wallBox)
		wall1.AddComponent(wallCollider2D)

		wall2 := engosdl.NewEntity("wall2")
		wall2.GetTransform().SetPosition(engosdl.NewVector(500, 100))
		wall2.SetTag("wall")
		wall2.AddComponent(components.NewBox("wall2/box", &sdl.Rect{W: 64, H: 64}, sdl.Color{G: 255}, true))
		wall2.AddComponent(components.NewCollider2D("wall2/collider-2d"))

		scene.AddEntity(h.player)
		scene.AddEntity(wall1)
		scene.AddEntity(wall2)

		scene.SetCollisionMode(engosdl.ModeBox)
		return true
	}
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
