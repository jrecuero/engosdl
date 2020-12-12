package main

import (
	"fmt"

	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
	"github.com/veandco/go-sdl2/sdl"
)

// GameManager is the application game manager.
type GameManager struct {
	*engosdl.GameManager
}

var _ engosdl.IGameManager = (*GameManager)(nil)

// NewGameManager created a new game manager instance.
func NewGameManager(name string) *GameManager {
	engosdl.Logger.Trace().Str("game-manager", name).Msg("new game-manager")
	return &GameManager{
		GameManager: engosdl.NewGameManager(name),
	}
}

// CreateAssets creates all application assets and resources. it is called
// before game engine starts in order to create all required assets and
// resources.
func (h *GameManager) CreateAssets() {
	scene := engosdl.NewScene("play", "play")
	scene.SetSceneCode(h.createScene())
	engosdl.GetEngine().AddScene(scene)
}

func (h *GameManager) createBox(name string, position *engosdl.Vector, speed *engosdl.Vector) engosdl.IEntity {
	box := engosdl.NewEntity(name)
	box.SetCache("counter", 0)
	box.SetTag("box")
	box.GetTransform().SetPositionXY(position.X, position.Y)
	box.AddComponent(components.NewBox("box/box", &engosdl.Rect{W: 50, H: 50}, sdl.Color{B: 255}, true))
	box.AddComponent(components.NewBody("box/body", false))
	box.AddComponent(components.NewMoveTo("box/move-to", speed))
	box.GetComponent(&components.MoveTo{}).AddDelegateToRegister(nil, nil, &components.Body{}, func(params ...interface{}) bool {
		entity := params[0].(engosdl.IEntity)
		c := entity.GetComponent(&components.MoveTo{}).(*components.MoveTo)
		outAt := params[2].(int)
		switch outAt {
		case engosdl.Left, engosdl.Right:
			c.SetSpeed(engosdl.NewVector(c.GetSpeed().X*-1, c.GetSpeed().Y))
		case engosdl.Up, engosdl.Down:
			c.SetSpeed(engosdl.NewVector(c.GetSpeed().X, c.GetSpeed().Y*-1))
		}
		return true
	})
	box.GetComponent(&components.MoveTo{}).AddDelegateToRegister(engosdl.GetDelegateManager().GetCollisionDelegate(), nil, nil, func(params ...interface{}) bool {
		c := box.GetComponent(&components.MoveTo{}).(*components.MoveTo)
		if me, other, err := engosdl.EntitiesInCollision(c.GetEntity(), params...); err == nil && other.GetTag() == "box" {
			if counter, err := me.GetCache("counter"); err == nil {
				me.SetCache("counter", counter.(int)+1)
			}
			collisionRect := params[2].(*engosdl.Rect)
			if collisionRect.W > collisionRect.H {
				if (collisionRect.Y > me.GetTransform().GetRect().Y && c.GetSpeed().Y > 0) ||
					(collisionRect.Y <= me.GetTransform().GetRect().Y && c.GetSpeed().Y < 0) {
					c.SetSpeed(engosdl.NewVector(c.GetSpeed().X, c.GetSpeed().Y*-1))
				}
			} else {
				if (collisionRect.X > me.GetTransform().GetRect().X && c.GetSpeed().X > 0) ||
					(collisionRect.X <= me.GetTransform().GetRect().X && c.GetSpeed().X < 0) {
					c.SetSpeed(engosdl.NewVector(c.GetSpeed().X*-1, c.GetSpeed().Y))
				}
			}
		}
		return true
	})
	return box
}

func (h *GameManager) createScene() func(engine *engosdl.Engine, scene engosdl.IScene) bool {
	return func(engine *engosdl.Engine, scene engosdl.IScene) bool {
		box1 := h.createBox("box1", engosdl.NewVector(50, 50), engosdl.NewVector(10, 10))
		box2 := h.createBox("box2", engosdl.NewVector(50, 350), engosdl.NewVector(9, -9))
		box7 := h.createBox("box7", engosdl.NewVector(50, 200), engosdl.NewVector(8, -8))
		box3 := h.createBox("box3", engosdl.NewVector(750, 50), engosdl.NewVector(-7, 7))
		box4 := h.createBox("box4", engosdl.NewVector(750, 350), engosdl.NewVector(-6, 6))
		box8 := h.createBox("box8", engosdl.NewVector(750, 200), engosdl.NewVector(-5, 5))
		box5 := h.createBox("box5", engosdl.NewVector(400, 50), engosdl.NewVector(-4, 4))
		box6 := h.createBox("box6", engosdl.NewVector(400, 350), engosdl.NewVector(3, -3))
		box9 := h.createBox("box9", engosdl.NewVector(400, 200), engosdl.NewVector(2, -2))

		boxes := []engosdl.IEntity{box1, box2, box3, box4, box5, box6, box7, box8, box9}
		counter := engosdl.NewEntity("counter")
		counter.AddComponent(components.NewTimer("counter/timer", 100, -1))
		counterDisplay := engosdl.NewComponent("counter/display")
		counterDisplay.AddDelegateToRegister(nil, nil, &components.Timer{}, func(params ...interface{}) bool {
			fmt.Println()
			for _, box := range boxes {
				if c, err := box.GetCache("counter"); err == nil {
					fmt.Printf("%s : %d\n", box.GetName(), c)
				}
			}
			return true
		})
		counter.AddComponent(counterDisplay)

		scene.AddEntity(box1)
		scene.AddEntity(box2)
		scene.AddEntity(box3)
		scene.AddEntity(box4)
		scene.AddEntity(box5)
		scene.AddEntity(box6)
		scene.AddEntity(box7)
		scene.AddEntity(box8)
		scene.AddEntity(box9)
		scene.AddEntity(counter)

		scene.SetCollisionMode(engosdl.ModeBox)
		return true
	}
}

// DoFrameEnd is called at the end of every engine frame.
func (h *GameManager) DoFrameEnd() {
}

// DoFrameStart is called at the start of the game frame.
func (h *GameManager) DoFrameStart() {
}

// DoInit initializes internal game manager resources.
func (h *GameManager) DoInit() {
}
