package main

import (
	"strconv"

	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
)

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

// onOutOfBounds is called when any enemy goes out of bounds.
func (c *EnemyController) onOutOfBounds(params ...interface{}) bool {
	if enemy, ok := params[0].(engosdl.IEntity); ok {
		// fmt.Println("[Controller] enemy " + enemy.GetName() + " is out of bounds")
		engosdl.GetDelegateHandler().TriggerDelegate(c.GetDelegate(), false, enemy)
	}
	return true
}

// OnAwake is called when component is added to an entity. It creates all
// non-dependent resources.
func (c *EnemyController) OnAwake() {
	engosdl.Logger.Trace().Str("component", "enemy-controller").Msg("OnAwake")
	c.SetDelegate(engosdl.GetDelegateHandler().CreateDelegate(c, "enemy-controller"))
	c.Component.OnAwake()
}

// onDestroy is called when any entity is removed from the scenario.
func (c *EnemyController) onDestroy(params ...interface{}) bool {
	entity := params[0].(engosdl.IEntity)
	// fmt.Printf("Entity %s has been destroyed\n", entity.GetName())
	if entity.GetTag() == "enemy" {
		c.totalEnemies--
		if c.totalEnemies == 0 {
			engosdl.GetEngine().GetSceneHandler().SetActiveFirstScene()
		}
	}
	return true
}

// onLoad is called when any entity is loaded in the scene.
func (c *EnemyController) onLoad(params ...interface{}) bool {
	// entity := params[0].(engosdl.IEntity)
	// fmt.Printf("Entity %s has been loaded\n", entity.GetName())
	return true
}

// OnStart is called first time the component is enabled.
func (c *EnemyController) OnStart() {
	c.AddDelegateToRegister(engosdl.GetDelegateHandler().GetLoadDelegate(), nil, nil, c.onLoad)
	c.AddDelegateToRegister(engosdl.GetDelegateHandler().GetDestroyDelegate(), nil, nil, c.onDestroy)
	// enemies := createEnemies(engosdl.GetEngine(), 2, c.GetEntity())
	// for _, enemy := range enemies {
	// 	c.AddDelegateToRegister(nil, enemy, &components.OutOfBounds{}, c.onOutOfBounds)
	// 	c.GetEntity().GetScene().AddEntity(enemy)
	// }
	c.Component.OnStart()
}

type enemySpriteT struct {
	*components.Sprite
	onCollisionF func(engosdl.ISprite) engosdl.TDelegateSignature
	hit          bool
}

func newEnemySprite(name string, filenames []string, numberOfSprites int) *enemySpriteT {
	result := &enemySpriteT{
		Sprite: components.NewSprite(name, filenames, numberOfSprites, engosdl.GetRenderer()),
		hit:    false,
	}
	result.AddDelegateToRegister(engosdl.GetDelegateHandler().GetCollisionDelegate(), nil, nil, result.onCollision)
	return result
}

// onCollision checks when there is a collision with other entity.
func (c *enemySpriteT) onCollision(params ...interface{}) bool {
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

// createEnemyController creates enemy controller entity.
func createEnemyController(totalEnemies int) engosdl.IEntity {
	enemyController := engosdl.NewEntity("enemy-controller")
	enemyController.AddComponent(NewEnemyController("enemy-controller", totalEnemies))
	return enemyController
}

// createEnemies creates all enemies instances.
func createEnemies(engine *engosdl.Engine, maxEnemies int, enemyController engosdl.IEntity) []engosdl.IEntity {
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
	enemySprite.DefaultAddDelegateToRegister()
	// enemySprite.AddDelegateToRegister(nil, enemy, &components.OutOfBounds{}, func(params ...interface{}) bool {
	// 	speed := enemyMove.GetSpeed()
	// 	enemyMove.SetSpeed(engosdl.NewVector(speed.X*-1, speed.Y*-1))
	// 	return true
	// })
	enemyStats := components.NewEntityStats("enemy-stats", 20)
	enemyStats.DefaultAddDelegateToRegister()
	enemyCollider := components.NewCollider2D("enemy-collider-2D")
	enemyCollider.DefaultAddDelegateToRegister()

	enemy.AddComponent(enemyMove)
	enemy.AddComponent(enemyOutOfBounds)
	enemy.AddComponent(enemySprite)
	enemy.AddComponent(enemyCollider)
	enemy.AddComponent(enemyStats)

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
