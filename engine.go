package engosdl

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
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

var gameEngine *Engine

// GetEngine returns the singleton game engine.
func GetEngine() *Engine {
	return gameEngine
}

// GetEventHandler returns the event handler.
func GetEventHandler() IEventHandler {
	if engine := GetEngine(); engine != nil {
		return engine.GetEventHandler()
	}
	return nil
}

// GetRenderer returns the engine renderr.
func GetRenderer() *sdl.Renderer {
	if engine := GetEngine(); engine != nil {
		return engine.GetRenderer()
	}
	return nil
}

// GetSceneHandler returns the engine scene handler.
func GetSceneHandler() ISceneHandler {
	if engine := GetEngine(); engine != nil {
		return engine.GetSceneHandler()
	}
	return nil
}

// NewEngine creates a new engine instance.
func NewEngine(name string, w, h int32) *Engine {
	Logger.Trace().Str("engine", name).Msg("new engine")
	if GetEngine() == nil {
		gameEngine = &Engine{
			name:         name,
			width:        w,
			height:       h,
			sceneHandler: NewSceneHandler("engine-scene-handler"),
			eventHandler: NewEventHandler("engine-event-handler"),
		}
	}
	return gameEngine
}

// AddScene adds a new scene to the engine.
func (engine *Engine) AddScene(scene IScene) bool {
	Logger.Trace().Str("engine", engine.name).Msg("AddScene")
	return engine.GetSceneHandler().AddScene(scene)
}

// DestroyEntity removes the given entity from the game.
func (engine *Engine) DestroyEntity(entity IEntity) bool {
	scene := entity.GetScene()
	scene.DeleteEntity(entity)
	entity.SetActive(false)
	Logger.Trace().Str("engine", engine.name).Str("scene", scene.GetName()).Str("entity", entity.GetName()).Msg("destroy entity")
	return true
}

// DoCleanup clean-ups all graphical resources created by teh engine.
func (engine *Engine) DoCleanup() {
	Logger.Trace().Str("engine", engine.name).Msg("end engine")
	defer sdl.Quit()
	defer engine.window.Destroy()
	defer engine.renderer.Destroy()
}

// DoCycleEnd calls all methods to run at the end of a tick cycle.
func (engine *Engine) DoCycleEnd() {
	engine.GetSceneHandler().DoCycleEnd()
}

// DoCycleStart calls all methods to run at the start of a tick cycle.
func (engine *Engine) DoCycleStart() {
	engine.GetSceneHandler().DoCycleStart()
}

// DoInit initializes basic engine resources.
func (engine *Engine) DoInit() {
	Logger.Trace().Str("engine", engine.name).Msg("DoInit")
	gameEngine.DoInitSdl()
	gameEngine.DoInitResources()
}

// DoInitResources initializes all internal resources, like scene handler and
// event handler.
func (engine *Engine) DoInitResources() {
	Logger.Trace().Str("engine", engine.name).Msg("init resources")
	engine.GetEventHandler().OnStart()
	engine.GetSceneHandler().OnStart()
}

// DoInitSdl initializes all engine sdl structures.
func (engine *Engine) DoInitSdl() {
	var err error

	ttf.Init()

	Logger.Trace().Str("engine", engine.name).Msg("init sdl module")
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

// DoRun runs the engine.
func (engine *Engine) DoRun() {
	Logger.Trace().Str("engine", engine.name).Msg("run engine")

	for engine.active {

		// Execute everything required at the start of a tick cycle.
		engine.DoCycleStart()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				Logger.Trace().Str("engine", engine.name).Msg("exit engine")
				engine.active = false
				break
			}
		}

		// Execute all update calls.
		engine.GetSceneHandler().OnUpdate()
		// Call update for delegate handler.
		engine.GetEventHandler().GetDelegateHandler().OnUpdate()
		// Execute any post updates behavior.
		engine.GetSceneHandler().OnAfterUpdate()

		engine.renderer.SetDrawColor(255, 255, 255, 255)
		engine.renderer.Clear()

		// Execute all draw calls.
		engine.GetSceneHandler().OnDraw()

		engine.renderer.Present()

		// Execute everything required at the end of the tick cycle.
		engine.DoCycleEnd()
		sdl.Delay(30)
	}
}

// DoStart starts the game engine. At this point all scenes and entities have
// been already added to the engine.
func (engine *Engine) DoStart() {
	Logger.Trace().Str("engine", engine.name).Msg("DoStart")
	engine.active = true

	// Set first scene as the active by default.
	engine.GetSceneHandler().SetActiveFirstScene()
}

// GetEventHandler returns the engine event handler.
func (engine *Engine) GetEventHandler() IEventHandler {
	return engine.eventHandler
}

// GetHeight returns engine window height
func (engine *Engine) GetHeight() int32 {
	return engine.height
}

// GetRenderer returns the engine renderer.
func (engine *Engine) GetRenderer() *sdl.Renderer {
	return engine.renderer
}

// GetSceneHandler returns the engine scene handler.
func (engine *Engine) GetSceneHandler() ISceneHandler {
	return engine.sceneHandler
}

// GetWidth returns engine window width.
func (engine *Engine) GetWidth() int32 {
	return engine.width
}

// RunEngine runs the game engine.
func (engine *Engine) RunEngine() bool {
	engine.DoStart()
	engine.DoRun()
	engine.DoCleanup()
	return true
}
