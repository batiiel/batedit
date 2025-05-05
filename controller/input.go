package controller

import (
	"batedit/model"

	"github.com/nsf/termbox-go"
)

type EventHandler struct {
	doc   *model.Documnet
	event *termbox.Event
}

func New(ev *termbox.Event, d *model.Documnet) *EventHandler {
	return &EventHandler{
		doc:   d,
		event: ev,
	}
}

// func (eh *EventHandler) Run() {
// 	switch eh.event.Type {
// 	case termbox.EventKey:
// 		eh.HandleKeyPress()
// 	case termbox.EventResize:
// 		ChangeWin()
// 	}
// }
// func (e *EventHandler) HandleKeyPress() {
// 	key := e.event.Key
// 	ch := eh.event.Ch
// 	switch key {
// 	case termbox.KeyCtrlQ, termbox.KeyEsc:
// 		termbox.Close()
// 		os.Exit(0)
// 	case termbox.KeyCtrlS:
// 		eh.doc.SaveToFile()
// 	case termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight, termbox.KeyArrowUp:
// 		e.moveCursor(key)
// 	case termbox.KeyPgup, termbox.KeyPgdn:
// 		for time := e.win.ROWS; time != 0; time-- {
// 			if key == termbox.KeyPgup {
// 				e.moveCursor(termbox.KeyArrowUp)
// 			} else {
// 				e.moveCursor(termbox.KeyArrowDown)
// 			}
// 		}
// 	case termbox.KeyHome:
// 		e.doc.Cursor.X = 0
// 	case termbox.KeyEnd:
// 		e.doc.Cursor.X = len(e.doc.TextBuffer[e.doc.Cursor.Y]) - 1
// 	case termbox.KeyTab:
// 		e.doc.InsertChar('\t', e.win.ROWS)
// 	case termbox.KeySpace:
// 		e.doc.InsertChar(' ', e.win.ROWS)
// 	case termbox.KeyBackspace, termbox.KeyBackspace2:
// 		e.doc.DeleteChar()
// 	case termbox.KeyEnter:
// 		e.doc.Enter()
// 	default:
// 		e.doc.InsertChar(ch, e.win.ROWS)
// 	}
// }

// func (e *EventHandler) moveCursor(key termbox.Key) {
// 	switch key {
// 	case termbox.KeyArrowUp:
// 		if e.doc.Cursor.Y > 0 { // Можно двигаться вверх?
// 			e.doc.Cursor.Y-- // Перемещаемся на строку выше
// 			// Корректируем X, чтобы не выйти за границы новой строки

// 			maxX := len(e.doc.TextBuffer[e.doc.Cursor.Y])
// 			if e.cx_temp < maxX {
// 				e.doc.Cursor.X = e.cx_temp
// 			} else if e.doc.Cursor.X > maxX {
// 				e.doc.Cursor.X = maxX
// 			}
// 		}
// 	case termbox.KeyArrowDown:
// 		if e.doc.Cursor.Y < len(e.doc.TextBuffer)-1 { // Можно двигаться вниз?
// 			e.doc.Cursor.Y++ // Перемещаемся на строку ниже
// 			// Корректируем X, чтобы не выйти за границы новой строки

// 			maxX := len(e.doc.TextBuffer[e.doc.Cursor.Y])
// 			if e.cx_temp < maxX {
// 				e.doc.Cursor.X = e.cx_temp
// 			} else if e.doc.Cursor.X > maxX {
// 				e.doc.Cursor.X = maxX
// 			}
// 		}
// 	case termbox.KeyArrowRight:
// 		if e.doc.Cursor.X < len(e.doc.TextBuffer[e.doc.Cursor.Y]) {
// 			e.doc.Cursor.X++
// 		} else {
// 			if e.doc.Cursor.Y < len(e.doc.TextBuffer)-1 {
// 				e.doc.Cursor.Y++
// 				e.doc.Cursor.X = 0
// 			}
// 		}
// 		e.cx_temp = e.doc.Cursor.X
// 	case termbox.KeyArrowLeft:
// 		if e.doc.Cursor.X > 0 {
// 			e.doc.Cursor.X--
// 		} else {
// 			if e.doc.Cursor.Y > 0 {
// 				e.doc.Cursor.Y--
// 				e.doc.Cursor.X = len(e.doc.TextBuffer[e.doc.Cursor.Y])
// 			}
// 		}
// 		e.cx_temp = e.doc.Cursor.X
// 	}
// }

func ChangeWin() {}

var keyBindings = map[termbox.Key]func(){}
