package engosdl

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/veandco/go-sdl2/sdl"
)

var logger zerolog.Logger

func init() {
	file, err := os.Create("engosdl.log")
	if err != nil {
		panic(err)
	}
	//"2006-01-02T15:04:05.999999999Z07:00"
	zerolog.TimeFieldFormat = time.RFC3339Nano
	logger = zerolog.New(file).With().Timestamp().Logger()
}

// Engine represents the main game engine in charge of running the game.
type Engine struct {
	name     string
	width    int32
	height   int32
	active   bool
	window   *sdl.Window
	renderer *sdl.Renderer
}

// NewEngine creates a new engine instance.
func NewEngine(name string, w, h int32) *Engine {
	logger.Trace().Str("engine", name).Msg("new engine")
	return &Engine{
		name:   name,
		width:  w,
		height: h,
	}
}

// DoStart starts the game engine.
func (e *Engine) DoStart() {
	var err error
	logger.Trace().Str("engine", e.name).Msg("start engine")
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		logger.Error().Err(err)
		panic(err)
	}

	e.window, err = sdl.CreateWindow(e.name,
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		e.width, e.height,
		sdl.WINDOW_SHOWN)
	if err != nil {
		logger.Error().Err(err)
		panic(err)
	}

	e.renderer, err = sdl.CreateRenderer(e.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		logger.Error().Err(err)
		panic(err)
	}

	e.active = true
}

// DoRun runs the engine.
func (e *Engine) DoRun() {
	logger.Trace().Str("engine", e.name).Msg("run engine")
	for e.active {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				logger.Trace().Str("engine", e.name).Msg("exit engine")
				e.active = false
				break
			}
		}
		e.renderer.SetDrawColor(255, 255, 255, 255)
		e.renderer.Clear()
		e.renderer.Present()
		sdl.Delay(30)
	}
}

//DoCleanup clean-ups all graphical resources created by teh engine.
func (e *Engine) DoCleanup() {
	logger.Trace().Str("engine", e.name).Msg("end engine")
	defer sdl.Quit()
	defer e.window.Destroy()
	defer e.renderer.Destroy()
}
