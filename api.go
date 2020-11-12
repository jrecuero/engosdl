package engosdl

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/veandco/go-sdl2/sdl"
)

// Logger is the system logger to be used by the application.
var Logger zerolog.Logger

// gameEngine is the game engine singlenton instance.
var gameEngine *Engine

// componentManager is the component manager in charge of tracking all
// components registered in the application.
var componentManager *ComponentManager

func init() {
	file, err := os.Create("engosdl.log")
	if err != nil {
		panic(err)
	}
	//"2006-01-02T15:04:05.999999999Z07:00"
	zerolog.TimeFieldFormat = time.RFC3339Nano
	Logger = zerolog.New(file).With().Timestamp().Logger()
	Logger.Info().Msg("start engosdl")
	Logger.Info().Msg("create component manager")
	componentManager = &ComponentManager{
		Object:     NewObject("component-manager"),
		components: []IComponent{},
	}
}

const (
	_fps   uint32 = 30
	_delay uint32 = 1000 / _fps
)

// GetComponentManager gets the component manager.
func GetComponentManager() *ComponentManager {
	return componentManager
}

// GetEngine returns the singleton game engine.
func GetEngine() *Engine {
	return gameEngine
}

// GetDelegateHandler returns the delegate handler.
func GetDelegateHandler() IDelegateHandler {
	if engine := GetEngine(); engine != nil {
		return engine.GetDelegateHandler()
	}
	return nil
}

// GetEventHandler returns the event handler.
func GetEventHandler() IEventHandler {
	if engine := GetEngine(); engine != nil {
		return engine.GetEventHandler()
	}
	return nil
}

// GetFontHandler returns the font handler.
func GetFontHandler() IFontHandler {
	if engine := GetEngine(); engine != nil {
		return engine.GetFontHandler()
	}
	return nil
}

// GetRenderer returns the engine renderer.
func GetRenderer() *sdl.Renderer {
	if engine := GetEngine(); engine != nil {
		return engine.GetRenderer()
	}
	return nil
}

// GetResourceHandler returns the engine resource handler.
func GetResourceHandler() IResourceHandler {
	if engine := GetEngine(); engine != nil {
		return engine.GetResourceHandler()
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

// GetSoundHandler returns the engine sound handler.
func GetSoundHandler() ISoundHandler {
	if engine := GetEngine(); engine != nil {
		return engine.GetSoundHandler()
	}
	return nil
}
