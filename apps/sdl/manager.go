package main

import (
	"fmt"

	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
	"github.com/veandco/go-sdl2/sdl"
)

// GameManager is the application game manager.
type GameManager struct {
	*engosdl.GameManager
	engine        *engosdl.Engine
	player        *engosdl.Entity
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
	playScene.SetSceneCode(h.createScenePlay)
	playSceneTwo := engosdl.NewScene("play-scene-2", "battle")
	playSceneTwo.SetSceneCode(h.createScenePlayTwo)
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
	playerSprite := components.NewSprite("player-sprite", []string{"images/player.bmp"}, 1, h.engine.GetRenderer())
	playerSprite.SetDestroyOnOutOfBounds(false)
	playerSprite.DefaultAddDelegateToRegister()
	playerKeyboard := components.NewKeyboard("player-keyboard")
	playerKeyboard.DefaultAddDelegateToRegister()
	playerKeyShooter := components.NewKeyShooter("player-key-shooter", sdl.SCANCODE_SPACE)
	playerKeyShooter.DefaultAddDelegateToRegister()
	playerShootBullet := components.NewShootBullet("player-shoot-bullet")
	playerShootBullet.DefaultAddDelegateToRegister()
	playerOutOfBounds := components.NewOutOfBounds("player-out-of-bounds", true)
	playerOutOfBounds.DefaultAddDelegateToRegister()
	playerMoveIt := components.NewMoveIt("player-move-it", engosdl.NewVector(5, 0))
	playerMoveIt.DefaultAddDelegateToRegister()
	if obj := h.player.GetComponent(&components.EntityStats{}); obj != nil {
		if stats, ok := obj.(*components.EntityStats); ok {
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

	h.player.DoDump()

	return h.player
}

// createScenePlay creates the play scene.
func (h *GameManager) createScenePlay(engine *engosdl.Engine, scene engosdl.IScene) bool {
	h.winnerCreated = false
	maxEnemies := 2

	// Entities
	bg := createEntityBackground(engine, "images/space.bmp")
	h.createEntityPlayer()
	enemyController := createEnemyController(maxEnemies)
	enemies := createEnemies(engine, maxEnemies, enemyController)
	score := createEntityScore(engine)
	sceneController := engosdl.NewEntity("scene-controller")
	sceneControllerKeyboard := components.NewKeyboard("scene-controller-keyboard")
	sceneControllerKeyboard.DefaultAddDelegateToRegister()
	sceneControllerComponent := engosdl.NewComponent("scene-controller-controller")
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
	scene.AddEntity(score)
	scene.AddEntity(sceneController)
	// scene.AddEntity(winner)

	return true
}

// createScenePlayTwo creates the play scene.
func (h *GameManager) createScenePlayTwo(engine *engosdl.Engine, scene engosdl.IScene) bool {
	h.winnerCreated = false
	maxEnemies := 1

	// Entities
	bg := createEntityBackground(engine, "images/space2.bmp")
	h.createEntityPlayer()
	enemyController := createEnemyController(maxEnemies)
	enemies := createEnemies(engine, maxEnemies, enemyController)
	score := createEntityScore(engine)
	sceneController := engosdl.NewEntity("scene-controller")
	sceneControllerKeyboard := components.NewKeyboard("scene-controller-keyboard")
	sceneControllerKeyboard.DefaultAddDelegateToRegister()
	sceneControllerComponent := engosdl.NewComponent("scene-controller-controller")
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
	scene.AddEntity(score)
	scene.AddEntity(sceneController)
	// scene.AddEntity(winner)

	return true
}

// createSceneStats creates stats scene
func (h *GameManager) createSceneStats(engine *engosdl.Engine, scene engosdl.IScene) bool {
	message := engosdl.NewEntity("message")
	message.GetTransform().SetPosition(engosdl.NewVector(50, 100))
	messageString := "player stats"
	if playerStats, ok := h.player.GetComponent(&components.EntityStats{}).(*components.EntityStats); ok {
		messageString = fmt.Sprintf("player stats\nlife: %d\nexp: %d", playerStats.Life, playerStats.Experience)
	}
	messageText := components.NewText("message-text", "fonts/lato.ttf", 16, sdl.Color{R: 0, G: 255, B: 0}, messageString, engine.GetRenderer())
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
	playerStats := components.NewEntityStats("player stats", 100)
	playerStats.SetRemoveOnDestroy(false)

	h.player.AddComponent(playerStats)
}
