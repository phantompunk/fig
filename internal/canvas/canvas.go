package canvas

import (
	"strings"

	fig "github.com/phantompunk/fig/internal/font"
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
			if cx >= len(c.cells[y]) {
				break
			}
			c.cells[y][xOffset+x] = ch
		}
	}
}

func (c *Canvas) FindOverlap(glyph [][]rune, maxOverlap int, rules []fig.SmushRule, hb rune) int {
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

	for overlap := maxOverlap; overlap >= 1; overlap-- {
		if overlap > rightEdge {
			continue
		}
		canSmush := true
	rowLoop:
		for rowIdx, glyphRow := range glyph {
			if rowIdx >= c.height {
				break
			}
			for i := 0; i < overlap; i++ {
				canvasX := rightEdge - overlap + i
				if i >= len(glyphRow) {
					continue
				}
				left := c.cells[rowIdx][canvasX]
				right := glyphRow[i]
				if left == 0 || right == 0 {
					continue
				}
				if _, ok := smushCell(left, right, rules, hb); !ok {
					canSmush = false
					break rowLoop
				}
			}
		}
		if canSmush {
			return overlap
		}
	}
	return 0
}

func (c *Canvas) StampSmush(glyph [][]rune, overlap int, rules []fig.SmushRule, hb rune) {
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

func (c *Canvas) String(hb rune) string {
	var sb strings.Builder
	for _, row := range c.cells {
		// Find last non-zero cell to trim trailing empty cells.
		end := len(row)
		for end > 0 && row[end-1] == 0 {
			end--
		}
		for _, ch := range row[:end] {
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

func smushCell(l, r rune, rules []fig.SmushRule, hb rune) (rune, bool) {
	if l != hb && r != hb {
		if r == ' ' {
			return l, true
		}
		if l == ' ' {
			return r, true
		}
	}

	for _,rule := range rules {
		if result := rule(l, r); result.Ok {
			return result.Char, true
		}
	}
	return 0, false
}

