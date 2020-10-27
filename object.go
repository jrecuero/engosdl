package engosdl

// IObject represents the interface for any object in the game.
type IObject interface {
	GetID() string
	GetName() string
	SetName(string) IObject
}

// Object is the default implementation for IObject interface.
type Object struct {
	id   string
	name string
}

var _ IObject = (*Object)(nil)

// GetID returns the object id.
func (obj *Object) GetID() string {
	return obj.id
}

// GetName returns the object name
func (obj *Object) GetName() string {
	return obj.name
}

// SetName sets the object name
func (obj *Object) SetName(name string) IObject {
	obj.name = name
	return obj
}

// NewObject returns a new object instance.
func NewObject(name string) *Object {
	return &Object{
		id:   nextIder(),
		name: name,
	}
}
