package engosdl

import (
	"fmt"
	"math"
)

// Scene layer constants.
const (
	LayerTop        int = 3
	LayerMiddle     int = 2
	LayerBottom     int = 1
	LayerBackground int = 0
	maxLayers           = 4
)

// TSceneCodeSignature represents the method to be called in order to create
// a scene.
type TSceneCodeSignature func(*Engine, IScene) bool

// IScene represents the interface for any game scene
type IScene interface {
	IObject
	AddEntity(IEntity) bool
	AuditEntities()
	DeleteEntity(IEntity) bool
	DoDestroy()
	DoFrameEnd()
	DoFrameStart()
	DoLoad()
	DoSwapFrom()
	DoSwapBack()
	DoUnLoad()
	GetEntities() []IEntity
	GetEntitiesByTag(string) []IEntity
	GetEntity(string) IEntity
	GetEntityByName(string) IEntity
	GetSceneCode() TSceneCodeSignature
	GetTag() string
	OnAfterUpdate()
	OnRender()
	OnEnable()
	OnStart()
	OnUpdate()
	SetSceneCode(TSceneCodeSignature)
	SetTag(string)
}

// Scene is the default implementation for IScene interface.
type Scene struct {
	*Object
	entities            []IEntity
	toDeleteEntities    []IEntity
	loadedEntities      []IEntity
	unloadedEntities    []IEntity
	layers              [][]IEntity
	collisionCollection []ICollider
	sceneCode           TSceneCodeSignature
	tag                 string
}

var _ IScene = (*Scene)(nil)

// NewScene creates a new scene instance
func NewScene(name string, tag string) *Scene {
	Logger.Trace().Str("scene", name).Msg("new scene")
	scene := &Scene{
		Object:           NewObject(name),
		entities:         []IEntity{},
		toDeleteEntities: []IEntity{},
		loadedEntities:   []IEntity{},
		unloadedEntities: []IEntity{},
		layers:           make([][]IEntity, maxLayers),
		sceneCode:        nil,
		tag:              tag,
	}
	return scene
}

// AddEntity adds a new entity to the scene.
func (scene *Scene) AddEntity(entity IEntity) bool {
	Logger.Trace().Str("scene", scene.GetName()).Str("Entity", entity.GetName()).Msg("add entity")
	scene.entities = append(scene.entities, entity)
	scene.unloadedEntities = append(scene.unloadedEntities, entity)
	entity.SetScene(scene)
	return true
}

// AuditEntities displays all entities for audit purposes.
func (scene *Scene) AuditEntities() {
	for i, entity := range scene.GetEntities() {
		fmt.Printf("%d entity: [%s] %s\n", i, entity.GetID(), entity.GetName())
	}
}

// DeleteEntity deletes a entity from the scene.
func (scene *Scene) DeleteEntity(entity IEntity) bool {
	Logger.Trace().Str("scene", scene.GetName()).Str("Entity", entity.GetName()).Msg("delete entity")
	for _, traverseObj := range scene.entities {
		if traverseObj == entity {
			// Entity to be deleted in OnAfterUpdate method.
			// scene.Entities = append(scene.Entities[:i], scene.Entities[i+1:]...)
			scene.toDeleteEntities = append(scene.toDeleteEntities, entity)
			// Remove collider from the collision collection, so there is not
			// more checks between this collider and other other one.
			if index, ok := scene.getIndexInCollisionCollectionByEntity(entity); ok {
				scene.collisionCollection = append(scene.collisionCollection[:index], scene.collisionCollection[index+1:]...)
			}
			// Trigger destroy delegate
			destroyDelegate := GetDelegateManager().GetDestroyDelegate()
			GetDelegateManager().TriggerDelegate(destroyDelegate, true, entity)
			return true
		}
	}
	return false
}

// DoDestroy calls all methods to clean up scene.
func (scene *Scene) DoDestroy() {
	Logger.Trace().Str("scene", scene.GetName()).Msg("DoDestroy")
	scene.SetLoaded(false)
	for _, entity := range scene.loadedEntities {
		entity.DoDestroy()
	}
	scene.entities = []IEntity{}
	scene.loadedEntities = []IEntity{}
	scene.unloadedEntities = []IEntity{}
	scene.collisionCollection = []ICollider{}
	scene.layers = make([][]IEntity, maxLayers)
}

// DoFrameEnd calls all methods to run at the end of a tick frame.
func (scene *Scene) DoFrameEnd() {
}

// DoFrameStart calls all methods to run at the start of a tick frame.
func (scene *Scene) DoFrameStart() {
	scene.loadUnloadedEntities()
	for _, entity := range scene.loadedEntities {
		entity.DoFrameStart()
	}
}

