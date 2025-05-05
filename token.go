package main

import "github.com/nsf/termbox-go"

var keywords = map[string]termbox.Attribute{
	"func":    termbox.ColorBlue,
	"package": termbox.ColorBlue,
	"import":  termbox.ColorLightRed,
}
