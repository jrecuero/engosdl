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
	*Object
}

var _ IComponent = (*Component)(nil)

// OnAwake is called first time the component is created.
func (c *Component) OnAwake() {
	Logger.Trace().Str("component", c.name).Msg("OnAwake")
}

// OnStart is called first time the component is enabled.
func (c *Component) OnStart() {
	Logger.Trace().Str("component", c.name).Msg("OnStart")
}

// OnEnable is called every time the component is enabled.
func (c *Component) OnEnable() {
	Logger.Trace().Str("component", c.name).Msg("OnEnable")
}

// OnUpdate is called for every updata tick.
func (c *Component) OnUpdate() {
	Logger.Trace().Str("component", c.name).Msg("OnUpdate")
}

// OnDraw is called for every draw tick.
func (c *Component) OnDraw() {
	Logger.Trace().Str("component", c.name).Msg("OnDraw")
}

// NewComponent creates a new component instance.
func NewComponent(name string) *Component {
	Logger.Trace().Str("component", name).Msg("new component")
	return &Component{
		Object: NewObject(name),
	}
}
