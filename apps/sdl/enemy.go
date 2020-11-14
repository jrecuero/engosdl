package main

import (
	"strconv"

	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
)

// createEnemyController creates enemy controller entity.
func createEnemyController(totalEnemies int) engosdl.IEntity {
	enemyController := engosdl.NewEntity("enemy-controller")
	enemyControllerComponent := NewEnemyController("enemy-controller", totalEnemies)
	enemyControllerComponent.DefaultAddDelegateToRegister()
	enemyController.AddComponent(enemyControllerComponent)
	return enemyController
}

// createEnemies creates all enemies instances.
func createEnemies(engine *engosdl.Engine, maxEnemies int, enemyController engosdl.IEntity) []engosdl.IEntity {
	engosdl.GetComponentManager().RegisterComponent(&EnemyController{})
	engosdl.GetComponentManager().RegisterComponent(&enemySpriteT{})
	enemies := []engosdl.IEntity{}
	var x float64 = 10
	for i := 0; i < maxEnemies; i++ {
		enemy := createEnemy(engine, i, engosdl.NewVector(x, 10), enemyController)
		enemies = append(enemies, enemy)
		x += 150
	}
	return enemies
}

// createEnemy creates a single enemy instance.
func createEnemy(engine *engosdl.Engine, index int, position *engosdl.Vector, enemyController engosdl.IEntity) engosdl.IEntity {
	enemy := engosdl.NewEntity("enemy-" + strconv.Itoa(index))
	enemy.SetTag("enemy")
	enemy.GetTransform().SetPosition(position)
	// enemy.GetTransform().SetScale(engosdl.NewVector(0.5, 0.5))

	enemyOutOfBounds := components.NewOutOfBounds("enemy-out-of-bounds", true)
	enemyOutOfBounds.DefaultAddDelegateToRegister()
	enemyMove := components.NewMoveTo("enemy-move", engosdl.NewVector(5, 0))
	// enemySprite := components.NewSprite("enemy-sprite", "images/basic_enemy.bmp", engine.GetRenderer())
	// enemySprite := components.NewMultiSprite("enemy-sprite", []string{"images/basic_enemy.bmp"}, engine.GetRenderer())
	// enemySprite := components.NewSpriteSheet("enemy-sprite", []string{"images/enemies.bmp"}, 3, engine.GetRenderer())
	enemySprite := newEnemySprite("enemy-sprite", []string{"images/enemies.bmp"}, 3)
	enemySprite.SetDestroyOnOutOfBounds(false)
	// enemySprite.DefaultAddDelegateToRegister()
	enemySprite.AddDelegateToRegister(engosdl.GetDelegateManager().GetCollisionDelegate(), nil, nil, enemySprite.DefaultOnCollision)
	enemySprite.AddDelegateToRegister(nil, enemy, &components.Timer{}, func(params ...interface{}) bool {
		x, y := enemySprite.GetEntity().GetTransform().GetPosition().Get()
		enemySprite.GetEntity().GetTransform().SetPosition(engosdl.NewVector(x, y+10))
		return true
	})
	enemyStats := components.NewEntityStats("enemy-stats", 50)
	enemyStats.DefaultAddDelegateToRegister()
	enemyCollider := components.NewCollider2D("enemy-collider-2D")
	enemyCollider.DefaultAddDelegateToRegister()
	enemyTimer := components.NewTimer("enemy-timer", 100)

	enemy.AddComponent(enemyMove)
	enemy.AddComponent(enemyOutOfBounds)
	enemy.AddComponent(enemySprite)
	enemy.AddComponent(enemyCollider)
	enemy.AddComponent(enemyStats)
	enemy.AddComponent(enemyTimer)

	if controller, ok := enemyController.GetComponent(&EnemyController{}).(*EnemyController); ok {
		controller.AddDelegateToRegister(nil, enemy, &components.OutOfBounds{}, controller.onOutOfBounds)
		enemySprite.AddDelegateToRegister(nil, controller.GetEntity(), controller, func(params ...interface{}) bool {
			speed := enemyMove.GetSpeed()
			enemyMove.SetSpeed(engosdl.NewVector(speed.X*-1, speed.Y*-1))
			enemy := enemySprite.GetEntity()
			position := enemy.GetTransform().GetPosition()
			enemy.GetTransform().SetPosition(engosdl.NewVector(position.X-speed.X, position.Y-speed.Y))
			_ = enemy
			return true
		})
	}

	return enemy
}
