package model

import (
	"bufio"
	"fmt"
	"os"
)

type Document struct {
	Name          string
	TextBuffer    [][]rune
	Cursor        struct{ X, Y int }
	CountNumLines int
}

func NewDocument() *Document {
	d := &Document{
		Name:       "",
		TextBuffer: [][]rune{{}},
		Cursor: struct {
			X int
			Y int
		}{X: 0, Y: 0},
	}
	d.CountNumLines = countNumber(len(d.TextBuffer))
	return d
}

func (d *Document) UpdateCountNum() {
	d.CountNumLines = countNumber(len(d.TextBuffer))
}

/* INSERT CHAR - start */
func (d *Document) InsertChar(ch rune, countRows int) {
	if ch == 0 {
		return
	}
	if d.Cursor.Y == countRows {
		d.TextBuffer = append(d.TextBuffer, []rune{})
	}
	lineInsertChar(&d.TextBuffer[d.Cursor.Y], d.Cursor.X, rune(ch))
	d.Cursor.X++
	d.UpdateCountNum()
}

func lineInsertChar(line *[]rune, at int, ch rune) {
	if at < 0 || at > len(*line) {
		at = len(*line)
	}
	*line = append(*line, 0)
	copy((*line)[at+1:], (*line)[at:])
	(*line)[at] = ch
}

/* INSERT CHAR - end */

/* DELETE CHAR - start */
func (d *Document) DeleteChar() {
	if d.Cursor.Y == len(d.TextBuffer) {
		return
	}
	if d.Cursor.X > 0 {
		d.Cursor.X--
		lineDeleteChar(&d.TextBuffer[d.Cursor.Y], d.Cursor.X)
	} else if d.Cursor.X == 0 {
		d.raiseLine()
	}
	d.UpdateCountNum()
}

func lineDeleteChar(line *[]rune, at int) {
	if at < 0 || at >= len(*line) {
		return
	}
	*line = append((*line)[:at], (*line)[at+1:]...)
}

func (d *Document) raiseLine() {
	if d.Cursor.Y == 0 {
		return
	}
	new_line := []rune{}
	new_line = append(new_line, d.TextBuffer[d.Cursor.Y-1]...)
	new_line = append(new_line, d.TextBuffer[d.Cursor.Y]...)
	d.Cursor.X = len(d.TextBuffer[d.Cursor.Y-1])

	d.TextBuffer = append(d.TextBuffer[:d.Cursor.Y], d.TextBuffer[d.Cursor.Y+1:]...)
	d.TextBuffer[d.Cursor.Y-1] = new_line

	d.Cursor.Y--
	d.UpdateCountNum()
}

/* DELETE CHAR - end */

func (d *Document) Enter() {
	rightLine := d.TextBuffer[d.Cursor.Y][d.Cursor.X:]
	leftLine := d.TextBuffer[d.Cursor.Y][:d.Cursor.X]

	d.TextBuffer = append(d.TextBuffer, []rune{})
	copy(d.TextBuffer[d.Cursor.Y+1:], d.TextBuffer[d.Cursor.Y:])

	d.TextBuffer[d.Cursor.Y] = leftLine
	d.TextBuffer[d.Cursor.Y+1] = rightLine

	d.Cursor.Y++
	d.Cursor.X = 0
	d.UpdateCountNum()
}

// READ FILE
func (d *Document) ReadFile(filename string) {
	d.Name = filename

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for line_num := 0; scanner.Scan(); line_num++ {
		if len(d.TextBuffer) <= line_num {
			d.TextBuffer = append(d.TextBuffer, []rune{})
		}
		for _, ch := range scanner.Text() {
			d.TextBuffer[line_num] = append(d.TextBuffer[line_num], rune(ch))
		}
	}
	d.UpdateCountNum()
}

// WRITE FILE
func (d *Document) SaveToFile() {
	if d.Name == "" {
		d.Name = "testFile.txt"
	}
	file, err := os.Create(d.Name)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for row, line := range d.TextBuffer {
		new_line := "\n"
		if row == len(d.TextBuffer) {
			new_line = ""
		}
		write_line := string(line) + new_line
		_, err := writer.WriteString(write_line)
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
	writer.Flush()
	d.UpdateCountNum()
}

func countNumber(l int) int {
	count := 0
	for l > 0 {
		count++
		l /= 10
	}
	return count
}
