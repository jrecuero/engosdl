package engosdl

// ITransform represents the interface for any game object transformation.
type ITransform interface {
	GetPosition() *Vector
	SetPosition(*Vector) ITransform
	GetRotation() float64
	SetRotation(float64) ITransform
	GetScale() *Vector
	SetScale(*Vector) ITransform
}

// Transform is the default implementation for ITransform interface.
type Transform struct {
	position *Vector
	rotation float64
	scale    *Vector
}

// GetPosition returns the transform position.
func (t *Transform) GetPosition() *Vector {
	return t.position
}

// SetPosition sets the transform position.
func (t *Transform) SetPosition(v *Vector) ITransform {
	t.position = v
	return t
}

// GetRotation returns the transform rotation.
func (t *Transform) GetRotation() float64 {
	return t.rotation
}

// SetRotation sets the transform rotation.
func (t *Transform) SetRotation(r float64) ITransform {
	t.rotation = r
	return t
}

// GetScale returns the transform scale.
func (t *Transform) GetScale() *Vector {
	return t.scale
}

// SetScale sets the transform scale.
func (t *Transform) SetScale(v *Vector) ITransform {
	t.scale = v
	return t
}

// NewTransform creates a new transform instance.
func NewTransform() *Transform {
	return &Transform{
		position: NewVector(0, 0),
		rotation: 0,
		scale:    NewVector(1, 1),
	}
}
