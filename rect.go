package engosdl

import "github.com/veandco/go-sdl2/sdl"

// Rect represents any 2-Dimensional rectangle
type Rect struct {
	X float64
	Y float64
	W float64
	H float64
}

// NewRect creates a new rectangle instance.
func NewRect(x, y, w, h float64) *Rect {
	return &Rect{
		X: x,
		Y: y,
		W: w,
		H: h,
	}
}

// HasIntersection returns if two rectangle have intersection.
func (r *Rect) HasIntersection(other *Rect) bool {
	meRect := &sdl.Rect{
		X: int32(r.X),
		Y: int32(r.Y),
		W: int32(r.W),
		H: int32(r.H),
	}
	otherRect := &sdl.Rect{
		X: int32(other.X),
		Y: int32(other.Y),
		W: int32(other.W),
		H: int32(other.H),
	}
	return meRect.HasIntersection(otherRect)
}

// Intersect returns the rectangle intersection result.
func (r *Rect) Intersect(other *Rect) (*Rect, bool) {
	meRect := &sdl.Rect{
		X: int32(r.X),
		Y: int32(r.Y),
		W: int32(r.W),
		H: int32(r.H),
	}
	otherRect := &sdl.Rect{
		X: int32(other.X),
		Y: int32(other.Y),
		W: int32(other.W),
		H: int32(other.H),
	}
	resultRect, ok := meRect.Intersect(otherRect)
	return &Rect{
		X: float64(resultRect.X),
		Y: float64(resultRect.Y),
		W: float64(resultRect.W),
		H: float64(resultRect.H),
	}, ok
}
