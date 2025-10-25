package fig

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFont(t *testing.T) {
	_, err := loadFont("demo")
	assert.Error(t, err)

	font, err := loadFont("standard")
	assert.Nil(t, err)
	assert.Equal(t, font.name, "standard")
}

func TestNewParseHeader(t *testing.T) {
	invalid := "flf2$ 6 5 20 15 13"
	_, err := newParser(invalid)
	assert.Error(t, err)

	invalid = "flf2a$ 6 5 20"
	_, err = newParser(invalid)
	assert.Error(t, err)

	given := "flf2a$ 6 5 20 15 13"
	_, err = newParser(given)
	assert.Nil(t, err)
}

func TestParseHeader(t *testing.T) {
	given := "flf2a$ 6 5 20 15 A"
	parser, _ := newParser(given)

	assert.Equal(t, parser.parseInt(1, "height"), 6)
	assert.Equal(t, parser.parseInt(2, "baseline"), 5)
	assert.Equal(t, parser.parseInt(3, "max_length"), 20)
	assert.Equal(t, parser.parseInt(4, "old_layout"), 15)
	assert.Equal(t, parser.parseInt(5, "comment_lines"), 0)
	assert.Error(t, parser.err)
}

func TestParseGlyph(t *testing.T) {
	input := `  __ #
 /  |#
/_/ |#
  | |#
  | |#
  |_|#
     ##`
	expected := strings.ReplaceAll(input, "#", "")
	scanner := bufio.NewScanner(strings.NewReader(input))
	g, err := readCharacter(scanner, 7)
	assert.Equal(t, expected, strings.Join(g.lines, "\n"))
	assert.Nil(t, err)

	input = `     _    @
    / \   @
   / _ \  @
  / ___ \ @
 /_/   \_\@
          @@`
	expected = strings.ReplaceAll(input, "@", "")
	scanner = bufio.NewScanner(strings.NewReader(input))
	g, err = readCharacter(scanner, 6)
	assert.Equal(t, expected, strings.Join(g.lines, "\n"))
	assert.Nil(t, err)
}
