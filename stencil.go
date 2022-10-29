package stencil

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Stencil struct{}

// var charA = `
//
//	__
//
// |__|
// |  |
// `
var charA = `
    
|\/|
|  |
`
var m = `     <>  [] `
var stars = "************"
var dashes = `------------`
var hashes = `############`

var fontSimple = map[string]string{
	"a": ` __ |__||  |`,
	"b": `    |__ |__)`,
	"c": ` __ |  '|__;`,
	"e": ` __ | _)|__,`,
	"r": ` __ |'_|| \ `,
	"i": `     <>  [] `,
	"l": `    |   |__ `,
	"m": `    |\/||  |`,
	"*": stars,
	"#": dashes,
	"-": hashes,
}

var fontTest = map[string]string{
	"*": stars,
	"-": dashes,
	"#": hashes,
}

//  __  __  __      __		Row 0
// |__||__)|  ¡|\/|| _|		Row 1
// |  ||__)|__!|  ||__;		Row 2

// __ 		Row 0
//|__|		Row 1
//|  |		Row 2

func NewStencil() *Stencil {
	return &Stencil{}
}

func (st *Stencil) DrawText(charArray []string) {
	for _, part := range charArray {
		fmt.Println(part)
	}
}

func (st *Stencil) AssignChars(phrase, font string) ([]string, error) {
	charArray := make([]string, 3) // [1,2,3]
	for index, char := range phrase {
		logrus.Debugf("Index:%d, %c, %s\n", index, char, fontSimple[string(char)])
		start := 0
		end := 4
		for part := 0; part < 3; part++ {
			charArray[part] += fontSimple[string(char)][start:end]
			start += 4
			end += 4
		}
	}
	return charArray, nil
}
