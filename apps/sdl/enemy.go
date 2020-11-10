package main

import (
	"fmt"
	"strconv"

	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
)

// EnemyController represents a component that control enemy entities.
type EnemyController struct {
	*engosdl.Component
}

// NewEnemyController creates new instance for component enemy controller.
func NewEnemyController(name string) *EnemyController {
	return &EnemyController{
		Component: engosdl.NewComponent(name),
	}
}

// onOutOfBounds is called when any enemy goes out of bounds.
func (c *EnemyController) onOutOfBounds(params ...interface{}) bool {
	if enemy, ok := params[0].(engosdl.IEntity); ok {
		fmt.Println("[Controller] enemy " + enemy.GetName() + " is out of bounds")
		engosdl.GetEngine().GetEventHandler().GetDelegateHandler().TriggerDelegate(c.GetDelegate(), false, enemy)
	}
	return true
}

// OnAwake is called when component is added to an entity. It creates all
// non-dependent resources.
func (c *EnemyController) OnAwake() {
	engosdl.Logger.Trace().Str("component", "enemy-controller").Msg("OnAwake")
	c.SetDelegate(engosdl.GetEngine().GetEventHandler().GetDelegateHandler().CreateDelegate(c, "enemy-controller"))
}

type enemySpriteT struct {
	*components.SpriteSheet
	onCollisionF func(engosdl.ISprite) engosdl.TDelegateSignature
	hit          bool
}

func newEnemySprite(name string, filenames []string, numberOfSprites int) *enemySpriteT {
	result := &enemySpriteT{
		SpriteSheet: components.NewSpriteSheet(name, filenames, numberOfSprites, engosdl.GetRenderer()),
		hit:         false,
	}
	// result.onCollisionF = func(instance engosdl.ISprite) engosdl.TDelegateSignature {
	// 	sprite := instance.(*enemySpriteT)
	// 	return sprite.onCollision
	// 	// return func(params ...interface{}) bool {
	// 	// 	collisionEntityOne := params[0].(*engosdl.Entity)
	// 	// 	collisionEntityTwo := params[1].(*engosdl.Entity)
	// 	// 	if collisionEntityOne.GetTag() == "bullet" || collisionEntityTwo.GetTag() == "bullet" {
	// 	// 		sprite := instance.(*enemySpriteT)
	// 	// 		sprite.hit = true
	// 	// 	}
	// 	// 	return true
	// 	// }
	// }
	result.AddDelegateToRegister(engosdl.GetEngine().GetEventHandler().GetDelegateHandler().GetCollisionDelegate(), nil, nil, result.onCollision)
	return result
}

// onCollision checks when there is a collision with other entity.
func (c *enemySpriteT) onCollision(params ...interface{}) bool {
	collisionEntityOne := params[0].(*engosdl.Entity)
	collisionEntityTwo := params[1].(*engosdl.Entity)
	if collisionEntityOne.GetTag() == "bullet" || collisionEntityTwo.GetTag() == "bullet" {
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

// // OnStart is called first time the component is enabled.
// func (c *enemySpriteT) OnStart() {
// 	// Register to: "on-collision" and "out-of-bounds"
// 	engosdl.Logger.Trace().Str("component", "sprite").Str("sprite", c.GetName()).Msg("OnStart")
// 	if c.CanRegisterTo(engosdl.CollisionName) {
// 		delegate := engosdl.GetEngine().GetEventHandler().GetDelegateHandler().GetCollisionDelegate()
// 		// c.AddDelegateToRegister(delegate, nil, nil, c.onCollision)
// 		c.AddDelegateToRegister(delegate, nil, nil, c.onCollisionF(c))
// 	}
// 	c.Component.OnStart()
// }

// createEnemyController creates enemy controller entity.
func createEnemyController() engosdl.IEntity {
	enemyController := engosdl.NewEntity("enemy-controller")
	enemyController.AddComponent(NewEnemyController("enemy-controller"))
	return enemyController
}

// createEnemies creates all enemies instances.
func createEnemies(engine *engosdl.Engine, maxEnemies int, enemyController engosdl.IEntity) []engosdl.IEntity {
	enemies := []engosdl.IEntity{}
	var x float64 = 10
	for i := 0; i < maxEnemies; i++ {
		enemy := createEnemy(engine, i, engosdl.NewVector(x, 10), enemyController)
		enemies = append(enemies, enemy)
		x += 110
	}
	return enemies
}

// createEnemy creates a single enemy instance.
func createEnemy(engine *engosdl.Engine, index int, position *engosdl.Vector, enemyController engosdl.IEntity) engosdl.IEntity {
	enemy := engosdl.NewEntity("enemy-" + strconv.Itoa(index))
	enemy.SetTag("enemy")
	enemy.GetTransform().SetPosition(position)
	enemy.GetTransform().SetScale(engosdl.NewVector(0.5, 0.5))

	enemyOutOfBounds := components.NewOutOfBounds("enemy-out-of-bounds", true)
	enemyMove := components.NewMoveTo("enemy-move", engosdl.NewVector(5, 0))
	enemyMove.SetRegisterOnStart(map[string]bool{engosdl.OutOfBoundsName: false})
	// enemySprite := components.NewSprite("enemy-sprite", "images/basic_enemy.bmp", engine.GetRenderer())
	// enemySprite := components.NewMultiSprite("enemy-sprite", []string{"images/basic_enemy.bmp"}, engine.GetRenderer())
	// enemySprite := components.NewSpriteSheet("enemy-sprite", []string{"images/enemies.bmp"}, 3, engine.GetRenderer())
	enemySprite := newEnemySprite("enemy-sprite", []string{"images/enemies.bmp"}, 3)
	enemySprite.SetDestroyOnOutOfBounds(false)
	// enemySprite.AddDelegateToRegister(nil, enemy, &components.OutOfBounds{}, func(params ...interface{}) bool {
	// 	speed := enemyMove.GetSpeed()
	// 	enemyMove.SetSpeed(engosdl.NewVector(speed.X*-1, speed.Y*-1))
	// 	return true
	// })
	enemyStats := components.NewEntityStats("enemy-stats", 100)

	enemy.AddComponent(enemyMove)
	enemy.AddComponent(enemyOutOfBounds)
	enemy.AddComponent(enemySprite)
	enemy.AddComponent(components.NewCollider2D("enemy-collider-2D"))
	enemy.AddComponent(enemyStats)

	if controller, ok := enemyController.GetComponent(&EnemyController{}).(*EnemyController); ok {
		controller.AddDelegateToRegister(nil, enemy, &components.OutOfBounds{}, controller.onOutOfBounds)
		enemySprite.AddDelegateToRegister(controller.GetDelegate(), nil, nil, func(params ...interface{}) bool {
			speed := enemyMove.GetSpeed()
			enemyMove.SetSpeed(engosdl.NewVector(speed.X*-1, speed.Y*-1))
			// if enemy, ok := params[0].(engosdl.IEntity); ok {
			// if enemy.GetID() == enemySprite.GetEntity().GetID() {
			enemy := enemySprite.GetEntity()
			position := enemy.GetTransform().GetPosition()
			enemy.GetTransform().SetPosition(engosdl.NewVector(position.X-speed.X, position.Y-speed.Y))
			_ = enemy
			// }
			// }
			return true
		})
	}

	return enemy
}
