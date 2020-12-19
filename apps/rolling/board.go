package main

import (
	"fmt"
	"reflect"

	"github.com/jrecuero/engosdl"
)

// ComponentNameBoard is the name to refer board component.
var ComponentNameBoard string = reflect.TypeOf(&Board{}).String()

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor(ComponentNameBoard, CreateBoard)
	}
}

// Position represents a row and column in a board.
type Position struct {
	Row int
	Col int
}

// None sets the position to an invalid position in the board.
func (p *Position) None() {
	p.Row = -1
	p.Col = -1
}

// EqualTo checks if the given row and columns match the position.
func (p *Position) EqualTo(row int, col int) bool {
	return p.Row == row && p.Col == col
}

// TCellResultSignature is the callback to execute a cell action.
type TCellResultSignature func(*Board, engosdl.IEntity, *Position, ...interface{}) (string, error)

// Cell represents any cell in a board component.
type Cell struct {
	EnterDialog     string
	ExitDialog      string
	ActionAndResult map[string]TCellResultSignature
	Position        *Position
}

// NewCell creates a new cell instance.
func NewCell(row int, column int) *Cell {
	return &Cell{
		EnterDialog:     "",
		ExitDialog:      "",
		ActionAndResult: make(map[string]TCellResultSignature),
		Position:        &Position{Row: row, Col: column},
	}
}

// GetActions returns all actions to be performed in a given cell.
func (c *Cell) GetActions() []string {
	actions := make([]string, len(c.ActionAndResult))
	i := 0
	for k := range c.ActionAndResult {
		actions[i] = k
		i++
	}
	return actions
}

// Board represents a component.
type Board struct {
	*engosdl.Component
	Entities       [][]engosdl.IEntity
	Columns        int
	Rows           int
	Origin         *engosdl.Vector
	CellSize       int
	Cells          [][]*Cell
	dirty          bool
	playerPosition *Position
}

// var _ engosdl.IBoard = (*Board)(nil)

// NewBoard create a new Board instance.
func NewBoard(name string, columns int, rows int, origin *engosdl.Vector, cellSize int) *Board {
	engosdl.Logger.Trace().Str("component", "Board").Str("Board", name).Msg("new Board")
	result := &Board{
		Component:      engosdl.NewComponent(name),
		Columns:        columns,
		Rows:           rows,
		Origin:         origin,
		CellSize:       cellSize,
		Entities:       [][]engosdl.IEntity{},
		Cells:          [][]*Cell{},
		dirty:          true,
		playerPosition: &Position{},
	}
	result.playerPosition.None()
	result.Entities = make([][]engosdl.IEntity, rows)
	result.Cells = make([][]*Cell, rows)
	for i := 0; i < rows; i++ {
		result.Entities[i] = make([]engosdl.IEntity, columns)
		result.Cells[i] = make([]*Cell, columns)
	}
	return result
}

// CreateBoard implements Board constructor used by component
// manager.
func CreateBoard(params ...interface{}) engosdl.IComponent {
	if len(params) == 5 {
		return NewBoard(params[0].(string), params[1].(int), params[2].(int), params[3].(*engosdl.Vector), params[4].(int))
	}
	return NewBoard("", 0, 0, nil, 0)
}

// AddEntityAt adds an entity to the given position.
func (c *Board) AddEntityAt(entity engosdl.IEntity, row int, col int, player bool) error {
	if c.Entities[row][col] == nil {
		c.Entities[row][col] = entity
		pos, _ := c.GetPositionFromCell(row, col)
		entity.GetTransform().SetPosition(pos)
		if player {
			c.dirty = true
			c.playerPosition = &Position{Row: row, Col: col}
		}
		return nil
	}
	return fmt.Errorf("cell at row %d col %d is not free", row, col)
}

