package engosdl

import (
	"io/ioutil"

	"github.com/veandco/go-sdl2/mix"
)

// ISoundResource represents any sound resource to be handled by the sound manager.
type ISoundResource interface {
	IObject
	Clear()
	Delete() int
	GetFilename() string
	GetFormat() int
	GetResource() (interface{}, int)
	New()
}

// SoundResource is the default implementation for the sound interface.
type SoundResource struct {
	*Object
	filename string
	format   int
	counter  int
	sound    *mix.Music
	chunk    *mix.Chunk
}

var _ ISoundResource = (*SoundResource)(nil)

// NewSound creates a new source instance.
func NewSound(name string, filename string, format int) *SoundResource {
	var err error
	Logger.Trace().Str("sound", name).Str("filename", filename).Msg("new sound")
	result := &SoundResource{
		Object:   NewObject(name),
		filename: filename,
		format:   format,
		counter:  0,
		sound:    nil,
		chunk:    nil,
	}
	switch format {
	case SoundMP3:
		if result.sound, err = mix.LoadMUS(filename); err != nil {
			Logger.Error().Err(err).Str("filename", filename).Msg("load MP3 error")
			panic(err)
		}
		break
	case SoundWAV:
		var data []byte
		if data, err = ioutil.ReadFile(filename); err != nil {
			Logger.Error().Err(err).Str("filename", filename).Msg("load WAV file error")
			panic(err)
		}
		if result.chunk, err = mix.QuickLoadWAV(data); err != nil {
			Logger.Error().Err(err).Str("filename", filename).Msg("load WAV error")
			panic(err)
		}
		break
	}
	return result
}

// Clear deletes sound even if counter is not zero.
func (s *SoundResource) Clear() {
	Logger.Trace().Str("sound", s.GetName()).Str("filename", s.GetFilename()).Msg("clear sound")
	s.counter = 1
	s.Delete()
}

// Delete deletes sound and release all memory.
func (s *SoundResource) Delete() int {
	Logger.Trace().Str("source", s.GetName()).Str("filename", s.GetFilename()).Msg("delete source")
	s.counter--
	if s.counter == 0 {
		if s.sound != nil {
			s.sound.Free()
		} else if s.chunk != nil {
			s.chunk.Free()
		}
	}
	return s.counter
}

// GetFilename returns sound filename.
func (s *SoundResource) GetFilename() string {
	return s.filename
}

// GetFormat returns sound format.
func (s *SoundResource) GetFormat() int {
	return s.format
}

// GetResource returns the sound resource and sound type
func (s *SoundResource) GetResource() (interface{}, int) {
	if s.sound != nil {
		return s.sound, s.format
	} else if s.chunk != nil {
		return s.chunk, s.format
	}
	return nil, -1
}

// New increases the number of times this sound is being used.
func (s *SoundResource) New() {
	s.counter++
}

// ISoundManager represents the handler that is in charge of all sounds.
type ISoundManager interface {
	IObject
	Clear()
	CreateSound(string, string, int) ISoundResource
	DeleteSound(ISoundResource) bool
	DoInit()
	GetSound(string) ISoundResource
	GetSoundByFilename(string) ISoundResource
	GetSoundByName(string) ISoundResource
	GetSounds() []ISoundResource
	OnStart()
}

// SoundManager is the default implementation for the Sound handler.
type SoundManager struct {
	*Object
	sounds []ISoundResource
}

var _ ISoundManager = (*SoundManager)(nil)

// NewSoundManager creates a new Sound handler instance.
func NewSoundManager(name string) *SoundManager {
	Logger.Trace().Str("Sound-manager", name).Msg("new sound-manager")
	return &SoundManager{
		Object: NewObject(name),
		sounds: []ISoundResource{},
	}
}

// Clear removes all sounds from the sound manager.
func (h *SoundManager) Clear() {
	Logger.Trace().Str("sound-manager", h.GetName()).Msg("Clear")
	for _, s := range h.sounds {
		s.Clear()
	}
	h.sounds = []ISoundResource{}
}

// CreateSound creates a new sound. If the same sound has already been created
// with the same filename, existing sound is returned.
func (h *SoundManager) CreateSound(name string, filename string, format int) ISoundResource {
	Logger.Trace().Str("sound-manager", h.GetName()).Str("name", name).Str("filename", filename).Msg("CreateSound")
	for _, sound := range h.sounds {
		if sound.GetFilename() == filename {
			sound.New()
			return sound
		}
	}
	sound := NewSound(name, filename, format)
	h.sounds = append(h.sounds, sound)
	return sound
}

// DeleteSound deleted sound from the sound manager. Memory resources are
// released from the given resource.
func (h *SoundManager) DeleteSound(sound ISoundResource) bool {
	Logger.Trace().Str("sound-manager", h.GetName()).Str("name", sound.GetName()).Str("filename", sound.GetFilename()).Msg("DeleteSound")
	for i := len(h.sounds) - 1; i >= 0; i-- {
		s := h.sounds[i]
		if s.GetID() == sound.GetID() {
			if result := s.Delete(); result == 0 {
				h.sounds = append(h.sounds[:i], h.sounds[i+1:]...)
			}
			return true
		}
	}
	return false
}

// DoInit initializes all sound manager resources
func (h *SoundManager) DoInit() {
	Logger.Trace().Str("Sound-manager", h.GetName()).Msg("DoInit")
}

// GetSound returns a sound with the given ID.
func (h *SoundManager) GetSound(id string) ISoundResource {
	for _, sound := range h.sounds {
		if sound.GetID() == id {
			return sound
		}
	}
	return nil
}

// GetSoundByFilename returns the sound with the given filename.
func (h *SoundManager) GetSoundByFilename(filename string) ISoundResource {
	for _, sound := range h.sounds {
		if sound.GetFilename() == filename {
			return sound
		}
	}
	return nil
}

// GetSoundByName returns the sound with the given name.
func (h *SoundManager) GetSoundByName(name string) ISoundResource {
	for _, sound := range h.sounds {
		if sound.GetName() == name {
			return sound
		}
	}
	return nil
}

// GetSounds returns all sounds.
func (h *SoundManager) GetSounds() []ISoundResource {
	return h.sounds
}

// OnStart initializes all sound manager structure.
func (h *SoundManager) OnStart() {
	Logger.Trace().Str("Sound-manager", h.GetName()).Msg("OnStart")
}
