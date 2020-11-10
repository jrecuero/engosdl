package components

import (
	"fmt"

	"github.com/jrecuero/engosdl"
)

// EntityStats represents a component that contains generic entity stats.
type EntityStats struct {
	*engosdl.Component
	life int
}

// NewEntityStats creates a new entity stats instance.
// It registers to on-collision delegate.
func NewEntityStats(name string, life int) *EntityStats {
	engosdl.Logger.Trace().Str("component", "entity-stats").Str("entity-stats", name).Msg("new entity-stats")
	result := &EntityStats{
		Component: engosdl.NewComponent(name),
		life:      life,
	}
	result.AddDelegateToRegister(engosdl.GetDelegateHandler().GetCollisionDelegate(), nil, nil, result.onCollision)
	return result
}

// OnAwake should create all component resources that don't have any dependency
// with any other component or entity.
func (c *EntityStats) OnAwake() {
	// Create new delegate "entity-stats"
	engosdl.Logger.Trace().Str("component", "entity-stats").Str("entity-stats", c.GetName()).Msg("OnAwake")
	c.SetDelegate(engosdl.GetDelegateHandler().CreateDelegate(c, "entity-stats"))

}

// onCollision checks when there is a collision with other entity.
func (c *EntityStats) onCollision(params ...interface{}) bool {
	collisionEntityOne := params[0].(*engosdl.Entity)
	collisionEntityTwo := params[1].(*engosdl.Entity)
	if c.GetEntity().GetID() == collisionEntityOne.GetID() || c.GetEntity().GetID() == collisionEntityTwo.GetID() {
		c.life -= 10
		engosdl.GetDelegateHandler().TriggerDelegate(c.GetDelegate(), true, c.life)
		fmt.Printf("%s [live %d] onCollision %s with %s\n",
			c.GetEntity().GetName(),
			c.life,
			collisionEntityOne.GetName(),
			collisionEntityTwo.GetName())
		if c.life == 0 {
			engosdl.GetEngine().DestroyEntity(c.GetEntity())
		}
	}
	return true
}

// OnStart is called first time the component is enabled.
func (c *EntityStats) OnStart() {
	engosdl.Logger.Trace().Str("component", "entity-stats").Str("entity-stats", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}
