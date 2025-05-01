package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/nsf/termbox-go"
)

type Editor struct {
	source_file string
	text_buffer [][]rune

	saveMode bool
	//undo_buffer [][]rune

	msg string

	cx_temp int

	cx int
	cy int

	win window
}

func (e *Editor) Init() {
	e.source_file = ""
	e.text_buffer = [][]rune{{}}
	e.saveMode = false

	e.msg = ""

	e.cx = 0
	e.cy = 0

	e.win.initWindow()
}

func (e *Editor) moveCursor(key termbox.Key) {
	switch key {
	case termbox.KeyArrowUp:
		if e.cy > 0 { // Можно двигаться вверх?
			e.cy-- // Перемещаемся на строку выше
			// Корректируем X, чтобы не выйти за границы новой строки

			maxX := len(e.text_buffer[e.cy])
			if e.cx_temp < maxX {
				e.cx = e.cx_temp
			} else if e.cx > maxX {
				e.cx = maxX
			}
		}
	case termbox.KeyArrowDown:
		if e.cy < len(e.text_buffer)-1 { // Можно двигаться вниз?
			e.cy++ // Перемещаемся на строку ниже
			// Корректируем X, чтобы не выйти за границы новой строки

			maxX := len(e.text_buffer[e.cy])
			if e.cx_temp < maxX {
				e.cx = e.cx_temp
			} else if e.cx > maxX {
				e.cx = maxX
			}
		}
	case termbox.KeyArrowRight:
		if e.cx < len(e.text_buffer[e.cy]) {
			e.cx++
		} else {
			if e.cy < len(e.text_buffer)-1 {
				e.cy++
				e.cx = 0
			}
		}
		e.cx_temp = e.cx
	case termbox.KeyArrowLeft:
		if e.cx > 0 {
			e.cx--
		} else {
			if e.cy > 0 {
				e.cy--
				e.cx = len(e.text_buffer[e.cy])
			}
		}
		e.cx_temp = e.cx
	}
}

func (e *Editor) editorScroll() {
	if e.saveMode {
		return
	}
	// Если курсор выше видимой области
	if e.cy < e.win.offset_row {
		e.win.offset_row = e.cy
	}
	// Если курсор ниже видимой области
	if e.cy >= e.win.offset_row+e.win.ROWS-1 {
		e.win.offset_row = e.cy - e.win.ROWS + 2
	}
	lineLength := len(e.text_buffer[e.cy])
	// Если курсор левее видимой области
	if e.cx < e.win.offset_col {
		e.win.offset_col = e.cx
	}
	// Если курсор правее видимой области
	if e.cx >= e.win.offset_col+e.win.COLS {
		e.win.offset_col = e.cx - e.win.COLS + 1
	}
	// Не скроллим дальше конца строки
	if e.win.offset_col > lineLength-e.win.COLS {
		e.win.offset_col = max(0, lineLength-e.win.COLS)
	}
}

// READ FILE
func (e *Editor) ReadFile(filename string) {
	e.source_file = filename

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for line_num := 0; scanner.Scan(); line_num++ {
		if len(e.text_buffer) <= line_num {
			e.text_buffer = append(e.text_buffer, []rune{})
		}
		for _, ch := range scanner.Text() {
			e.text_buffer[line_num] = append(e.text_buffer[line_num], rune(ch))
		}
		//e.text_buffer[line_num] = append(e.text_buffer[line_num], rune('\n'))

	}
}

