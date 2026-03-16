package fig

import (
	"strings"
)

type Renderer struct {
	font *FigFont
}

type RenderOptions struct {
	Justify Justification
}

type Justification int

const (
	JustifyLeft Justification = iota
	JustifyCenter
	JustifyRight
)

func New(font *FigFont) *Renderer {
	return &Renderer{font: font}
}

func (r *Renderer) Render(text string) string {
	return r.linesToString(r.render(text))
}

func (r *Renderer) Lines(text string) []string {
	return r.render(text)
}

func (r *Renderer) render(text string) []string {
	lines := make([]string, r.font.metadata.height)

	first := true
	for _, char := range text {
		g := r.font.getGlyph(char)

		if first {
			g = *r.leftFlush(&g)
			copy(lines, g.lines)
			first = false
			continue
		}

		overlap := r.computeOverlap(lines, &g)
		lines = r.merge(lines, &g, overlap)
	}
	return lines
}

func (r *Renderer) computeOverlap(leftLines []string, right *Glyph) int {
	if r.font.metadata.layoutMode.FullWidth {
		return 0
	}

	minDist := r.font.metadata.maxLength

	for row := range len(leftLines) {
		leftRunes := []rune(leftLines[row])
		rightRunes := []rune(right.lines[row])

		leftEnd := len(leftRunes)
		for leftEnd > 0 && leftRunes[leftEnd-1] == ' ' {
			leftEnd--
		}
		rightStart := 0
		for rightStart < len(rightRunes) && rightRunes[rightStart] == ' ' {
			rightStart++
		}

		curDist := (len(leftRunes) - leftEnd) + rightStart
		if leftEnd > 0 && rightStart < len(rightRunes) {
			if r.isSmushable(leftRunes[leftEnd-1], rightRunes[rightStart]) {
				curDist++
			}
		}

		if curDist < minDist {
			minDist = curDist
		}
	}

	if minDist < 0 {
		minDist = 0
	}

	return minDist
}

func (r *Renderer) merge(lines []string, right *Glyph, overlap int) []string {
	result := make([]string, r.font.metadata.height)

	for i := range r.font.metadata.height {
		leftRunes := []rune(lines[i])
		rightRunes := []rune(right.lines[i])

		// Cap overlap to the length of rightRunes to prevent out-of-bounds slice
		safeOverlap := overlap
		if safeOverlap < 0 {
			safeOverlap = 0
		}
		if safeOverlap > len(rightRunes) {
			safeOverlap = len(rightRunes)
		}

		for j := range safeOverlap {
			li := len(leftRunes) - safeOverlap + j
			ri := j

			if li < 0 || ri >= len(rightRunes) {
				continue
			}

			lChar := leftRunes[li]
			rChar := rightRunes[ri]
			if lChar == ' ' {
				leftRunes[li] = rChar
			} else if rChar != ' ' {
				res := r.applyRules(lChar, rChar)
				if res.Ok {
					leftRunes[li] = res.Char
				}
			}
		}
		result[i] = string(leftRunes) + string(rightRunes[safeOverlap:])
	}
	return result
}

func (r *Renderer) isSmushable(a, b rune) bool {
	return r.applyRules(a, b).Ok
}

func (r *Renderer) applyRules(a, b rune) SmushResult {
	for _, rule := range r.font.rules {
		if res := rule(a, b); res.Ok {
			return res
		}
	}
	return SmushResult{Char: 0, Ok: false}
}

func (r *Renderer) leftFlush(glyph *Glyph) *Glyph {
	var flush int = glyph.width
	for _, row := range glyph.lines {
		for i, char := range row {
			if char != ' ' && char != r.font.metadata.hardBlank {
				flush = min(i, flush)
				break
			}
		}
	}

	// Ensure flush is within valid bounds
	if flush < 0 {
		flush = 0
	}

	for i, row := range glyph.lines {
		// Cap flush to the length of the row to prevent out-of-bounds slice
		safeFlush := flush
		if safeFlush > len(row) {
			safeFlush = len(row)
		}
		glyph.lines[i] = row[safeFlush:]
	}
	glyph.width -= flush
	return glyph
}

func (r *Renderer) linesToString(lines []string) string {
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		row := strings.ReplaceAll(line, string(r.font.metadata.hardBlank), " ")
		result = append(result, row)
	}
	return strings.Join(result, "\n") + "\n"
}

type SmushResult struct {
	Char rune
	Ok   bool
}

type SmushRule func(a, b rune) SmushResult

func EqualCharsRule(hardblank rune) SmushRule {
	return func(a, b rune) SmushResult {
		if a == b && a != hardblank {
			return SmushResult{Char: a, Ok: true}
		}
		return SmushResult{Char: 0, Ok: false}
	}
}

// EqualChars is the legacy fixed-sentinel version kept for direct test use.
func EqualChars(a, b rune) SmushResult {
	return EqualCharsRule('$')(a, b)
}

func BigX(a, b rune) SmushResult {
	if a == '/' && b == '\\' {
		return SmushResult{Char: '|', Ok: true}
	} else if a == '\\' && b == '/' {
		return SmushResult{Char: 'Y', Ok: true}
	} else if a == '>' && b == '<' {
		return SmushResult{Char: 'X', Ok: true}
	}
	return SmushResult{Char: 0, Ok: false}
}

func Underscore(a, b rune) SmushResult {
	pairs := map[rune]struct{}{
		'[':  {},
		']':  {},
		'{':  {},
		'}':  {},
		'(':  {},
		')':  {},
		'|':  {},
		'/':  {},
		'\\': {},
		'<':  {},
		'>':  {},
	}

	if _, found := pairs[a]; found && b == '_' {
		return SmushResult{Char: a, Ok: true}
	} else if _, found := pairs[b]; found && a == '_' {
		return SmushResult{Char: b, Ok: true}
	}
	return SmushResult{Char: 0, Ok: false}
}

func OppositePair(a, b rune) SmushResult {
	pairs := map[rune]rune{
		'[': ']',
		']': '[',
		'{': '}',
		'}': '{',
		'(': ')',
		')': '(',
	}
	if pairs[a] == b {
		return SmushResult{Char: '|', Ok: true}
	}
	return SmushResult{Char: 0, Ok: false}
}

func Heirarchy(a, b rune) SmushResult {
	ch := map[rune]int{
		'|':  1,
		'/':  2,
		'\\': 2,
		'[':  3,
		']':  3,
		'{':  4,
		'}':  4,
		'(':  5,
		')':  5,
		'<':  6,
		'>':  6,
	}

	if ch[a] > ch[b] {
		return SmushResult{Char: a, Ok: true}
	} else if ch[a] < ch[b] {
		return SmushResult{Char: b, Ok: true}
	}
	return SmushResult{Char: 0, Ok: false}
}

func HardblankRule(hardblank rune) SmushRule {
	return func(a, b rune) SmushResult {
		if a == hardblank && b == hardblank {
			return SmushResult{Char: hardblank, Ok: true}
		}
		return SmushResult{Char: 0, Ok: false}
	}
}

// Hardblank is the legacy fixed-sentinel version kept for direct test use.
func Hardblank(a, b rune) SmushResult {
	return HardblankRule('$')(a, b)
}
