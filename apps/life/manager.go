package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jrecuero/engosdl"
	"github.com/veandco/go-sdl2/sdl"
)

// GameManager is the application game manager.
type GameManager struct {
	*engosdl.GameManager
}

var _ engosdl.IGameManager = (*GameManager)(nil)

// Pixel represents any child from an alive entity.
type Pixel struct {
	*engosdl.Entity
}

// NewPixel creates a new pixel instance.
func NewPixel(name string) *Pixel {
	return &Pixel{
		Entity: engosdl.NewEntity(name),
	}
}

// NewGameManager created a new game manager instance.
func NewGameManager(name string) *GameManager {
	engosdl.Logger.Trace().Str("game-manager", name).Msg("new game-manager")
	rand.Seed(time.Now().UTC().UnixNano())
	return &GameManager{
		GameManager: engosdl.NewGameManager(name),
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

func (h *GameManager) createAlive(name string, row int, col int, color sdl.Color) engosdl.IEntity {
	alive := engosdl.NewEntity(name)
	alive.AddComponent(NewAlive(fmt.Sprintf("%s/alive", name), 10, row, col, 10, color))
	return alive
}

func (h *GameManager) createScenePlay() func(engine *engosdl.Engine, scene engosdl.IScene) bool {
	return func(engine *engosdl.Engine, scene engosdl.IScene) bool {
		controller := engosdl.NewEntity("controller")
		board := NewBoard("board", 80, 80, engosdl.NewVector(0, 0), 10)

		controller.AddComponent(board)
		alive1 := h.createAlive("alive-1", 30, 40, sdl.Color{R: 255})
		alive2 := h.createAlive("alive-1", 30, 50, sdl.Color{B: 255})
		alive3 := h.createAlive("alive-1", 40, 40, sdl.Color{G: 255})
		alive4 := h.createAlive("alive-1", 40, 50, sdl.Color{})

		scene.AddEntity(controller)
		scene.AddEntity(alive1)
		scene.AddEntity(alive2)
		scene.AddEntity(alive3)
		scene.AddEntity(alive4)
		scene.SetCollisionCheck(false)
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
