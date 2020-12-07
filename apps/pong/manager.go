package main

import (
	"strconv"

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
		speed := 4.0

		net := engosdl.NewEntity("net")
		net.GetTransform().SetPositionXY(395, 0)
		net.AddComponent(components.NewBox("net/box", &sdl.Rect{W: 10, H: 400}, sdl.Color{B: 255}, true))

		score1 := engosdl.NewEntity("score1")
		score1.GetTransform().SetPositionXY(50, 10)
		text1 := components.NewText("score1/text", "fonts/fira.ttf", 32, sdl.Color{R: 255}, strconv.Itoa(h.score1))
		score1.AddComponent(text1)

		score2 := engosdl.NewEntity("score2")
		score2.GetTransform().SetPositionXY(750, 10)
		text2 := components.NewText("score2/text", "fonts/fira.ttf", 32, sdl.Color{R: 255}, strconv.Itoa(h.score1))
		score2.AddComponent(text2)

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
		ball.SetCache("last", "")
		ball.GetTransform().SetPositionXY(100, 180)
		ballPong := engosdl.NewEntity("ball-pong")
		ballPong.AddComponent(components.NewSound("ball-pong/sound", "sounds/pong.mp3", engosdl.SoundMP3))
		ballOut := engosdl.NewEntity("ball-out")
		ballOut.AddComponent(components.NewSound("ball-out/sound", "sounds/out.mp3", engosdl.SoundMP3))

		box3 := components.NewBox("ball/box", &sdl.Rect{W: 10, H: 10}, sdl.Color{}, true)
		body3 := components.NewBody("ball/body", false)
		moveTo3 := components.NewMoveTo("ball/move-it", engosdl.NewVector(speed, speed))
		moveTo3.AddDelegateToRegister(nil, nil, &components.Body{}, func(params ...interface{}) bool {
			c := moveTo3
			outAt := params[2].(int)
			switch outAt {
			case engosdl.Left:
				// c.SetSpeed(engosdl.NewVector(c.GetSpeed().X*-1, c.GetSpeed().Y))
				c.GetEntity().GetTransform().SetPositionXY(100, 180)
				c.SetSpeed(engosdl.NewVector(speed, speed))
				h.score2++
				text2.SetMessage(strconv.Itoa(h.score2))
				ball.SetCache("last", "")
				if child := c.GetEntity().GetChildByName("ball-out"); child != nil {
					if sound := child.GetComponent(&components.Sound{}); sound != nil {
						sound.(*components.Sound).Play(1)
					}
				}
				break
			case engosdl.Right:
				// c.SetSpeed(engosdl.NewVector(c.GetSpeed().X*-1, c.GetSpeed().Y))
				c.GetEntity().GetTransform().SetPositionXY(700, 180)
				c.SetSpeed(engosdl.NewVector(speed*-1, speed))
				h.score1++
				text1.SetMessage(strconv.Itoa(h.score1))
				ball.SetCache("last", "")
				if child := c.GetEntity().GetChildByName("ball-out"); child != nil {
					if sound := child.GetComponent(&components.Sound{}); sound != nil {
						sound.(*components.Sound).Play(1)
					}
				}
				break
			case engosdl.Up:
				c.SetSpeed(engosdl.NewVector(c.GetSpeed().X, c.GetSpeed().Y*-1))
				break
			case engosdl.Down:
				c.SetSpeed(engosdl.NewVector(c.GetSpeed().X, c.GetSpeed().Y*-1))
				break
			}
			return true
		})
		moveTo3.AddDelegateToRegister(engosdl.GetDelegateManager().GetCollisionDelegate(), nil, nil, func(params ...interface{}) bool {
			c := moveTo3
			if _, other, err := engosdl.EntitiesInCollision(c.GetEntity(), params...); err == nil && other.GetTag() == "player" {
				if last, err := c.GetEntity().GetCache("last"); err == nil && last != other.GetID() {
					if child := c.GetEntity().GetChildByName("ball-pong"); child != nil {
						if sound := child.GetComponent(&components.Sound{}); sound != nil {
							sound.(*components.Sound).Play(1)
						}
					}
					c.SetSpeed(engosdl.NewVector(c.GetSpeed().X*-1, c.GetSpeed().Y))
					c.GetEntity().SetCache("last", other.GetID())
				}
			}
			return true
		})
		ball.AddChild(ballPong)
		ball.AddChild(ballOut)
		ball.AddComponent(box3)
		ball.AddComponent(body3)
		ball.AddComponent(moveTo3)

		scene.AddEntity(net)
		scene.AddEntity(player1)
		scene.AddEntity(player2)
		scene.AddEntity(ball)
		scene.AddEntity(score1)
		scene.AddEntity(score2)
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
