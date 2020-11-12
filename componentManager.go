package engosdl

import (
	"fmt"
	"reflect"
)

// ComponentManager is in charge of storing all components registered.
type ComponentManager struct {
	*Object
	components []IComponent
}

// RegisterComponent registers a new component to the manager.
func (cm *ComponentManager) RegisterComponent(component IComponent) bool {
	Logger.Trace().Msg(fmt.Sprintf("register component %s", reflect.TypeOf(component)))
	cm.components = append(cm.components, component)
	return true
}
