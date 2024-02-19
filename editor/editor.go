package editor

import "github.com/gdamore/tcell/v2"

// Type for representing the cursor position
type Cursor struct {
	x    int
	y    int
	memX int
}

// Type for the most common data related to the window
type EditorWindow struct {
	screen   tcell.Screen
	startRow int
	endRow   int
	startCol int
	cursor   Cursor
}

func (e *EditorWindow) Clear() {
	e.screen.Clear()
}
