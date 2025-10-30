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
	lines := make([]string, r.font.metadata.height)

	var prev *Glyph
	for _, char := range text {
		g := r.font.getGlyph(char)

		if prev == nil {
			prev = r.leftFlush(&g)
			copy(lines, prev.lines)
			prev = &g
			continue
		}

		overlap := r.computeOverlap(prev, &g)
		lines = r.merge(lines, &g, overlap)
		prev = &g
	}
	return r.linesToString(lines)
}

func (r *Renderer) computeOverlap(left, right *Glyph) int {
	minDist := r.font.metadata.maxLength

	for row := range len(left.lines) {
		leftRunes := []rune(left.lines[row])
		rightRunes := []rune(right.lines[row])

		dist := 0
		ri := 0
		li := len(leftRunes) - 1
		for li > 0 && leftRunes[li] == ' ' {
			li--
		}
		for ri < len(rightRunes) && rightRunes[ri] == ' ' {
			ri++
		}

		if li < 0 || ri >= len(rightRunes) {
			continue
		}

		for {
			leftPos := li - dist
			rightPos := ri + dist

			if leftPos < 0 || rightPos >= len(rightRunes) {
				break
			}

			lChar := leftRunes[leftPos]
			rChar := rightRunes[rightPos]

			// fmt.Printf("Dist: %d, Row %d: %q + %q, Chars: %q + %q, Sms: %v\n", dist, row, left.lines[row], right.lines[row], string(lChar), string(rChar), r.isSmushable(lChar, rChar))
			if !r.isSmushable(lChar, rChar) {
				break
			}
			dist++
		}

		if dist < minDist {
			minDist = dist
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

		for j := range overlap {
			li := len(leftRunes) - overlap + j
			ri := j

			if li < 0 || ri >= len(rightRunes) {
				continue
			}

			res := r.applyRules(leftRunes[li], rightRunes[ri])
			if res.allowed {
				leftRunes[li] = res.char
			}
		}
		result[i] = string(leftRunes) + string(rightRunes[overlap:])
	}
	return result
}

func (r *Renderer) isSmushable(a, b rune) bool {
	if a == ' ' || b == ' ' {
		return true
	}
	result := r.applyRules(a, b)
	return result.allowed
}

func (r *Renderer) applyRules(a, b rune) SmushResult {
	for _, rule := range r.font.rules {
		if res := rule(a, b); res.allowed {
			return res
		}
	}
	return SmushResult{char: 0, allowed: false}
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

	for i, row := range glyph.lines {
		glyph.lines[i] = row[flush:]
	}
	glyph.width -= flush
	return glyph
}

func (r *Renderer) linesToString(lines []string) string {
	for i, line := range lines {
		lines[i] = strings.ReplaceAll(line, string(r.font.metadata.hardBlank), " ")
	}
	return strings.Join(lines, "\n")
}

type SmushResult struct {
	char    rune
	allowed bool
}

type SmushRule func(a, b rune) SmushResult

func EqualChars(a, b rune) SmushResult {
	if a == b && a != '$' {
		return SmushResult{char: a, allowed: true}
	}
	return SmushResult{char: 0, allowed: false}
}

func BigX(a, b rune) SmushResult {
	if a == '/' && b == '\\' {
		return SmushResult{char: '|', allowed: true}
	} else if a == '\\' && b == '/' {
		return SmushResult{char: 'Y', allowed: true}
	} else if a == '>' && b == '<' {
		return SmushResult{char: 'X', allowed: true}
	}
	return SmushResult{char: 0, allowed: false}
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
		return SmushResult{char: a, allowed: true}
	} else if _, found := pairs[b]; found && a == '_' {
		return SmushResult{char: b, allowed: true}
	}
	return SmushResult{char: 0, allowed: false}
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
		return SmushResult{char: '|', allowed: true}
	}
	return SmushResult{char: 0, allowed: false}
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
		return SmushResult{char: a, allowed: true}
	} else if ch[a] < ch[b] {
		return SmushResult{char: b, allowed: true}
	}
	return SmushResult{char: 0, allowed: false}
}

func Hardblank(a, b rune) SmushResult {
	if a == '$' || b == '$' {
		return SmushResult{char: a, allowed: true}
	}
	return SmushResult{char: 0, allowed: false}
}
