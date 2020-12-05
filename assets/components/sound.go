package components

import (
	"reflect"

	"github.com/jrecuero/engosdl"
)

// ComponentNameSound is the name to refer sound component.
var ComponentNameSound string = reflect.TypeOf(&Sound{}).String()

func init() {
	if componentManager := engosdl.GetComponentManager(); componentManager != nil {
		componentManager.RegisterConstructor(ComponentNameSound, CreateSound)
	}
}

var _ engosdl.ISound = (*Sound)(nil)

// Sound represents a component that play a sound/music.
type Sound struct {
	*engosdl.Component
	Filename string `json:"filename"`
	Format   int    `json:"format"`
	resource engosdl.ISoundResource
}

// NewSound creates a new sound instance.
func NewSound(name string, filename string, format int) *Sound {
	engosdl.Logger.Trace().Str("component", "sound").Str("sound", name).Msg("new sound")
	return &Sound{
		Component: engosdl.NewComponent(name),
		Filename:  filename,
		Format:    format,
		resource:  nil,
	}
}

// CreateSound implements sound constructor used by component manager.
func CreateSound(params ...interface{}) engosdl.IComponent {
	if len(params) == 3 {
		return NewSound(params[0].(string), params[1].(string), params[2].(int))
	}
	return NewSound("", "", engosdl.SoundMP3)
}

// DoDestroy calls all methods to clean up sound.
func (c *Sound) DoDestroy() {
	engosdl.Logger.Trace().Str("component", "sound").Str("sound", c.GetName()).Msg("DoDestroy")
	c.resource.Delete()
	c.Component.DoDestroy()
}

// DoUnLoad is called when component is unloaded, so all resources have
// to be released.
func (c *Sound) DoUnLoad() {
	engosdl.Logger.Trace().Str("component", "sound").Str("sound", c.GetName()).Msg("DoUnLoad")
	c.Component.DoUnLoad()
}

// GetFilename returns filename used for the sound.
func (c *Sound) GetFilename() string {
	return c.Filename
}

// GetFormat returns sound format for the sound.
func (c *Sound) GetFormat() int {
	return c.Format
}

// LoadSound loads the sound from the filename.
func (c *Sound) LoadSound() {
	engosdl.Logger.Trace().Str("component", "sound").Str("sound", c.GetName()).Msg("LoadSound")
	c.resource = engosdl.GetSoundManager().CreateSound(c.GetName(), c.GetFilename(), c.GetFormat())
}

// OnAwake should create all component resources that don't have any dependency
// with any other component or entity.
func (c *Sound) OnAwake() {
	engosdl.Logger.Trace().Str("component", "sound").Str("sound", c.GetName()).Msg("OnAwake")
	c.LoadSound()
	c.Component.OnAwake()
}

// OnStart is called first time the component is enabled.
func (c *Sound) OnStart() {
	// Register to: "on-collision" and "out-of-bounds"
	engosdl.Logger.Trace().Str("component", "sound").Str("sound", c.GetName()).Msg("OnStart")
	c.Component.OnStart()
}

// Play plays the sound.
func (c *Sound) Play(times int) {
	// if err := c.resource.GetResource().Play(times); err != nil {
	if err := c.resource.GetResource().FadeIn(times, 2500); err != nil {
		engosdl.Logger.Error().Err(err).Msg("play mix resource error")
		panic(err)
	}
}

// Unmarshal takes information from a ComponentToUnmarshal instance and
//  creates a new component instance.
func (c *Sound) Unmarshal(data map[string]interface{}) {
	c.Component.Unmarshal(data)
	c.Filename = data["format"].(string)
	c.Format = int(data["format"].(float64))
}
