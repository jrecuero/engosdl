package engosdl

import "github.com/veandco/go-sdl2/sdl"

// ITransform represents the interface for any entity transformation.
type ITransform interface {
	GetDim() *Vector
	GetPosition() *Vector
	GetRect() *sdl.Rect
	GetRectExt() (float64, float64, float64, float64)
	GetRotation() float64
	GetScale() *Vector
	SetDim(*Vector) ITransform
	SetPosition(*Vector) ITransform
	SetRotation(float64) ITransform
	SetScale(*Vector) ITransform
}

// Transform is the default implementation for ITransform interface.
type Transform struct {
	Position *Vector `json:"position"`
	Rotation float64 `json:"rotation"`
	Scale    *Vector `json:"scale"`
	Dim      *Vector `json:"dimension"`
}

// NewTransform creates a new transform instance.
func NewTransform() *Transform {
	return &Transform{
		Position: NewVector(0, 0),
		Rotation: 0,
		Scale:    NewVector(1, 1),
		Dim:      NewVector(0, 0),
	}
}

//GetDim returns the transform original dimensions.
func (t *Transform) GetDim() *Vector {
	return t.Dim
}

// GetPosition returns the transform position.
func (t *Transform) GetPosition() *Vector {
	return t.Position
}

// GetRect returns a rectangle with real position and dimensions.
// Real dimensions are affected by the scale value.
func (t *Transform) GetRect() *sdl.Rect {
	return &sdl.Rect{
		X: int32(t.GetPosition().X),
		Y: int32(t.GetPosition().Y),
		W: int32(t.GetDim().X * t.GetScale().X),
		H: int32(t.GetDim().Y * t.GetScale().Y),
	}
}

// GetRectExt returns rectangle coordinates as x, y, w, and h.
func (t *Transform) GetRectExt() (float64, float64, float64, float64) {
	return t.GetPosition().X, t.GetPosition().Y, t.GetDim().X * t.GetScale().X, t.GetDim().Y * t.GetScale().Y
}

// GetRotation returns the transform rotation.
func (t *Transform) GetRotation() float64 {
	return t.Rotation
}

// GetScale returns the transform scale.
func (t *Transform) GetScale() *Vector {
	return t.Scale
}

// SetDim sets the transform original dimensions.
func (t *Transform) SetDim(v *Vector) ITransform {
	t.Dim = v
	return t
}

// SetPosition sets the transform position.
func (t *Transform) SetPosition(v *Vector) ITransform {
	t.Position = v
	return t
}

// SetRotation sets the transform rotation.
func (t *Transform) SetRotation(r float64) ITransform {
	t.Rotation = r
	return t
}

// SetScale sets the transform scale.
func (t *Transform) SetScale(v *Vector) ITransform {
	t.Scale = v
	return t
}
