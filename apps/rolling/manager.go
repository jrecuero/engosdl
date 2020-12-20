package main

import (
	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
	"github.com/veandco/go-sdl2/sdl"
)

// GameManager is the application game manager.
type GameManager struct {
	*engosdl.GameManager
}

var _ engosdl.IGameManager = (*GameManager)(nil)

// NewGameManager created a new game manager instance.
func NewGameManager(name string) *GameManager {
	engosdl.Logger.Trace().Str("game-manager", name).Msg("new game-manager")
	return &GameManager{
		GameManager: engosdl.NewGameManager(name),
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

func (h *GameManager) createBoard() *Board {
	var cell *Cell
	board := NewBoard("dungeon", 10, 1, engosdl.NewVector(5, 5), 32)

	cell = NewCell(0, 0)
	cell.EnterDialog = "You enter in the dungeon"
	cell.ExitDialog = "Perils follow from this point"
	cell.ActionAndResult["move"] = func(board *Board, entity engosdl.IEntity, position *Position, params ...interface{}) (string, error) {
		newPos := &Position{Row: position.Row, Col: position.Col + 1}
		board.DeleteEntityAt(position.Row, position.Col)
		board.AddEntityAt(entity, newPos.Row, newPos.Col, true)
		return "Move to the next cell", nil
	}
	board.Cells[0][0] = cell

	cell = NewCell(0, 1)
	cell.EnterDialog = "There is an enemy here"
	cell.ExitDialog = "You beat the enemy. You can continue"
	cell.ActionAndResult["look"] = func(board *Board, entity engosdl.IEntity, position *Position, params ...interface{}) (string, error) {
		return "Enemy is ready to fight", nil
	}
	cell.ActionAndResult["move"] = func(board *Board, entity engosdl.IEntity, position *Position, params ...interface{}) (string, error) {
		newPos := &Position{Row: position.Row, Col: position.Col + 1}
		board.DeleteEntityAt(position.Row, position.Col)
		board.AddEntityAt(entity, newPos.Row, newPos.Col, true)
		return "Move to the next cell", nil
	}
	cell.ActionAndResult["attack"] = func(board *Board, entity engosdl.IEntity, position *Position, params ...interface{}) (string, error) {
		return "Attack enemy", nil
	}
	board.Cells[0][1] = cell

	cell = NewCell(0, 2)
	cell.EnterDialog = "You reach your destination"
	cell.ExitDialog = "This is the end"
	cell.ActionAndResult["look"] = func(board *Board, entity engosdl.IEntity, position *Position, params ...interface{}) (string, error) {
		return "You made it!", nil
	}
	board.Cells[0][2] = cell
	return board
}

func (h *GameManager) createScenePlay() func(*engosdl.Engine, engosdl.IScene) bool {
	return func(engine *engosdl.Engine, scene engosdl.IScene) bool {
		sceneController := engosdl.NewEntity("scene-controller")
		sceneController.AddComponent(NewSceneController("scene-controller/component"))

		// board := engosdl.NewEntity("board")
		board := sceneController.GetComponent(&SceneController{}).(*SceneController).Board
		board.GetTransform().SetPositionXY(5, 5)
		board.AddComponent(h.createBoard())

		// player := NewPlayer("player")
		player := sceneController.GetComponent(&SceneController{}).(*SceneController).Player
		player.SetTag("player")
		mouse := components.NewMouse("player/mouse", true)
		player.AddComponent(mouse)
		player.AddComponent(components.NewBox("player/box", &engosdl.Rect{W: 32, H: 32}, sdl.Color{B: 125, A: 255}, true))
		row, col := 0, 0
		board.GetComponent(&Board{}).(*Board).AddEntityAt(player, row, col, true)

		playerController := engosdl.NewComponent("player/controller")
		player.AddComponentExt(playerController, player)
		// playerController.AddDelegateToRegister(nil, board, &Board{}, func(c engosdl.IComponent) func(params ...interface{}) bool {
		// 	return func(params ...interface{}) bool {
		// 		actions := params[0].([]string)
		// 		c.GetEntity().(*Player).UpdateActions(actions)
		// 		return true
		// 	}
		// }(playerController))

		// console := engosdl.NewEntity("console")
		console := sceneController.GetComponent(&SceneController{}).(*SceneController).Console
		console.GetTransform().SetPositionXY(10, 200)
		// console.AddComponent(components.NewBox("console/box", &engosdl.Rect{}, sdl.Color{}, false))
		consoleText := components.NewText("console/text", "fonts/fira2.ttf", 12, sdl.Color{}, " ")
		console.AddComponent(consoleText)
		// message := ""

		// lookButton := engosdl.NewEntity("look")
		// moveButton := engosdl.NewEntity("move")
		// attackButton := engosdl.NewEntity("attack")
		lookButton := player.GetChildByName("look")
		moveButton := player.GetChildByName("move")
		attackButton := player.GetChildByName("attack")

		lookButton.GetTransform().SetPositionXY(10, 50)
		lookButton.AddComponent(components.NewButton("loo/button", "fonts/fira.ttf", 32, sdl.Color{B: 255}, "LOOK", &engosdl.Rect{}, sdl.Color{B: 255}, false))
		// lookButton.GetComponent(&components.Button{}).AddDelegateToRegister(nil, player, &components.Mouse{}, func(c engosdl.IComponent) func(params ...interface{}) bool {
		// 	return func(params ...interface{}) bool {
		// 		mousePos := engosdl.NewVector(float64(params[0].(int32)), float64(params[1].(int32)))
		// 		if c.GetEntity().IsInside(mousePos) {
		// 			if c.GetEnabled() {
		// 				if output, err := board.GetComponent(&Board{}).(*Board).ExecuteAtPlayerPos("look"); err == nil {
		// 					message += output + "\n"
		// 					consoleText.SetMessage(message)
		// 				}
		// 			}
		// 		}
		// 		return true
		// 	}
		// }(lookButton.GetComponent(&components.Button{})))
		player.AddChild(lookButton)

		moveButton.GetTransform().SetPositionXY(100, 50)
		moveButton.AddComponent(components.NewButton("loo/button", "fonts/fira.ttf", 32, sdl.Color{B: 255}, "MOVE", &engosdl.Rect{}, sdl.Color{B: 255}, false))
		// moveButton.GetComponent(&components.Button{}).AddDelegateToRegister(nil, player, &components.Mouse{}, func(c engosdl.IComponent) func(params ...interface{}) bool {
		// 	return func(params ...interface{}) bool {
		// 		mousePos := engosdl.NewVector(float64(params[0].(int32)), float64(params[1].(int32)))
		// 		if c.GetEntity().IsInside(mousePos) {
		// 			if c.GetEnabled() {
		// 				if output, err := board.GetComponent(&Board{}).(*Board).ExecuteAtPlayerPos("move"); err == nil {
		// 					message += output + "\n"
		// 					consoleText.SetMessage(message)
		// 				}
		// 			}
		// 		}
		// 		return true
		// 	}
		// }(moveButton.GetComponent(&components.Button{})))
		player.AddChild(moveButton)

		attackButton.GetTransform().SetPositionXY(210, 50)
		attackButton.AddComponent(components.NewButton("attack/button", "fonts/fira.ttf", 32, sdl.Color{B: 255}, "ATTACK", &engosdl.Rect{}, sdl.Color{B: 255}, false))
		// attackButton.GetComponent(&components.Button{}).AddDelegateToRegister(nil, player, &components.Mouse{}, func(c engosdl.IComponent) func(params ...interface{}) bool {
		// 	return func(params ...interface{}) bool {
		// 		mousePos := engosdl.NewVector(float64(params[0].(int32)), float64(params[1].(int32)))
		// 		if c.GetEntity().IsInside(mousePos) {
		// 			if c.GetEnabled() {
		// 				if output, err := board.GetComponent(&Board{}).(*Board).ExecuteAtPlayerPos("attack"); err == nil {
		// 					message += output + "\n"
		// 					consoleText.SetMessage(message)
		// 				}
		// 			}
		// 		}
		// 		return true
		// 	}
		// }(attackButton.GetComponent(&components.Button{})))
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
