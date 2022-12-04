package font

import (
	"bufio"
	"bytes"
	"fmt"
	"path"
	"strings"
	"testing"
)

var testFonts = []struct {
	name      string
	signature string
	height    int
	comments  int
	hardblank string
}{
	{
		name:      "drpepper",
		signature: "flf2a$ 5 4 20 0 16",
		height:    5,
		comments:  16,
		hardblank: "$",
	},
	{
		name:      "starwars",
		signature: "flf2a$ 7 6 22 15 4",
		height:    7,
		comments:  4,
		hardblank: "$",
	},
	{
		name:      "short",
		signature: "flf2a$ 3 2 8 -1 4",
		height:    3,
		comments:  4,
		hardblank: "$",
	},
}

func TestNewFonts(t *testing.T) {
	for _, test := range testFonts {
		t.Run(fmt.Sprintf("Font:%s", test.name), func(t *testing.T) {
			actual := NewFont(test.name)
			if actual.height != test.height {
				t.Errorf("Expected: %d, actual: %d", test.height, actual.height)
			}
			if actual.comments != test.comments {
				t.Errorf("Expected: %d, actual: %d", test.comments, actual.comments)
			}
			if actual.hardblank != test.hardblank {
				t.Errorf("Expected: %s, actual: %s", test.hardblank, actual.hardblank)
			}
		})
	}
}

func TestSetFontAttributes(t *testing.T) {
	// Given
	fontMetadata := "flf2a$ 5 4 20 0 16"

	// When
	font := &Font{}
	scanner := bufio.NewScanner(strings.NewReader(fontMetadata))
	font.SetAttributes(scanner)

	// Then
	if font.height != 5 {
		t.Fatalf("Font height is not correct got: %d, want: %d", font.height, 5)
	}
	if font.comments != 16 {
		t.Fatalf("Font comment lines is not correct got: %d, want: %d", font.comments, 16)
	}
	if font.hardblank != "$" {
		t.Fatalf("Font hardblank is not correct got: %s, want: %s", font.hardblank, "$")
	}
}

func TestSetFontChar(t *testing.T) {
	fontName := "drpepper"
	firstChar := "19"
	fmt.Print(rune(firstChar[0]))
	font := NewFont(fontName)
	// font.ParseChar(firstChar)
	if len(font.name) > 0 {
		t.Log("Pass")
	}
}

func TestSetFontLetters(t *testing.T) {
	font := &Font{comments: 16, height: 5}
	fontBytes, _ := Asset(path.Join("fonts", "drpepper.flf"))
	bytesReader := bytes.NewReader(fontBytes)
	scanner := bufio.NewScanner(bytesReader)

	font.SetLetters(scanner)

	// fmt.Println(font.letters[0])
	fmt.Println(font.letters[52][0])
	fmt.Println(font.letters[52][1])
	fmt.Println(font.letters[52][2])
	fmt.Println(font.letters[52][3])
	fmt.Println(font.letters[52][4])
	upperCaseT := 52
	if font.letters[upperCaseT][1] != ` ___ @` {
		t.Fail()
	}
	// if font.letters[upperCaseT][4] != `     @@` {
	// 	t.Fail()
	// }
}
