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
	enemy.GetTransform().SetPosition(position)
	enemy.GetTransform().SetScale(engosdl.NewVector(0.5, 0.5))

	enemyOutOfBounds := components.NewOutOfBounds("enemy-out-of-bounds", true)
	enemyMove := components.NewMoveTo("enemy-move", engosdl.NewVector(5, 0))
	enemyMove.SetRegisterOnStart(map[string]bool{engosdl.OutOfBoundsName: false})
	enemySprite := components.NewSprite("enemy-sprite", "images/basic_enemy.bmp", engine.GetRenderer())
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
