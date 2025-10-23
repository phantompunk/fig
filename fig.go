package fig

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed fonts/*.flf
var figFonts embed.FS

type GlyphDict map[rune]Glyph

type FigFont struct {
	name     string
	metadata Metadata
	glyphs   GlyphDict
}

// fontloader takes an FS?
// parse []byte
func Font(name string) (*FigFont, error) {
	data, err := figFonts.ReadFile("fonts/" + name + ".flf")
	if err != nil {
		return nil, err
	}

	// fmt.Println("Content: ", string(data))
	font, err := parseFont(name, data)
	return font, nil
}

type Metadata struct {
	signature      string
	hardBlank      rune
	height         int
	baseline       int
	maxLength      int
	oldLayout      int
	commentLines   int
	printDirection int
	fullLayout     int
	codeTag        int
	comments       string
}

type Glyph struct {
	lines []string
	width int
}

func parseFont(name string, data []byte) (*FigFont, error) {
	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)

	if !scanner.Scan() {
		return nil, fmt.Errorf("empty font file")
	}

	// flf2a$ 6 4 6 -1 4
	header := scanner.Text()
	if !strings.HasPrefix(header, "flf2a") {
		return nil, fmt.Errorf("unsupported font file")
	}

	font := &FigFont{name: name, glyphs: make(GlyphDict)}
	meta := Metadata{}

	parts := strings.Fields(header[5:])
	if len(parts) < 6 {
		return nil, fmt.Errorf("invalid header format")
	}

	meta.hardBlank = rune(header[5])
	meta.height, _ = strconv.Atoi(parts[1])
	meta.baseline, _ = strconv.Atoi(parts[2])
	meta.maxLength, _ = strconv.Atoi(parts[3])
	meta.oldLayout, _ = strconv.Atoi(parts[4])
	meta.commentLines, _ = strconv.Atoi(parts[5])

	if len(parts) > 6 {
		meta.printDirection, _ = strconv.Atoi(parts[6])
	}

	if len(parts) > 7 {
		meta.fullLayout, _ = strconv.Atoi(parts[7])
	}

	if len(parts) > 8 {
		meta.codeTag, _ = strconv.Atoi(parts[8])
	}

	for range meta.commentLines {
		if !scanner.Scan() {
			return nil, fmt.Errorf("unexpected eof in comments")
		}
	}

	for charCode := 32; charCode <= 126; charCode++ {
		char, err := readCharacter(scanner, meta.height)
		if err != nil {
			return nil, fmt.Errorf("failed to read character %d: %w", charCode, err)
		}
		font.glyphs[rune(charCode)] = char
	}

	font.metadata = meta
	return font, nil
}

func readCharacter(scanner *bufio.Scanner, height int) (Glyph, error) {
	lines := make([]string, height)
	width := 0

	for i := range height {
		if !scanner.Scan() {
			return Glyph{}, fmt.Errorf("unexpected end of file while reading characters")
		}

		line := scanner.Text()
		line = strings.TrimRight(line, "@")
		// fmt.Println("Line:", line)

		width = max(width, len(line))
		lines[i] = line
	}

	return Glyph{
		lines: lines,
		width: width,
	}, nil
}

func (f *FigFont) Render(text string) string {
	letters := make([]Glyph, len(text))

	for i, r := range text {
		letters[i] = f.glyphs[r]
	}

	// letters[0].lines[0] + letters[1].lines[0] +	letters[2].lines[0] +
	// letters[0].lines[1] + letters[1].lines[1] +	letters[2].lines[1] +
	var result strings.Builder
	for h := range f.metadata.height {
		for _, l := range letters {
			result.WriteString(l.lines[h])
		}
		if h < f.metadata.height-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}
