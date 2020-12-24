package main

import (
	"fmt"
	"reflect"

	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
)

// ComponentNameSceneController is the name to refer SceneController component.
var ComponentNameSceneController string = reflect.TypeOf(&SceneController{}).String()

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor(ComponentNameSceneController, CreateSceneController)
	}
}

// SceneController represents a component.
type SceneController struct {
	*engosdl.Component
	Player  *Player
	Board   engosdl.IEntity
	Console engosdl.IEntity
	Enemies []*Player
}

// var _ engosdl.ISceneController = (*SceneController)(nil)

// NewSceneController create a new SceneController instance.
func NewSceneController(name string, player *Player) *SceneController {
	engosdl.Logger.Trace().Str("component", "SceneController").Str("SceneController", name).Msg("new SceneController")
	result := &SceneController{
		Component: engosdl.NewComponent(name),
		Player:    player,
		Board:     engosdl.NewEntity("board"),
		Console:   engosdl.NewEntity("console"),
		Enemies:   []*Player{},
	}
	result.Player.AddChild(engosdl.NewEntity("look"))
	result.Player.AddChild(engosdl.NewEntity("move"))
	result.Player.AddChild(engosdl.NewEntity("attack"))
	result.Player.SetCache("sheet", NewCharacterSheet(NewAbility(18, 16, 12, 10, 8, 10)))
	enemy := NewPlayer("goblin")
	enemy.SetCache("sheet", NewCharacterSheet(NewAbility(10, 8, 8, 6, 6, 6)))
	result.Enemies = append(result.Enemies, enemy)
	result.Board.AddComponent(result.createBoard())
	result.Console.SetCache("message", "")
	return result
}

// CreateSceneController implements SceneController constructor used by component
// manager.
func CreateSceneController(params ...interface{}) engosdl.IComponent {
	if len(params) == 1 {
		return NewSceneController(params[0].(string), params[1].(*Player))
	}
	return NewSceneController("", nil)
}

func (c *SceneController) addDelegateToRegisterToButton(name string) {
	component := c.Player.GetChildByName(name).GetComponent(&components.Button{})
	component.AddDelegateToRegister(nil, c.Player, &components.Mouse{}, func(params ...interface{}) bool {
		mousePos := engosdl.NewVector(float64(params[0].(int32)), float64(params[1].(int32)))
		if component.GetEntity().IsInside(mousePos) {
			if component.GetEnabled() {
				if output, err := c.Board.GetComponent(&Board{}).(*Board).ExecuteAtPlayerPos(name); err == nil {
					if obj, error := c.Console.GetCache("message"); error == nil {
						message := obj.(string) + output + "\n"
						c.Console.SetCache("message", message)
						c.Console.GetComponent(&components.Text{}).(*components.Text).SetMessage(message)
					}
				}
			}
		}
		return true
	})
}

