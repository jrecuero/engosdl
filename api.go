package engosdl

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/veandco/go-sdl2/sdl"
)

// Logger is the system logger to be used by the application.
var Logger zerolog.Logger
var gameEngine *Engine

func init() {
	file, err := os.Create("engosdl.log")
	if err != nil {
		panic(err)
	}
	//"2006-01-02T15:04:05.999999999Z07:00"
	zerolog.TimeFieldFormat = time.RFC3339Nano
	Logger = zerolog.New(file).With().Timestamp().Logger()
	Logger.Info().Msg("start engosdl")
}

const (
	_fps   uint32 = 30
	_delay uint32 = 1000 / _fps
)

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
