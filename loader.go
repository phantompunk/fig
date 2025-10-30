package fig

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
)

//go:embed fonts/*.flf
var figFS embed.FS

// ListFonts returns a list of available font names.
func ListFonts() []string {
	files, err := figFS.ReadDir("fonts")
	if err != nil {
		return nil
	}

	names := []string{}
	for _, file := range files {
		name := file.Name()[:len(file.Name())-4]
		names = append(names, name)
	}

	return names
}

func loadFont(name string) (*FigFont, error) {
	fileName := filepath.Join("fonts", name+".flf")

	data, err := figFS.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return parseFont(name, bytes.NewReader(data))
}

func parseFont(name string, data io.Reader) (*FigFont, error) {
	scanner := bufio.NewScanner(data)
	if !scanner.Scan() {
		return nil, fmt.Errorf("empty font file")
	}

	// flf2a$ 6 4 6 -1 4
	header := scanner.Text()
	meta, err := parseHeader(header)
	if err != nil {
		return nil, err
	}

	// skipCommentLines
	for range meta.commentLines {
		if !scanner.Scan() {
			return nil, fmt.Errorf("unexpected eof in comments")
		}
	}

	font := NewFigFont(name, meta)
	// parseCharacters
	for charCode := 32; charCode <= 126; charCode++ {
		char, err := readCharacter(scanner, meta.height, meta.hardBlank)
		if err != nil {
			return nil, fmt.Errorf("failed to read character %d: %w", charCode, err)
		}
		font.glyphs[rune(charCode)] = char
	}

	return font, nil
}

type headerParser struct {
	fields []string
	size   int
	err    error
}

func newParser(header string) (*headerParser, error) {
	if len(header) < 6 {
		return nil, fmt.Errorf("invalid header: %q", header)
	}

	if !strings.HasPrefix(header, "flf2a") {
		return nil, fmt.Errorf("invalid header prefix: %q", header)
	}

	fields := strings.Fields(header[5:])
	if len(fields) < 6 {
		return nil, fmt.Errorf("invalid header format: expected at least 6 fields, got %d", len(fields))
	}

	return &headerParser{fields: fields, size: len(fields)}, nil
}

func (p *headerParser) parseInt(i int, name string) int {
	if p.err != nil {
		return 0
	}

	if i >= len(p.fields) {
		p.err = fmt.Errorf("missing field: %s", name)
		return 0
	}

	v, err := strconv.Atoi(p.fields[i])
	if err != nil {
		p.err = fmt.Errorf("invalid field %s: %v", name, err)
		return 0
	}

	return v
}

func parseSmushMode(layout int) SmushMode {
	if layout < 0 {
		return SmushMode{Enabled: true}
	}

	return SmushMode{
		Enabled:      layout > 0,
		EqualChar:    layout&16 != 0,
		Underscore:   layout&32 != 0,
		Hierarchy:    layout&64 != 0,
		OppositePair: layout&128 != 0,
		BigX:         layout&256 != 0,
		Hardblank:    layout&512 != 0,
		// TODO: support vertical smushing
	}
}

func parseHeader(header string) (Metadata, error) {
	meta := Metadata{}

	parser, err := newParser(header)
	if err != nil {
		return Metadata{}, err
	}

	meta.hardBlank = rune(header[5])
	meta.height = parser.parseInt(1, "height")
	meta.baseline = parser.parseInt(2, "baseline")
	meta.maxLength = parser.parseInt(3, "max_length")
	meta.oldLayout = parser.parseInt(4, "old_layout")
	meta.commentLines = parser.parseInt(5, "comment_lines")

	if parser.size > 6 {
		meta.printDirection = parser.parseInt(6, "print_direction")
	}

	if parser.size > 7 {
		meta.fullLayout = parser.parseInt(7, "full_layout")
	} else {
		meta.fullLayout = meta.oldLayout
	}

	if parser.size > 8 {
		meta.codeTag = parser.parseInt(8, "code_tag")
	}

	meta.smushMode = parseSmushMode(meta.fullLayout)
	return meta, parser.err
}

func readCharacter(scanner *bufio.Scanner, height int, hardBlank rune) (Glyph, error) {
	lines := make([]string, height)
	width := 0

	for i := range height {
		if !scanner.Scan() {
			return Glyph{}, fmt.Errorf("unexpected end of file while reading characters")
		}

		line := scanner.Text()

		if i == height-1 {
			line = line[:len(line)-2]
		} else {
			line = line[:len(line)-1]
		}

		width = max(width, len(line))
		lines[i] = line
	}

	return Glyph{
		lines: lines,
		width: width,
	}, nil
}
