package components

import (
	"fmt"

	"github.com/jrecuero/engosdl"
)

// EntityStats represemts a component that contains generic entity stats.
type EntityStats struct {
	*engosdl.Component
	life int
}

// onCollision checks when there is a collision with other entity.
func (c *EntityStats) onCollision(params ...interface{}) bool {
	collisionEntityOne := params[0].(*engosdl.Entity)
	collisionEntityTwo := params[1].(*engosdl.Entity)
	if c.GetEntity().GetID() == collisionEntityOne.GetID() || c.GetEntity().GetID() == collisionEntityTwo.GetID() {
		c.life -= 10
		fmt.Printf("%s [live %d] onCollision %s with %s\n", c.GetEntity().GetName(), c.life, collisionEntityOne.GetName(), collisionEntityTwo.GetName())
		if c.life == 0 {
			engosdl.GetEngine().DestroyEntity(c.GetEntity())
		}
	}
	return true
}

// OnStart is called first time the component is enabled.
func (c *EntityStats) OnStart() {
	engosdl.Logger.Trace().Str("component", "entity-stats").Str("entity-stats", c.GetName()).Msg("OnStart")
	delegateHandler := engosdl.GetEngine().GetEventHandler().GetDelegateHandler()
	delegate := delegateHandler.GetCollisionDelegate()
	delegateHandler.CreateDelegate(c, "entity-stats")
	delegateHandler.RegisterToDelegate(delegate, c.onCollision)
}

// NewEntityStats creates a new entity stats instance.
func NewEntityStats(name string, life int) *EntityStats {
	engosdl.Logger.Trace().Str("component", "entity-stats").Str("entity-stats", name).Msg("new entity-stats")
	return &EntityStats{
		Component: engosdl.NewComponent(name),
		life:      life,
	}
}
