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
	score1 int
	score2 int
}

var _ engosdl.IGameManager = (*GameManager)(nil)

// NewGameManager created a new game manager instance.
func NewGameManager(name string) *GameManager {
	engosdl.Logger.Trace().Str("game-manager", name).Msg("new game-manager")
	return &GameManager{
		GameManager: engosdl.NewGameManager(name),
		score1:      0,
		score2:      0,
	}
}

// CreateAssets creates all application assets and resources. it is called
// before game engine starts in order to create all required assets and
// resources.
func (h *GameManager) CreateAssets() {
	playScene := engosdl.NewScene("play", "play")
	playScene.SetSceneCode(h.createScenePlay())
	engosdl.GetEngine().AddScene(playScene)
}

func (h *GameManager) createScenePlay() func(engine *engosdl.Engine, scene engosdl.IScene) bool {
	return func(engine *engosdl.Engine, scene engosdl.IScene) bool {
		net := engosdl.NewEntity("net")
		net.GetTransform().SetPositionXY(395, 0)
		net.AddComponent(components.NewBox("net/box", &sdl.Rect{W: 10, H: 400}, sdl.Color{B: 255}, true))

		player1 := engosdl.NewEntity("player1")
		player1.SetTag("player")
		player1.GetTransform().SetPositionXY(750, 180)
		box1 := components.NewBox("player1/box", &sdl.Rect{W: 10, H: 40}, sdl.Color{B: 255}, true)
		body1 := components.NewBody("player1/body", false)
		keyboard1 := components.NewKeyboard("player1/keyboard", map[int]bool{sdl.SCANCODE_LEFT: true, sdl.SCANCODE_RIGHT: true})
		moveIt1 := components.NewMoveIt("player1/move-it", engosdl.NewVector(0, 5))
		moveIt1.AddDelegateToRegister(nil, nil, &components.Keyboard{}, func(params ...interface{}) bool {
			c := moveIt1
			position := c.GetEntity().GetTransform().GetPosition()
			key := params[0].(int)
			switch key {
			case sdl.SCANCODE_LEFT:
				position.Y -= c.Speed.Y
				c.LastMove.Y = -1 * c.Speed.Y
				c.LastMove.X = 0
				break
			case sdl.SCANCODE_RIGHT:
				position.Y += c.Speed.Y
				c.LastMove.Y = c.Speed.Y
				c.LastMove.X = 0
				break
			}
			return true
		})
		player1.AddComponent(box1)
		player1.AddComponent(body1)
		player1.AddComponent(keyboard1)
		player1.AddComponent(moveIt1)

		player2 := engosdl.NewEntity("player2")
		player2.SetTag("player")
		player2.GetTransform().SetPositionXY(50, 180)
		box2 := components.NewBox("player2/box", &sdl.Rect{W: 10, H: 40}, sdl.Color{B: 255}, true)
		body2 := components.NewBody("player2/body", false)
		keyboard2 := components.NewKeyboard("player2/keyboard", map[int]bool{sdl.SCANCODE_A: true, sdl.SCANCODE_S: true})
		moveIt2 := components.NewMoveIt("player2/move-it", engosdl.NewVector(0, 5))
		moveIt2.AddDelegateToRegister(nil, nil, &components.Keyboard{}, func(params ...interface{}) bool {
			c := moveIt2
			position := c.GetEntity().GetTransform().GetPosition()
			key := params[0].(int)
			switch key {
			case sdl.SCANCODE_A:
				position.Y -= c.Speed.Y
				c.LastMove.Y = -1 * c.Speed.Y
				c.LastMove.X = 0
				break
			case sdl.SCANCODE_S:
				position.Y += c.Speed.Y
				c.LastMove.Y = c.Speed.Y
				c.LastMove.X = 0
				break
			}
			return true
		})
		player2.AddComponent(box2)
		player2.AddComponent(body2)
		player2.AddComponent(keyboard2)
		player2.AddComponent(moveIt2)

		ball := engosdl.NewEntity("ball")
		ball.GetTransform().SetPositionXY(100, 180)
		box3 := components.NewBox("ball/box", &sdl.Rect{W: 10, H: 10}, sdl.Color{}, true)
		body3 := components.NewBody("ball/body", true)
		moveTo3 := components.NewMoveTo("ball/move-it", engosdl.NewVector(5, 0))
		moveTo3.AddDelegateToRegister(nil, nil, &components.Body{}, func(params ...interface{}) bool {
			c := moveTo3
			entity := params[0].(engosdl.IEntity)
			if outAt := params[2].(int); outAt == engosdl.Left || outAt == engosdl.Right {
				if entity.GetID() == c.GetEntity().GetID() {
					if outAt == engosdl.Left {
						h.score2++
					} else {
						h.score1++
					}
					fmt.Printf("p1: %d - p2: %d\n", h.score1, h.score2)
					engosdl.GetEngine().DestroyEntity(c.GetEntity())
				}
			}
			return true
		})
		moveTo3.AddDelegateToRegister(engosdl.GetDelegateManager().GetCollisionDelegate(), nil, nil, func(params ...interface{}) bool {
			c := moveTo3
			if _, other, err := engosdl.EntitiesInCollision(c.GetEntity(), params...); err == nil && other.GetTag() == "player" {
				c.SetSpeed(engosdl.NewVector(c.GetSpeed().X*-1, c.GetSpeed().Y*-1))
			}
			return true
		})
		ball.AddComponent(box3)
		ball.AddComponent(body3)
		ball.AddComponent(moveTo3)

		scene.AddEntity(net)
		scene.AddEntity(player1)
		scene.AddEntity(player2)
		scene.AddEntity(ball)
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
