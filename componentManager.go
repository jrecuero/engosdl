package engosdl

// ComponentManager is in charge of storing all components registered.
type ComponentManager struct {
	*Object
	// components   []IComponent
	Constructors map[string]func(...interface{}) IComponent
}

// NewComponentManager create a new component manager instance
func NewComponentManager(name string) *ComponentManager {
	Logger.Trace().Str("component-manager", name).Msg("new component-manager")
	return &ComponentManager{
		Object: NewObject(name),
		// components:   []IComponent{},
		Constructors: make(map[string]func(...interface{}) IComponent),
	}
}

// // GetComponents returns all components in the component manager.
// func (cm *ComponentManager) GetComponents() []IComponent {
// 	return cm.components
// }

// // RegisterComponent registers a new component to the manager.
// func (cm *ComponentManager) RegisterComponent(component IComponent) bool {
// 	Logger.Trace().Msg(fmt.Sprintf("register component %s", reflect.TypeOf(component)))
// 	cm.components = append(cm.components, component)
// 	return true
// }

// RegisterConstructor registers a new constructor.
func (cm *ComponentManager) RegisterConstructor(componentName string, constructor func(...interface{}) IComponent) {
	Logger.Trace().Str("component-manager", cm.GetName()).Str("component", componentName).Msg("register constructor")
	cm.Constructors[componentName] = constructor
}
