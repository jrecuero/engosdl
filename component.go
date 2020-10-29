package engosdl

// IComponent represents the interface for any component to be added to any
// GameObject
type IComponent interface {
	IObject
	GetActive() bool
	GetDelegates() []IDelegate
	GetGameObject() *GameObject
	Load()
	SetActive(bool)
	OnAwake()
	OnCycleEnd()
	OnCycleStart()
	OnDraw()
	OnEnable()
	OnStart()
	OnUpdate()
	Unload()
}

// Component represents the default IComponent implementation.
type Component struct {
	*Object
	gameObject *GameObject
	active     bool
	loaded     bool
}

var _ IComponent = (*Component)(nil)

// NewComponent creates a new component instance.
func NewComponent(name string, gobj *GameObject) *Component {
	Logger.Trace().Str("component", name).Msg("new component")
	return &Component{
		Object:     NewObject(name),
		gameObject: gobj,
		active:     true,
		loaded:     false,
	}
}

// GetActive returns if component is active or not
func (c *Component) GetActive() bool {
	return c.active
}

// GetDelegates returns all delegates registered to the component.
func (c *Component) GetDelegates() []IDelegate {
	return nil
}

// GetGameObject return the component game object parent.
func (c *Component) GetGameObject() *GameObject {
	return c.gameObject
}

// Load is called when component is loaded by the game object.
func (c *Component) Load() {
	c.loaded = true
	c.OnStart()
}

// OnAwake is called first time the component is created.
func (c *Component) OnAwake() {
	Logger.Trace().Str("component", c.name).Msg("OnAwake")
}

// OnCycleEnd calls all methods to run at the end of a tick cycle.
func (c *Component) OnCycleEnd() {
}

// OnCycleStart calls all methods to run at the start of a tick cycle.
func (c *Component) OnCycleStart() {
}

// OnDraw is called for every draw tick.
func (c *Component) OnDraw() {
	// Logger.Trace().Str("component", c.name).Msg("OnDraw")
}

// OnEnable is called every time the component is enabled.
func (c *Component) OnEnable() {
	Logger.Trace().Str("component", c.name).Msg("OnEnable")
}

// OnStart is called first time the component is enabled.
func (c *Component) OnStart() {
	Logger.Trace().Str("component", c.name).Msg("OnStart")
}

// OnUpdate is called for every update tick.
func (c *Component) OnUpdate() {
	// Logger.Trace().Str("component", c.name).Msg("OnUpdate")
}

// SetActive sets component active attribute
func (c *Component) SetActive(active bool) {
	c.active = active
}

// Unload is called when component is unloaded by the game object.
func (c *Component) Unload() {
	c.loaded = false
}
