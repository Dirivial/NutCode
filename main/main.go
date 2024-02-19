package main

import (
	"NutCode/rope"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
)

const (
	NORMAL = iota
	INSERT
	COMMAND
)

// Type for representing the cursor position
type Cursor struct {
	x    int
	y    int
	memX int
}

// Type for the most common data related to the window
type EditorWindow struct {
	screen          tcell.Screen
	style           tcell.Style
	cursor          Cursor
	endRow          int
	startRow        int
	startCol        int
	lineNumberWidth int
	contentOffset   int
}

// Completely redraw the screen
func drawFull(ew EditorWindow, content, fileName string, unsavedChanges bool, mode int) {
	ew.screen.Clear()
	drawContent(ew, content)
	drawLineNumbers(ew)
	drawStatus(ew, mode, fileName, unsavedChanges)
	ew.screen.ShowCursor(ew.cursor.x+ew.contentOffset, ew.cursor.y)
}

// Draw line numbers
func drawLineNumbers(ew EditorWindow) {
	_, height := ew.screen.Size()
	style := tcell.StyleDefault

	for i := 0; i < height; i++ {
		str := fmt.Sprint(i)
		off := ew.lineNumberWidth - len(str)
		for r := range str {
			byte := rune(str[r])
			ew.screen.SetContent(r+off, i, byte, nil, style)
		}
	}
}

// Draw the content to the screen
func drawContent(ew EditorWindow, content string) {
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
func drawStatus(ew EditorWindow, mode int, filename string, unsavedChanges bool) {
	style := tcell.StyleDefault.Background(tcell.Color18).Foreground(tcell.ColorReset)
	w, h := ew.screen.Size()

	// Draw information
	curEnd := drawMode(ew.screen, mode, h, style)
	curEnd = drawCursorPositionStatus(ew.screen, ew.cursor.x, ew.cursor.y, curEnd, h, style)
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

func main() {

	filename := flag.String("filename", "", "the name of the file to read from or write to")
	// Parse command-line arguments

	flag.Parse()

	// Check if the filename flag is provided
	if *filename == "" {
		fmt.Println("Please provide a filename")
		os.Exit(1)
	}
	// Open the file
	file, err := os.OpenFile(*filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer file.Close()

	// Get the file size
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}
	fileSize := fileInfo.Size()

	// Read the entire file into a byte slice
	data := make([]byte, fileSize)
	_, err = file.Read(data)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	content := rope.New(string(data))
	charCount := len(data)

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.SetCursorStyle(tcell.CursorStyleBlinkingBar)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	tabSize := 4
	unsavedChanges := false
	mode := NORMAL
	c := 0
	_, h := s.Size()
	cursor := Cursor{
		x:    0,
		y:    0,
		memX: 0,
	}
	editor := EditorWindow{
		screen:          s,
		startRow:        0,
		endRow:          h,
		startCol:        0,
		lineNumberWidth: 5,
		contentOffset:   7,
		cursor:          cursor,
		style:           defStyle,
	}
	drawFull(editor, content.GetContent(), *filename, unsavedChanges, mode)

	// Event loop
	s.ShowCursor(editor.contentOffset, editor.cursor.y)
	for {
		// Update screen
		s.Show()

		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else if ev.Key() == tcell.KeyCtrlS {
				// Save into file
				err := os.WriteFile(*filename, []byte(content.GetContent()), 0644)
				if err != nil {
					fmt.Println("Error writing to file: ", err)
					return
				}
				unsavedChanges = false

			} else if ev.Key() == tcell.KeyRight {
				nextRune := content.Index(c + 1)
				if nextRune != "\n" && nextRune != "" {
					c++
					editor.cursor.x++
				}
			} else if ev.Key() == tcell.KeyLeft {
				if editor.cursor.x > 0 {
					c--
					editor.cursor.x--
				}
			} else if ev.Key() == tcell.KeyDown {
				// Move cursor depending on line length

				minMove := content.SearchChar('\n', c+1)
				if minMove != -1 {
					// Note: this moves us after the newline
					c = minMove
					//x = 0
					editor.cursor.y++
					// Check if we can move the pointer foward to the old x position
					lineEnd := content.SearchChar('\n', c+1)
					if lineEnd != -1 {
						// Compute length of the line we move to
						diff := lineEnd - minMove

						if diff > 0 {
							if diff <= editor.cursor.x {
								// Move x to the end of the line (-1 for the newline)
								editor.cursor.x = diff - 1
							}
							c += editor.cursor.x
						} else {
							editor.cursor.x = 0
						}
					} else {
						editor.cursor.x = 0
					}
				}
			} else if ev.Key() == tcell.KeyUp {
				if editor.cursor.y > 0 {
					// Find end of last row
					lineEnd, err := content.SearchCharReverse('\n', c)
					if lineEnd != -1 && err == nil {
						// Move up
						editor.cursor.y--
						// Find start of last row
						lineStart, err := content.SearchCharReverse('\n', lineEnd-1)
						if err == nil {
							if lineStart == -1 {
								lineStart = 1
							}
							c = lineStart
							// Compute length of the line we move to
							diff := lineEnd - lineStart

							if diff > 0 {
								if diff <= editor.cursor.x {
									// Move x to the end of the line (-1 for the newline)
									editor.cursor.x = diff - 1
								}
								c += editor.cursor.x
							} else {
								editor.cursor.x = 0
							}
						} else {
							editor.cursor.x = 0
						}
					}
				} else {
					// Move to the beginning of the file
					editor.cursor.x = 0
					c = 0
				}
			} else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBS || ev.Key() == tcell.KeyBackspace2 {
				// Make sure there is something to delete
				if c > 0 {
					content = content.Delete(c-1, 1)
					charCount--
					c--
					unsavedChanges = true
					// Move cursor
					if editor.cursor.x > 0 {
						editor.cursor.x--
					} else {

						// Find start of last row
						lineStart, err := content.SearchCharReverse('\n', c)
						if err == nil {
							if lineStart == -1 {
								lineStart = 0
							}
							editor.cursor.y--
							// Compute length of the line we move to
							diff := c - lineStart
							editor.cursor.x = diff
						}
					}
				}
				// Move cursor depending on line length
			} else if ev.Key() == tcell.KeyEnter {
				// Insert a newline and move to the next line
				content = content.Insert(c, string('\n'))
				charCount++
				editor.cursor.x = 0
				editor.cursor.y++
				c++
				unsavedChanges = true

			} else if ev.Key() == tcell.KeyTab || ev.Key() == tcell.KeyTAB {
				// Tab key, replaces with tabSize number of spaces
				content = content.Insert(c, strings.Repeat(" ", tabSize))
				charCount += tabSize
				editor.cursor.x += tabSize
				c += tabSize
				unsavedChanges = true

			} else {
				// Catch-all for remaining characters,
				// adding them to the content at the current cursor position
				content = content.Insert(c, string(ev.Rune()))
				charCount++
				editor.cursor.x++
				c++
			}

			// === Draw ===
			drawFull(editor, content.GetContent(), *filename, unsavedChanges, mode)
		}
	}
}