// DoLoad is call when scene is loaded in the scene handler.
func (scene *Scene) DoLoad() {
	Logger.Trace().Str("scene", scene.GetName()).Msg("DoLoad")
	scene.SetLoaded(true)
	scene.loadUnloadedEntities()
}

// DoSwapBack is called when a scene is swap to.
func (scene *Scene) DoSwapBack() {
	Logger.Trace().Str("scene", scene.GetName()).Msg("DoResume")
	// scene.SetLoaded(true)
	scene.loadUnloadedEntities()
}

// DoSwapFrom is called when scene is swap from, but it is not unloaded.
func (scene *Scene) DoSwapFrom() {
	Logger.Trace().Str("scene", scene.GetName()).Msg("DoPause")
	// scene.SetLoaded(false)
	for _, entity := range scene.loadedEntities {
		entity.DoUnLoad()
	}
	scene.loadedEntities = []IEntity{}
	scene.unloadedEntities = []IEntity{}
	for _, entity := range scene.GetEntities() {
		scene.unloadedEntities = append(scene.unloadedEntities, entity)
	}
	scene.collisionCollection = []ICollider{}
	scene.layers = make([][]IEntity, maxLayers)
}

// DoUnLoad is called when scene is unloaded from the scene handler.
func (scene *Scene) DoUnLoad() {
	Logger.Trace().Str("scene", scene.GetName()).Msg("DoUnLoad")
	scene.SetLoaded(false)
	for _, entity := range scene.loadedEntities {
		entity.DoUnLoad()
	}
	scene.entities = []IEntity{}
	scene.loadedEntities = []IEntity{}
	scene.unloadedEntities = []IEntity{}
	// for _, entity := range scene.GetEntities() {
	// 	scene.unloadedEntities = append(scene.unloadedEntities, entity)
	// }
	scene.collisionCollection = []ICollider{}
	scene.layers = make([][]IEntity, maxLayers)
}

// getEntity returns entity and index for the given name.
func (scene *Scene) getEntity(id string) (IEntity, int) {
	for i, entity := range scene.entities {
		if entity.GetID() == id {
			return entity, i
		}
	}
	return nil, -1
}

// getEntityByName returns entity and index for the given name.
func (scene *Scene) getEntityByName(name string) (IEntity, int) {
	for i, entity := range scene.entities {
		if entity.GetName() == name {
			return entity, i
		}
	}
	return nil, -1
}

// GetEntities returns all Entities in the scene.
func (scene *Scene) GetEntities() []IEntity {
	return scene.entities
}

// GetEntitiesByTag returns all entities in the scene with the given tag.
func (scene *Scene) GetEntitiesByTag(tag string) []IEntity {
	result := []IEntity{}
	for _, entity := range scene.GetEntities() {
		if entity.GetTag() == tag {
			result = append(result, entity)
		}
	}
	return result
}

// GetEntity returns a entity for the entity ID.
func (scene *Scene) GetEntity(id string) IEntity {
	for _, entity := range scene.entities {
		if entity.GetID() == id {
			return entity
		}
	}
	return nil
}

// GetEntityByName returns a entity for the given name.
func (scene *Scene) GetEntityByName(name string) IEntity {
	for _, entity := range scene.entities {
		if entity.GetName() == name {
			return entity
		}
	}
	return nil
}

// getIndexInCollisionCollectionByEntity returns the index for the collider in
// the collision collection that belongs to the given entity.
func (scene *Scene) getIndexInCollisionCollectionByEntity(entity IEntity) (int, bool) {
	for i, collider := range scene.collisionCollection {
		if collider.GetEntity().GetID() == entity.GetID() {
			return i, true
		}
	}
	return -1, false
}

// getIndexInLoadedEntity return the index for the given entity in
// loadedEntity array.
func (scene Scene) getIndexInLoadedEntity(entity IEntity) (int, bool) {
	for i, obj := range scene.loadedEntities {
		if obj.GetID() == entity.GetID() {
			return i, true
		}
	}
	return -1, false
}

// getIndexInUnloadedEntity return the index for the given entity in
// unloadedEntity array.
func (scene Scene) getIndexInUnloadedEntity(entity IEntity) (int, bool) {
	for i, obj := range scene.unloadedEntities {
		if obj.GetName() == entity.GetName() {
			return i, true
		}
	}
	return -1, false
}

// GetSceneCode returns the scene code.
func (scene *Scene) GetSceneCode() TSceneCodeSignature {
	return scene.sceneCode
}

// GetTag returns the scene tag.
func (scene *Scene) GetTag() string {
	return scene.tag
}

