package components

import (
	"fmt"

	"github.com/jrecuero/engosdl"
)

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterComponent(&EntityStats{})
	}
}

// EntityStats represents a component that contains generic entity stats.
type EntityStats struct {
	*engosdl.Component
	Life       int `json:"life"`
	Experience int `json:"experience"`
}

// NewEntityStats creates a new entity stats instance.
// It creates delegate "on-enemy-stats".
// It registers to "on-collision" delegate.
func NewEntityStats(name string, life int) *EntityStats {
	engosdl.Logger.Trace().Str("component", "entity-stats").Str("entity-stats", name).Msg("new entity-stats")
	result := &EntityStats{
		Component: engosdl.NewComponent(name),
		Life:      life,
	}
	return result
}

// DefaultAddDelegateToRegister will proceed to add default delegate to
// register for the component.
func (c *EntityStats) DefaultAddDelegateToRegister() {
	c.AddDelegateToRegister(engosdl.GetDelegateManager().GetCollisionDelegate(), nil, nil, c.onCollision)
}

// OnAwake should create all component resources that don't have any dependency
// with any other component or entity.
func (c *EntityStats) OnAwake() {
	// Create new delegate "entity-stats"
	engosdl.Logger.Trace().Str("component", "entity-stats").Str("entity-stats", c.GetName()).Msg("OnAwake")
	c.SetDelegate(engosdl.GetDelegateManager().CreateDelegate(c, "on-entity-stats"))
	c.Component.OnAwake()

}

// onCollision checks when there is a collision with other entity.
func (c *EntityStats) onCollision(params ...interface{}) bool {
	collisionEntityOne := params[0].(*engosdl.Entity)
	collisionEntityTwo := params[1].(*engosdl.Entity)
	var me, other *engosdl.Entity
	if c.GetEntity().GetID() == collisionEntityOne.GetID() {
		me = collisionEntityOne
		other = collisionEntityTwo
	} else if c.GetEntity().GetID() == collisionEntityTwo.GetID() {
		me = collisionEntityTwo
		other = collisionEntityOne
	}
	if me != nil && other != nil {
		tag := c.GetEntity().GetTag()
		var otherTag string
		if parent := other.GetParent(); parent != nil {
			otherTag = parent.GetTag()
		}
		if tag != otherTag {
			c.Life -= 10
			engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), true, c.Life)
			fmt.Printf("%s [live %d] onCollision %s with %s\n",
				c.GetEntity().GetName(),
				c.Life,
				collisionEntityOne.GetName(),
				collisionEntityTwo.GetName())
			if c.Life == 0 {
				engosdl.GetEngine().DestroyEntity(c.GetEntity())
			}
			engosdl.GetEngine().DestroyEntity(other)
		}
	}
	return true
}

// OnStart is called first time the component is enabled.
func (c *EntityStats) OnStart() {
	engosdl.Logger.Trace().Str("component", "entity-stats").Str("entity-stats", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}
