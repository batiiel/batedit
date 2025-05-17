package editor

import (
	"batedit/model"
	"batedit/screen"
	"fmt"
	"os"
	"strconv"

	"github.com/nsf/termbox-go"
)

type EditorMode int

const (
	ModeEdit EditorMode = iota
	ModeSave
)

type Editor struct {
	Doc       *model.Document
	Screen    *screen.ScreenBuffer
	Mode      EditorMode
	oldCursor struct{ X, Y int }
}

func New() *Editor {
	err := termbox.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	e := &Editor{}
	e.Mode = ModeEdit
	e.Doc = model.NewDocument()
	e.Screen = screen.NewScreenBuffer()
	return e
}

func (e *Editor) Draw() {
	e.Screen.Clear()
	e.DrawStatusBar()

	if e.Mode == ModeEdit {
		e.Screen.ScrollDocument(e.Doc)
	}

	e.Screen.RenderDocument(e.Doc)

	if e.Mode == ModeEdit {
		termbox.SetCursor(e.Doc.Cursor.X-e.Screen.OffsetCol+e.Doc.CountNumLines+1, e.Doc.Cursor.Y-e.Screen.OffsetRow)
	} else if e.Mode == ModeSave {
		termbox.SetCursor(e.Doc.Cursor.X, e.Doc.Cursor.Y)
	}

	termbox.Flush()
}

func (e *Editor) DrawStatusBar() {
	if e.Mode == ModeEdit {
		e.StatusBar()
	} else if e.Mode == ModeSave {
		e.SaveBar()
	}
}

func (e *Editor) SaveBar() {
	text := "save file:" + e.Doc.Name
	x := len([]rune(text))
	e.Doc.Cursor.X = x
	e.Doc.Cursor.Y = e.Screen.Height - 1
	e.printBar(text)
}

func (e *Editor) StatusBar() {
	len_bar := 0
	text := "filename:" + e.Doc.Name
	pos_cur := "ln:" + strconv.Itoa(e.Doc.Cursor.Y) + " col:" + strconv.Itoa(e.Doc.Cursor.X) + " "
	len_bar += len([]rune(text)) + len([]rune(pos_cur))
	if len_bar >= e.Screen.Width {
		e.printBar("...")
		return
	}
	e.printBar(text)
	e.printPosCurBar(pos_cur)
}

func (e *Editor) printPosCurBar(posText string) {
	buf := []rune(posText)
	x := e.Screen.Width - len(buf)
	for _, ch := range buf {
		termbox.SetCell(x, e.Screen.Height-1, ch, termbox.ColorBlue, termbox.ColorDarkGray)
		x++
	}
}

func (e *Editor) printBar(text string) {
	buffer_line := text
	for x := 0; x < e.Screen.Width; x++ {
		ch := ' '
		if x < len(buffer_line) {
			ch = rune(buffer_line[x])
		}
		termbox.SetCell(x, e.Screen.Height-1, ch, termbox.ColorBlue, termbox.ColorDarkGray)
	}
}

func (e *Editor) HandlerEvent() {
	event := termbox.PollEvent()

	switch event.Type {
	case termbox.EventKey:
		if e.Mode == ModeEdit {
			e.editorHandlerKeypress(event.Key, event.Ch)
		} else if e.Mode == ModeSave {
			e.saveHandlerKeypress(event.Key, event.Ch)
		}
	case termbox.EventResize:
		e.Screen.ReSize(event.Width, event.Height)
	}
}

func (e *Editor) saveHandlerKeypress(key termbox.Key, ch rune) {
	switch key {
	case termbox.KeyCtrlQ, termbox.KeyEsc:
		termbox.Close()
		os.Exit(0)
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		e.saveDelChar()
	case termbox.KeyEnter:
		e.Doc.SaveToFile()
		e.Mode = ModeEdit
		e.Doc.Cursor.X = e.oldCursor.X
		e.Doc.Cursor.Y = e.oldCursor.Y
	default:
		e.saveInserChar(ch)
	}
}

