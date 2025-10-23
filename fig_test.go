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

	val := standard.glyphs['!']

	assert.NilError(t, err)
	assert.Equal(t, len(val.lines), 6)
	assert.Equal(t, val.width, 4)
	assert.Equal(t, val.lines[0], "  _ ")
	assert.Equal(t, val.lines[1], " | |")
	assert.Equal(t, val.lines[2], " | |")
	assert.Equal(t, val.lines[3], " |_|")
	assert.Equal(t, val.lines[4], " (_)")
	assert.Equal(t, val.lines[5], "    ")
}

func TestRender(t *testing.T) {
	standard, _ := Font("standard")
	expected := `     _    
    / \   
   / _ \  
  / ___ \ 
 /_/   \_\
          `

	actual := standard.Render("A")
	assert.Equal(t, actual, expected)

	expected = `    _    ____   ____ 
   / \  | __ ) / ___|
  / _ \ |  _ \| |    
 / ___ \| |_) | |___ 
/_/   \_\____/ \____|
                     `
	actual = standard.Render("ABC")
	assert.Equal(t, actual, expected)
}
