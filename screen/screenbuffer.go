package screen

import (
	"batedit/model"
	"batedit/syntax"
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

func (buf *ScreenBuffer) RenderDocument(doc *model.Document) {
	for row := 0; row < buf.Height-1; row++ {
		bufRow := row + buf.OffsetRow
		if bufRow >= len(doc.TextBuffer) {
			continue
		}

		// Отрисовка номеров строк (как у вас было)
		colorNumLine := termbox.ColorDarkGray
		if bufRow == doc.Cursor.Y {
			colorNumLine = termbox.ColorLightGray
		}
		lineNumber := bufRow + 1
		numStr := strconv.Itoa(lineNumber)
		for i, ch := range numStr {
			xPos := doc.CountNumLines - len(numStr) + i
			termbox.SetCell(xPos, row, ch, colorNumLine, termbox.ColorBlack)
		}

		// Токенизация строки
		tokens := syntax.Tokenize(doc.TextBuffer[bufRow])
		//fmt.Println("End")
		// Отрисовка токенов с подсветкой
		colOffset := doc.CountNumLines + 1 - buf.OffsetCol
		for _, token := range tokens {
			for _, ch := range token.Content {
				fgColor := token.Color
				termbox.SetCell(colOffset, row, ch, fgColor, termbox.ColorBlack)
				colOffset++
			}
		}

	}
}

// func (buf *ScreenBuffer) RenderDocument(doc *model.Document) {
// 	for row := 0; row <= buf.Height-2; row++ {
// 		bufRow := row + buf.OffsetRow

// 		coloNumLine := termbox.ColorDarkGray
// 		if bufRow == doc.Cursor.Y {
// 			coloNumLine = termbox.ColorLightGray
// 		}
// 		if bufRow < len(doc.TextBuffer) {
// 			lineNumber := bufRow + 1

// 			numSt := strconv.Itoa(lineNumber)

// 			for i, ch := range numSt {
// 				xPos := doc.CountNumLines - len(numSt) + i
// 				termbox.SetCell(xPos, row, ch, coloNumLine, termbox.ColorBlack)
// 			}
// 		}
// 		for col := 0; col <= buf.Width; col++ {
// 			bufCol := col + buf.OffsetCol
// 			if bufRow >= len(doc.TextBuffer) {
// 				termbox.SetCell(0, row, '*', termbox.ColorBlue, termbox.ColorBlack)
// 			} else if bufCol < len(doc.TextBuffer[bufRow]) {
// 				ch := doc.TextBuffer[bufRow][bufCol]
// 				offsetLine := col + doc.CountNumLines + 1

// 				// coutnNum := countNumber(bufRow)
// 				// bf := bufRow

// 				switch ch {
// 				case '\t':
// 					termbox.SetCell(offsetLine, row, '→', termbox.ColorDarkGray, termbox.ColorBlack)
// 				case ' ':
// 					termbox.SetCell(offsetLine, row, '•', termbox.ColorDarkGray, termbox.ColorBlack)
// 				case '(', ')', '[', ']', '{', '}':
// 					termbox.SetCell(offsetLine, row, ch, termbox.ColorYellow, termbox.ColorBlack)
// 				default:
// 					termbox.SetCell(offsetLine, row, ch, termbox.ColorWhite, termbox.ColorBlack)
// 				}
// 			}
// 		}
// 	}
// }

func (buf *ScreenBuffer) ScrollDocument(doc *model.Document) {

	endLine := 2 //для каректоровки курсора при скрлинге

	// Если курсор выше видимой области
	if doc.Cursor.Y < buf.OffsetRow {
		buf.OffsetRow = doc.Cursor.Y
	}
	// Если курсор ниже видимой области
	if doc.Cursor.Y > buf.OffsetRow+buf.Height-endLine {
		buf.OffsetRow = doc.Cursor.Y - buf.Height + endLine
	}
	// Если курсор левее видимой области
	if doc.Cursor.X < buf.OffsetCol {
		buf.OffsetCol = doc.Cursor.X
	}
	// Если курсор правее видимой области
	if doc.Cursor.X > buf.OffsetCol+buf.Width-doc.CountNumLines-endLine {
		buf.OffsetCol = doc.Cursor.X - buf.Width + endLine + doc.CountNumLines
	}
	// Не скроллим дальше конца строки
	lineLength := len(doc.TextBuffer[doc.Cursor.Y])
	if buf.OffsetCol > lineLength-buf.Width+endLine+doc.CountNumLines {
		buf.OffsetCol = max(0, lineLength-buf.Width+endLine+doc.CountNumLines)
	}
}

func (buf *ScreenBuffer) ReSize(w, h int) {
	buf.Width = w
	buf.Height = h
}

// func countNumber(l int) int {
// 	count := 0
// 	for l > 0 {
// 		count++
// 		l /= 10
// 	}
// 	return count
// }
