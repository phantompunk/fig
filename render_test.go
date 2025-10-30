package fig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanSmushEqualChars(t *testing.T) {
	testcases := []struct {
		name     string
		left     rune
		right    rune
		expected bool
	}{
		{"Same letters", 'X', 'X', true},
		{"Different letters", 'X', 'Y', false},
		{"Same numbers", '5', '5', true},
		{"Different numbers", '5', '4', false},
	}

	renderer := NewSmushRenderer(EqualChars)
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, renderer.isSmushable(tc.left, tc.right))
		})
	}
}

func TestCanSmushUnderscores(t *testing.T) {
	testcases := []struct {
		name     string
		left     rune
		right    rune
		expected bool
	}{
		{"underscore with pipe", '_', '|', true},
		{"underscore with brace", '_', '{', true},
		{"underscore with angle", '<', '_', true},
		{"underscore with letter", 'l', '_', false},
		{"two underscores", '_', '_', false},
	}

	renderer := NewSmushRenderer(Underscore)
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, renderer.isSmushable(tc.left, tc.right))
		})
	}
}

// A hierarchy of six classes is used: "|", "/\", "[]", "{}", "()", and "<>". When two smushing sub-characters are from different classes, the one from the latter class will be used.
func TestCanSmushHierarchy(t *testing.T) {
	testcases := []struct {
		name     string
		left     rune
		right    rune
		expected bool
	}{
		{"pipe vs slash", '|', '/', true},
		{"bracket vs pipe", '[', '|', true},
		{"paren vs bracket", '(', '[', true},
		{"same hierarchy", '[', ']', true},
		{"non-hierarchy chars", 'A', 'B', false},
	}

	renderer := NewSmushRenderer(Heirarchy)
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, renderer.isSmushable(tc.left, tc.right))
		})
	}
}

func TestCanSmushOpposite(t *testing.T) {
	testcases := []struct {
		name     string
		left     rune
		right    rune
		expected bool
	}{
		{"square brackets", '[', ']', true},
		{"reverse square brackets", ']', '[', true},
		{"curly braces", '{', '}', true},
		{"parentheses", '(', ')', true},
		{"same bracket", '[', '[', false},
		{"non-brackets", 'A', 'B', false},
	}

	renderer := NewSmushRenderer(OppositePair)
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, renderer.isSmushable(tc.left, tc.right))
		})
	}
}

func TestCanSmushBigX(t *testing.T) {
	testcases := []struct {
		name     string
		left     rune
		right    rune
		expected bool
	}{
		{"forward slash + backslash", '/', '\\', true},
		{"backslash + forward slash", '\\', '/', true},
		{"two forward slashes", '/', '/', false},
		{"two backslashes", '\\', '\\', false},
		{"non-slashes", 'X', 'Y', false},
	}

	renderer := NewSmushRenderer(BigX)
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, renderer.isSmushable(tc.left, tc.right))
		})
	}
}

func TestCanSmushHardblank(t *testing.T) {
	testcases := []struct {
		name     string
		left     rune
		right    rune
		expected bool
	}{
		{"pipe vs slash", '|', '/', true},
		{"bracket vs pipe", '[', '|', true},
		{"paren vs bracket", '(', '[', true},
		{"same hierarchy", '[', ']', true},
		{"non-hierarchy chars", 'A', 'B', false},
	}

	renderer := NewSmushRenderer(Heirarchy)
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, renderer.isSmushable(tc.left, tc.right))
		})
	}
}

