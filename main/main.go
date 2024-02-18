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

// Completely redraw the screen
func drawFull(s tcell.Screen, lineNumberRoom, offset, cx, cy int, style tcell.Style, content, fileName string, unsavedChanges bool, mode int) {
	s.Clear()
	drawContent(s, offset, style, content)
	drawLineNumbers(s, lineNumberRoom)
	drawStatus(s, mode, cy, cx, fileName, unsavedChanges)
}

// Draw line numbers
func drawLineNumbers(s tcell.Screen, offset int) {
	_, height := s.Size()
	style := tcell.StyleDefault

	for i := 0; i < height; i++ {
		str := fmt.Sprint(i)
		off := offset - len(str)
		for r := range str {
			byte := rune(str[r])
			s.SetContent(r+off, i, byte, nil, style)
		}
	}
}

// Draw the content to the screen
func drawContent(s tcell.Screen, offset int, style tcell.Style, content string) {
	row := 0
	col := offset
	for _, r := range content {
		if r == '\n' {
			row++
			col = offset
			s.SetContent(col, row, r, nil, style)
		} else {
			s.SetContent(col, row, r, nil, style)
			col++
		}
	}
}

// Draw a statusbar showing line:col numbers, filename, mode and if there are unsaved changes
func drawStatus(s tcell.Screen, mode, lineNr, colNr int, filename string, unsavedChanges bool) {
	style := tcell.StyleDefault.Background(tcell.Color18).Foreground(tcell.ColorReset)
	w, h := s.Size()

	// Draw information
	curEnd := drawMode(s, mode, h, style)
	curEnd = drawCursorPositionStatus(s, lineNr, colNr, curEnd, h, style)
	curEnd = drawFileStatus(s, filename, unsavedChanges, curEnd, h, style)

	// Fill the rest of the row
	for i := curEnd; i < w; i++ {
		s.SetContent(i, h-1, rune(' '), nil, style)
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
	lineNumRoom := 5
	contentStart := lineNumRoom + 2
	unsavedChanges := false
	drawFull(s, lineNumRoom, contentStart, 0, 0, defStyle, content.GetContent(), "myfile", unsavedChanges, NORMAL)

	// Event loop
	x, y, c := 0, 0, 0
	s.ShowCursor(contentStart, y)
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
					x++
				}
			} else if ev.Key() == tcell.KeyLeft {
				if x > 0 {
					c--
					x--
				}
			} else if ev.Key() == tcell.KeyDown {
				// Move cursor depending on line length

				minMove := content.SearchChar('\n', c+1)
				if minMove != -1 {
					// Note: this moves us after the newline
					c = minMove
					//x = 0
					y++
					// Check if we can move the pointer foward to the old x position
					lineEnd := content.SearchChar('\n', c+1)
					if lineEnd != -1 {
						// Compute length of the line we move to
						diff := lineEnd - minMove

						if diff > 0 {
							if diff <= x {
								// Move x to the end of the line (-1 for the newline)
								x = diff - 1
							}
							c += x
						} else {
							x = 0
						}
					} else {
						x = 0
					}
				}
			} else if ev.Key() == tcell.KeyUp {
				if y > 0 {
					// Find end of last row
					lineEnd, err := content.SearchCharReverse('\n', c)
					if lineEnd != -1 && err == nil {
						// Move up
						y--
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
								if diff <= x {
									// Move x to the end of the line (-1 for the newline)
									x = diff - 1
								}
								c += x
							} else {
								x = 0
							}
						} else {
							x = 0
						}
					}
				} else {
					// Move to the beginning of the file
					x = 0
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
					if x > 0 {
						x--
					} else {

						// Find start of last row
						lineStart, err := content.SearchCharReverse('\n', c)
						if err == nil {
							if lineStart == -1 {
								lineStart = 0
							}
							y--
							// Compute length of the line we move to
							diff := c - lineStart
							x = diff
						}
					}
				}
				// Move cursor depending on line length
			} else if ev.Key() == tcell.KeyEnter {
				// Insert a newline and move to the next line
				content = content.Insert(c, string('\n'))
				charCount++
				x = 0
				y++
				c++
				unsavedChanges = true

			} else if ev.Key() == tcell.KeyTab || ev.Key() == tcell.KeyTAB {
				// Tab key, replaces with tabSize number of spaces
				content = content.Insert(c, strings.Repeat(" ", tabSize))
				charCount += tabSize
				x += tabSize
				c += tabSize
				unsavedChanges = true

			} else {
				// Catch-all for remaining characters,
				// adding them to the content at the current cursor position
				content = content.Insert(c, string(ev.Rune()))
				charCount++
				x++
				c++
			}

			// === Draw ===
			drawFull(s, lineNumRoom, contentStart, x, y, defStyle, content.GetContent(), *filename, unsavedChanges, NORMAL)
			s.ShowCursor(x+contentStart, y)
		}
	}
}
