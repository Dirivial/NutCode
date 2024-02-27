package editor

import (
	"fmt"
	"math"

	"github.com/gdamore/tcell/v2"
)

const (
	NORMAL = iota
	INSERT
	COMMAND
)

// Type for representing the cursor position
type Cursor struct {
	X int
	Y int
}

// Type for the most common data related to the window
type EditorWindow struct {
	screen          tcell.Screen
	Cursor          *Cursor
	style           tcell.Style
	NumRows         int
	height          int
	width           int
	startRow        int
	StartCol        int
	lineNumberWidth int
	contentOffset   int
}

func New(s tcell.Screen, startRow, StartCol, lineNumberWidth, contentOffset int, style tcell.Style) *EditorWindow {
	cursor := Cursor{
		X: 0,
		Y: 0,
	}
	w, h := s.Size()
	return &EditorWindow{
		screen:          s,
		Cursor:          &cursor,
		height:          h,
		width:           w,
		startRow:        startRow,
		StartCol:        StartCol,
		lineNumberWidth: lineNumberWidth,
		contentOffset:   contentOffset,
		style:           style,
	}
}

func (ew *EditorWindow) ComputeNumRows(content string) {
	for _, r := range content {
		if r == '\n' {
			ew.NumRows++
		}
	}
}

func (ew *EditorWindow) MoveY(numRows int) {
	if numRows > 0 {
		// Check if we should move the window down
		if ew.Cursor.Y+numRows+5 >= ew.height && ew.startRow+numRows+ew.height-1 <= ew.NumRows+1 {
			ew.startRow = min(ew.startRow+numRows, ew.NumRows+1)
		} else {
			if ew.Cursor.Y+numRows == ew.height-1 {
				ew.startRow += numRows
			} else {
				ew.Cursor.Y += numRows

			}
		}
	} else {
		// Check if we should move the window up
		if ew.Cursor.Y+numRows-5 <= 0 && ew.startRow+numRows >= 0 {
			ew.startRow = max(ew.startRow+numRows, 0)
		} else {
			ew.Cursor.Y = max(numRows+ew.Cursor.Y, 0)
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (ew *EditorWindow) MoveX(numCols int) {
	if numCols > 0 {
		if ew.Cursor.X+numCols+5+ew.contentOffset+ew.lineNumberWidth >= ew.width {
			ew.StartCol += numCols
		} else {
			ew.Cursor.X += numCols
		}
	} else {
		if ew.Cursor.X+numCols <= ew.StartCol {
			if ew.StartCol > -1*numCols {
				ew.StartCol += numCols
			} else {
				ew.Cursor.X = max(ew.StartCol+numCols, 0)
				ew.StartCol = 0
			}
		} else {
			ew.Cursor.X = max(ew.Cursor.X+numCols, 0)
		}
	}
}

func (ew *EditorWindow) SetX(col int) {

	windowSize := ew.width - ew.contentOffset
	if col > windowSize {
		leftSpace := int(math.Floor(float64(windowSize/3) * 2))
		ew.StartCol = col - leftSpace
		ew.Cursor.X = col - ew.StartCol
	} else {
		ew.StartCol = 0
		ew.Cursor.X = col
	}
}

func (ew *EditorWindow) ResetX() {
	ew.StartCol = 0
	ew.Cursor.X = 0
}

// Completely redraw the screen
func (ew *EditorWindow) DrawFull(content, fileName string, unsavedChanges bool) {
	ew.screen.Clear()
	ew.DrawContent(content)
	ew.DrawLineNumbers()
	ew.DrawStatus(fileName, unsavedChanges)
	ew.screen.ShowCursor(ew.Cursor.X+ew.contentOffset, ew.Cursor.Y)
}

// Draw line numbers
func (ew *EditorWindow) DrawLineNumbers() {
	_, height := ew.screen.Size()
	style := tcell.StyleDefault.Foreground(tcell.Color140)
	activeRow := tcell.StyleDefault.Foreground(tcell.ColorReset)

	for i := 0; i < height; i++ {
		if i < ew.Cursor.Y {
			str := fmt.Sprint(ew.Cursor.Y - i)
			off := ew.lineNumberWidth - len(str)
			for j, r := range str {
				ew.screen.SetContent(j+off, i, r, nil, style)
			}
		} else if i > ew.Cursor.Y {
			str := fmt.Sprint(i - ew.Cursor.Y)
			off := ew.lineNumberWidth - len(str)
			for j, r := range str {
				ew.screen.SetContent(j+off, i, r, nil, style)
			}
		} else {
			str := fmt.Sprint(i + ew.startRow)
			off := ew.lineNumberWidth - len(str)
			for j, r := range str {
				ew.screen.SetContent(j+off, i, r, nil, activeRow)
			}
		}
	}
}

// Draw the content to the screen
func (ew *EditorWindow) DrawContent(content string) {
	row := 0
	col := ew.contentOffset
	epicCol := col
	minCol := ew.contentOffset + ew.StartCol
	activeRow := tcell.StyleDefault.Background(tcell.Color24).Foreground(tcell.ColorReset)
	for _, r := range content {
		if row > ew.NumRows {
			// Ignore rest of file
			break
		}
		if r == '\n' {
			row++
			if row >= ew.startRow {
				ew.screen.SetContent(col, row-ew.startRow, r, nil, ew.style)
				col = ew.contentOffset
			}
		} else {
			if row >= ew.startRow && row <= ew.startRow+ew.height {
				if row-ew.startRow == ew.Cursor.Y {
					if col >= minCol {
						ew.screen.SetContent(col-ew.StartCol, row-ew.startRow, r, nil, activeRow)
					}
					epicCol = col
				} else {
					if col >= minCol {
						ew.screen.SetContent(col-ew.StartCol, row-ew.startRow, r, nil, ew.style)
					}
				}
				col++
			}
		}
	}
	// Fill rest of activeRow
	for i := epicCol + 1 - ew.StartCol; i < ew.width; i++ {
		ew.screen.SetContent(i, ew.Cursor.Y, ' ', nil, activeRow)
	}
}

// Draw a statusbar showing line:col numbers, filename and if there are unsaved changes
func (ew *EditorWindow) DrawStatus(filename string, unsavedChanges bool) {
	style := tcell.StyleDefault.Background(tcell.Color18).Foreground(tcell.ColorReset)
	w, h := ew.screen.Size()

	// Draw information
	curEnd := drawCursorPositionStatus(ew.screen, ew.Cursor.X+ew.StartCol, ew.Cursor.Y+ew.startRow, 0, h, style)
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
