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

	maxEnemies := 2

	// Scenes
	scene := engosdl.NewScene("main scene")

	// Entities
	bg := createBackground(engine)
	player := createPlayer(engine)
	enemyController := createEnemyController()
	// enemy := createEnemy(engine, engosdl.NewVector(200, 10))
	enemies := createEnemies(engine, maxEnemies, enemyController)
	score := createScore(engine)

	// Add entities to scene
	scene.AddEntity(bg)
	scene.AddEntity(player)
	// scene.AddEntity(enemy)
	scene.AddEntity(enemyController)
	for _, enemy := range enemies {
		scene.AddEntity(enemy)
	}
	scene.AddEntity(score)

	// Add scenes to engine
	engine.AddScene(scene)
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
	playerKeyboard := components.NewKeyboard("player-keyboard", engosdl.NewVector(5, 10))
	playerKeyboard.DefaultAddDelegateToRegister()
	playerKeyShooter := components.NewKeyShooter("player-key-shooter", sdl.SCANCODE_SPACE)
	playerKeyShooter.DefaultAddDelegateToRegister()
	playerShootBullet := components.NewShootBullet("player-shoot-bullet")
	playerShootBullet.DefaultAddDelegateToRegister()
	playerOutOfBounds := components.NewOutOfBounds("player-out-of-bounds", true)
	playerOutOfBounds.DefaultAddDelegateToRegister()

	player.AddComponent(playerSprite)
	player.AddComponent(playerKeyboard)
	player.AddComponent(playerKeyShooter)
	player.AddComponent(playerShootBullet)
	player.AddComponent(playerOutOfBounds)

	return player
}

// createScore creates all text entities.
func createScore(engine *engosdl.Engine) engosdl.IEntity {
	scoreValue := 0
	score := engosdl.NewEntity("score")
	score.GetTransform().SetPosition(engosdl.NewVector(10, 560))

	scoreText := components.NewText("score-text", "fonts/lato.ttf", 24, sdl.Color{R: 255, G: 0, B: 0}, "Score: 0000", engine.GetRenderer())
	scoreText.DefaultAddDelegateToRegister()
	destroyDelegate := engosdl.GetDelegateHandler().GetDestroyDelegate()
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
	engine := engosdl.NewEngine("engosdl app", 400, 600)
	engine.DoInit()
	if engine != nil {
		createAssets(engine)
		engine.RunEngine()
	}
}
