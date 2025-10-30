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
		char, err := readCharacter(scanner, meta.height)
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

const (
	BitEqualChar       = 1 << 0  // 1
	BitUnderscore      = 1 << 1  // 2
	BitHierarchy       = 1 << 2  // 4
	BitOppositePair    = 1 << 3  // 8
	BitBigX            = 1 << 4  // 16
	BitHardblank       = 1 << 5  // 32
	BitKern            = 1 << 6  // 64
	BitHSmush          = 1 << 7  // 128
	BitVEqualChar      = 1 << 8  // 256
	BitVUnderscore     = 1 << 9  // 512
	BitVHierarchy      = 1 << 10 // 1024
	BitHLine           = 1 << 11 // 2048
	BitVline           = 1 << 12 // 4096
	BitVKern           = 1 << 13 // 8192
	BitVSmush          = 1 << 14 // 16384
	BitLayoutKerning   = 1 << 0  // 1
	BitLayoutUniversal = 1 << 1  // 2
	BitLayoutVKerning  = 1 << 13 // 8192
	BitLayoutVSmushing = 1 << 14 // 16384
)

func parseSmushMode(mask int) SmushMode {
	if mask < 0 {
		return SmushMode{Enabled: false}
	}

	return SmushMode{
		// Enabled is true if the bitmask is greater than 0
		Enabled: mask > 0,

		// Horizontal Smushing Bits
		EqualChar:    mask&BitEqualChar != 0,
		Underscore:   mask&BitUnderscore != 0,
		Hierarchy:    mask&BitHierarchy != 0,
		OppositePair: mask&BitOppositePair != 0,
		BigX:         mask&BitBigX != 0,
		Hardblank:    mask&BitHardblank != 0,

		// Vertical Smushing Bits
		VEqualChar:  mask&BitVEqualChar != 0,
		VUnderscore: mask&BitVUnderscore != 0,
		VHierarchy:  mask&BitVHierarchy != 0,
		HLine:       mask&BitHLine != 0,
		Vline:       mask&BitVline != 0,
	}
}

func parseLayoutMode(mask int) LayoutMode {
	if mask < 0 {
		return LayoutMode{FullWidth: true}
	}

	isKerning := mask >= 0
	isUniversal := mask&BitLayoutUniversal != 0

	mode := LayoutMode{
		// Kerning is set if Bit 0 is active
		Kerning: isKerning,

		// Universal Smushing is set if Bit 1 is active (this is often the 'Smushing' flag)
		Universal: isUniversal,

		// Smushing is enabled if Universal Smushing is active
		Smushing: isUniversal,

		// FullWidth is true only if neither Kerning nor Universal Smushing bits are set
		// FullWidth: !isKerning && !isUniversal,

		// Vertical Modes (using dedicated higher bits)
		VKerning:  mask&BitLayoutVKerning != 0,
		VSmushing: mask&BitLayoutVSmushing != 0,
	}

	return mode
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

func readCharacter(scanner *bufio.Scanner, height int) (Glyph, error) {
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
