package engosdl

// ISoundHandler represents the handler that is in charge of all sounds.
type ISoundHandler interface {
	IObject
	OnStart()
}

// SoundHandler is the default implementation for the Sound handler.
type SoundHandler struct {
	*Object
}

// NewSoundHandler creates a new Sound handler instance.
func NewSoundHandler(name string) *SoundHandler {
	Logger.Trace().Str("Sound-handler", name).Msg("new sound-handler")
	return &SoundHandler{
		Object: NewObject(name),
	}

}

// OnStart initializes all Sound handler structure.
func (h *SoundHandler) OnStart() {
	Logger.Trace().Str("Sound-handler", h.GetName()).Msg("OnStart")
}
