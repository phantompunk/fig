package render

import (
	"strings"

	"github.com/phantompunk/fig/internal/font"
)

type Canvas struct {
	height int
	cells  [][]rune
}

func NewCanvas(height, width int) *Canvas {
	cells := make([][]rune, height)
	for i := range cells {
		cells[i] = make([]rune, width)
	}
	return &Canvas{height: height, cells: cells}
}

func (c *Canvas) Stamp(glyph [][]rune, xOffset int) {
	for y, row := range glyph {
		for x, ch := range row {
			cx := xOffset + x
			if cx < 0 {
				continue
			}
			if cx >= len(c.cells[y]) {
				break
			}
			c.cells[y][cx] = ch
		}
	}
}

func (c *Canvas) FindOverlap(glyph [][]rune, maxOverlap int, rules []font.SmushRule, hb rune) int {
	// Find the global right edge of the canvas (max across all rows).
	rightEdge := 0
	for _, row := range c.cells {
		for x := len(row) - 1; x >= 0; x-- {
			if row[x] != 0 {
				if x+1 > rightEdge {
					rightEdge = x + 1
				}
				break
			}
		}
	}

	// FIGlet gap algorithm: for each row compute the overlap as
	//   trailing_empty(canvas_row) + leading_empty(glyph_row) [+1 if touching chars smush]
	// then take the minimum across all rows.
	overlap := maxOverlap
	for rowIdx, glyphRow := range glyph {
		if rowIdx >= c.height {
			break
		}
		canvasRow := c.cells[rowIdx]

		// Rightmost non-zero in this canvas row.
		canvasRight := 0
		for x := len(canvasRow) - 1; x >= 0; x-- {
			if canvasRow[x] != 0 {
				canvasRight = x + 1
				break
			}
		}

		// Trailing empty columns relative to the global right edge.
		trailingEmpty := rightEdge - canvasRight

		// Leftmost non-zero in the glyph row.
		glyphLeft := len(glyphRow) // sentinel: row is all empty
		for x, ch := range glyphRow {
			if ch != 0 {
				glyphLeft = x
				break
			}
		}

		gap := trailingEmpty + glyphLeft

		// One extra column of overlap when the touching chars can smush.
		if canvasRight > 0 && glyphLeft < len(glyphRow) {
			if _, ok := smushCell(canvasRow[canvasRight-1], glyphRow[glyphLeft], rules, hb); ok {
				gap++
			}
		}

		if gap < overlap {
			overlap = gap
		}
	}

	return overlap
}

func (c *Canvas) StampSmush(glyph [][]rune, overlap int, rules []font.SmushRule, hb rune) {
	// Find the rightmost non-zero column across all rows.
	rightEdge := 0
	for _, row := range c.cells {
		for x := len(row) - 1; x >= 0; x-- {
			if row[x] != 0 {
				if x+1 > rightEdge {
					rightEdge = x + 1
				}
				break
			}
		}
	}

	if overlap == 0 {
		c.Stamp(glyph, rightEdge)
		return
	}

	startX := rightEdge - overlap

	// First pass: check whether every cell in the overlap zone can be smushed.
	for rowIdx, glyphRow := range glyph {
		if rowIdx >= c.height {
			break
		}
		for i := 0; i < overlap && i < len(glyphRow); i++ {
			canvasX := startX + i
			if canvasX < 0 || canvasX >= len(c.cells[rowIdx]) {
				continue
			}
			left := c.cells[rowIdx][canvasX]
			right := glyphRow[i]
			if left == 0 || right == 0 {
				continue // one side is empty, no conflict
			}
			if _, ok := smushCell(left, right, rules, hb); !ok {
				// Smush failed — fall back to placing glyph right after existing content.
				c.Stamp(glyph, rightEdge)
				return
			}
		}
	}

	// Second pass: apply smushed overlap zone and write the non-overlap tail.
	for rowIdx, glyphRow := range glyph {
		if rowIdx >= c.height {
			break
		}
		for i, ch := range glyphRow {
			canvasX := startX + i
			if canvasX < 0 || canvasX >= len(c.cells[rowIdx]) {
				break
			}
			if i < overlap {
				left := c.cells[rowIdx][canvasX]
				if left == 0 {
					c.cells[rowIdx][canvasX] = ch
				} else if ch != 0 {
					result, _ := smushCell(left, ch, rules, hb)
					c.cells[rowIdx][canvasX] = result
				}
			} else {
				c.cells[rowIdx][canvasX] = ch
			}
		}
	}
}

func (c *Canvas) String(hb rune, minWidth int) string {
	// Find the widest content column across all rows so every row is the same
	// width, matching figlet's fixed-width output. minWidth sets a floor (used
	// for full-width fonts where trailing glyph whitespace must be preserved).
	maxEnd := minWidth
	for _, row := range c.cells {
		end := len(row)
		for end > 0 && row[end-1] == 0 {
			end--
		}
		if end > maxEnd {
			maxEnd = end
		}
	}

	var sb strings.Builder
	for _, row := range c.cells {
		for i := 0; i < maxEnd; i++ {
			var ch rune
			if i < len(row) {
				ch = row[i]
			}
			if ch == 0 || ch == hb {
				sb.WriteRune(' ')
			} else {
				sb.WriteRune(ch)
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func smushCell(l, r rune, rules []font.SmushRule, hb rune) (rune, bool) {
	if l != hb && r != hb {
		if r == ' ' {
			return l, true
		}
		if l == ' ' {
			return r, true
		}
	}

	for _, rule := range rules {
		if result := rule(l, r); result.Ok {
			return result.Char, true
		}
	}
	return 0, false
}
