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
		board := engosdl.NewEntity("board")
		board.GetTransform().SetPositionXY(5, 5)
		board.AddComponent(h.createBoard())

		// player := engosdl.NewEntity("player")
		player := NewPlayer("player")
		player.SetTag("player")
		mouse := components.NewMouse("player/mouse", true)
		player.AddComponent(mouse)
		player.AddComponent(components.NewBox("player/box", &engosdl.Rect{W: 32, H: 32}, sdl.Color{B: 125}, true))
		row, col := 0, 0
		board.GetComponent(&Board{}).(*Board).AddEntityAt(player, row, col, true)

		playerController := engosdl.NewComponent("player/controller")
		player.AddComponentExt(playerController, player)
		// playerController.SetEntity(player)
		playerController.AddDelegateToRegister(nil, board, &Board{}, func(c engosdl.IComponent) func(params ...interface{}) bool {
			return func(params ...interface{}) bool {
				actions := params[0].([]string)
				c.GetEntity().(*Player).UpdateActions(actions)
				return true
			}
		}(playerController))

		console := engosdl.NewEntity("console")
		console.GetTransform().SetPositionXY(10, 200)
		console.AddComponent(components.NewBox("console/box", &engosdl.Rect{}, sdl.Color{}, false))
		consoleText := components.NewText("console/text", "fonts/fira2.ttf", 12, sdl.Color{}, " ")
		console.AddComponent(consoleText)
		message := ""

		lookButton := engosdl.NewEntity("look")
		moveButton := engosdl.NewEntity("move")

		lookButton.GetTransform().SetPositionXY(10, 50)
		lookButton.AddComponent(components.NewBox("look/box", &engosdl.Rect{}, sdl.Color{R: 255}, false))
		lookButton.GetComponent(&components.Box{}).SetCustomOnUpdate(func(c engosdl.IComponent) {
			entity := c.GetEntity()
			x, y, _ := sdl.GetMouseState()
			box := entity.GetComponent(&components.Box{}).(*components.Box)
			if entity.IsInside(engosdl.NewVector(float64(x), float64(y))) {
				box.Color = sdl.Color{G: 255}
				cursor := sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_HAND)
				sdl.SetCursor(cursor)
			} else {
				box.Color = sdl.Color{R: 255}
				cursor := sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_ARROW)
				sdl.SetCursor(cursor)
			}

		})
		lookText := components.NewText("player/look-text", "fonts/fira.ttf", 32, sdl.Color{B: 255}, "LOOK")
		lookText.AddDelegateToRegister(nil, player, &components.Mouse{}, func(c engosdl.IComponent) func(params ...interface{}) bool {
			// enabled := false
			return func(params ...interface{}) bool {
				mousePos := engosdl.NewVector(float64(params[0].(int32)), float64(params[1].(int32)))
				if c.GetEntity().IsInside(mousePos) {
					if c.GetEnabled() {
						// message += "You clicked inside LOOK button\n"
						// consoleText.SetMessage(message)
						// fmt.Println("you clicked inside LOOK button")
						// 	moveText := moveButton.GetComponent(&components.Text{}).(*components.Text)
						// 	moveText.SetEnabled(enabled)
						// 	if enabled {
						// 		moveText.SetColor(sdl.Color{R: 255, A: 255})
						// 	} else {
						// 		moveText.SetColor(sdl.Color{R: 255, A: 100})
						// 	}
						// 	enabled = !enabled
						if output, err := board.GetComponent(&Board{}).(*Board).ExecuteAtPlayerPos("look"); err == nil {
							message += output + "\n"
							consoleText.SetMessage(message)
						}
					}
					// } else {
					// 	fmt.Println("you clicked outside LOOK button")
				}
				return true
			}
		}(lookText))
		lookButton.AddComponent(lookText)
		player.AddChild(lookButton)

		moveButton.GetTransform().SetPositionXY(100, 50)
		// moveButton.AddComponent(components.NewBox("move/box", &engosdl.Rect{}, sdl.Color{B: 255}, true))
		moveText := components.NewText("player/move-text", "fonts/fira.ttf", 32, sdl.Color{R: 255}, "MOVE")
		moveText.AddDelegateToRegister(nil, player, &components.Mouse{}, func(c engosdl.IComponent) func(params ...interface{}) bool {
			return func(params ...interface{}) bool {
				mousePos := engosdl.NewVector(float64(params[0].(int32)), float64(params[1].(int32)))
				if c.GetEntity().IsInside(mousePos) {
					if c.GetEnabled() {
						// message += "You clicked inside MOVE button\n"
						// consoleText.SetMessage(message)
						// fmt.Println("you clicked inside MOVE button")
						if output, err := board.GetComponent(&Board{}).(*Board).ExecuteAtPlayerPos("move"); err == nil {
							message += output + "\n"
							consoleText.SetMessage(message)
						}
						// } else {
						// 	message += "Component MOVE is not enabled\n"
						// 	consoleText.SetMessage(message)
					}
					// } else {
					// 	fmt.Println("you clicked outside MOVE button")
				}
				return true
			}
		}(moveText))
		moveButton.AddComponent(moveText)
		player.AddChild(moveButton)

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
