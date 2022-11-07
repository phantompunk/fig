package stencil

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
	height    int
	comments  int
	hardblank string
}{
	{
		name:      "drpepper",
		height:    5,
		comments:  16,
		hardblank: "$",
	},
}

func TestNewFonts(t *testing.T) {
	for _, testCase := range testFonts {
		steps := NewFont(testCase.name)
		if steps.name != testCase.name {
			t.Fatal("Failed name")
		}
		if steps.height != testCase.height {
			t.Fatal("Failed height")
		}
		if steps.comments != testCase.comments {
			t.Fatal("Failed comments")
		}
		if steps.hardblank != testCase.hardblank {
			t.Fatal("Failed hardblank")
		}

		t.Logf("PASS:")
	}
}

func TestSetFontAttributes(t *testing.T) {
	// Given
	fontMetadata := "flf2a$ 5 4 20 0 16"

	// When
	font := &font{}
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
	font := &font{comments: 16, height: 5}
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
	if font.letters[upperCaseT][0] != ` ___ @` {
		t.Fail()
	}
	if font.letters[upperCaseT][4] != `     @@` {
		t.Fail()
	}
}
