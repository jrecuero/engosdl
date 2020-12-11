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

// Board represents a component.
type Board struct {
	*engosdl.Component
	Space    [][]engosdl.IEntity
	Width    int
	Height   int
	Origin   *engosdl.Vector
	CellSize int
}

// var _ engosdl.IBoard = (*Board)(nil)

// NewBoard create a new Board instance.
func NewBoard(name string, w int, h int, origin *engosdl.Vector, cellSize int) *Board {
	engosdl.Logger.Trace().Str("component", "Board").Str("Board", name).Msg("new Board")
	result := &Board{
		Component: engosdl.NewComponent(name),
		Width:     w,
		Height:    h,
		Origin:    origin,
		CellSize:  cellSize,
		Space:     [][]engosdl.IEntity{},
	}
	result.Space = make([][]engosdl.IEntity, h)
	for i := 0; i < h; i++ {
		result.Space[i] = make([]engosdl.IEntity, w)
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
func (c *Board) AddEntityAt(entity engosdl.IEntity, row int, col int) error {
	if c.Space[row][col] == nil {
		c.Space[row][col] = entity
		return nil
	}
	return fmt.Errorf("cell at row %d col %d is not free", row, col)
}

// GetPositionFromCell returns display position for the given cell.
func (c *Board) GetPositionFromCell(row int, col int) (*engosdl.Vector, error) {
	x := c.Origin.X + float64(row*c.CellSize)
	y := c.Origin.Y + float64(col*c.CellSize)
	return engosdl.NewVector(x, y), nil
}

// DefaultAddDelegateToRegister will proceed to add default delegates to
// register to the component.
func (c *Board) DefaultAddDelegateToRegister() {
	// c.AddDelegateToRegister(<DELEGATE>, nil, <OTHER-COMPONENT>, <SIGNATURE>)
}

// DeleteEntityAt deletes entity at given position.
func (c *Board) DeleteEntityAt(row int, col int) error {
	c.Space[row][col] = nil
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

// GetEntityAt returns the entity at the given position.
func (c *Board) GetEntityAt(row int, col int) engosdl.IEntity {
	return c.Space[row][col]
}

// IsCellFree checks if a cell in the board is free or not.
func (c *Board) IsCellFree(row int, col int) bool {
	if row >= 0 && row < c.Height && col >= 0 && col < c.Width {
		return c.Space[row][col] == nil
	}
	return false
}

// UseCell checks if cell in the board is free to be used, if it is already
// used, if it is already used, remove the entity without assigning the
// new one.
func (c *Board) UseCell(row int, col int) bool {
	if row < 0 || row >= c.Height || col < 0 || col >= c.Width {
		return false
	} else if c.Space[row][col] != nil {
		entity := c.Space[row][col].(*Pixel)
		parent := entity.GetParent()
		parent.DeleteChild(entity.GetID())
		parent.GetScene().DeleteEntity(entity)
		c.Space[row][col] = nil
		fmt.Printf("delete pixel at %d, %d\n", row, col)
		return false
	}
	return true
}

// OnAwake is called when component is first loaded into the scene and all
// component resources have to be created. No resources dependent with other
// components or entities can be created at this point.
func (c *Board) OnAwake() {
	engosdl.Logger.Trace().Str("component", "Board").Str("Board", c.GetName()).Msg("OnAwake")
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
	c.Component.OnUpdate()
}

// Unmarshal takes information from a ComponentToUnmarshal instance and
// creates a new component instance.
func (c *Board) Unmarshal(data map[string]interface{}) {
	c.Component.Unmarshal(data)
}
