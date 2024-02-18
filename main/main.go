package main

import (
	"NutCode/rope"
	"fmt"
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// Completely redraw the screen
func drawFull(s tcell.Screen, lineNumberRoom, offset, cx, cy int, style tcell.Style, content string) {
	s.Clear()
	drawContent(s, cx, cy, offset, style, content)
	drawLineNumbers(s, lineNumberRoom)
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
func drawContent(s tcell.Screen, cx, cy, offset int, style tcell.Style, content string) {
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

func main() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	// TODO: Remove and implement reading of files
	testContent := string("This is some text content. I wonder how this will be displayed.\n    Bruhmode.engaged\n,,\n\n,\nBrusch")
	content := rope.New(testContent)
	charCount := len(testContent)

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
	drawLineNumbers(s, lineNumRoom)
	drawContent(s, 0, 0, contentStart, defStyle, content.GetContent())

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
					// === Draw ===
					drawFull(s, lineNumRoom, contentStart, 0, 0, defStyle, content.GetContent())
				}
				// Move cursor depending on line length
			} else if ev.Key() == tcell.KeyEnter {
				// Insert a newline and move to the next line
				content = content.Insert(c, string('\n'))
				charCount++
				x = 0
				y++
				c++

				// === Draw ===
				drawFull(s, lineNumRoom, contentStart, 0, 0, defStyle, content.GetContent())
			} else if ev.Key() == tcell.KeyTab || ev.Key() == tcell.KeyTAB {
				// Tab key, replaces with tabSize number of spaces
				content = content.Insert(c, strings.Repeat(" ", tabSize))
				charCount += tabSize
				x += tabSize
				c += tabSize

				// === Draw ===
				drawFull(s, lineNumRoom, contentStart, 0, 0, defStyle, content.GetContent())
			} else {
				// Catch-all for remaining characters,
				// adding them to the content at the current cursor position
				content = content.Insert(c, string(ev.Rune()))
				charCount++
				x++
				c++

				// === Draw ===
				drawFull(s, lineNumRoom, contentStart, 0, 0, defStyle, content.GetContent())
			}

			s.ShowCursor(x+contentStart, y)
		}
	}
}