// GetPositionFromCell returns display position for the given cell.
func (c *Board) GetPositionFromCell(row int, col int) (*engosdl.Vector, error) {
	x := c.Origin.X + float64(col*c.CellSize)
	y := c.Origin.Y + float64(row*c.CellSize)
	return engosdl.NewVector(x, y), nil
}

// DefaultAddDelegateToRegister will proceed to add default delegates to
// register to the component.
func (c *Board) DefaultAddDelegateToRegister() {
	// c.AddDelegateToRegister(<DELEGATE>, nil, <OTHER-COMPONENT>, <SIGNATURE>)
}

// DeleteEntityAt deletes entity at given position.
func (c *Board) DeleteEntityAt(row int, col int) error {
	c.Entities[row][col] = nil
	if c.playerPosition.EqualTo(row, col) {
		c.dirty = true
		c.playerPosition.None()
	}
	return nil
}

// DoDestroy should destroy all component resources. This is called when
// component is removed from the scene and resources are not anymore
// required.
func (c *Board) DoDestroy() {
	engosdl.Logger.Trace().Str("component", "Board").Str("Board", c.GetName()).Msg("DoDestroy")
	c.Component.DoDestroy()
}

// DoUnLoad is called when component is unloaded from scene.
func (c *Board) DoUnLoad() {
	engosdl.Logger.Trace().Str("component", "Board").Str("Board", c.GetName()).Msg("DoUnLoad")
	c.Component.DoUnLoad()
}

// ExecuteAtPlayerPos executes action at player position.
func (c *Board) ExecuteAtPlayerPos(action string) (string, error) {
	cell := c.Cells[c.playerPosition.Row][c.playerPosition.Col]
	if actionSignature := cell.ActionAndResult[action]; actionSignature != nil {
		player := c.Entities[c.playerPosition.Row][c.playerPosition.Col]
		c.dirty = true
		return actionSignature(c, player, c.playerPosition)

	}
	return "", nil
}

// GetEntityAt returns the entity at the given position.
func (c *Board) GetEntityAt(row int, col int) engosdl.IEntity {
	return c.Entities[row][col]
}

// IsCellFree checks if a cell in the board is free or not.
func (c *Board) IsCellFree(row int, col int) bool {
	if row >= 0 && row < c.Rows && col >= 0 && col < c.Columns {
		return c.Entities[row][col] == nil
	}
	return false
}

// OnAwake is called when component is first loaded into the scene and all
// component resources have to be created. No resources dependent with other
// components or entities can be created at this point.
func (c *Board) OnAwake() {
	engosdl.Logger.Trace().Str("component", "Board").Str("Board", c.GetName()).Msg("OnAwake")
	// Create new delegate "on-board"
	name := fmt.Sprintf("on-board/%s", c.GetName())
	c.SetDelegate(engosdl.GetDelegateManager().CreateDelegate(c, name))
	c.Component.OnAwake()
}

// OnRender is called every engine frame when component has to be rendered.
func (c *Board) OnRender() {
}

// OnStart is called at the end of the component being loaded by the scene.
// Any component resource dependent from other entities or components has
// to be created at this point.
func (c *Board) OnStart() {
	engosdl.Logger.Trace().Str("component", "Board").Str("Board", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}

// OnUpdate is called every engine frame in order to update any component
// resource.
func (c *Board) OnUpdate() {
	if c.dirty && c.playerPosition.Row != -1 && c.playerPosition.Col != -1 {
		cell := c.Cells[c.playerPosition.Row][c.playerPosition.Col]
		engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), false, cell.GetActions())
		// for _, action := range cell.GetActions() {
		// 	player := c.Entities[c.playerPosition.Row][c.playerPosition.Col]
		// 	fmt.Println(player.GetName(), action)
		// }
		c.dirty = false
	}
	c.Component.OnUpdate()
}

// Unmarshal takes information from a ComponentToUnmarshal instance and
// creates a new component instance.
func (c *Board) Unmarshal(data map[string]interface{}) {
	c.Component.Unmarshal(data)
}
