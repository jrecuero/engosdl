package engosdl

// IComponent represents the interface for any component to be added to any
// GameObject
type IComponent interface {
	IObject
	OnAwake()
	OnStart()
	OnUpdate()
	OnEnable()
	OnDraw()
}

// Component represents the default IComponent implementation.
type Component struct {
	id   string
	name string
}

var _ IComponent = (*Component)(nil)

// GetID returns the component id.
func (c *Component) GetID() string {
	return c.id
}

// GetName returns the component name
func (c *Component) GetName() string {
	return c.name
}

// SetName sets the component name
func (c *Component) SetName(name string) IObject {
	c.name = name
	return c
}

// OnAwake is called first time the component is created.
func (c *Component) OnAwake() {
}

// OnStart is called first time the component is enabled.
func (c *Component) OnStart() {
}

// OnEnable is called every time the component is enabled.
func (c *Component) OnEnable() {
}

// OnUpdate is called for every updata tick.
func (c *Component) OnUpdate() {
}

// OnDraw is called for every draw tick.
func (c *Component) OnDraw() {
}

// NewComponent creates a new component instance.
func NewComponent(name string) *Component {
	logger.Trace().Str("component", name).Msg("new component")
	return &Component{
		id:   "",
		name: name,
	}
}
