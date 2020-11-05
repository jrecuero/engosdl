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

// IScene represents the interface for any game scene
type IScene interface {
	IObject
	AddEntity(IEntity) bool
	DeleteEntity(IEntity) bool
	DoCycleEnd()
	DoCycleStart()
	DoLoad()
	DoUnLoad()
	GetEntities() []IEntity
	GetEntity(string) IEntity
	GetEntityByName(string) IEntity
	OnAfterUpdate()
	OnDraw()
	OnEnable()
	OnStart()
	OnUpdate()
}

// Scene is the default implementation for IScene interface.
type Scene struct {
	*Object
	entities            []IEntity
	toDeleteEntities    []IEntity
	loadedEntities      []IEntity
	unloadedEntities    []IEntity
	layers              [][]IEntity
	loaded              bool
	collisionCollection []ICollider
}

var _ IScene = (*Scene)(nil)

// NewScene creates a new scene instance
func NewScene(name string) *Scene {
	Logger.Trace().Str("scene", name).Msg("new scene")
	scene := &Scene{
		Object:           NewObject(name),
		entities:         []IEntity{},
		toDeleteEntities: []IEntity{},
		loadedEntities:   []IEntity{},
		unloadedEntities: []IEntity{},
		layers:           make([][]IEntity, maxLayers),
		loaded:           false,
		// layers:           [][]IEntity{{}, {}, {}, {}},
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
			destroyDelegate := GetEngine().GetEventHandler().GetDelegateHandler().GetDestroyDelegate()
			GetEngine().GetEventHandler().GetDelegateHandler().TriggerDelegate(destroyDelegate, entity)
			return true
		}
	}
	return false
}

// DoCycleEnd calls all methods to run at the end of a tick cycle.
func (scene *Scene) DoCycleEnd() {
}

// DoCycleStart calls all methods to run at the start of a tick cycle.
func (scene *Scene) DoCycleStart() {
	scene.loadUnloadedEntities()
	for _, entity := range scene.loadedEntities {
		entity.DoCycleStart()
	}
}

// DoLoad is call when scene is loaded in the scene handler.
func (scene *Scene) DoLoad() {
	scene.loaded = true
	scene.loadUnloadedEntities()
}

// DoUnLoad is called when scene is unloaded from the scene handler.
func (scene *Scene) DoUnLoad() {
	scene.loaded = false
	for _, entity := range scene.loadedEntities {
		entity.DoUnLoad()
	}
	scene.loadedEntities = []IEntity{}
	scene.unloadedEntities = []IEntity{}
	scene.collisionCollection = []ICollider{}
	scene.layers = make([][]IEntity, maxLayers)
}

// getEntity returns entity and index for the given name.
func (scene *Scene) getEntity(name string) (IEntity, int) {
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
		if obj.GetName() == entity.GetName() {
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
		loadDelegate := GetEngine().GetEventHandler().GetDelegateHandler().GetLoadDelegate()
		GetEngine().GetEventHandler().GetDelegateHandler().TriggerDelegate(loadDelegate, entity)
	}
	scene.unloadedEntities = unloaded
}

// OnAfterUpdate calls executed after all DoUpdates have been executed and
// before OnDraw.
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
		}
		scene.toDeleteEntities = []IEntity{}
	}
}

// OnDraw calls all Entities OnDraw methods. It call active entities using
// layers struct, calling from background to top layer.
func (scene *Scene) OnDraw() {
	for _, layer := range scene.layers {
		for _, entity := range layer {
			if entity.GetActive() {
				entity.OnDraw()
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
	for _, entity := range scene.entities {
		entity.OnStart()
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
				delegate := GetEngine().GetEventHandler().GetDelegateHandler().GetCollisionDelegate()
				GetEngine().GetEventHandler().GetDelegateHandler().TriggerDelegate(delegate, entityI, entityJ)
			}
		}
	}
	for _, entity := range scene.loadedEntities {
		if entity.GetActive() {
			entity.OnUpdate()
		}
	}
}
