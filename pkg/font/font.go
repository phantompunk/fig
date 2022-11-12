package font

import (
	"bufio"
	"bytes"
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

func NewFont(fontName string) *Font {
	if len(fontName) == 0 {
		fontName = DEFAULT_FONT
	}
	ft := &Font{}
	ft.name = fontName

	// fontPath := filepath.Join("..", "fonts", fontName+".flf")
	// if _, err := os.Stat(fontPath); err != nil {
	// 	logrus.Fatal("Font not found", err)
	// }
	fontBytes, _ := Asset(path.Join("fonts", fontName+".flf"))

	bytesReader := bytes.NewReader(fontBytes)
	// f, err := os.Open(fontPath)
	// check(err)
	// defer f.Close()

	// reader := read
	scanner := bufio.NewScanner(bytesReader)
	scanner.Split(bufio.ScanLines)
	// var lines []string
	// for scanner.Scan() {
	// 	lines = append(lines, scanner.Text())
	// }

	ft.SetAttributes(scanner)
	ft.SetLetters(scanner)
	// _ @
	// / |@
	// | |@
	// |_|@
	//    @@
	// fmt.Print(lines[102])
	// fmt.Println(lines[142])
	// fmt.Print(lines[103])
	// fmt.Println(lines[143])
	// fmt.Print(lines[104])
	// fmt.Println(lines[144])
	// fmt.Print(lines[105])
	// fmt.Println(lines[145])
	// fmt.Print(lines[106])
	// fmt.Println(lines[146])
	return ft
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
	// letter := make([][]string, f.height)
	// charValue := int(rune(char[0]))
	// start := 1 + f.comments + (charValue-32)*5

	// // letter = append(letter, char)
	// // fmt.Printf("%d", char[0])
	// logrus.Infof("Parsing Char: %c, %d", char[0], char[0])
	// for _, letter := range char {
	// 	logrus.Infof("%c", letter)
	// }

}

func (f *Font) SetLetters(scanner *bufio.Scanner) {
	f.letters = append(f.letters, make([]string, f.height))
	// fmt.Print(f.letters)
	// f.letters[0] = []string{` ____@`}
	// f.letters[1] = []string{`|_  /@`}
	// f.letters[2] = []string{` / / @`}
	// f.letters[3] = []string{`/___|@`}
	// f.letters[4] = []string{`    @@`}

	for i := 0; i <= f.comments-1; i++ {
		scanner.Scan()
	}
	f.letters = append(f.letters, make([]string, f.height, f.height))
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
