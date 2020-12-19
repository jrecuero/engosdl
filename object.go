package engosdl

// IObject represents the interface for any object in the game.
type IObject interface {
	GetID() string
	GetDirty() bool
	GetEnabled() bool
	GetLoaded() bool
	GetName() string
	GetStarted() bool
	SetDirty(bool)
	SetEnabled(bool)
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
	enabled bool
	dirty   bool
}

var _ IObject = (*Object)(nil)

// NewObject returns a new object instance.
func NewObject(name string) *Object {
	return &Object{
		ID:      nextIder(),
		Name:    name,
		loaded:  false,
		started: false,
		enabled: true,
		dirty:   true,
	}
}

// GetID returns the object id.
func (obj *Object) GetID() string {
	return obj.ID
}

// GetDirty returns if the object is dirty or not.
func (obj *Object) GetDirty() bool {
	return obj.dirty
}

// GetEnabled  returns if the object is enabled or not.
func (obj *Object) GetEnabled() bool {
	return obj.enabled
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

// SetDirty sets if the object is dirty or not.
func (obj *Object) SetDirty(dirty bool) {
	obj.dirty = dirty
}

// SetEnabled sets if the object is enabled or not.
func (obj *Object) SetEnabled(enabled bool) {
	obj.enabled = enabled
	obj.SetDirty(true)
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
