package engosdl

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

// IResource represents any graphical resource to be handled by the resource
// handler.
type IResource interface {
	IObject
	Clear()
	Delete() int
	GetFilename() string
	GetFormat() int
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
	format   int
}

var _ IResource = (*Resource)(nil)

// NewResource creates a new resource instance.
func NewResource(name string, filename string, format int) *Resource {
	var err error
	Logger.Trace().Str("resource", name).Str("filename", filename).Msg("new resource")
	result := &Resource{
		Object:   NewObject(name),
		filename: filename,
		counter:  0,
		format:   format,
	}
	switch format {
	case FormatBMP:
		result.surface, err = sdl.LoadBMP(filename)
		break
	case FormatPNG:
		result.surface, err = img.Load(filename)
		break
	case FormatJPG:
		result.surface, err = img.Load(filename)
		break
	default:
		err := fmt.Errorf("unknown format %d", format)
		Logger.Error().Err(err)
		panic(err)
	}
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

// GetFilename returns resource filename.
func (r *Resource) GetFilename() string {
	return r.filename
}

// GetFormat returns resource format.
func (r *Resource) GetFormat() int {
	return r.format
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

// IResourceManager represents the handler that is in charge of all graphical
// resources.
type IResourceManager interface {
	IObject
	Clear()
	CreateResource(string, string, int) IResource
	DeleteResource(IResource) bool
	DoInit()
	GetResource(string) IResource
	GetResourceByFilename(string) IResource
	GetResourceByName(string) IResource
	GetResources() []IResource
	OnStart()
}

// ResourceManager is the default implementation for the resource handler.
type ResourceManager struct {
	*Object
	resources []IResource
}

var _ IResourceManager = (*ResourceManager)(nil)

// NewResourceManager creates a new resource handler instance.
func NewResourceManager(name string) *ResourceManager {
	Logger.Trace().Str("resource-manager", name).Msg("new resource-manager")
	return &ResourceManager{
		Object:    NewObject(name),
		resources: []IResource{},
	}

}

// Clear removes all resources from the resource handler.
func (h *ResourceManager) Clear() {
	Logger.Trace().Str("resource-manager", h.GetName()).Msg("Clear")
	for _, r := range h.resources {
		r.Clear()
	}
	h.resources = []IResource{}
}

// CreateResource creates a new resource. If the same resource has already
// been created with the same filename, existing resource is returned.
func (h *ResourceManager) CreateResource(name string, filename string, format int) IResource {
	Logger.Trace().Str("resource-manager", h.GetName()).Str("name", name).Str("filename", filename).Msg("CreateResource")
	for _, resource := range h.resources {
		if resource.GetFilename() == filename {
			resource.New()
			return resource
		}
	}
	resource := NewResource(name, filename, format)
	h.resources = append(h.resources, resource)
	return resource
}

// DeleteResource deletes resource from the handler. Memory resources are
// released from the given resource.
func (h *ResourceManager) DeleteResource(resource IResource) bool {
	Logger.Trace().Str("resource-manager", h.GetName()).Str("name", resource.GetName()).Str("filename", resource.GetFilename()).Msg("DeleteResource")
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

// DoInit initializes all resource manager resources.
func (h *ResourceManager) DoInit() {
	Logger.Trace().Str("resource-manager", h.GetName()).Msg("DoInit")
}

// GetResource returns a resource with the given resource ID.
func (h *ResourceManager) GetResource(id string) IResource {
	for _, resource := range h.resources {
		if resource.GetID() == id {
			return resource
		}
	}
	return nil
}

// GetResourceByFilename returns the resource with the given filename.
func (h *ResourceManager) GetResourceByFilename(filename string) IResource {
	for _, resource := range h.resources {
		if resource.GetFilename() == filename {
			return resource
		}
	}
	return nil
}

// GetResourceByName returns the resource with the given name.
func (h *ResourceManager) GetResourceByName(name string) IResource {
	for _, resource := range h.resources {
		if resource.GetName() == name {
			return resource
		}
	}
	return nil
}

// GetResources returns all resources.
func (h *ResourceManager) GetResources() []IResource {
	return h.resources
}

// OnStart initializes all resource handler structure.
func (h *ResourceManager) OnStart() {
	Logger.Trace().Str("resource-manager", h.GetName()).Msg("OnStart")
}
