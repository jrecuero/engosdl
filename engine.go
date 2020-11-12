package engosdl

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Engine represents the main game engine in charge of running the game.
type Engine struct {
	name            string
	width           int32
	height          int32
	active          bool
	window          *sdl.Window
	renderer        *sdl.Renderer
	delegateManager IDelegateManager
	eventManager    IEventManager
	fontManager     IFontManager
	resourceManager IResourceManager
	sceneManager    ISceneManager
	soundManager    ISoundManager
	gameManager     IGameManager
}

// NewEngine creates a new engine instance.
func NewEngine(name string, w, h int32) *Engine {
	Logger.Trace().Str("engine", name).Msg("new engine")
	if GetEngine() == nil {
		gameEngine = &Engine{
			name:            name,
			width:           w,
			height:          h,
			delegateManager: NewDelegateManager("engine-delegate-handler"),
			eventManager:    NewEventManager("engine-event-handler"),
			fontManager:     NewFontManager("engine-font-handler"),
			resourceManager: NewResourceManager("engine-resource-handler"),
			sceneManager:    NewSceneManager("engine-scene-handler"),
			soundManager:    NewSoundManager("engine-sound-handler"),
		}
	}
	return gameEngine
}

// AddScene adds a new scene to the engine.
func (engine *Engine) AddScene(scene IScene) bool {
	Logger.Trace().Str("engine", engine.name).Msg("AddScene")
	return engine.GetSceneManager().AddScene(scene)
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

// DoFrameEnd calls all methods to run at the end of a tick frame.
func (engine *Engine) DoFrameEnd() {
	engine.GetSceneManager().DoFrameEnd()
}

// DoFrameStart calls all methods to run at the start of a tick frame.
func (engine *Engine) DoFrameStart() {
	engine.GetSceneManager().DoFrameStart()
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
	engine.GetEventManager().OnStart()
	engine.GetDelegateManager().OnStart()
	engine.GetResourceManager().OnStart()
	engine.GetFontManager().OnStart()
	engine.GetSoundManager().OnStart()
	engine.GetSceneManager().OnStart()
}

// DoInitSdl initializes all engine sdl structures.
func (engine *Engine) DoInitSdl() {
	var err error

	ttf.Init()

	Logger.Trace().Str("engine", engine.name).Msg("init sdl module")
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		Logger.Error().Err(err).Msg("sdl.Init error")
		panic(err)
	}

	engine.window, err = sdl.CreateWindow(engine.name,
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		engine.width, engine.height,
		sdl.WINDOW_SHOWN)
	if err != nil {
		Logger.Error().Err(err).Msg("CreateWindow error")
		panic(err)
	}

	engine.renderer, err = sdl.CreateRenderer(engine.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		Logger.Error().Err(err).Msg("CreateRenderer error")
		panic(err)
	}
}

// DoRun runs the engine.
func (engine *Engine) DoRun() {
	Logger.Trace().Str("engine", engine.name).Msg("run engine")

	for engine.active {

		frameStart := sdl.GetTicks()

		// Execute everything required at the start of a tick frame.
		engine.DoFrameStart()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				Logger.Trace().Str("engine", engine.name).Msg("exit engine")
				engine.active = false
				break
			}
		}

		// Execute all update calls.
		engine.GetSceneManager().OnUpdate()
		// Call update for delegate handler.
		engine.GetDelegateManager().OnUpdate()
		// Execute any post updates behavior.
		engine.GetSceneManager().OnAfterUpdate()

		engine.renderer.SetDrawColor(255, 255, 255, 255)
		engine.renderer.Clear()

		// Execute all render calls.
		engine.GetSceneManager().OnRender()

		engine.renderer.Present()

		// Execute everything required at the end of the tick frame.
		engine.DoFrameEnd()

		frameTime := sdl.GetTicks() - frameStart

		if frameTime < _delay {
			sdl.Delay(_delay - frameTime)
		}
	}
}

// DoStart starts the game engine. At this point all scenes and entities have
// been already added to the engine.
func (engine *Engine) DoStart(scene IScene) {
	Logger.Trace().Str("engine", engine.name).Msg("DoStart")
	engine.active = true

	// Set first scene as the active by default.
	if scene != nil {
		engine.GetSceneManager().SetActiveScene(scene)
	} else {
		engine.GetSceneManager().SetActiveFirstScene()
	}
}

// GetDelegateManager returns the engine delegate handler.
func (engine *Engine) GetDelegateManager() IDelegateManager {
	return engine.delegateManager
}

// GetFontManager returns the engine font handler.
func (engine *Engine) GetFontManager() IFontManager {
	return engine.fontManager
}

// GetEventManager returns the engine event handler.
func (engine *Engine) GetEventManager() IEventManager {
	return engine.eventManager
}

// GetHeight returns engine window height
func (engine *Engine) GetHeight() int32 {
	return engine.height
}

// GetRenderer returns the engine renderer.
func (engine *Engine) GetRenderer() *sdl.Renderer {
	return engine.renderer
}

// GetResourceManager returns the engine resource handler.
func (engine *Engine) GetResourceManager() IResourceManager {
	return engine.resourceManager
}

// GetSceneManager returns the engine scene handler.
func (engine *Engine) GetSceneManager() ISceneManager {
	return engine.sceneManager
}

// GetSoundManager returns the engine sound handler.
func (engine *Engine) GetSoundManager() ISoundManager {
	return engine.soundManager
}

// GetWidth returns engine window width.
func (engine *Engine) GetWidth() int32 {
	return engine.width
}

// RunEngine runs the game engine.
func (engine *Engine) RunEngine(scene IScene) bool {
	engine.DoStart(scene)
	engine.DoRun()
	engine.DoCleanup()
	return true
}
