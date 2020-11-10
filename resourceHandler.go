package engosdl

import "github.com/veandco/go-sdl2/sdl"

// IResource represents any graphical resource to be handled by the resource
// handler.
type IResource interface {
	IObject
	Clear()
	Delete() int
	GetFilename() string
	GetSurface() *sdl.Surface
	GetTextureFromSurface() *sdl.Texture
	New()
}

// Resource is the default implementation for the resource interface.
type Resource struct {
	*Object
	filename string
	surface  *sdl.Surface
	counter  int
}

var _ IResource = (*Resource)(nil)

// NewResource creates a new resource instance.
func NewResource(name string, filename string) *Resource {
	var err error
	Logger.Trace().Str("resource", name).Str("filename", filename).Msg("new resource")
	result := &Resource{
		Object:   NewObject(name),
		filename: filename,
		counter:  0,
	}
	result.surface, err = sdl.LoadBMP(filename)
	if err != nil {
		Logger.Error().Err(err).Msg("LoadBMP error")
		panic(err)
	}
	return result
}

// Clear deletes resource even if counter is not zero.
func (r *Resource) Clear() {
	Logger.Trace().Str("resource", r.GetName()).Str("filename", r.GetFilename()).Msg("clear resource")
	r.counter = 1
	r.Delete()
}

// Delete deletes resource and relese all memory.
func (r *Resource) Delete() int {
	Logger.Trace().Str("resource", r.GetName()).Str("filename", r.GetFilename()).Msg("delete resource")
	r.counter--
	if r.counter == 0 {
		r.surface.Free()
	}
	return r.counter
}

// GetFilename returns resource filename
func (r *Resource) GetFilename() string {
	return r.filename
}

// GetSurface returns resource surface.
func (r *Resource) GetSurface() *sdl.Surface {
	return r.surface
}

// GetTextureFromSurface returns a texture from the resource surface.
func (r *Resource) GetTextureFromSurface() *sdl.Texture {
	Logger.Trace().Str("resource", r.GetName()).Str("filename", r.GetFilename()).Msg("get texture from surface")
	texture, err := GetRenderer().CreateTextureFromSurface(r.surface)
	if err != nil {
		Logger.Error().Err(err).Msg("CreateTextureFromSurface error")
		panic(err)
	}
	return texture
}

// New increases the number of times this resource is being used.
func (r *Resource) New() {
	r.counter++
}

// IResourceHandler represents the handler that is in charge of all graphical
// resources.
type IResourceHandler interface {
	IObject
	Clear()
	CreateResource(string, string) IResource
	DeleteResource(IResource) bool
	GetResource(string) IResource
	GetResourceByFilename(string) IResource
	GetResourceByName(string) IResource
	GetResources() []IResource
	OnStart()
}

// ResourceHandler is the default implementation for the resource handler.
type ResourceHandler struct {
	*Object
	resources []IResource
}

var _ IResourceHandler = (*ResourceHandler)(nil)

// NewResourceHandler creates a new resource handler instance.
func NewResourceHandler(name string) *ResourceHandler {
	Logger.Trace().Str("resource-handler", name).Msg("new resource-handler")
	return &ResourceHandler{
		Object:    NewObject(name),
		resources: []IResource{},
	}

}

// Clear removes all resources from the resource handler.
func (h *ResourceHandler) Clear() {
	Logger.Trace().Str("resource-handler", h.GetName()).Msg("Clear")
	for _, r := range h.resources {
		r.Clear()
	}
	h.resources = []IResource{}
}

// CreateResource creates a new resource. If the same resource has already
// been created with the same filename, existing resource is returned.
func (h *ResourceHandler) CreateResource(name string, filename string) IResource {
	Logger.Trace().Str("resource-handler", h.GetName()).Str("name", name).Str("filename", filename).Msg("CreateResource")
	for _, resource := range h.resources {
		if resource.GetFilename() == filename {
			resource.New()
			return resource
		}
	}
	resource := NewResource(name, filename)
	h.resources = append(h.resources, resource)
	return resource
}

// DeleteResource deletes resource from the handler. Memory resources are
// released from the given resource.
func (h *ResourceHandler) DeleteResource(resource IResource) bool {
	Logger.Trace().Str("resource-handler", h.GetName()).Str("name", resource.GetName()).Str("filename", resource.GetFilename()).Msg("DeleteResource")
	for i := len(h.resources) - 1; i >= 0; i-- {
		r := h.resources[i]
		if r.GetID() == resource.GetID() {
			if result := r.Delete(); result == 0 {
				h.resources = append(h.resources[:i], h.resources[i+1:]...)
			}
			return true
		}
	}
	return false
}

// GetResource returns a resource with the given resource ID.
func (h *ResourceHandler) GetResource(id string) IResource {
	for _, resource := range h.resources {
		if resource.GetID() == id {
			return resource
		}
	}
	return nil
}

// GetResourceByFilename returns the resource with the given filename.
func (h *ResourceHandler) GetResourceByFilename(filename string) IResource {
	for _, resource := range h.resources {
		if resource.GetFilename() == filename {
			return resource
		}
	}
	return nil
}

// GetResourceByName returns the resource with the given name.
func (h *ResourceHandler) GetResourceByName(name string) IResource {
	for _, resource := range h.resources {
		if resource.GetName() == name {
			return resource
		}
	}
	return nil
}

// GetResources returns all resources.
func (h *ResourceHandler) GetResources() []IResource {
	return h.resources
}

// OnStart initializes all resource handler structure.
func (h *ResourceHandler) OnStart() {
	Logger.Trace().Str("resource-handler", h.GetName()).Msg("OnStart")
}
