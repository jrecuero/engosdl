package engosdl

// Vector represents any 2-Dimentional variable (position, rotation, ...)
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
