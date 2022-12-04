package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/phantompunk/stencil/pkg/stencil"
)

func main() {
	params := os.Args[1:]

	fontname := ""
	phrase := params[0]
	if len(params) >= 2 {
		fontname = params[1]
	}
	st, err := stencil.NewStencil(phrase, fontname)
	if err != nil {
		errAndExit(err.Error())
	}
	st.DrawText()

}

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprint(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.Usage()
	os.Exit(0)
}

func errAndExit(msg string) {
	fmt.Fprint(os.Stderr, msg, "\n")
	os.Exit(1)
}
