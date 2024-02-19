package editor

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

const (
	NORMAL = iota
	INSERT
	COMMAND
)

// Type for representing the cursor position
type Cursor struct {
	X    int
	Y    int
	MemX int
}

// Type for the most common data related to the window
type EditorWindow struct {
	screen          tcell.Screen
	Cursor          *Cursor
	style           tcell.Style
	endRow          int
	startRow        int
	startCol        int
	lineNumberWidth int
	contentOffset   int
}

func New(s tcell.Screen, startRow, endRow, startCol, lineNumberWidth, contentOffset int, style tcell.Style) *EditorWindow {
	cursor := Cursor{
		X:    0,
		Y:    0,
		MemX: 0,
	}
	return &EditorWindow{
		screen:          s,
		Cursor:          &cursor,
		endRow:          endRow,
		startRow:        startRow,
		startCol:        startCol,
		lineNumberWidth: lineNumberWidth,
		contentOffset:   contentOffset,
		style:           style,
	}
}

// Completely redraw the screen
func (ew *EditorWindow) DrawFull(content, fileName string, unsavedChanges bool, mode int) {
	ew.screen.Clear()
	ew.DrawContent(content)
	ew.DrawLineNumbers()
	ew.DrawStatus(mode, fileName, unsavedChanges)
	ew.screen.ShowCursor(ew.Cursor.X+ew.contentOffset, ew.Cursor.Y)
}

// Draw line numbers
func (ew *EditorWindow) DrawLineNumbers() {
	_, height := ew.screen.Size()
	style := tcell.StyleDefault

	for i := 0; i < height; i++ {
		str := fmt.Sprint(i)
		off := ew.lineNumberWidth - len(str)
		for j, r := range str {
			ew.screen.SetContent(j+off, i, r, nil, style)
		}
	}
}

// Draw the content to the screen
func (ew *EditorWindow) DrawContent(content string) {
	row := 0
	col := ew.contentOffset
	for _, r := range content {
		if r == '\n' {
			row++
			col = ew.contentOffset
			ew.screen.SetContent(col, row, r, nil, ew.style)
		} else {
			ew.screen.SetContent(col, row, r, nil, ew.style)
			col++
		}
	}
}

// Draw a statusbar showing line:col numbers, filename, mode and if there are unsaved changes
func (ew *EditorWindow) DrawStatus(mode int, filename string, unsavedChanges bool) {
	style := tcell.StyleDefault.Background(tcell.Color18).Foreground(tcell.ColorReset)
	w, h := ew.screen.Size()

	// Draw information
	curEnd := drawMode(ew.screen, mode, h, style)
	curEnd = drawCursorPositionStatus(ew.screen, ew.Cursor.X, ew.Cursor.Y, curEnd, h, style)
	curEnd = drawFileStatus(ew.screen, filename, unsavedChanges, curEnd, h, style)

	// Fill the rest of the row
	for i := curEnd; i < w; i++ {
		ew.screen.SetContent(i, h-1, rune(' '), nil, style)
	}
}

// Draw file name & if the changes the user has made are saved
func drawFileStatus(s tcell.Screen, fileName string, unsavedChanges bool, startAt, height int, style tcell.Style) int {
	info := ""
	if unsavedChanges {
		info = "*"
	}
	info += fileName

	for i, r := range info {
		s.SetContent(i+startAt, height-1, r, nil, style)
	}
	s.SetContent(len(info)+startAt, height-1, rune(' '), nil, style)
	return startAt + len(info)

}

// Draw the current line and col number in the status bar
func drawCursorPositionStatus(s tcell.Screen, lineNr, colNr, startAt, height int, style tcell.Style) int {
	info := fmt.Sprintf(" %d:%d ", lineNr+1, colNr)
	for i, r := range info {
		s.SetContent(i+startAt, height-1, r, nil, style)
	}
	return startAt + len(info)
}

// Draw the current mode in the status bar
func drawMode(s tcell.Screen, mode, height int, style tcell.Style) int {

	modeString := ""
	switch mode {
	case NORMAL:
		modeString = "NORMAL"
	case INSERT:
		modeString = "INSERT"
	case COMMAND:
		modeString = "COMMAND"
	default:
		modeString = "unknown"
	}
	s.SetContent(0, height-1, rune(' '), nil, style)
	for i, r := range modeString {
		s.SetContent(i+1, height-1, r, nil, style)
	}
	s.SetContent(len(modeString)+1, height-1, rune(' '), nil, style)
	return len(modeString) + 2
}
