package stencil

import (
	"fmt"
	"regexp"

	"github.com/phantompunk/stencil/pkg/font"
)

type Stencil struct {
	phrase string
	font   *font.Font
}

func NewStencil(phrase, fontName string) *Stencil {
	font := font.NewFont(fontName)
	return &Stencil{
		phrase: phrase,
		font:   font,
	}
}

func (st *Stencil) DrawText() {
	regexp, _ := regexp.Compile(`@`)
	for p := 0; p <= st.font.GetHeight()-1; p++ {
		for _, c := range st.phrase {
			charValue := int(rune(c)) - 31
			fmt.Print(regexp.ReplaceAllString(st.font.GetLetters()[charValue][p], ""))
		}
		fmt.Println()
	}
}

// func assignChars(phrase, font string) []string {
// 	var charArray = []string // [1,2,3]
// 	for index, char := range phrase {
// 		logrus.Debugf("Index:%d, %c, %s\n", index, char, fontSimple[string(char)])
// 		start := 0
// 		end := 4
// 		for part := 0; part < 3; part++ {
// 			charArray[part] += fontSimple[string(char)][start:end]
// 			start += 4
// 			end += 4
// 		}
// 	}
// 	return charArray
// }
