package engosdl

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
	OnAwake()
	OnDraw()
	OnEnable()
	OnStart()
	OnUpdate()
}

// Scene is the default implementation for IScene interface.
type Scene struct {
	*Object
	Entities         []IEntity
	toDeleteEntities []IEntity
	loadedEntities   []IEntity
	unloadedEntities []IEntity
	loaded           bool
}

var _ IScene = (*Scene)(nil)

// NewScene creates a new scene instance
func NewScene(name string) *Scene {
	Logger.Trace().Str("scene", name).Msg("new scene")
	return &Scene{
		Object:           NewObject(name),
		Entities:         []IEntity{},
		toDeleteEntities: []IEntity{},
		loadedEntities:   []IEntity{},
		unloadedEntities: []IEntity{},
		loaded:           false,
	}
}

// AddEntity adds a new entity to the scene.
func (scene *Scene) AddEntity(entity IEntity) bool {
	Logger.Trace().Str("scene", scene.GetName()).Str("Entity", entity.GetName()).Msg("add entity")
	scene.Entities = append(scene.Entities, entity)
	scene.unloadedEntities = append(scene.unloadedEntities, entity)
	entity.SetScene(scene)
	return true
}

// DeleteEntity deletes a entity from the scene.
func (scene *Scene) DeleteEntity(entity IEntity) bool {
	Logger.Trace().Str("scene", scene.GetName()).Str("Entity", entity.GetName()).Msg("delete entity")
	for _, travObj := range scene.Entities {
		if travObj == entity {
			// Entity to be deleted in OnAfterUpdate method.
			// scene.Entities = append(scene.Entities[:i], scene.Entities[i+1:]...)
			scene.toDeleteEntities = append(scene.toDeleteEntities, entity)
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
}

// getEntity returns entity and index for the given name.
func (scene *Scene) getEntity(name string) (IEntity, int) {
	for i, entity := range scene.Entities {
		if entity.GetName() == name {
			return entity, i
		}
	}
	return nil, -1
}

// GetEntities returns all Entities in the scene.
func (scene *Scene) GetEntities() []IEntity {
	return scene.Entities
}

// GetEntity returns a entity for the entity ID.
func (scene *Scene) GetEntity(id string) IEntity {
	for _, entity := range scene.Entities {
		if entity.GetID() == id {
			return entity
		}
	}
	return nil
}

// GetEntityByName returns a entity for the given name.
func (scene *Scene) GetEntityByName(name string) IEntity {
	for _, entity := range scene.Entities {
		if entity.GetName() == name {
			return entity
		}
	}
	return nil
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
		} else {
			unloaded = append(unloaded, entity)
		}
	}
	scene.unloadedEntities = unloaded
}

// OnAfterUpdate calls executed after all DoUpdates have been executed and
// before OnDraw.
func (scene *Scene) OnAfterUpdate() {
	// Delete all Entities being marked to be deleted
	if len(scene.toDeleteEntities) != 0 {
		for _, entity := range scene.toDeleteEntities {
			if _, i := scene.getEntity(entity.GetName()); i != -1 {
				scene.Entities = append(scene.Entities[:i], scene.Entities[i+1:]...)
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

// OnAwake calls all entity OnAwake methods.
func (scene *Scene) OnAwake() {
	for _, entity := range scene.Entities {
		entity.OnAwake()
	}
}

// OnDraw calls all Entities OnDraw methods.
func (scene *Scene) OnDraw() {
	for _, entity := range scene.loadedEntities {
		if entity.GetActive() {
			entity.OnDraw()
		}
	}
}

// OnEnable calls all entity OnEnable methods.
func (scene *Scene) OnEnable() {
	for _, entity := range scene.Entities {
		entity.OnEnable()
	}
}

// OnStart calls all Entities OnStart methods.
func (scene *Scene) OnStart() {
	for _, entity := range scene.Entities {
		entity.OnStart()
	}
}

// OnUpdate calls all Entities OnUpdate methods.
func (scene *Scene) OnUpdate() {
	for _, entity := range scene.loadedEntities {
		if entity.GetActive() {
			entity.OnUpdate()
		}
	}
}
