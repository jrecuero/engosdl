package main

import (
	"fmt"
	"strconv"

	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
	"github.com/veandco/go-sdl2/sdl"
)

// GameManager is the application game manager.
type GameManager struct {
	*engosdl.GameManager
	engine        *engosdl.Engine
	player        *engosdl.Entity
	score         *engosdl.Entity
	scoreTotal    int
	winnerCreated bool
}

var _ engosdl.IGameManager = (*GameManager)(nil)

// NewGameManager creates a new game manager instance.
func NewGameManager(name string) *GameManager {
	engosdl.Logger.Trace().Str("game-manager", name).Msg("new game-manager")
	return &GameManager{
		GameManager: engosdl.NewGameManager(name),
	}
}

// CreateAssets is called before game engine starts in order to create all
// required assets.
func (h *GameManager) CreateAssets() {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	// Scenes
	titleScene := engosdl.NewScene("title-scene", "title")
	titleScene.SetSceneCode(createSceneTitle)
	playScene := engosdl.NewScene("play-scene-1", "battle")
	// playScene.SetSceneCode(h.createScenePlay)
	playScene.SetSceneCode(h.createScenePlay("images/space.bmp", 3))
	playSceneTwo := engosdl.NewScene("play-scene-2", "battle")
	// playSceneTwo.SetSceneCode(h.createScenePlayTwo)
	playSceneTwo.SetSceneCode(h.createScenePlay("images/space2.bmp", 4))
	statsScene := engosdl.NewScene("stats-scene", "stats")
	statsScene.SetSceneCode(h.createSceneStats)

	// Add scenes to engine
	h.engine.AddScene(titleScene)
	h.engine.AddScene(playScene)
	h.engine.AddScene(playSceneTwo)
	h.engine.AddScene(statsScene)
}

// createEntityPlayer creates the player entity
func (h *GameManager) createEntityPlayer() engosdl.IEntity {
	h.player.GetTransform().SetPosition(engosdl.NewVector(float64(h.engine.GetWidth())/2, float64(h.engine.GetHeight()-125)))
	playerSprite := components.NewSprite("player-sprite", []string{"images/player.bmp"}, 1, engosdl.FormatBMP)
	playerSprite.DefaultAddDelegateToRegister()
	playerKeyboard := components.NewKeyboard("player-keyboard")
	playerKeyboard.DefaultAddDelegateToRegister()
	playerKeyShooter := components.NewKeyShooter("player-key-shooter", sdl.SCANCODE_SPACE)
	playerKeyShooter.DefaultAddDelegateToRegister()
	playerShootBullet := components.NewShootBullet("player-shoot-bullet", engosdl.NewVector(0, -5))
	playerShootBullet.DefaultAddDelegateToRegister()
	playerOutOfBounds := components.NewOutOfBounds("player-out-of-bounds", true)
	playerOutOfBounds.DefaultAddDelegateToRegister()
	playerMoveIt := components.NewMoveIt("player-move-it", engosdl.NewVector(5, 0))
	playerMoveIt.DefaultAddDelegateToRegister()
	playerCollider := components.NewCollider2D("player-collider")
	if obj := h.player.GetComponent(&components.EntityStats{}); obj != nil {
		if stats, ok := obj.(*components.EntityStats); ok {
			stats.DefaultAddDelegateToRegister()
			stats.AddDelegateToRegister(engosdl.GetDelegateManager().GetDestroyDelegate(), nil, nil, func(params ...interface{}) bool {
				entity := params[0].(engosdl.IEntity)
				if entity.GetTag() == "enemy" {
					stats.Experience += 10
				}
				return true
			})
		}
	}

	h.player.AddComponent(playerSprite)
	h.player.AddComponent(playerKeyboard)
	h.player.AddComponent(playerKeyShooter)
	h.player.AddComponent(playerShootBullet)
	h.player.AddComponent(playerOutOfBounds)
	h.player.AddComponent(playerMoveIt)
	h.player.AddComponent(playerCollider)

	return h.player
}

// createEntityScore creates the score entity.
func (h *GameManager) createEntityScore() engosdl.IEntity {
	h.score.GetTransform().SetPosition(engosdl.NewVector(10, 560))
	if obj := h.score.GetComponent(&components.Text{}); obj != nil {
		if text, ok := obj.(*components.Text); ok {
			text.DefaultAddDelegateToRegister()
			destroyDelegate := engosdl.GetDelegateManager().GetDestroyDelegate()
			text.AddDelegateToRegister(destroyDelegate, nil, nil, func(params ...interface{}) bool {
				entity := params[0].(engosdl.IEntity)
				if entity.GetTag() == "enemy" {
					h.scoreTotal += 10
					text.SetMessage("Score: " + strconv.Itoa(h.scoreTotal))
				}
				return true
			})
		}
	}
	return h.score
}

