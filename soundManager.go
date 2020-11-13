package engosdl

// ISoundManager represents the handler that is in charge of all sounds.
type ISoundManager interface {
	IObject
	DoInit()
	OnStart()
}

// SoundManager is the default implementation for the Sound handler.
type SoundManager struct {
	*Object
}

var _ ISoundManager = (*SoundManager)(nil)

// NewSoundManager creates a new Sound handler instance.
func NewSoundManager(name string) *SoundManager {
	Logger.Trace().Str("Sound-manager", name).Msg("new sound-manager")
	return &SoundManager{
		Object: NewObject(name),
	}

}

// DoInit initializes all sound manager resources
func (h *SoundManager) DoInit() {
	Logger.Trace().Str("Sound-manager", h.GetName()).Msg("DoInit")
}

// OnStart initializes all sound manager structure.
func (h *SoundManager) OnStart() {
	Logger.Trace().Str("Sound-manager", h.GetName()).Msg("OnStart")
}
