package main

import (
	"batedit/editor"
	"os"
)

// var e Editor

func main() {

	e := editor.New()
	if len(os.Args) > 1 {
		e.Doc.ReadFile(os.Args[1])
	}
	for {
		e.Draw()
		e.HandlerEvent()
	}
}
