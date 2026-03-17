package render

import (
	"os"
	"strings"

	"github.com/charmbracelet/x/term"
)

// Alignment controls how rendered output is horizontally positioned.
type Alignment int

const (
	AlignLeft   Alignment = iota // default, no padding
	AlignCenter                  // center each line within the terminal width
	AlignRight                   // right-align each line within the terminal width
)

// alignOutput pads every line in s (a newline-terminated string from
// canvas.String) so the block sits at the requested alignment within width
// columns. Lines are padded relative to the widest line in the block so
// the whole block moves as a unit.
func alignOutput(s string, align Alignment, width int) string {
	if align == AlignLeft || width <= 0 {
		return s
	}

	// Trim the trailing newline, split, then restore it at the end.
	trimmed := strings.TrimRight(s, "\n")
	lines := strings.Split(trimmed, "\n")

	// Content width = widest line in the block.
	contentWidth := 0
	for _, l := range lines {
		if len(l) > contentWidth {
			contentWidth = len(l)
		}
	}

	var sb strings.Builder
	for _, l := range lines {
		pad := 0
		switch align {
		case AlignCenter:
			pad = (width - contentWidth) / 2
		case AlignRight:
			pad = width - contentWidth
		}
		if pad < 0 {
			pad = 0
		}
		sb.WriteString(strings.Repeat(" ", pad))
		sb.WriteString(l)
		sb.WriteByte('\n')
	}
	return sb.String()
}

// terminalWidth returns the current terminal column count, falling back to 80
// when stdout is not a TTY (pipes, CI, tests).
func terminalWidth() int {
	w, _, err := term.GetSize(os.Stdout.Fd())
	if err != nil || w <= 0 {
		return 80
	}
	return w
}
