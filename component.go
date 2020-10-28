package engosdl

// IComponent represents the interface for any component to be added to any
// GameObject
type IComponent interface {
	IObject
	GetActive() bool
	GetDelegates() []IDelegate
	GetGameObject() *GameObject
	OnAwake()
	OnDraw()
	OnEnable()
	OnStart()
	OnUpdate()
	SetActive(bool)
}

// Component represents the default IComponent implementation.
type Component struct {
	*Object
	gameObject *GameObject
	active     bool
}

var _ IComponent = (*Component)(nil)

// GetGameObject return the component game object parent.
func (c *Component) GetGameObject() *GameObject {
	return c.gameObject
}

// GetDelegates returns all delegates registered to the component.
func (c *Component) GetDelegates() []IDelegate {
	return nil
}

// SetActive sets component active attribute
func (c *Component) SetActive(active bool) {
	c.active = active
}

// GetActive returns if component is active or not
func (c *Component) GetActive() bool {
	return c.active
}

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

// OnUpdate is called for every update tick.
func (c *Component) OnUpdate() {
	// Logger.Trace().Str("component", c.name).Msg("OnUpdate")
}

// OnDraw is called for every draw tick.
func (c *Component) OnDraw() {
	// Logger.Trace().Str("component", c.name).Msg("OnDraw")
}

// NewComponent creates a new component instance.
func NewComponent(name string, gobj *GameObject) *Component {
	Logger.Trace().Str("component", name).Msg("new component")
	return &Component{
		Object:     NewObject(name),
		gameObject: gobj,
		active:     true,
	}
}
