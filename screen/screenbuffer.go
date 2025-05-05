package screen

import (
	"batedit/model"
	"strconv"

	"github.com/nsf/termbox-go"
)

type ScreenBuffer struct {
	Width, Height        int
	OffsetRow, OffsetCol int
	// cells                []termbox.Cell
}

func NewScreenBuffer() *ScreenBuffer {
	w, h := termbox.Size()
	buf := &ScreenBuffer{Width: w, Height: h, OffsetRow: 0, OffsetCol: 0}
	return buf
}

func (buf *ScreenBuffer) Clear() {
	termbox.Clear(termbox.ColorLightGray, termbox.ColorBlack)
}
func (buf *ScreenBuffer) DrawDocument(doc *model.Documnet) {
	buf.scrollDocument(doc)
	buf.renderDocumnet(doc)
	termbox.SetCursor(doc.Cursor.X-buf.OffsetCol+doc.CoutnNumLines, doc.Cursor.Y-buf.OffsetRow)
	termbox.Flush()
}

func (buf *ScreenBuffer) renderDocumnet(doc *model.Documnet) {
	for row := 0; row <= buf.Height-2; row++ {
		bufRow := row + buf.OffsetRow

		coloNumLine := termbox.ColorDarkGray
		if bufRow == doc.Cursor.Y {
			coloNumLine = termbox.ColorLightGray
		}
		if bufRow < len(doc.TextBuffer) {
			lineNumber := bufRow + 1

			numSt := strconv.Itoa(lineNumber)

			for i, ch := range numSt {
				xPos := doc.CoutnNumLines - len(numSt) + i
				termbox.SetCell(xPos, row, ch, coloNumLine, termbox.ColorBlack)
			}
		}
		for col := 0; col <= buf.Width; col++ {
			bufCol := col + buf.OffsetCol
			if bufRow >= len(doc.TextBuffer) {
				termbox.SetCell(0, row, '*', termbox.ColorBlue, termbox.ColorBlack)
			} else if bufCol < len(doc.TextBuffer[bufRow]) {
				ch := doc.TextBuffer[bufRow][bufCol]
				offsetLine := col + doc.CoutnNumLines

				// coutnNum := countNumber(bufRow)
				// bf := bufRow

				switch ch {
				case '\t':
					termbox.SetCell(offsetLine, row, '→', termbox.ColorDarkGray, termbox.ColorBlack)
				case ' ':
					termbox.SetCell(offsetLine, row, '•', termbox.ColorDarkGray, termbox.ColorBlack)
				case '(', ')', '[', ']', '{', '}':
					termbox.SetCell(offsetLine, row, ch, termbox.ColorYellow, termbox.ColorBlack)
				default:
					termbox.SetCell(offsetLine, row, ch, termbox.ColorWhite, termbox.ColorBlack)
				}
			}
		}
	}
}

func (buf *ScreenBuffer) scrollDocument(doc *model.Documnet) {

	// Если курсор выше видимой области
	if doc.Cursor.Y < buf.OffsetRow {
		buf.OffsetRow = doc.Cursor.Y
	}
	// Если курсор ниже видимой области
	if doc.Cursor.Y > buf.OffsetRow+buf.Height-2 {
		buf.OffsetRow = doc.Cursor.Y - buf.Height + 2
	}
	// Если курсор левее видимой области
	if doc.Cursor.X < buf.OffsetCol {
		buf.OffsetCol = doc.Cursor.X
	}
	// Если курсор правее видимой области
	if doc.Cursor.X > buf.OffsetCol+buf.Width-doc.CoutnNumLines-1 {
		buf.OffsetCol = doc.Cursor.X - buf.Width + 1 + doc.CoutnNumLines
	}
	lineLength := len(doc.TextBuffer[doc.Cursor.Y])
	// Не скроллим дальше конца строки
	if buf.OffsetCol > lineLength-buf.Width+1+doc.CoutnNumLines {
		buf.OffsetCol = max(0, lineLength-buf.Width)
	}
}

func (buf *ScreenBuffer) ReSize(w, h int) {
	buf.Width = w
	buf.Height = h
}

func countNumber(l int) int {
	count := 0
	for l > 0 {
		count++
		l /= 10
	}
	return count
}
