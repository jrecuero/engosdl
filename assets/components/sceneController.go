package components

import (
	"reflect"

	"github.com/jrecuero/engosdl"
)

// ComponentNameSceneController is the name to refer scene controller component.
var ComponentNameSceneController string = reflect.TypeOf(&SceneController{}).String()

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor(ComponentNameSceneController, CreateSceneController)
	}
}

// SceneController represents a component, unique per scene that provides a
// way to control the whole scene. This component should be added to other
// unique per scene entity.
type SceneController struct {
	*engosdl.Component
}

var _ engosdl.IComponent = (*SceneController)(nil)

// NewSceneController creates a new scene controller instance.
func NewSceneController(name string) *SceneController {
	engosdl.Logger.Trace().Str("component", "scene-controller").Str("scene-controller", name).Msg("new scene-controller")
	return &SceneController{
		Component: engosdl.NewComponent(name),
	}
}

// CreateSceneController implements scene controller constructor used by component
// manager.
func CreateSceneController(params ...interface{}) engosdl.IComponent {
	if len(params) == 1 {
		return NewSceneController(params[0].(string))
	}
	return NewSceneController("")
}