// createScenePlay creates the play scene.
func (h *GameManager) createScenePlay(bgFilename string, maxEnemies int) func(engine *engosdl.Engine, scene engosdl.IScene) bool {
	return func(engine *engosdl.Engine, scene engosdl.IScene) bool {
		h.winnerCreated = false

		// Entities
		bg := createEntityBackground(engine, bgFilename)
		h.createEntityPlayer()
		enemyController := createEnemyController(maxEnemies)
		enemies := createEnemies(engine, maxEnemies, enemyController)
		h.createEntityScore()
		sceneController := engosdl.NewEntity("scene-controller")
		sceneControllerKeyboard := components.NewKeyboard("scene-controller-keyboard")
		sceneControllerKeyboard.DefaultAddDelegateToRegister()
		// sceneControllerComponent := engosdl.NewComponent("scene-controller-controller")
		sceneControllerComponent := components.NewSceneController("scene-controller-component")
		sceneControllerComponent.AddDelegateToRegister(nil, nil, &components.Keyboard{}, func(params ...interface{}) bool {
			key := params[0].(int)
			if key == sdl.SCANCODE_N {
				// if key == sdl.SCANCODE_RETURN {
				engosdl.GetEngine().GetSceneManager().SwapFromSceneTo(engine.GetSceneManager().GetSceneByName("stats-scene"))
			}
			return true
		})
		sceneController.AddComponent(sceneControllerKeyboard)
		sceneController.AddComponent(sceneControllerComponent)

		// Add entities to scene
		scene.AddEntity(bg)
		scene.AddEntity(h.player)
		scene.AddEntity(enemyController)
		for _, enemy := range enemies {
			scene.AddEntity(enemy)
		}
		scene.AddEntity(h.score)
		scene.AddEntity(sceneController)
		// scene.AddEntity(winner)

		return true
	}
}

// createSceneStats creates stats scene
func (h *GameManager) createSceneStats(engine *engosdl.Engine, scene engosdl.IScene) bool {
	message := engosdl.NewEntity("message")
	message.GetTransform().SetPosition(engosdl.NewVector(50, 100))
	messageString := "player stats"
	if playerStats, ok := h.player.GetComponent(&components.EntityStats{}).(*components.EntityStats); ok {
		messageString = fmt.Sprintf("player stats\nlife: %d\nexp: %d", playerStats.Life, playerStats.Experience)
	}
	messageText := components.NewText("message-text", "fonts/lato.ttf", 16, sdl.Color{R: 0, G: 255, B: 0}, messageString)
	messageText.DefaultAddDelegateToRegister()
	messageText.AddDelegateToRegister(nil, nil, &components.Keyboard{}, func(params ...interface{}) bool {
		key := params[0].(int)
		// if key == sdl.SCANCODE_P {
		if key == sdl.SCANCODE_RETURN {
			engosdl.GetEngine().GetSceneManager().SwapBack()
		}
		return true
	})
	messageKeyboard := components.NewKeyboard("title-keyboard")
	messageKeyboard.DefaultAddDelegateToRegister()

	message.AddComponent(messageText)
	message.AddComponent(messageKeyboard)

	scene.AddEntity(message)
	return true
}

// DoFrameEnd is called at the end of the engine frame.
func (h *GameManager) DoFrameEnd() {
}

// DoFrameStart is called at the start of engine frame.
func (h *GameManager) DoFrameStart() {
	// Count number of enemies at the end of a frame. Takes some functionality
	// from the enemy controller.
	activeScene := engosdl.GetSceneManager().GetActiveScene()
	if activeScene.GetTag() == "battle" {
		enemies := activeScene.GetEntitiesByTag("enemy")
		// fmt.Printf("there are %d enemies\n", len(enemies))
		if len(enemies) == 0 && !h.winnerCreated {
			h.winnerCreated = true
			winner := createWinner(engosdl.GetEngine())
			activeScene.AddEntity(winner)
		}
	}
}

// DoInit initializes game manager resources.
func (h *GameManager) DoInit() {
	h.engine = engosdl.GetEngine()
	h.player = engosdl.NewEntity("player")
	h.player.SetTag("player")
	playerStats := components.NewEntityStats("player stats", 100)
	playerStats.SetRemoveOnDestroy(false)
	h.player.AddComponent(playerStats)

	h.score = engosdl.NewEntity("score")
	scoreText := components.NewText("score-text", "fonts/lato.ttf", 24, sdl.Color{R: 255, G: 0, B: 0}, "Score: 0000")
	scoreText.SetRemoveOnDestroy(false)
	h.score.AddComponent(scoreText)
}
