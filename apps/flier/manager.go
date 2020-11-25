package main

import (
	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
	"github.com/jrecuero/engosdl/assets/entities"
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
		playerSprite := components.NewSprite("player-sprite", []string{"images/plane.png"}, 1, engosdl.FormatPNG)
		// playerSprite.DefaultAddDelegateToRegister()
		playerSprite.AddDelegateToRegister(nil, nil, &components.OutOfBounds{}, playerSprite.DefaultOnOutOfBounds)
		playerSprite.AddDelegateToRegister(engosdl.GetDelegateManager().GetCollisionDelegate(), nil, nil, func(params ...interface{}) bool {
			c := playerSprite
			collisionEntityOne := params[0].(*engosdl.Entity)
			collisionEntityTwo := params[1].(*engosdl.Entity)
			if c.GetEntity().GetID() == collisionEntityOne.GetID() || c.GetEntity().GetID() == collisionEntityTwo.GetID() {
				if collisionEntityOne.GetTag() == "wall" || collisionEntityTwo.GetTag() == "wall" {
					engosdl.GetEngine().DestroyEntity(c.GetEntity())
				}
			}
			return true
		})
		playerKeyboard := components.NewKeyboard("player-keyboard")
		playerKeyboard.DefaultAddDelegateToRegister()
		playerMoveIt := components.NewMoveIt("player-move-it", engosdl.NewVector(5, 5))
		playerMoveIt.DefaultAddDelegateToRegister()
		playerCollider2D := components.NewCollider2D("player-collider")
		h.player.AddComponent(playerSprite)
		h.player.AddComponent(playerKeyboard)
		h.player.AddComponent(playerMoveIt)
		h.player.AddComponent(playerCollider2D)

		obstacle1 := engosdl.NewEntity("obstacle")
		obstacle1.SetTag("wall")
		obstacle1.GetTransform().SetPosition(engosdl.NewVector(600, 150))
		// obstacle1.GetTransform().SetScale(engosdl.NewVector(1, 4))
		obstacleSprite := components.NewSprite("obstacle/`-sprite", []string{"images/cube.bmp"}, 1, engosdl.FormatBMP)
		obstacleSprite.AddDelegateToRegister(nil, nil, &components.OutOfBounds{}, obstacleSprite.DefaultOnOutOfBounds)
		obstacleOutOfBounds := components.NewOutOfBounds("obstacle/1-out-of-bounds", false)
		obstacleOutOfBounds.DefaultAddDelegateToRegister()
		obstacleMoveTo := components.NewMove("obstacle/1-move-to", func() components.NextMoveT {
			counter := 0
			speed := 5.0
			return func(c engosdl.IComponent) *engosdl.Vector {
				position := engosdl.NewVector(0, 0)
				mark := counter % 200
				if mark%2 == 0 {
				} else if mark < 50 {
					position.X = -1 * speed
				} else if mark < 100 {
					position.Y = -1 * speed
				} else if mark < 150 {
					position.X = speed
				} else {
					position.Y = speed
				}
				counter++
				return position
			}
		}())
		obstacleMoveTo.DefaultAddDelegateToRegister()
		obstacleCollider2D := components.NewCollider2D("obstacle/1-collider-2D")
		obstacle1.AddComponent(obstacleSprite)
		obstacle1.AddComponent(obstacleOutOfBounds)
		obstacle1.AddComponent(obstacleMoveTo)
		obstacle1.AddComponent(obstacleCollider2D)

		obstacle2 := entities.NewBody2D("obstable/2",
			[]string{"images/cube.bmp"},
			1,
			engosdl.FormatBMP,
			false,
			engosdl.NewVector(-2, 0))
		obstacle2.SetTag("wall")
		obstacle2.GetTransform().SetPosition(engosdl.NewVector(400, 200))
		// obstacle.GetTransform().SetScale(engosdl.NewVector(4, 4))

		scene.AddEntity(h.player)
		scene.AddEntity(obstacle1)
		scene.AddEntity(obstacle2)
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
