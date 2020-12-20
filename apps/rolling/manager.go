package main

import (
	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
	"github.com/veandco/go-sdl2/sdl"
)

// GameManager is the application game manager.
type GameManager struct {
	*engosdl.GameManager
	Player *Player
}

var _ engosdl.IGameManager = (*GameManager)(nil)

// NewGameManager created a new game manager instance.
func NewGameManager(name string) *GameManager {
	engosdl.Logger.Trace().Str("game-manager", name).Msg("new game-manager")
	return &GameManager{
		GameManager: engosdl.NewGameManager(name),
		Player:      NewPlayer("player"),
	}
}

// CreateAssets creates all application assets and resources. it is called
// before game engine starts in order to create all required assets and
// resources.
func (h *GameManager) CreateAssets() {
	scenePlay := engosdl.NewScene("play-scene", "play")
	scenePlay.SetSceneCode(h.createScenePlay())
	engosdl.GetEngine().AddScene(scenePlay)
}

func (h *GameManager) createScenePlay() func(*engosdl.Engine, engosdl.IScene) bool {
	return func(engine *engosdl.Engine, scene engosdl.IScene) bool {
		sceneController := engosdl.NewEntity("scene-controller")
		sceneController.AddComponent(NewSceneController("scene-controller/component", h.Player))

		board := sceneController.GetComponent(&SceneController{}).(*SceneController).Board
		board.GetTransform().SetPositionXY(5, 5)

		player := sceneController.GetComponent(&SceneController{}).(*SceneController).Player
		player.SetTag("player")
		mouse := components.NewMouse("player/mouse", true)
		player.AddComponent(mouse)
		player.AddComponent(components.NewBox("player/box", &engosdl.Rect{W: 32, H: 32}, sdl.Color{B: 125, A: 255}, true))
		row, col := 0, 0
		board.GetComponent(&Board{}).(*Board).AddEntityAt(player, row, col, true)

		playerController := engosdl.NewComponent("player/controller")
		player.AddComponentExt(playerController, player)

		console := sceneController.GetComponent(&SceneController{}).(*SceneController).Console
		console.GetTransform().SetPositionXY(10, 200)
		consoleText := components.NewText("console/text", "fonts/fira2.ttf", 12, sdl.Color{}, " ")
		console.AddComponent(consoleText)

		lookButton := player.GetChildByName("look")
		lookButton.GetTransform().SetPositionXY(10, 50)
		lookButton.AddComponent(components.NewButton("loo/button", "fonts/fira.ttf", 32, sdl.Color{B: 255}, "LOOK", &engosdl.Rect{}, sdl.Color{B: 255}, false))
		player.AddChild(lookButton)

		moveButton := player.GetChildByName("move")
		moveButton.GetTransform().SetPositionXY(100, 50)
		moveButton.AddComponent(components.NewButton("loo/button", "fonts/fira.ttf", 32, sdl.Color{B: 255}, "MOVE", &engosdl.Rect{}, sdl.Color{B: 255}, false))
		player.AddChild(moveButton)

		attackButton := player.GetChildByName("attack")
		attackButton.GetTransform().SetPositionXY(210, 50)
		attackButton.AddComponent(components.NewButton("attack/button", "fonts/fira.ttf", 32, sdl.Color{B: 255}, "ATTACK", &engosdl.Rect{}, sdl.Color{B: 255}, false))
		player.AddChild(attackButton)

		sceneController.GetComponent(&SceneController{}).(*SceneController).SetupResources()

		scene.AddEntity(sceneController)
		scene.AddEntity(board)
		scene.AddEntity(player)
		scene.AddEntity(console)
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
}
