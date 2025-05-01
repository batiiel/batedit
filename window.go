package main

import (
	"github.com/nsf/termbox-go"
)

type window struct {
	COLS int
	ROWS int

	offset_col int
	offset_row int
}

func (win *window) initWindow() {
	win.offset_col = 0
	win.offset_row = 0

	win.COLS, win.ROWS = termbox.Size()
}

func (win *window) setSize(width, height int) {
	win.COLS = width
	win.ROWS = height
}

func (win *window) getSize() (col, row int) {
	return win.COLS, win.ROWS
}
