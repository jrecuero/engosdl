package main

import (
	"fmt"
	"log"
	"strconv"

	"net/http"
	_ "net/http/pprof"

	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
	"github.com/veandco/go-sdl2/sdl"
)

func createAssets(engine *engosdl.Engine) {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Scenes
	titleScene := engosdl.NewScene("title-scene", "title")
	titleScene.SetSceneCode(createSceneTitle)
	playScene := engosdl.NewScene("play-scene", "battle")
	playScene.SetSceneCode(createScenePlay)
	statsScene := engosdl.NewScene("stats-scene", "stats")
	statsScene.SetSceneCode(createSceneStats)

	// Add scenes to engine
	engine.AddScene(titleScene)
	engine.AddScene(playScene)
	engine.AddScene(statsScene)
}

// createBackground creates the background
func createBackground(engine *engosdl.Engine) engosdl.IEntity {
	bg := engosdl.NewEntity("background")
	bg.SetLayer(engosdl.LayerBackground)
	bgSprite := components.NewScrollSprite("bg-sprite", "images/space.bmp", engine.GetRenderer())
	bgSprite.DefaultAddDelegateToRegister()
	bgSprite.SetScroll(engosdl.NewVector(0, -1))
	// bgSprite.SetCamera(&sdl.Rect{X: 0, Y: 0, W: 400, H: 800})
	bgMoveTo := components.NewMoveTo("bg-move", engosdl.NewVector(0, -5))
	bgMoveTo.DefaultAddDelegateToRegister()
	bg.AddComponent(bgSprite)
	bg.AddComponent(bgMoveTo)
	return bg
}

// createPlayer creates the player entity
func createPlayer(engine *engosdl.Engine) engosdl.IEntity {
	player := engosdl.NewEntity("player")
	player.GetTransform().SetPosition(engosdl.NewVector(float64(engine.GetWidth())/2, float64(engine.GetHeight()-125)))

	// playerSprite := components.NewSprite("player-sprite", "images/player.bmp", engine.GetRenderer())
	playerSprite := components.NewSprite("player-sprite", []string{"images/player.bmp"}, 1, engine.GetRenderer())
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

	player.AddComponent(playerSprite)
	player.AddComponent(playerKeyboard)
	player.AddComponent(playerKeyShooter)
	player.AddComponent(playerShootBullet)
	player.AddComponent(playerOutOfBounds)
	player.AddComponent(playerMoveIt)

	player.DoDump()

	return player
}

// createScenePlay creates the play scene.
func createScenePlay(engine *engosdl.Engine, scene engosdl.IScene) bool {
	maxEnemies := 2

	// Entities
	bg := createBackground(engine)
	player := createPlayer(engine)
	enemyController := createEnemyController(maxEnemies)
	enemies := createEnemies(engine, maxEnemies, enemyController)
	score := createScore(engine)
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
	// winner := engosdl.NewEntity("winner")
	// winnerKeyboard := components.NewKeyboard("winner-keyboard")
	// winnerKeyboard.DefaultAddDelegateToRegister()
	// winnerText := components.NewText("winner-text", "fonts/lato.ttf", 24, sdl.Color{R: 0, G: 0, B: 255}, "You Won...type any key", engosdl.GetRenderer())
	// winnerText.DefaultAddDelegateToRegister()
	// winnerText.AddDelegateToRegister(nil, nil, &components.Keyboard{}, func(params ...interface{}) bool {
	// 	key := params[0].(int)
	// 	if key == sdl.SCANCODE_RETURN {
	// 		engosdl.GetEngine().GetSceneManager().SetActiveFirstScene()
	// 	}
	// 	return true
	// })
	// winner.AddComponent(winnerKeyboard)
	// winner.AddComponent(winnerText)
	// winner.SetActive(false)
	// winner.SetTag("winner")

	// Add entities to scene
	scene.AddEntity(bg)
	scene.AddEntity(player)
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
func createSceneStats(engine *engosdl.Engine, scene engosdl.IScene) bool {
	message := engosdl.NewEntity("message")
	message.GetTransform().SetPosition(engosdl.NewVector(175, 100))
	messageText := components.NewText("message-text", "fonts/lato.ttf", 16, sdl.Color{R: 0, G: 255, B: 0}, "player stats", engine.GetRenderer())
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

// createSceneTitle creates title scene.
func createSceneTitle(engine *engosdl.Engine, scene engosdl.IScene) bool {
	title := engosdl.NewEntity("title")
	title.GetTransform().SetPosition(engosdl.NewVector(175, 250))
	titleText := components.NewText("title-text", "fonts/lato.ttf", 32, sdl.Color{R: 0, G: 0, B: 255}, "PLAY", engine.GetRenderer())
	titleText.DefaultAddDelegateToRegister()
	titleText.AddDelegateToRegister(nil, nil, &components.Keyboard{}, func(params ...interface{}) bool {
		key := params[0].(int)
		if key == sdl.SCANCODE_RETURN {
			engosdl.GetEngine().GetSceneManager().SetActiveNextScene()
		}
		return true
	})
	titleKeyboard := components.NewKeyboard("title-keyboard")
	titleKeyboard.DefaultAddDelegateToRegister()
	titleOutOfBounds := components.NewOutOfBounds("title-out-of-bounds", true)
	titleOutOfBounds.DefaultAddDelegateToRegister()
	titleMoveIt := components.NewMoveIt("title-move-it", engosdl.NewVector(5, 0))
	titleMoveIt.DefaultAddDelegateToRegister()

	title.AddComponent(titleText)
	title.AddComponent(titleKeyboard)
	title.AddComponent(titleOutOfBounds)
	title.AddComponent(titleMoveIt)

	scene.AddEntity(title)
	return true
}

// createScore creates all text entities.
func createScore(engine *engosdl.Engine) engosdl.IEntity {
	scoreValue := 0
	score := engosdl.NewEntity("score")
	score.GetTransform().SetPosition(engosdl.NewVector(10, 560))

	scoreText := components.NewText("score-text", "fonts/lato.ttf", 24, sdl.Color{R: 255, G: 0, B: 0}, "Score: 0000", engine.GetRenderer())
	scoreText.DefaultAddDelegateToRegister()
	destroyDelegate := engosdl.GetDelegateManager().GetDestroyDelegate()
	scoreText.AddDelegateToRegister(destroyDelegate, nil, nil, func(params ...interface{}) bool {
		entity := params[0].(engosdl.IEntity)
		if entity.GetTag() == "enemy" {
			scoreValue += 10
			scoreText.SetMessage("Score: " + strconv.Itoa(scoreValue))
		}
		return true
	})
	score.AddComponent(scoreText)

	return score
}

func main() {
	fmt.Println("engosdl app")
	if engine := engosdl.NewEngine("engosdl app", 400, 600, NewGameManager("app-game-manager")); engine != nil {
		// engine.DoInit()
		// createAssets(engine)
		engine.RunEngine(nil)
	}
}
