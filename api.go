package engosdl

import (
	"fmt"
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
	componentManager = NewComponentManager("component-manager")
}

// Frames per second constants.
const (
	_fps   uint32 = 30
	_delay uint32 = 1000 / _fps
)

// Graphics format constants.
const (
	// FormatBMP identifies sprites in BMP format.
	FormatBMP int = 1
	// FormatPNG identifies sprites in PNG format.
	FormatPNG int = 2
	// FormatJPG identifies sprites in JPG format.
	FormatJPG int = 3
)

// Sound format constants.
const (
	// SoundMP3 identifies sound in MP3 format.
	SoundMP3 int = 1
	// SoundWAV identifies sound in WAV format.
	SoundWAV int = 2
)

// Movement constants.
const (
	// No Movement.
	MoveNo int = 0
	// Move up in the screen (negative Y axis).
	MoveUp int = 1
	// Move down in the screen (positive Y axis).
	MoveDown int = 2
	// Move left in the screen (negative X axis).
	MoveLeft int = 4
	// Move right in the screen (positive X axis).
	MoveRight int = 8
)

// Location constants.
const (
	// Up location.
	Up int = 1
	// Down location.
	Down int = 2
	// Left location.
	Left int = 3
	// Right location.
	Right int = 4
)

// Collision box modes.
const (
	// Box or rectangle mode. Collision box is a rectangle.
	ModeBox int = 1
	// Circle mode. Collision box is a circle.
	ModeCircle int = 2
)

// GetComponentManager gets the component manager.
func GetComponentManager() *ComponentManager {
	return componentManager
}

// GetEngine returns the singleton game engine.
func GetEngine() *Engine {
	return gameEngine
}

// GetCursorManager returns the engine cursor manager.
func GetCursorManager() ICursorManager {
	if engine := GetEngine(); engine != nil {
		return engine.GetCursorManager()
	}
	return nil
}

// GetDelegateManager returns the engine delegate manager.
func GetDelegateManager() IDelegateManager {
	if engine := GetEngine(); engine != nil {
		return engine.GetDelegateManager()
	}
	return nil
}

// GetEventManager returns the engine event manager.
func GetEventManager() IEventManager {
	if engine := GetEngine(); engine != nil {
		return engine.GetEventManager()
	}
	return nil
}

// GetFontManager returns the engine font manager.
func GetFontManager() IFontManager {
	if engine := GetEngine(); engine != nil {
		return engine.GetFontManager()
	}
	return nil
}

// GetGameManager returns the engine game manager.
func GetGameManager() IGameManager {
	if engine := GetEngine(); engine != nil {
		return engine.GetGameManager()
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

// EntitiesInCollision identifies entities being passed in a collision
// notification.
func EntitiesInCollision(entity IEntity, params ...interface{}) (IEntity, IEntity, error) {
	collisionEntityOne := params[0].(*Entity)
	collisionEntityTwo := params[1].(*Entity)
	var me, other IEntity
	var result error = nil
	if entity.GetID() == collisionEntityOne.GetID() {
		me = collisionEntityOne
		other = collisionEntityTwo
	} else if entity.GetID() == collisionEntityTwo.GetID() {
		me = collisionEntityTwo
		other = collisionEntityOne
	} else {
		me = collisionEntityOne
		other = collisionEntityTwo
		result = fmt.Errorf("entity %s not found in collision", entity.GetName())
	}
	return me, other, result
}
