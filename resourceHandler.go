package engosdl

// IResourceHandler represents the handler that is in charge of all graphical
// resources.
type IResourceHandler interface {
	IObject
	OnStart()
}

// ResourceHandler is the default implementation for the resource handler.
type ResourceHandler struct {
	*Object
}

// NewResourceHandler creates a new resource handler instance.
func NewResourceHandler(name string) *ResourceHandler {
	Logger.Trace().Str("resource-handler", name).Msg("new resource-handler")
	return &ResourceHandler{
		Object: NewObject(name),
	}

}

// OnStart initializes all resource handler structure.
func (h *ResourceHandler) OnStart() {
	Logger.Trace().Str("resource-handler", h.GetName()).Msg("OnStart")
}
