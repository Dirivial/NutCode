package main

import (
	"NutCode/editor"
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
	c := 0

	cachedContent := content.GetContent()
	editor := editor.New(s, 0, 0, 5, 7, defStyle)
	editor.ComputeNumRows(cachedContent)

	editor.DrawFull(cachedContent, *filename, unsavedChanges)

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
					editor.MoveX(1)
				}
			} else if ev.Key() == tcell.KeyLeft {
				if editor.Cursor.X > 0 {
					c--
					editor.MoveX(-1)
				}
			} else if ev.Key() == tcell.KeyDown {
				// Move cursor depending on line length

				minMove := content.SearchChar('\n', c+1)
				if minMove != -1 {
					// Note: this moves us after the newline
					c = minMove
					editor.MoveY(1)
					// Check if we can move the pointer foward to the old x position
					lineEnd := content.SearchChar('\n', c+1)
					if lineEnd != -1 {
						// Compute length of the line we move to
						diff := lineEnd - minMove

						if diff > 0 {
							if diff <= editor.Cursor.X+editor.StartCol {
								// Move x to the end of the line (-1 for the newline)
								editor.SetX(diff - 1)
							}
							c += editor.Cursor.X + editor.StartCol
						} else {
							editor.ResetX()
						}
					} else {
						editor.ResetX()
					}
				}
			} else if ev.Key() == tcell.KeyUp {
				if editor.Cursor.Y > 0 {
					// Find end of last row
					lineEnd, err := content.SearchCharReverse('\n', c)
					if lineEnd != -1 && err == nil {
						// Move up
						editor.MoveY(-1)
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
								if diff <= editor.Cursor.X+editor.StartCol {
									// Move x to the end of the line (-1 for the newline)
									editor.SetX(diff - 1)
								}
								c += editor.Cursor.X + editor.StartCol
							} else {
								editor.ResetX()
							}
						} else {
							editor.ResetX()
						}
					}
				} else {
					// Move to the beginning of the file
					c = 0
					editor.ResetX()
				}
			} else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBS || ev.Key() == tcell.KeyBackspace2 {
				// Make sure there is something to delete
				if c > 0 {
					content = content.Delete(c-1, 1)
					charCount--
					c--
					unsavedChanges = true
					cachedContent = content.GetContent()
					// Move cursor
					if editor.Cursor.X > 0 {
						editor.MoveX(-1)
					} else {

						// Find start of last row
						lineStart, err := content.SearchCharReverse('\n', c)
						if err == nil {
							if lineStart == -1 {
								lineStart = 0
							}
							editor.MoveY(-1)
							// Compute length of the line we move to
							diff := c - lineStart
							editor.SetX(diff)
							editor.NumRows--
						}
					}
				}
				// Move cursor depending on line length
			} else if ev.Key() == tcell.KeyEnter {
				// Insert a newline and move to the next line
				content = content.Insert(c, string('\n'))
				charCount++
				editor.NumRows++
				editor.ResetX()
				editor.MoveY(1)
				c++
				unsavedChanges = true
				cachedContent = content.GetContent()

			} else if ev.Key() == tcell.KeyTab || ev.Key() == tcell.KeyTAB {
				// Tab key, replaces with tabSize number of spaces
				content = content.Insert(c, strings.Repeat(" ", tabSize))
				charCount += tabSize
				editor.MoveX(tabSize)
				c += tabSize
				unsavedChanges = true
				cachedContent = content.GetContent()

			} else {
				// Catch-all for remaining characters,
				// adding them to the content at the current cursor position
				content = content.Insert(c, string(ev.Rune()))
				charCount++
				editor.MoveX(1)
				c++
				unsavedChanges = true
				cachedContent = content.GetContent()
			}

			// === Draw ===
			editor.DrawFull(cachedContent, *filename, unsavedChanges)
		}
	}
}
