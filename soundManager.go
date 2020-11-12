package engosdl

// ISoundManager represents the handler that is in charge of all sounds.
type ISoundManager interface {
	IObject
	OnStart()
}

// SoundManager is the default implementation for the Sound handler.
type SoundManager struct {
	*Object
}

var _ ISoundManager = (*SoundManager)(nil)

// NewSoundManager creates a new Sound handler instance.
func NewSoundManager(name string) *SoundManager {
	Logger.Trace().Str("Sound-handler", name).Msg("new sound-handler")
	return &SoundManager{
		Object: NewObject(name),
	}

}

// OnStart initializes all Sound handler structure.
func (h *SoundManager) OnStart() {
	Logger.Trace().Str("Sound-handler", h.GetName()).Msg("OnStart")
}
