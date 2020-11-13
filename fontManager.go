package engosdl

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// IFont represents any graphical font to be handled by the font
// handler.
type IFont interface {
	IObject
	Clear()
	Delete() int
	GetFilename() string
	GetFont() *ttf.Font
	GetTextureFromFont(string, sdl.Color) *sdl.Texture
	New()
}

// Font is the default implementation for the font interface.
type Font struct {
	*Object
	filename string
	fontSize int
	font     *ttf.Font
	counter  int
}

var _ IFont = (*Font)(nil)

// NewFont creates a new font instance.
func NewFont(name string, filename string, fontSize int) *Font {
	var err error
	Logger.Trace().Str("font", name).Str("filename", filename).Msg("new font")
	result := &Font{
		Object:   NewObject(name),
		filename: filename,
		fontSize: fontSize,
		counter:  0,
	}
	result.font, err = ttf.OpenFont(filename, fontSize)
	if err != nil {
		Logger.Error().Err(err).Msg("OpenFont error")
		panic(err)
	}
	return result
}

// Clear deletes font even if counter is not zero.
func (r *Font) Clear() {
	Logger.Trace().Str("font", r.GetName()).Str("filename", r.GetFilename()).Msg("clear font")
	r.counter = 1
	r.Delete()
}

// Delete deletes font and relese all memory.
func (r *Font) Delete() int {
	Logger.Trace().Str("font", r.GetName()).Str("filename", r.GetFilename()).Msg("delete font")
	r.counter--
	if r.counter == 0 {
		r.font.Close()
	}
	return r.counter
}

// GetFilename returns font filename
func (r *Font) GetFilename() string {
	return r.filename
}

// GetFont returns font.
func (r *Font) GetFont() *ttf.Font {
	return r.font
}

// GetTextureFromFont returns a texture from the font surface.
func (r *Font) GetTextureFromFont(message string, color sdl.Color) *sdl.Texture {
	Logger.Trace().Str("font", r.GetName()).Str("filename", r.GetFilename()).Msg("get texture from font")
	surface, err := r.font.RenderUTF8Solid(message, color)
	if err != nil {
		Logger.Error().Err(err).Msg("RenderUTF8Solid error")
		panic(err)
	}
	defer surface.Free()

	texture, err := GetRenderer().CreateTextureFromSurface(surface)
	if err != nil {
		Logger.Error().Err(err).Msg("CreateTextureFromSurface error")
		panic(err)
	}
	return texture
}

// New increases the number of times this font is being used.
func (r *Font) New() {
	r.counter++
}

// IFontManager represents the handler that is in charge of all graphical
// fonts.
type IFontManager interface {
	IObject
	Clear()
	CreateFont(string, string, int) IFont
	DoInit()
	DeleteFont(IFont) bool
	GetFont(string) IFont
	GetFontByFilename(string) IFont
	GetFontByName(string) IFont
	GetFonts() []IFont
	OnStart()
}

// FontManager is the default implementation for the font handler.
type FontManager struct {
	*Object
	fonts []IFont
}

var _ IFontManager = (*FontManager)(nil)

// NewFontManager creates a new font handler instance.
func NewFontManager(name string) *FontManager {
	Logger.Trace().Str("font-manager", name).Msg("new font-manager")
	return &FontManager{
		Object: NewObject(name),
		fonts:  []IFont{},
	}

}

// Clear removes all fonts from the font handler.
func (h *FontManager) Clear() {
	Logger.Trace().Str("font-manager", h.GetName()).Msg("Clear")
	for _, r := range h.fonts {
		r.Clear()
	}
	h.fonts = []IFont{}
}

// CreateFont creates a new font. If the same font has already
// been created with the same filename, existing font is returned.
func (h *FontManager) CreateFont(name string, filename string, fontSize int) IFont {
	Logger.Trace().Str("font-manager", h.GetName()).Str("name", name).Str("filename", filename).Msg("CreateFont")
	for _, font := range h.fonts {
		if font.GetFilename() == filename {
			font.New()
			return font
		}
	}
	font := NewFont(name, filename, fontSize)
	h.fonts = append(h.fonts, font)
	return font
}

// DeleteFont deletes font from the handler. Memory fonts are
// released from the given font.
func (h *FontManager) DeleteFont(font IFont) bool {
	Logger.Trace().Str("font-manager", h.GetName()).Str("name", font.GetName()).Str("filename", font.GetFilename()).Msg("DeleteFont")
	for i := len(h.fonts) - 1; i >= 0; i-- {
		r := h.fonts[i]
		if r.GetID() == font.GetID() {
			if result := r.Delete(); result == 0 {
				h.fonts = append(h.fonts[:i], h.fonts[i+1:]...)
			}
			return true
		}
	}
	return false
}

// DoInit initializes all font manager resources.
func (h *FontManager) DoInit() {
	Logger.Trace().Str("font-manager", h.GetName()).Msg("DoInit")
}

// GetFont returns a font with the given font ID.
func (h *FontManager) GetFont(id string) IFont {
	for _, font := range h.fonts {
		if font.GetID() == id {
			return font
		}
	}
	return nil
}

// GetFontByFilename returns the font with the given filename.
func (h *FontManager) GetFontByFilename(filename string) IFont {
	for _, font := range h.fonts {
		if font.GetFilename() == filename {
			return font
		}
	}
	return nil
}

// GetFontByName returns the font with the given name.
func (h *FontManager) GetFontByName(name string) IFont {
	for _, font := range h.fonts {
		if font.GetName() == name {
			return font
		}
	}
	return nil
}

// GetFonts returns all fonts.
func (h *FontManager) GetFonts() []IFont {
	return h.fonts
}

// OnStart initializes all font handler structure.
func (h *FontManager) OnStart() {
	Logger.Trace().Str("font-manager", h.GetName()).Msg("OnStart")
}
