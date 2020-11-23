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
		playerSprite.DefaultAddDelegateToRegister()
		playerKeyboard := components.NewKeyboard("player-keyboard")
		playerKeyboard.DefaultAddDelegateToRegister()
		playerMoveIt := components.NewMoveIt("player-move-it", engosdl.NewVector(5, 5))
		playerMoveIt.DefaultAddDelegateToRegister()
		h.player.AddComponent(playerSprite)
		h.player.AddComponent(playerKeyboard)
		h.player.AddComponent(playerMoveIt)

		// obstacle := engosdl.NewEntity("obstacle")
		// obstacle.GetTransform().SetPosition(engosdl.NewVector(400, 200))
		// obstacle.GetTransform().SetScale(engosdl.NewVector(1, 4))
		// obstacleSprite := components.NewSprite("obstacle-sprite", []string{"images/cube.bmp"}, 1, engosdl.FormatBMP)
		// obstacleSprite.AddDelegateToRegister(nil, nil, &components.OutOfBounds{}, obstacleSprite.DefaultOnOutOfBounds)
		// obstacleOutOfBounds := components.NewOutOfBounds("obstacle-out-of-bounds", false)
		// obstacleOutOfBounds.DefaultAddDelegateToRegister()
		// obstacleMoveTo := components.NewMoveTo("obstacle-move-to", engosdl.NewVector(-2, 0))
		// obstacleMoveTo.DefaultAddDelegateToRegister()
		// obstacle.AddComponent(obstacleSprite)
		// obstacle.AddComponent(obstacleOutOfBounds)
		// obstacle.AddComponent(obstacleMoveTo)
		obstacle := entities.NewBody2D("obstable",
			[]string{"images/cube.bmp"},
			1,
			engosdl.FormatBMP,
			false,
			engosdl.NewVector(-2, 0))
		obstacle.GetTransform().SetPosition(engosdl.NewVector(400, 200))
		obstacle.GetTransform().SetScale(engosdl.NewVector(1, 4))

		scene.AddEntity(h.player)
		scene.AddEntity(obstacle)
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
