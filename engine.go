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
	eventHandler IEventHandler
}

// NewEngine creates a new engine instance.
func NewEngine(name string, w, h int32) *Engine {
	Logger.Trace().Str("engine", name).Msg("new engine")
	return &Engine{
		name:         name,
		width:        w,
		height:       h,
		sceneHandler: NewSceneHandler("engine-scene-handler"),
		eventHandler: NewEventHandler("engine-event-handler"),
	}
}

// DoInitSdl initialiazes all engine sdl structures.
func (engine *Engine) DoInitSdl() {
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
}

// DoAwake awakes all engine structures.
func (engine *Engine) DoAwake() {
	engine.GetSceneHandler().OnAwake()
}

// DoStart starts the game engine.
func (engine *Engine) DoStart() {
	engine.active = true

	engine.GetSceneHandler().OnStart()
}

// DoRun runs the engine.
func (engine *Engine) DoRun() {
	Logger.Trace().Str("engine", engine.name).Msg("run engine")

	for engine.active {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				Logger.Trace().Str("engine", engine.name).Msg("exit engine")
				engine.active = false
				break
			}
		}

		engine.GetSceneHandler().OnUpdate()

		engine.renderer.SetDrawColor(255, 255, 255, 255)
		engine.renderer.Clear()

		engine.GetSceneHandler().OnDraw()

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

// GetRenderer returns the engine renderer.
func (engine *Engine) GetRenderer() *sdl.Renderer {
	return engine.renderer
}

// GetSceneHandler returns the engine scene handler.
func (engine *Engine) GetSceneHandler() ISceneHandler {
	return engine.sceneHandler
}

// GetEventHandler returns the engine event handler.
func (engine *Engine) GetEventHandler() IEventHandler {
	return engine.eventHandler
}

// AddScene adds a new scene to the engine.
func (engine *Engine) AddScene(scene IScene) bool {
	return engine.GetSceneHandler().AddScene(scene)
}
