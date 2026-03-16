package input

import "testing"

func TestNormalizeInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Surrounding whitespace", "  hello  ", "hello"},
		{"Windows line endings", "line1\r\nline2", "line1\nline2"},
		{"Mixed line endings", "line1\rline2\nline3", "line1\nline2\nline3"},
		{"Trailing whitespace per line", "row1   \nrow2 ", "row1\nrow2"},
		{"Complex case", "  \n tabbed \t \r\n next   ", "tabbed\nnext"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeInput(tc.input)
			if got != tc.expected {
				t.Errorf("normalizeInput() = %q, want %q", got, tc.expected)
			}
		})
	}
}
