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

// GetDelegateManager returns the delegate handler.
func GetDelegateManager() IDelegateManager {
	if engine := GetEngine(); engine != nil {
		return engine.GetDelegateManager()
	}
	return nil
}

// GetEventManager returns the event handler.
func GetEventManager() IEventManager {
	if engine := GetEngine(); engine != nil {
		return engine.GetEventManager()
	}
	return nil
}

// GetFontManager returns the font handler.
func GetFontManager() IFontManager {
	if engine := GetEngine(); engine != nil {
		return engine.GetFontManager()
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

// GetResourceManager returns the engine resource handler.
func GetResourceManager() IResourceManager {
	if engine := GetEngine(); engine != nil {
		return engine.GetResourceManager()
	}
	return nil
}

// GetSceneManager returns the engine scene handler.
func GetSceneManager() ISceneManager {
	if engine := GetEngine(); engine != nil {
		return engine.GetSceneManager()
	}
	return nil
}

// GetSoundManager returns the engine sound handler.
func GetSoundManager() ISoundManager {
	if engine := GetEngine(); engine != nil {
		return engine.GetSoundManager()
	}
	return nil
}
