package main

import (
	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
)

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor("*main.EnemyController", CreateEnemyController)
		componentManager.RegisterConstructor("*main.EnemySprite", CreateEnemySprite)
	}
}

// EnemyController represents a component that control enemy entities.
type EnemyController struct {
	*engosdl.Component
	totalEnemies int
}

// NewEnemyController creates new instance for component enemy controller.
func NewEnemyController(name string, totalEnemies int) *EnemyController {
	return &EnemyController{
		Component:    engosdl.NewComponent(name),
		totalEnemies: totalEnemies,
	}
}

// CreateEnemyController is used by component manager to create a new
// enemy controller instance.
func CreateEnemyController(params ...interface{}) engosdl.IComponent {
	if len(params) == 2 {
		return NewEnemyController(params[0].(string), params[1].(int))
	}
	return NewEnemyController("", 0)
}

// DefaultAddDelegateToRegister will proceed to add default delegate to
// register for the component.
func (c *EnemyController) DefaultAddDelegateToRegister() {
	// c.AddDelegateToRegister(engosdl.GetDelegateManager().GetLoadDelegate(), nil, nil, c.onLoad)
	// c.AddDelegateToRegister(engosdl.GetDelegateManager().GetDestroyDelegate(), nil, nil, c.onDestroy)
}

// onOutOfBounds is called when any enemy goes out of bounds.
func (c *EnemyController) onOutOfBounds(params ...interface{}) bool {
	if enemy, ok := params[0].(engosdl.IEntity); ok {
		// fmt.Println("[Controller] enemy " + enemy.GetName() + " is out of bounds")
		engosdl.GetDelegateManager().TriggerDelegate(c.GetDelegate(), false, enemy)
	}
	return true
}

// OnAwake is called when component is added to an entity. It creates all
// non-dependent resources.
func (c *EnemyController) OnAwake() {
	engosdl.Logger.Trace().Str("component", "enemy-controller").Msg("OnAwake")
	c.SetDelegate(engosdl.GetDelegateManager().CreateDelegate(c, "enemy-controller"))
	c.Component.OnAwake()
}

// onDestroy is called when any entity is removed from the scenario.
func (c *EnemyController) onDestroy(params ...interface{}) bool {
	entity := params[0].(engosdl.IEntity)
	// fmt.Printf("Entity %s has been destroyed\n", entity.GetName())
	if entity.GetTag() == "enemy" {
		c.totalEnemies--
		if c.totalEnemies == 0 {
			engosdl.GetEngine().GetSceneManager().SetActiveFirstScene()
		}
	}
	return true
}

// onLoad is called when any entity is loaded in the scene.
func (c *EnemyController) onLoad(params ...interface{}) bool {
	return true
}

// OnStart is called first time the component is enabled.
func (c *EnemyController) OnStart() {
	c.Component.OnStart()
}

type enemySpriteT struct {
	*components.Sprite
	onCollisionF func(engosdl.ISprite) engosdl.TDelegateSignature
	hit          bool
}

func newEnemySprite(name string, filenames []string, numberOfSprites int) *enemySpriteT {
	result := &enemySpriteT{
		Sprite: components.NewSprite(name, filenames, numberOfSprites, engosdl.FormatBMP),
		hit:    false,
	}
	// result.AddDelegateToRegister(engosdl.GetDelegateManager().GetCollisionDelegate(), nil, nil, result.onCollision)
	return result
}

// CreateEnemySprite is used by component manager to create a new enemy sprite
// instance.
func CreateEnemySprite(params ...interface{}) engosdl.IComponent {
	if len(params) == 3 {
		return newEnemySprite(params[0].(string), params[1].([]string), params[2].(int))
	}
	return newEnemySprite("", []string{}, 0)
}

// DefaultOnCollision checks when there is a collision with other entity.
func (c *enemySpriteT) DefaultOnCollision(params ...interface{}) bool {
	collisionEntityOne := params[0].(*engosdl.Entity)
	collisionEntityTwo := params[1].(*engosdl.Entity)
	if (collisionEntityOne.GetTag() == "bullet" || collisionEntityTwo.GetTag() == "bullet") &&
		(collisionEntityOne.GetID() == c.GetEntity().GetID() || collisionEntityTwo.GetID() == c.GetEntity().GetID()) {
		c.hit = true
	}
	return true
}

// OnUpdate is called for every update tick.
func (c *enemySpriteT) OnUpdate() {
	if c.hit {
		c.NextSprite()
		if c.GetSpriteIndex() == 0 {
			c.hit = false
		}
	}
}
