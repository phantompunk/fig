package main

import (
	"os"

	"github.com/phantompunk/stencil/stencil"
)

func main() {
	params := os.Args[1:]

	fontname := ""
	phrase := params[0]
	if len(params) >= 2 {
		fontname = params[1]
	}
	st := stencil.NewStencil(phrase, fontname)
	st.DrawText()

}
