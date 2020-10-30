package engosdl

// ITransform represents the interface for any entity transformation.
type ITransform interface {
	GetDim() *Vector
	GetPosition() *Vector
	GetRotation() float64
	GetScale() *Vector
	SetDim(*Vector) ITransform
	SetPosition(*Vector) ITransform
	SetRotation(float64) ITransform
	SetScale(*Vector) ITransform
}

// Transform is the default implementation for ITransform interface.
type Transform struct {
	position *Vector
	rotation float64
	scale    *Vector
	dim      *Vector
}

//GetDim returns the transform original dimensions.
func (t *Transform) GetDim() *Vector {
	return t.dim
}

// GetPosition returns the transform position.
func (t *Transform) GetPosition() *Vector {
	return t.position
}

// GetRotation returns the transform rotation.
func (t *Transform) GetRotation() float64 {
	return t.rotation
}

// GetScale returns the transform scale.
func (t *Transform) GetScale() *Vector {
	return t.scale
}

// SetDim sets the transform original dimensions.
func (t *Transform) SetDim(v *Vector) ITransform {
	t.dim = v
	return t
}

// SetPosition sets the transform position.
func (t *Transform) SetPosition(v *Vector) ITransform {
	t.position = v
	return t
}

// SetRotation sets the transform rotation.
func (t *Transform) SetRotation(r float64) ITransform {
	t.rotation = r
	return t
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
		dim:      NewVector(0, 0),
	}
}
