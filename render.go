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
			continue
		}

		overlap := r.calculateOverlap(prev, &g)
		lines = r.merge(lines, &g, overlap)
		prev = &g
	}
	return r.linesToString(lines)
}

func (r *Renderer) calculateOverlap(left, right *Glyph) int {
	width := r.font.metadata.maxLength

	for i := 1; i < width; i++ {
		if r.canOverlap(left, right, i) {
			return i - 1
		}
	}
	return width - 1
}

func (r *Renderer) canOverlap(left, right *Glyph, offset int) bool {
	for i := range r.font.metadata.height {
		ll := []rune(left.lines[i])
		rl := []rune(right.lines[i])

		for j := range offset {
			li := len(ll) - offset + j
			ri := j

			if li < 0 || ri >= len(rl) {
				continue
			}

			if r.isSmushable(ll[li], rl[ri]) {
				continue
			}

			if ll[li] != ' ' && rl[ri] != ' ' {
				return false
			}
		}
	}
	return true
}

func (r *Renderer) merge(lines []string, right *Glyph, overlap int) []string {
	result := make([]string, r.font.metadata.height)

	for i := range r.font.metadata.height {
		left := []rune(lines[i])
		right := []rune(right.lines[i])

		for j := range overlap {
			li := len(left) - overlap + j
			ri := j

			if li < 0 || ri >= len(right) {
				continue
			}

			res := r.applyRules(left[li], right[ri])
			if res.found {
				left[li] = res.char
			}
		}
		result[i] = string(left) + string(right[overlap:])
	}
	return result
}

func (r *Renderer) isSmushable(a, b rune) bool {
	if a == ' ' || b == ' ' {
		return true
	}
	result := r.applyRules(a, b)
	return result.found
}

func (r *Renderer) applyRules(a, b rune) SmushResult {
	for _, rule := range r.font.rules {
		if res := rule(a, b); res.found {
			return res
		}
	}
	return SmushResult{char: 0, found: false}
}

func (r *Renderer) leftFlush(glyph *Glyph) *Glyph {
	var flush int = glyph.width
	for _, row := range glyph.lines {
		for i, char := range row {
			if char != ' ' && char != r.font.metadata.hardBlank {
				// fmt.Printf("Found visible char %s at %d, %d\n", string(char), j, i)
				flush = min(i, flush)
				break
			}
		}
	}

	for i, row := range glyph.lines {
		glyph.lines[i] = row[flush:]
	}
	return glyph
}

func (r *Renderer) linesToString(lines []string) string {
	for i, line := range lines {
		lines[i] = strings.ReplaceAll(line, string(r.font.metadata.hardBlank), " ")
	}
	return strings.Join(lines, "\n")
}

type SmushResult struct {
	char  rune
	found bool
}

type SmushRule func(a, b rune) SmushResult

func EqualChars(a, b rune) SmushResult {
	if a == b && a != '$' {
		return SmushResult{char: a, found: true}
	}
	return SmushResult{char: 0, found: false}
}

func BigX(a, b rune) SmushResult {
	if a == '/' && b == '\\' {
		return SmushResult{char: '|', found: true}
	} else if a == '\\' && b == '/' {
		return SmushResult{char: 'Y', found: true}
	} else if a == '>' && b == '<' {
		return SmushResult{char: 'X', found: true}
	}
	return SmushResult{char: 0, found: false}
}

func Underscore(a, b rune) SmushResult {
	pairs := map[rune]struct{}{
		'[': {},
		']': {},
		'{': {},
		'}': {},
		'(': {},
		')': {},
	}

	if _, found := pairs[a]; found && b == '_' {
		return SmushResult{char: a, found: true}
	} else if _, found := pairs[b]; found && a == '_' {
		return SmushResult{char: b, found: true}
	}
	return SmushResult{char: 0, found: false}
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
		return SmushResult{char: '|', found: true}
	}
	return SmushResult{char: 0, found: false}
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
		return SmushResult{char: a, found: true}
	} else if ch[a] < ch[b] {
		return SmushResult{char: b, found: true}
	}
	return SmushResult{char: 0, found: false}
}

func Hardblank(a, b rune) SmushResult {
	if a == '$' && b == '$' {
		return SmushResult{char: a, found: true}
	}
	return SmushResult{char: 0, found: false}
}
