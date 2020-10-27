package engosdl

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/veandco/go-sdl2/sdl"
)

// Logger is the system logger to be used by the application.
var Logger zerolog.Logger

func init() {
	file, err := os.Create("engosdl.log")
	if err != nil {
		panic(err)
	}
	//"2006-01-02T15:04:05.999999999Z07:00"
	zerolog.TimeFieldFormat = time.RFC3339Nano
	Logger = zerolog.New(file).With().Timestamp().Logger()
}

// Engine represents the main game engine in charge of running the game.
type Engine struct {
	name         string
	width        int32
	height       int32
	active       bool
	window       *sdl.Window
	renderer     *sdl.Renderer
	sceneHandler ISceneHandler
}

// NewEngine creates a new engine instance.
func NewEngine(name string, w, h int32) *Engine {
	Logger.Trace().Str("engine", name).Msg("new engine")
	return &Engine{
		name:         name,
		width:        w,
		height:       h,
		sceneHandler: NewSceneHandler("engine-scene-handler"),
	}
}

// DoInit initialiazes all engine required structures.
func (engine *Engine) DoInit() {
}

// DoStart starts the game engine.
func (engine *Engine) DoStart() {
	var err error
	Logger.Trace().Str("engine", engine.name).Msg("start engine")
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		Logger.Error().Err(err)
		panic(err)
	}

	engine.window, err = sdl.CreateWindow(engine.name,
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		engine.width, engine.height,
		sdl.WINDOW_SHOWN)
	if err != nil {
		Logger.Error().Err(err)
		panic(err)
	}

	engine.renderer, err = sdl.CreateRenderer(engine.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		Logger.Error().Err(err)
		panic(err)
	}

	engine.active = true

	for _, scene := range engine.sceneHandler.Scenes() {
		scene.OnAwake()
	}
}

// DoRun runs the engine.
func (engine *Engine) DoRun() {
	Logger.Trace().Str("engine", engine.name).Msg("run engine")
	for _, scene := range engine.sceneHandler.Scenes() {
		scene.OnStart()
	}
	for engine.active {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				Logger.Trace().Str("engine", engine.name).Msg("exit engine")
				engine.active = false
				break
			}
		}
		engine.renderer.SetDrawColor(255, 255, 255, 255)
		engine.renderer.Clear()
		engine.renderer.Present()
		sdl.Delay(30)
	}
}

// DoCleanup clean-ups all graphical resources created by teh engine.
func (engine *Engine) DoCleanup() {
	Logger.Trace().Str("engine", engine.name).Msg("end engine")
	defer sdl.Quit()
	defer engine.window.Destroy()
	defer engine.renderer.Destroy()
}

// AddScene adds a new scene to the engine.
func (engine *Engine) AddScene(scene IScene) bool {
	return engine.sceneHandler.AddScene(scene)
}