func (e *Editor) editorHandlerKeypress(key termbox.Key, ch rune) {
	switch key {
	case termbox.KeyCtrlQ, termbox.KeyEsc:
		termbox.Close()
		os.Exit(0)
	case termbox.KeyCtrlS:
		if e.Doc.Name == "" {
			e.Mode = ModeSave
			e.oldCursor.X = e.Doc.Cursor.X
			e.oldCursor.Y = e.Doc.Cursor.Y
			return
		}
		e.Doc.SaveToFile()
	case termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight, termbox.KeyArrowUp:
		e.moveCursor(key)
	case termbox.KeyPgup, termbox.KeyPgdn:
		for time := e.Screen.Height; time != 0; time-- {
			if key == termbox.KeyPgup {
				e.moveCursor(termbox.KeyArrowUp)
			} else {
				e.moveCursor(termbox.KeyArrowDown)
			}
		}
	case termbox.KeyHome:
		e.Doc.Cursor.X = 0
	case termbox.KeyEnd:
		e.Doc.Cursor.X = len(e.Doc.TextBuffer[e.Doc.Cursor.Y])
	case termbox.KeyTab:
		e.Doc.InsertChar('\t', e.Screen.Height)
	case termbox.KeySpace:
		e.Doc.InsertChar(' ', e.Screen.Height)
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		e.Doc.DeleteChar()
	case termbox.KeyEnter:
		e.Doc.Enter()
	default:
		e.Doc.InsertChar(ch, e.Screen.Height)
	}
}

func (e *Editor) moveCursor(key termbox.Key) {
	switch key {
	case termbox.KeyArrowUp:
		if e.Doc.Cursor.Y > 0 { // Можно двигаться вверх?
			e.Doc.Cursor.Y-- // Перемещаемся на строку выше
			// Корректируем X, чтобы не выйти за границы новой строки

			maxX := len(e.Doc.TextBuffer[e.Doc.Cursor.Y])
			if e.oldCursor.X < maxX {
				e.Doc.Cursor.X = e.oldCursor.X
			} else if e.Doc.Cursor.X > maxX {
				e.Doc.Cursor.X = maxX
			}
		}
	case termbox.KeyArrowDown:
		if e.Doc.Cursor.Y < len(e.Doc.TextBuffer)-1 { // Можно двигаться вниз?
			e.Doc.Cursor.Y++ // Перемещаемся на строку ниже
			// Корректируем X, чтобы не выйти за границы новой строки

			maxX := len(e.Doc.TextBuffer[e.Doc.Cursor.Y])
			if e.oldCursor.X < maxX {
				e.Doc.Cursor.X = e.oldCursor.X
			} else if e.Doc.Cursor.X > maxX {
				e.Doc.Cursor.X = maxX
			}
		}
	case termbox.KeyArrowRight:
		if e.Doc.Cursor.X < len(e.Doc.TextBuffer[e.Doc.Cursor.Y]) {
			e.Doc.Cursor.X++
		} else if e.Doc.Cursor.Y < len(e.Doc.TextBuffer)-1 {
			e.Doc.Cursor.Y++
			e.Doc.Cursor.X = 0
		}
		e.oldCursor.X = e.Doc.Cursor.X
	case termbox.KeyArrowLeft:
		if e.Doc.Cursor.X > 0 {
			e.Doc.Cursor.X--
		} else if e.Doc.Cursor.Y > 0 {
			e.Doc.Cursor.Y--
			e.Doc.Cursor.X = len(e.Doc.TextBuffer[e.Doc.Cursor.Y])
		}
		e.oldCursor.X = e.Doc.Cursor.X
	}
}

func (e *Editor) saveDelChar() {
	l := len([]rune(e.Doc.Name))
	e.Doc.Name = e.Doc.Name[:l-1]
}
func (e *Editor) saveInserChar(ch rune) {
	e.Doc.Name = e.Doc.Name + string(ch)
}