func TestSmushModes(t *testing.T) {
	testcases := []struct {
		name   string
		input  []rune
		result SmushResult
		rule   SmushRule
	}{
		{
			name:   "EqualChars match",
			input:  []rune{'|', '|'},
			result: SmushResult{'|', true},
			rule:   EqualChars,
		},
		{
			name:   "EqualChars hardblank",
			input:  []rune{'$', '$'},
			result: SmushResult{0, false},
			rule:   EqualChars,
		},
		{
			name:   "BigX |",
			input:  []rune{'/', '\\'},
			result: SmushResult{'|', true},
			rule:   BigX,
		},
		{
			name:   "BigX Y",
			input:  []rune{'\\', '/'},
			result: SmushResult{'Y', true},
			rule:   BigX,
		},
		{
			name:   "BigX X",
			input:  []rune{'>', '<'},
			result: SmushResult{'X', true},
			rule:   BigX,
		},
		{
			name:   "Hardblank",
			input:  []rune{'$', '$'},
			result: SmushResult{'$', true},
			rule:   Hardblank,
		},
		{
			name:   "OppositePair {}",
			input:  []rune{'{', '}'},
			result: SmushResult{'|', true},
			rule:   OppositePair,
		},
		{
			name:   "OppositePair )(",
			input:  []rune{')', '('},
			result: SmushResult{'|', true},
			rule:   OppositePair,
		},
		{
			name:   "Underscore _{",
			input:  []rune{'_', '{'},
			result: SmushResult{'{', true},
			rule:   Underscore,
		},
		{
			name:   "hierachy >|",
			input:  []rune{'>', '|'},
			result: SmushResult{'>', true},
			rule:   Heirarchy,
		},
		{
			name:   "hierachy {)",
			input:  []rune{'{', ')'},
			result: SmushResult{')', true},
			rule:   Heirarchy,
		},
		{
			name:   "whitespace )",
			input:  []rune{' ', ')'},
			result: SmushResult{')', true},
			rule:   Heirarchy,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.rule(tc.input[0], tc.input[1])
			assert.Equal(t, tc.result, result)
		})
	}
}

func TestIsSmushable(t *testing.T) {
	testcases := []struct {
		name     string
		fontName string
		input    []rune
		expected bool
	}{
		{
			name:     "3x5 AB",
			fontName: "3x5",
			input:    []rune{'A', 'B'},
			expected: false,
		},
		{
			name:     "standard HH",
			fontName: "standard",
			input:    []rune{'|', '_'},
			expected: true,
		},
		{
			name:     "standard space",
			fontName: "standard",
			input:    []rune{' ', 'H'},
			expected: true,
		},
		{
			name:     "standard $",
			fontName: "standard",
			input:    []rune{' ', '$'},
			expected: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			font, _ := loadFont(tc.fontName)
			render := New(font)
			assert.Equal(t, tc.expected, render.isSmushable(tc.input[0], tc.input[1]))
		})
	}
}

func TestStandardSmushable(t *testing.T) {
	testcases := []struct {
		name     string
		fontName string
		input    []rune
		expected bool
	}{
		{
			name:     "standard HH",
			fontName: "standard",
			input:    []rune{'|', '_'},
			expected: true,
		},
		{
			name:     "standard space",
			fontName: "standard",
			input:    []rune{' ', 'H'},
			expected: true,
		},
		{
			name:     "standard spaces",
			fontName: "standard",
			input:    []rune{' ', ' '},
			expected: true,
		},
		{
			name:     "standard $",
			fontName: "standard",
			input:    []rune{' ', '$'},
			expected: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			font, _ := loadFont(tc.fontName)
			render := New(font)
			assert.Equal(t, tc.expected, render.isSmushable(tc.input[0], tc.input[1]))
		})
	}
}

func TestOverlap(t *testing.T) {
	testcases := []struct {
		name     string
		fontName string
		input    []rune
		expected int
	}{
		// {
		// 	name:     "3x5 AB",
		// 	fontName: "3x5",
		// 	input:    []rune{'A', 'B'},
		// 	expected: 0,
		// },
		{
			name:     "standard Hspace",
			fontName: "standard",
			input:    []rune{'H', ' '},
			expected: 2,
		},
		{
			name:     "standard HH",
			fontName: "standard",
			input:    []rune{'H', 'i'},
			expected: 2,
		},
		// {
		// 	name:     "standard space",
		// 	fontName: "standard",
		// 	input:    []rune{' ', 'H'},
		// 	expected: 1,
		// },
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			font, _ := loadFont(tc.fontName)
			render := New(font)
			left := font.getGlyph(tc.input[0])
			right := font.getGlyph(tc.input[1])
			assert.Equal(t, tc.expected, render.computeOverlap(&left, &right))
		})
	}
}

func NewSmushRenderer(mode SmushRule) *Renderer {
	font := &FigFont{
		metadata: Metadata{
			hardBlank: '$',
		},
		rules: []SmushRule{mode},
	}
	return &Renderer{font: font}
}
