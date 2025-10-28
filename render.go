package fig

import (
	"fmt"
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
	letters := make([]Glyph, len(text))

	for i, char := range text {
		letters[i] = r.font.glyphs[char]
	}

	var result strings.Builder
	for h := range r.font.metadata.height {
		for _, l := range letters {
			result.WriteString(l.lines[h])
		}
		if h < r.font.metadata.height-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

func (r *Renderer) RenderExp(text string) string {
	lines := make([]string, r.font.metadata.height)

	var prev *Glyph
	for _, char := range text {
		g := r.font.getGlyph(char)

		if prev == nil {
			copy(lines, g.lines)
			prev = &g
			continue
		}

		overlap := r.calculateOverlap(prev, &g)
		lines = r.merge(lines, &g, overlap)
		prev = &g
	}
	return strings.Join(lines, "\n")
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

		if len(ll) < offset {
			continue
		}

		if r.isSmushable(ll[len(ll)-offset], rl[0]) {
			continue
		}

		if ll[len(ll)-offset] != ' ' && rl[0] != ' ' {
			return false
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
		fmt.Println("Applying rule:", rule)
		if res := rule(a, b); res.found {
			return res
		}
	}
	return SmushResult{char: 0, found: false}
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
