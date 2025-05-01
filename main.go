package main

import (
	"fmt"
	"os"

	"github.com/nsf/termbox-go"
)

var e Editor

func processEvent() {
	event := termbox.PollEvent()

	switch event.Type {
	case termbox.EventKey:
		if e.saveMode {
			saveProcessKeypress(event.Key, event.Ch)
		} else {
			editorProcessKeypress(event.Key, event.Ch)
		}
	case termbox.EventResize:
		e.win.setSize(event.Width, event.Height)
	}
}
func saveProcessKeypress(key termbox.Key, ch rune) {
	switch key {
	case termbox.KeyCtrlQ, termbox.KeyEsc:
		termbox.Close()
		os.Exit(0)
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		e.saveDelChar()
	case termbox.KeyEnter:
		e.saveMode = false
		e.WriteFile()
		e.cx = 0
		e.cy = 0
	default:
		e.saveInserChar(ch)
	}
}

func editorProcessKeypress(key termbox.Key, ch rune) {
	switch key {
	case termbox.KeyCtrlQ, termbox.KeyEsc:
		termbox.Close()
		os.Exit(0)
	case termbox.KeyCtrlS:
		e.WriteFile()
	case termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight, termbox.KeyArrowUp:
		e.moveCursor(key)
	case termbox.KeyPgup, termbox.KeyPgdn:
		for time := e.win.ROWS; time != 0; time-- {
			if key == termbox.KeyPgup {
				e.moveCursor(termbox.KeyArrowUp)
			} else {
				e.moveCursor(termbox.KeyArrowDown)
			}
		}
	case termbox.KeyHome:
		e.cx = 0
	case termbox.KeyEnd:
		e.cx = len(e.text_buffer[e.cy]) - 1
	case termbox.KeyTab:
		e.editorInsertChar('\t')
	case termbox.KeySpace:
		e.editorInsertChar(' ')
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		e.editorDelChar()
	case termbox.KeyEnter:
		e.editorInsertRow()
	default:
		e.editorInsertChar(ch)
	}
}

func run_editor() {
	err := termbox.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	e.Init()
	if len(os.Args) > 1 {
		e.ReadFile(os.Args[1])
	}
	for {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		e.drawStatusBar()
		e.editorScroll()
		e.draw()
		processEvent()
	}
}

func main() {
	run_editor()
}
