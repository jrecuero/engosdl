package engosdl

// IObject represents the interface for any object in the game.
type IObject interface {
	GetID() string
	GetLoaded() bool
	GetName() string
	GetStarted() bool
	SetLoaded(bool)
	SetName(string) IObject
	SetStarted(bool)
}

// Object is the default implementation for IObject interface.
type Object struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	loaded  bool
	started bool
}

var _ IObject = (*Object)(nil)

// NewObject returns a new object instance.
func NewObject(name string) *Object {
	return &Object{
		ID:      nextIder(),
		Name:    name,
		loaded:  false,
		started: false,
	}
}

// GetID returns the object id.
func (obj *Object) GetID() string {
	return obj.ID
}

// GetLoaded returns if object has been loaded.
func (obj *Object) GetLoaded() bool {
	return obj.loaded
}

// GetName returns the object name
func (obj *Object) GetName() string {
	return obj.Name
}

// GetStarted returns if object has been started.
func (obj *Object) GetStarted() bool {
	return obj.started
}

// SetLoaded sets if object has been loaded.
func (obj *Object) SetLoaded(loaded bool) {
	obj.loaded = loaded
}

// SetName sets the object name
func (obj *Object) SetName(name string) IObject {
	obj.Name = name
	return obj
}

// SetStarted sets if object has been started.
func (obj *Object) SetStarted(started bool) {
	obj.started = started
}
