package engosdl

import "github.com/veandco/go-sdl2/sdl"

// Vector represents any 2-Dimensional variable (position, rotation, ...)
type Vector struct {
	X float64
	Y float64
}

// NewVector creates a new vector instance
func NewVector(x, y float64) *Vector {
	return &Vector{
		X: x,
		Y: y,
	}
}

// Get returns vector values as a pair.
func (v *Vector) Get() (float64, float64) {
	return v.X, v.Y
}

// InRect returns if the vector is inside the given rectangle.
func (v *Vector) InRect(rect *Rect) bool {
	sdlPoint := &sdl.FPoint{X: float32(v.X), Y: float32(v.Y)}
	sdlRect := &sdl.FRect{
		X: float32(rect.X),
		Y: float32(rect.Y),
		W: float32(rect.W),
		H: float32(rect.H),
	}
	return sdlPoint.InRect(sdlRect)
}