// loadUnloadedEntities proceeds to load any unloaded entity
func (scene *Scene) loadUnloadedEntities() {
	unloaded := []IEntity{}
	for _, entity := range scene.unloadedEntities {
		if entity.GetActive() {
			entity.DoLoad()
			scene.loadedEntities = append(scene.loadedEntities, entity)
			layer := entity.GetLayer()
			scene.layers[layer] = append(scene.layers[layer], entity)
			for _, component := range entity.GetComponents() {
				if component.GetActive() {
					if collider, ok := interface{}(component).(ICollider); ok {
						scene.collisionCollection = append(scene.collisionCollection, collider)
					}
				}
			}
		} else {
			unloaded = append(unloaded, entity)
		}
		// Trigger load delegate
		loadDelegate := GetDelegateManager().GetLoadDelegate()
		GetDelegateManager().TriggerDelegate(loadDelegate, true, entity)
	}
	scene.unloadedEntities = unloaded
}

// OnAfterUpdate calls executed after all DoUpdates have been executed and
// before OnRender.
func (scene *Scene) OnAfterUpdate() {
	// Delete all Entities being marked to be deleted
	if len(scene.toDeleteEntities) != 0 {
		for _, entity := range scene.toDeleteEntities {
			if _, i := scene.getEntity(entity.GetID()); i != -1 {
				scene.entities = append(scene.entities[:i], scene.entities[i+1:]...)
			}
			if index, ok := scene.getIndexInLoadedEntity(entity); ok {
				scene.loadedEntities = append(scene.loadedEntities[:index], scene.loadedEntities[index+1:]...)
			}
			if index, ok := scene.getIndexInUnloadedEntity(entity); ok {
				scene.unloadedEntities = append(scene.unloadedEntities[:index], scene.unloadedEntities[index+1:]...)
			}
			entity.DoUnLoad()
		}
		scene.toDeleteEntities = []IEntity{}
	}
}

// OnRender calls all Entities OnRender methods. It call active entities using
// layers struct, calling from background to top layer.
func (scene *Scene) OnRender() {
	for _, layer := range scene.layers {
		for _, entity := range layer {
			if entity.GetActive() {
				entity.OnRender()
			}
		}
	}
}

// OnEnable calls all entity OnEnable methods.
func (scene *Scene) OnEnable() {
	for _, entity := range scene.entities {
		entity.OnEnable()
	}
}

// OnStart calls all Entities OnStart methods.
func (scene *Scene) OnStart() {
	Logger.Trace().Str("scene", scene.GetName()).Msg("OnStart")
	for _, entity := range scene.entities {
		if entity.GetActive() {
			entity.OnStart()
		}
	}
}

// OnUpdate calls all Entities OnUpdate methods. It does not use layers struct,
// but loadEntities struct. It calls to test collision in all entities active
// in the scene.
func (scene *Scene) OnUpdate() {
	for i := 0; i < len(scene.collisionCollection); i++ {
		colliderI := scene.collisionCollection[i]
		collisionBoxI := colliderI.GetCollisionBox()
		centerI := collisionBoxI.GetCenter()
		radiusI := collisionBoxI.GetRadius()
		// rectI := collisionBoxI.GetRect()
		entityI := colliderI.GetEntity()
		for j := i + 1; j < len(scene.collisionCollection); j++ {
			colliderJ := scene.collisionCollection[j]
			collisionBoxJ := colliderJ.GetCollisionBox()
			centerJ := collisionBoxJ.GetCenter()
			radiusJ := collisionBoxJ.GetRadius()
			// rectJ := collisionBoxJ.GetRect()
			entityJ := colliderJ.GetEntity()
			distance := math.Sqrt(math.Pow(centerI.X-centerJ.X, 2) + math.Pow(centerI.Y-centerJ.Y, 2))
			if distance < (radiusI + radiusJ) {
				// if rectI.HasIntersection(rectJ) {
				fmt.Printf("check collision %s with %s\n", entityI.GetName(), entityJ.GetName())
				delegate := GetDelegateManager().GetCollisionDelegate()
				GetDelegateManager().TriggerDelegate(delegate, true, entityI, entityJ)
			}
		}
	}
	for _, entity := range scene.loadedEntities {
		if entity.GetActive() {
			entity.OnUpdate()
		}
	}
}

// SetSceneCode sets the scene code.
func (scene *Scene) SetSceneCode(sceneCode TSceneCodeSignature) {
	scene.sceneCode = sceneCode
}

// SetTag sets the scene tag.
func (scene *Scene) SetTag(tag string) {
	scene.tag = tag
}
