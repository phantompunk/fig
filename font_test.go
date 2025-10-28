package fig

import (
	"testing"

	"github.com/phantompunk/fig/assert"
)

func TestNewFonts(t *testing.T) {
	// header flf2a$ 6 5 16 15 11 0 24463
	standard, err := Font("standard")
	meta := standard.metadata

	assert.NilError(t, err)
	assert.Equal(t, standard.name, "standard")
	assert.Equal(t, meta.hardBlank, '$')
	assert.Equal(t, meta.height, 6)
	assert.Equal(t, meta.baseline, 5)
	assert.Equal(t, meta.maxLength, 16)
	assert.Equal(t, meta.oldLayout, 15)
	assert.Equal(t, meta.commentLines, 11)
	assert.Equal(t, meta.printDirection, 0)
	assert.Equal(t, meta.fullLayout, 24463)
	assert.Equal(t, meta.codeTag, 0)
}

func TestRenderFonts(t *testing.T) {
	standard, err := Font("standard")
	char := standard.glyphs['!']

	assert.NilError(t, err)
	assert.Equal(t, len(char.lines), 6)
	assert.Equal(t, char.width, 4)
	assert.Equal(t, char.lines[0], "  _ ")
	assert.Equal(t, char.lines[1], " | |")
	assert.Equal(t, char.lines[2], " | |")
	assert.Equal(t, char.lines[3], " |_|")
	assert.Equal(t, char.lines[4], " (_)")
	assert.Equal(t, char.lines[5], "    ")
}