func (c *SceneController) createBoard() *Board {
	var cell *Cell
	board := NewBoard("dungeon", 10, 1, engosdl.NewVector(5, 5), 32)

	cell = NewCell(0, 0)
	cell.EnterDialog = "You enter in the dungeon"
	cell.ExitDialog = "Perils follow from this point"
	cell.ActionAndResult["move"] = func(board *Board, player engosdl.IEntity, entities []engosdl.IEntity, position *Position, params ...interface{}) (string, error) {
		newPos := &Position{Row: position.Row, Col: position.Col + 1}
		board.DeleteEntityAt(player, position.Row, position.Col)
		board.AddEntityAt(player, newPos.Row, newPos.Col, true)
		return "Move to the next cell", nil
	}
	board.Cells[0][0] = cell

	cell = NewCell(0, 1)
	cell.EnterDialog = "There is an enemy here"
	cell.ExitDialog = "You beat the enemy. You can continue"
	cell.ActionAndResult["look"] = func(board *Board, player engosdl.IEntity, entities []engosdl.IEntity, position *Position, params ...interface{}) (string, error) {
		text := "room is empty"
		for _, entity := range entities {
			if entity.GetID() != player.GetID() {
				text = fmt.Sprintf("you can see in the room: %s", entity.GetName())
			}
		}
		return text, nil
	}
	cell.ActionAndResult["move"] = func(board *Board, player engosdl.IEntity, entities []engosdl.IEntity, position *Position, params ...interface{}) (string, error) {
		newPos := &Position{Row: position.Row, Col: position.Col + 1}
		board.DeleteEntityAt(player, position.Row, position.Col)
		board.AddEntityAt(player, newPos.Row, newPos.Col, true)
		return "Move to the next cell", nil
	}
	cell.ActionAndResult["attack"] = func(board *Board, player engosdl.IEntity, entities []engosdl.IEntity, position *Position, params ...interface{}) (string, error) {
		result := "nothing happened"
		for _, entity := range entities {
			if entity.GetID() != player.GetID() {
				attack := RollDice(20)
				obj, _ := player.GetCache("sheet")
				playerStr := obj.(*CharacterSheet).Modifier.Strength
				obj, _ = entity.GetCache("sheet")
				enemyAC := obj.(*CharacterSheet).ArmorClass
				// if (attack + playerStr) >= enemyAC {
				result = fmt.Sprintf("%s Dice roll %d+%d against %s AC %d", player.GetName(), attack, playerStr, entity.GetName(), enemyAC)
				// }
			}
		}
		return result, nil
	}
	board.Cells[0][1] = cell

	cell = NewCell(0, 2)
	cell.EnterDialog = "You reach your destination"
	cell.ExitDialog = "This is the end"
	cell.ActionAndResult["look"] = func(board *Board, player engosdl.IEntity, entities []engosdl.IEntity, position *Position, params ...interface{}) (string, error) {
		return "You made it!", nil
	}
	board.Cells[0][2] = cell
	return board
}

// DefaultAddDelegateToRegister will proceed to add default delegates to
// register to the component.
func (c *SceneController) DefaultAddDelegateToRegister() {
	// c.AddDelegateToRegister(<DELEGATE>, nil, <OTHER-COMPONENT>, <SIGNATURE>)
}

// DoDestroy should destroy all component resources. This is called when
// component is removed from the scene and resources are not anymore
// required.
func (c *SceneController) DoDestroy() {
	engosdl.Logger.Trace().Str("component", "SceneController").Str("SceneController", c.GetName()).Msg("DoDestroy")
	c.Component.DoDestroy()
}

// DoUnLoad is called when component is unloaded from scene.
func (c *SceneController) DoUnLoad() {
	engosdl.Logger.Trace().Str("component", "SceneController").Str("SceneController", c.GetName()).Msg("DoUnLoad")
	c.Component.DoUnLoad()
}

// OnAwake is called when component is first loaded into the scene and all
// component resources have to be created. No resources dependent with other
// components or entities can be created at this point.
func (c *SceneController) OnAwake() {
	engosdl.Logger.Trace().Str("component", "SceneController").Str("SceneController", c.GetName()).Msg("OnAwake")
	c.Component.OnAwake()
}

// OnRender is called every engine frame when component has to be rendered.
func (c *SceneController) OnRender() {
}

// OnStart is called at the end of the component being loaded by the scene.
// Any component resource dependent from other entities or components has
// to be created at this point.
func (c *SceneController) OnStart() {
	engosdl.Logger.Trace().Str("component", "SceneController").Str("SceneController", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}

// OnUpdate is called every engine frame in order to update any component
// resource.
func (c *SceneController) OnUpdate() {
	c.Component.OnUpdate()
}

// SetupResources configures all scene controller resources.
func (c *SceneController) SetupResources() {
	c.AddDelegateToRegister(nil, c.Board, &Board{}, func(params ...interface{}) bool {
		actions := params[0].([]string)
		c.Player.UpdateActions(actions)
		return true
	})

	c.addDelegateToRegisterToButton("look")
	c.addDelegateToRegisterToButton("move")
	c.addDelegateToRegisterToButton("attack")
}

// Unmarshal takes information from a ComponentToUnmarshal instance and
// creates a new component instance.
func (c *SceneController) Unmarshal(data map[string]interface{}) {
	c.Component.Unmarshal(data)
}