// WRITE FILE
func (e *Editor) WriteFile() {
	if e.source_file == "" {
		e.saveMode = true
		return
	}
	file, err := os.Create(e.source_file)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for row, line := range e.text_buffer {
		new_line := "\n"
		if row == len(e.text_buffer) {
			new_line = ""
		}
		write_line := string(line) + new_line
		_, err := writer.WriteString(write_line)
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
	writer.Flush()
}

func (e *Editor) draw() {
	for row := 0; row < e.win.ROWS-1; row++ {
		bufRow := row + e.win.offset_row
		if bufRow >= len(e.text_buffer) {
			// Очистка пустых строк
			for col := 0; col < e.win.COLS; col++ {
				termbox.SetCell(0, row, '*', termbox.ColorDefault, termbox.ColorDefault)
			}
			continue
		}

		line := e.text_buffer[bufRow]
		col := 0
		for bufCol := e.win.offset_col; col < e.win.COLS; {
			if bufCol >= len(line) {
				termbox.SetCell(col+e.win.offset_col, row+e.win.offset_row, ' ', termbox.ColorDefault, termbox.ColorDefault)
				col++
				continue
			}

			ch := line[bufCol]
			if ch == '\t' {
				// Отрисовка табуляции
				//for i := 0; i < tabSpace; i++ {
				termbox.SetCell(col, row, '→', termbox.ColorDarkGray, termbox.ColorDefault)
				col++
				//}
				bufCol++
			} else if ch == ' ' {
				termbox.SetCell(col, row, '•', termbox.ColorDarkGray, termbox.ColorDefault)
				col++
				//}
				bufCol++
			} else {
				termbox.SetCell(col, row, ch, termbox.ColorDefault, termbox.ColorDefault)
				col++
				bufCol++
			}
		}
	}

	// Отрисовка курсора
	//e.countTab()
	termbox.SetCursor(e.cx-e.win.offset_col, e.cy-e.win.offset_row)
	termbox.Flush()

}
func (e *Editor) drawStatusBar() {
	if !e.saveMode {
		e.StatusBar()
	} else {
		e.SaveBar()
	}

}
func (e *Editor) SaveBar() {
	text := "filename:" + e.source_file
	x := len([]rune(text))
	e.cx = x
	e.cy = e.win.ROWS - 1
	e.printBar(text)
}

func (e *Editor) StatusBar() {
	len_bar := 0
	text := "filename:" + e.source_file
	pos_cur := "ln:" + strconv.Itoa(e.cy) + " col:" + strconv.Itoa(e.cx) + " "
	len_bar += len([]rune(text)) + len([]rune(pos_cur))
	if len_bar >= e.win.COLS {
		e.printBar("...")
		termbox.Flush()
		return
	}
	e.printBar(text)
	e.printPosCurBar(pos_cur)
	termbox.Flush()
}

func (e *Editor) printPosCurBar(posText string) {
	buf := []rune(posText)
	x := e.win.COLS - len(buf)
	for _, ch := range buf {
		termbox.SetCell(x, e.win.ROWS-1, ch, termbox.ColorLightGray, termbox.ColorDarkGray)
		x++
	}
}

func (e *Editor) printBar(text string) {
	buffer_line := text
	for x := 0; x < e.win.COLS; x++ {
		ch := ' '
		if x < len(buffer_line) {
			ch = rune(buffer_line[x])
		}
		termbox.SetCell(x, e.win.ROWS-1, ch, termbox.ColorLightGray, termbox.ColorDarkGray)
	}
}

func (e *Editor) editorRowInsertChar(line *[]rune, at int, ch rune) {
	if at < 0 || at > len(*line) {
		at = len(*line)
	}
	*line = append(*line, 0)
	copy((*line)[at+1:], (*line)[at:])
	(*line)[at] = ch
}

func (e *Editor) editorInsertChar(ch rune) {
	if ch == 0 {
		return
	}
	if e.cy == e.win.ROWS {
		e.text_buffer = append(e.text_buffer, []rune{})
	}
	e.editorRowInsertChar(&e.text_buffer[e.cy], e.cx, rune(ch))
	e.cx++
}

func (e *Editor) editorDelRowCahr(line *[]rune, at int) {
	if at < 0 || at >= len(*line) {
		return
	}
	*line = append((*line)[:at], (*line)[at+1:]...)
}

func (e *Editor) editorDelChar() {
	if e.cy == len(e.text_buffer) {
		return
	}
	if e.cx > 0 {
		e.cx--
		e.editorDelRowCahr(&e.text_buffer[e.cy], e.cx)
	} else if e.cx == 0 {
		e.editorDelRow()
	}
}

func (e *Editor) editorDelRow() {
	if e.cy == 0 {
		return
	}
	new_line := []rune{}
	new_line = append(new_line, e.text_buffer[e.cy-1]...)
	new_line = append(new_line, e.text_buffer[e.cy]...)
	e.cx = len(e.text_buffer[e.cy-1])

	e.text_buffer = append(e.text_buffer[:e.cy], e.text_buffer[e.cy+1:]...)
	e.text_buffer[e.cy-1] = new_line

	e.cy--
}

func (e *Editor) editorInsertRow() {
	rightLine := e.text_buffer[e.cy][e.cx:]
	leftLine := e.text_buffer[e.cy][:e.cx]

	e.text_buffer = append(e.text_buffer, []rune{})
	copy(e.text_buffer[e.cy+1:], e.text_buffer[e.cy:])

	e.text_buffer[e.cy] = leftLine
	e.text_buffer[e.cy+1] = rightLine

	e.cy++
	e.cx = 0
}

func (e *Editor) saveInserChar(ch rune) {
	if ch == 0 {
		return
	}
	e.source_file += string(ch)
	e.cx++
}

func (e *Editor) saveDelChar() {
	l := len([]rune(e.source_file))
	if l == 0 {
		return
	}
	e.source_file = e.source_file[:l-1]
}
