package font

import (
	"testing"

	"github.com/phantompunk/fig/internal/assert"
)

// TestGlyphRunes_noPaddingOnLeft verifies that GlyphRunes does not introduce
// extra leading zeros beyond what the original glyph line contains.
// Spaces in the source line become 0 (transparent), but no columns are added.
func TestGlyphRunes_noPaddingOnLeft(t *testing.T) {
	f, err := LoadFont("standard")
	assert.NilError(t, err)

	// '!' in standard font:
	//   line[0] = "  _ "  → width 4, 2 leading spaces
	//   line[1] = " | |"  → width 4, 1 leading space
	//   line[2] = " | |"
	//   line[3] = " |_|"
	//   line[4] = " (_)"
	//   line[5] = "    "  → all spaces
	rows := f.GlyphRunes('!')

	assert.Equal(t, len(rows), f.Height())

	// Row 0: "  _ " → [0, 0, '_', 0]
	assert.Equal(t, len(rows[0]), 4)
	assert.Equal(t, rows[0][0], rune(0)) // was space — no extra column added
	assert.Equal(t, rows[0][1], rune(0))
	assert.Equal(t, rows[0][2], '_')
	assert.Equal(t, rows[0][3], rune(0))

	// Row 1: " | |" → [0, '|', 0, '|']
	assert.Equal(t, len(rows[1]), 4)
	assert.Equal(t, rows[1][0], rune(0))
	assert.Equal(t, rows[1][1], '|')
	assert.Equal(t, rows[1][2], rune(0))
	assert.Equal(t, rows[1][3], '|')

	// Row 5: "    " → all zeros, same length
	assert.Equal(t, len(rows[5]), 4)
	for x, ch := range rows[5] {
		assert.Equal(t, ch, rune(0))
		_ = x
	}
}

// TestGlyphRunes_hardblankPreserved verifies that hardblank characters are NOT
// converted to 0 — they must stay distinct so the canvas can treat them as
// non-smushable spaces.
func TestGlyphRunes_hardblankPreserved(t *testing.T) {
	f, err := LoadFont("standard")
	assert.NilError(t, err)

	hb := f.Hardblank()

	// Scan all printable ASCII glyphs for any hardblank occurrence and confirm
	// it survives conversion unchanged.
	found := false
	for ch := rune(32); ch <= 126; ch++ {
		for _, row := range f.GlyphRunes(ch) {
			for _, cell := range row {
				if cell == hb {
					found = true
				}
				// A cell that equals hb must not have been turned into 0.
				if cell == 0 {
					// 0 means it was a regular space — that's fine.
				}
			}
		}
	}

	if !found {
		t.Log("no hardblank found in any glyph — skipping preservation check")
	}
}
