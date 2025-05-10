package token

import "github.com/nsf/termbox-go"

type Token struct {
	Content string
	Color   termbox.Attribute
}

var GoKeyWords = map[string]termbox.Attribute{
	//ColorBlue
	"func":      termbox.ColorBlue,
	"package":   termbox.ColorBlue,
	"interface": termbox.ColorBlue,
	"struct":    termbox.ColorBlue,
	"type":      termbox.ColorBlue,
	"var":       termbox.ColorBlue,
	"const":     termbox.ColorBlue,
	"chan":      termbox.ColorBlue,
	"map":       termbox.ColorBlue,
	//ColorMagenta
	"import":   termbox.ColorMagenta,
	"break":    termbox.ColorMagenta,
	"case":     termbox.ColorMagenta,
	"continue": termbox.ColorMagenta,
	"default":  termbox.ColorMagenta,
	"defer":    termbox.ColorMagenta,
	"else":     termbox.ColorMagenta,
	"for":      termbox.ColorMagenta,
	"go":       termbox.ColorMagenta,
	"goto":     termbox.ColorMagenta,
	"if":       termbox.ColorMagenta,
	"range":    termbox.ColorMagenta,
	"return":   termbox.ColorMagenta,
	//ColorCyan
	"switch":      termbox.ColorCyan,
	"fallthrough": termbox.ColorCyan,
}

func IsGoKeyword(word string) bool {
	_, ok := GoKeyWords[word]
	return ok
}
