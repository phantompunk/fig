package font

import (
	"bufio"
	"bytes"
	"fmt"
	"path"
	"strconv"
	"strings"
)

type Font struct {
	name      string
	hardblank string
	height    int
	comments  int
	letters   [][]string
}

const DEFAULT_FONT = "drpepper"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func NewFont(fontName string) (*Font, error) {
	if len(fontName) == 0 {
		fontName = DEFAULT_FONT
	}

	fontBytes, err := Asset(path.Join("fonts", fmt.Sprintf("%s.flf", fontName)))
	if err != nil {
		return nil, err
	}

	bytesReader := bytes.NewReader(fontBytes)
	scanner := bufio.NewScanner(bytesReader)
	scanner.Split(bufio.ScanLines)

	ft := &Font{}

	ft.name = fontName
	ft.SetAttributes(scanner)
	ft.SetLetters(scanner)
	return ft, nil
}

func setHeight(metadata string) int {
	dat := strings.Split(metadata, " ")[1]
	d, _ := strconv.Atoi(string(dat))
	return d
}

func (f *Font) GetHeight() int {
	return f.height
}

func setComments(metadata string) int {
	comments, _ := strconv.Atoi(strings.Split(metadata, " ")[5])
	return comments
}

func setHardBlank(metadata string) string {
	sig := strings.Split(metadata, " ")[0]
	blank := sig[len(sig)-1]
	return string(blank)
}

func (f *Font) SetAttributes(scanner *bufio.Scanner) {
	scanner.Scan()
	metadata := scanner.Text()
	f.hardblank = setHardBlank(metadata)
	f.height = setHeight(metadata)
	f.comments = setComments(metadata)
}

func (f *Font) ParseChar(char string) {
}

func (f *Font) SetLetters(scanner *bufio.Scanner) {
	f.letters = append(f.letters, make([]string, f.height))

	for i := 0; i <= f.comments-1; i++ {
		scanner.Scan()
	}
	f.letters = append(f.letters, make([]string, f.height))
	for i := range f.letters[0] {
		f.letters[0][i] = "  "
	}

	letter := 1
	part := 0
	for scanner.Scan() {
		line := scanner.Text()
		f.letters[letter] = append(f.letters[letter], line)
		part++
		if part > f.height-1 {
			part = 0
			letter++
			f.letters = append(f.letters, []string{})
		}
	}
}

func (f *Font) GetLetters() [][]string {
	return f.letters
}
