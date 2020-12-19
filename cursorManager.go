package engosdl

import "github.com/veandco/go-sdl2/sdl"

// CursorUpdate represents any cursor update per frame.  Cursor manager
// collects all cursor updates per frame in order to decide what should
// be the cursor to be used.
type CursorUpdate struct {
	inFocus  IEntity
	cursorID sdl.SystemCursor
}

// NewCursorUpdate creates a new cursor update instance.
func NewCursorUpdate(entity IEntity, cursorID sdl.SystemCursor) *CursorUpdate {
	return &CursorUpdate{
		inFocus:  entity,
		cursorID: cursorID,
	}
}

// ICursorManager represents the interface for the cursor manager.
type ICursorManager interface {
	IObject
	CursorUpdate(IEntity, sdl.SystemCursor)
	DoInit()
	OnStart()
	OnAfterUpdate()
}

// CursorManager represents the manager to handle cursor updates.
type CursorManager struct {
	*Object
	inFocus  IEntity
	cursor   *sdl.Cursor
	cursorID sdl.SystemCursor
	updates  []*CursorUpdate
}

var _ ICursorManager = (*CursorManager)(nil)

// NewCursorManager creates a new cursor manager instance.
func NewCursorManager(name string) *CursorManager {
	Logger.Trace().Str("cursor-manager", name).Msg("new cursor-manager")
	return &CursorManager{
		Object:   NewObject(name),
		inFocus:  nil,
		cursor:   nil,
		cursorID: sdl.SYSTEM_CURSOR_ARROW,
		updates:  []*CursorUpdate{},
	}
}

// CursorUpdate adds a new cursor update to the manager.
func (h *CursorManager) CursorUpdate(entity IEntity, cursorID sdl.SystemCursor) {
	h.updates = append(h.updates, NewCursorUpdate(entity, cursorID))
}

// DoInit initializes all cursor manager resources.
func (h *CursorManager) DoInit() {
	Logger.Trace().Str("cursor-manager", h.GetName()).Msg("DoInit")
}

// OnStart initializes all cursor handler structure.
func (h *CursorManager) OnStart() {
	Logger.Trace().Str("cursor-manager", h.GetName()).Msg("OnStart")
}

// OnAfterUpdate updates the cursor for the frame. Engine calls this after all
// updates have been executed.
func (h *CursorManager) OnAfterUpdate() {
	var inFocus IEntity = nil
	for _, cursorUpdate := range h.updates {
		if cursorUpdate.inFocus != nil {
			inFocus = cursorUpdate.inFocus
			h.cursorID = cursorUpdate.cursorID
			h.cursor = sdl.CreateSystemCursor(h.cursorID)
			sdl.SetCursor(h.cursor)
			break
		}
	}
	if inFocus == nil && h.inFocus != nil {
		h.inFocus = nil
		h.cursorID = sdl.SYSTEM_CURSOR_ARROW
		h.cursor = sdl.CreateSystemCursor(h.cursorID)
		sdl.SetCursor(h.cursor)
	}
	h.inFocus = inFocus
	h.updates = []*CursorUpdate{}
}
