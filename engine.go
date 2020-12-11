package engosdl

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
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
	debugServer     bool
}

// NewEngine creates a new engine instance.
func NewEngine(name string, w, h int32, gameManager IGameManager) *Engine {
	Logger.Trace().Str("engine", name).Msg("new engine")
	if GetEngine() == nil {
		if gameManager == nil {
			gameManager = NewGameManager("engine-game-manager")
		}
		gameEngine = &Engine{
			name:            name,
			width:           w,
			height:          h,
			delegateManager: NewDelegateManager("engine-delegate-manager"),
			eventManager:    NewEventManager("engine-event-manager"),
			fontManager:     NewFontManager("engine-font-manager"),
			resourceManager: NewResourceManager("engine-resource-manager"),
			sceneManager:    NewSceneManager("engine-scene-manager"),
			soundManager:    NewSoundManager("engine-sound-manager"),
			gameManager:     gameManager,
			debugServer:     false,
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
	defer ttf.Quit()
	defer img.Quit()
	defer mix.CloseAudio()
	defer mix.Quit()
	defer sdl.Quit()
	defer engine.window.Destroy()
	defer engine.renderer.Destroy()
}

// DoFrameEnd calls all methods to run at the end of a tick frame.
func (engine *Engine) DoFrameEnd() {
	engine.GetGameManager().DoFrameEnd()
	engine.GetSceneManager().DoFrameEnd()
}

// DoFrameStart calls all methods to run at the start of a tick frame.
func (engine *Engine) DoFrameStart() {
	engine.GetGameManager().DoFrameStart()
	engine.GetSceneManager().DoFrameStart()
}

// DoInit initializes basic engine resources.
func (engine *Engine) DoInit() {
	Logger.Trace().Str("engine", engine.name).Msg("DoInit")
	engine.DoInitSdl()
	engine.DoInitResources()
	engine.GetGameManager().DoInit()
}

// DoInitDebugServer initializes the debug server.
func (engine *Engine) DoInitDebugServer() {
	if !engine.debugServer {
		engine.debugServer = true
		fmt.Printf("init debug server %s\n", engine.name)
		go func() {
			router := mux.NewRouter().StrictSlash(true)
			router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "debug server %s\n", engine.name)
				fmt.Printf("debug server %s home page\n", engine.name)
			})
			router.HandleFunc("/scene", func(w http.ResponseWriter, r *http.Request) {
				scene := engine.GetSceneManager().GetActiveScene()
				fmt.Fprintf(w, "active scene: %s\n", scene.GetName())
			})
			router.HandleFunc("/entities", func(w http.ResponseWriter, r *http.Request) {
				scene := engine.GetSceneManager().GetActiveScene()
				for _, entity := range scene.GetEntities() {
					fmt.Fprintf(w, "entity %s: %t\n", entity.GetName(), entity.GetActive())
				}
			})
			Logger.Error().Err(http.ListenAndServe("localhost:6060", router))
		}()
	}
}

// DoInitResources initializes all internal resources, like scene handler and
// event handler.
func (engine *Engine) DoInitResources() {
	Logger.Trace().Str("engine", engine.name).Msg("init resources")
	engine.GetEventManager().DoInit()
	engine.GetDelegateManager().DoInit()
	engine.GetResourceManager().DoInit()
	engine.GetFontManager().DoInit()
	engine.GetSoundManager().DoInit()
	engine.GetSceneManager().DoInit()
	engine.GetGameManager().DoInit()
}

// DoInitSdl initializes all engine sdl structures.
func (engine *Engine) DoInitSdl() {
	var err error

	Logger.Trace().Str("engine", engine.name).Msg("init sdl module")
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		Logger.Error().Err(err).Msg("sdl.Init error")
		panic(err)
	}

	Logger.Trace().Str("engine", engine.name).Msg("init ttf module")
	if err = ttf.Init(); err != nil {
		Logger.Error().Err(err).Msg("ttf.Init error")
		panic(err)
	}

	Logger.Trace().Str("engine", engine.name).Msg("init img module")
	if err = img.Init(img.INIT_PNG); err != nil {
		Logger.Error().Err(err).Msg("img.Init error")
		panic(err)
	}

	Logger.Trace().Str("engine", engine.name).Msg("open audio module")
	if err = mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		Logger.Error().Err(err).Msg("mix.OpenAudio error")
		panic(err)
	}

	Logger.Trace().Str("engine", engine.name).Msg("init mix module")
	if err = mix.Init(mix.INIT_MP3); err != nil {
		Logger.Error().Err(err).Msg("mix.Init error")
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

		engine.DoUpdate()

		engine.renderer.SetDrawColor(255, 255, 255, 255)
		engine.renderer.Clear()

		engine.DoRender()

		engine.renderer.Present()

		// Execute everything required at the end of the tick frame.
		engine.DoFrameEnd()

		frameTime := sdl.GetTicks() - frameStart

		if frameTime < _delay {
			// fmt.Println(_delay - frameTime)
			sdl.Delay(_delay - frameTime)
		}
	}
}

// DoRender calls on OnRender methods to run.
func (engine *Engine) DoRender() {
	// Call game manager render.
	engine.GetGameManager().OnRender()
	// Execute all render calls.
	engine.GetSceneManager().OnRender()
}

// DoStart starts the game engine. At this point all scenes and entities have
// been already added to the engine.
func (engine *Engine) DoStart(scene IScene) {
	Logger.Trace().Str("engine", engine.name).Msg("DoStart")
	engine.active = true
	engine.GetEventManager().OnStart()
	engine.GetDelegateManager().OnStart()
	engine.GetResourceManager().OnStart()
	engine.GetFontManager().OnStart()
	engine.GetSoundManager().OnStart()
	engine.GetSceneManager().OnStart()
	engine.GetGameManager().OnStart()

	// Set first scene as the active by default.
	if scene != nil {
		engine.GetSceneManager().SetActiveScene(scene)
	} else {
		engine.GetSceneManager().SetActiveFirstScene()
	}
}

// DoUpdate calls all OnUpdate and OnAfterUpdate methods to run.
func (engine *Engine) DoUpdate() {
	// Call game manager update.
	engine.GetGameManager().OnUpdate()
	// Execute all update calls.
	engine.GetSceneManager().OnUpdate()
	// Call update for delegate handler.
	engine.GetDelegateManager().OnUpdate()
	// Execute any post updates behavior.
	engine.GetSceneManager().OnAfterUpdate()
	// Call game manager after update.
	engine.GetGameManager().OnAfterUpdate()
}

// GetDelegateManager returns the engine delegate handler.
func (engine *Engine) GetDelegateManager() IDelegateManager {
	return engine.delegateManager
}

// GetEventManager returns the engine event handler.
func (engine *Engine) GetEventManager() IEventManager {
	return engine.eventManager
}

// GetFontManager returns the engine font handler.
func (engine *Engine) GetFontManager() IFontManager {
	return engine.fontManager
}

// GetGameManager returns the engine game manager.
func (engine *Engine) GetGameManager() IGameManager {
	return engine.gameManager
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
	engine.DoInit()
	engine.GetGameManager().CreateAssets()
	engine.DoStart(scene)
	engine.DoInitDebugServer()
	engine.DoRun()
	engine.DoCleanup()
	return true
}
