package stencil

import (
	"fmt"
	"regexp"

	"github.com/phantompunk/stencil/pkg/font"
)

// var charA = `
//
//	__
//
// |__|
// |  |
// `
// var charA = `

// |\/|
// |  |
// `
// var m = `     <>  [] `
// var stars = "************"
// var dashes = `------------`
// var hashes = `############`

// var fontSimple = map[string]string{
// 	"a": ` __ |__||  |`,
// 	"b": `    |__ |__)`,
// 	"c": ` __ |  '|__;`,
// 	"e": ` __ | _)|__,`,
// 	"r": ` __ |'_|| \ `,
// 	"i": `     <>  [] `,
// 	"l": `    |   |__ `,
// 	"m": `    |\/||  |`,
// 	"*": stars,
// 	"#": dashes,
// 	"-": hashes,
// }

// var fontTest = map[string]string{
// 	"*": stars,
// 	"-": dashes,
// 	"#": hashes,
// }

//  __  __  __      __		Row 0
// |__||__)|  ¡|\/|| _|		Row 1
// |  ||__)|__!|  ||__;		Row 2

// __ 		Row 0
// |__|		Row 1
// |  |		Row 2
type Stencil struct {
	phrase string
	font   *font.Font
}

func NewStencil(phrase, fontName string) *Stencil {
	// if len(font) == 0 {
	// 	font = "drpepper"
	// }
	font := font.NewFont(fontName)
	return &Stencil{
		phrase: phrase,
		font:   font,
	}
}

func (st *Stencil) DrawText() {
	// logrus.Infof("Phase: %s", st.phrase)
	// logrus.Infof("Phase: %d", rune(st.phrase[0]))
	// fmt.Print(st.font.letters)
	// for _, c := range st.phrase {
	// 	charValue := int(rune(c)) - 31
	// 	for parts := 0; parts < st.font.height; parts++ {
	// 		fmt.Println(st.font.letters[charValue][parts])
	// 	}
	// }
	regexp, _ := regexp.Compile(`@`)
	for p := 0; p <= st.font.GetHeight()-1; p++ {
		for _, c := range st.phrase {
			charValue := int(rune(c)) - 31
			fmt.Print(regexp.ReplaceAllString(st.font.GetLetters()[charValue][p], ""))
		}
		fmt.Println()
	}
	// st.font.letters[letter][part]
	// for each letter of the phrase
	// check if letter has been parsed
	// otherwise parse and store letter
	// for _, part := range assignChars(st.phrase, st.font.name) {
	// fmt.Println(part)
	// }
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
