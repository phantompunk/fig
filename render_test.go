package fig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderEx(t *testing.T) {
	testcases := []struct {
		name     string
		fontName string
		input    string
		expected string
	}{
		{
			name:     "3x5 A",
			fontName: "3x5",
			input:    "A",
			expected: "    \n #  \n# # \n### \n# # \n# # ",
		},
		{
			name:     "3x5 AB",
			fontName: "3x5",
			input:    "AB",
			expected: "        \n #  ##  \n# # # # \n### ##  \n# # # # \n# # ##  ",
		},
		{
			name:     "standard A",
			fontName: "standard",
			input:    "A",
			expected: "    _    \n   / \\   \n  / _ \\  \n / ___ \\ \n/_/   \\_\\n         ",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			font, _ := loadFont(tc.fontName)
			render := New(font)
			assert.Equal(t, tc.expected, render.RenderExp(tc.input))
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
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.rule(tc.input[0], tc.input[1])
			assert.Equal(t, tc.result, result)
		})
	}
}
