package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
)

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

func drawContent(s tcell.Screen, cx, cy, offset int, style tcell.Style, content []rune) {
	row := 0
	col := offset
	for _, r := range content {
		if r == '\n' {
			row++
			col = offset
			s.SetContent(col, row, 'A', nil, style)
		} else {
			s.SetContent(col, row, r, nil, style)
			col++
		}
	}
}

func main() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	testContent := []rune("This is some text content.\nI wonder how this will be displayed\n\nBruhmode\tengaged.\n")
	content := make([]rune, 2048)
	copy(content, testContent)

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

	lineNumRoom := 5
	contentStart := lineNumRoom + 2
	drawLineNumbers(s, lineNumRoom)
	drawContent(s, 0, 0, contentStart, defStyle, content)

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
			} else if ev.Key() == tcell.KeyDown {
				y++
				// Move cursor depending on line length
			} else if ev.Key() == tcell.KeyUp {
				y--
				// Move cursor depending on line length
			} else if ev.Key() == tcell.KeyEnter {
				content[c] = '\n'
				s.Clear()
				drawContent(s, 0, 0, contentStart, defStyle, content)
				drawLineNumbers(s, lineNumRoom)
				x = 0
				y++
				c++
			} else {
				content[c] = ev.Rune()
				s.Clear()
				drawContent(s, 0, 0, contentStart, defStyle, content)
				drawLineNumbers(s, lineNumRoom)
				// s.SetContent(contentStart+x, y, ev.Rune(), nil, defStyle)
				x++
				c++
			}

			s.ShowCursor(x+contentStart, y)
		}
	}
}
