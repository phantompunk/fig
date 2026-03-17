package render

import (
	"strings"
	"testing"
)

// solidFLF builds a 1-row FLF where every printable ASCII glyph renders as
// a single '#' character. This makes output width exactly len(text) and lets
// alignment tests use simple arithmetic.
func solidFLF() []byte {
	var b strings.Builder
	b.WriteString("flf2a$ 1 1 2 0 0\n")
	for range 95 {
		b.WriteString("#@@\n")
	}
	return []byte(b.String())
}

func newAlignEngine(width int) *Engine {
	e := New(&stubFontLoader{fonts: map[string][]byte{"solid": solidFLF()}})
	e.TermWidth = func() int { return width }
	return e
}

func render1(t *testing.T, text string, align Alignment, width int) string {
	t.Helper()
	e := newAlignEngine(width)
	out, err := e.Render(text, RenderOptions{FontName: "solid", Align: align, Width: width})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	return out
}

func TestAlign_left_noPadding(t *testing.T) {
	out := render1(t, "Hi", AlignLeft, 20)
	line := strings.TrimRight(out, "\n")
	if strings.HasPrefix(line, " ") {
		t.Errorf("AlignLeft should not add leading spaces, got %q", line)
	}
}

func TestAlign_center(t *testing.T) {
	// "Hi" renders as "##" (2 chars). Width=20 → pad = (20-2)/2 = 9.
	out := render1(t, "Hi", AlignCenter, 20)
	line := strings.TrimRight(out, "\n")
	wantPad := (20 - 2) / 2
	wantPrefix := strings.Repeat(" ", wantPad)
	if !strings.HasPrefix(line, wantPrefix) {
		t.Errorf("AlignCenter: want %d leading spaces, got %q", wantPad, line)
	}
	if strings.HasPrefix(line, wantPrefix+" ") {
		t.Errorf("AlignCenter: too many leading spaces in %q", line)
	}
}

func TestAlign_right(t *testing.T) {
	// "Hi" renders as "##" (2 chars). Width=20 → pad = 20-2 = 18.
	out := render1(t, "Hi", AlignRight, 20)
	line := strings.TrimRight(out, "\n")
	wantPad := 20 - 2
	wantPrefix := strings.Repeat(" ", wantPad)
	if !strings.HasPrefix(line, wantPrefix) {
		t.Errorf("AlignRight: want %d leading spaces, got %q", wantPad, line)
	}
	if strings.HasPrefix(line, wantPrefix+" ") {
		t.Errorf("AlignRight: too many leading spaces in %q", line)
	}
}

func TestAlign_contentWiderThanTerminal_noPanic(t *testing.T) {
	// Content wider than terminal: pad should clamp to 0.
	out := render1(t, "Hello", AlignCenter, 3)
	line := strings.TrimRight(out, "\n")
	if strings.HasPrefix(line, " ") {
		t.Errorf("no padding expected when content wider than terminal, got %q", line)
	}
}

func TestAlign_trailingNewlinePreserved(t *testing.T) {
	out := render1(t, "Hi", AlignCenter, 20)
	if !strings.HasSuffix(out, "\n") {
		t.Errorf("trailing newline missing: %q", out)
	}
}

func TestAlign_widthFromTermWidth(t *testing.T) {
	// When Width=0, engine calls TermWidth. Inject a known value.
	e := New(&stubFontLoader{fonts: map[string][]byte{"solid": solidFLF()}})
	e.TermWidth = func() int { return 20 }

	out, err := e.Render("Hi", RenderOptions{FontName: "solid", Align: AlignCenter})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	line := strings.TrimRight(out, "\n")
	wantPad := (20 - 2) / 2
	if !strings.HasPrefix(line, strings.Repeat(" ", wantPad)) {
		t.Errorf("expected %d leading spaces from injected TermWidth, got %q", wantPad, line)
	}
}

func TestAlignOutput_multiline_uniformPad(t *testing.T) {
	// All lines in the block should be padded by the same amount, based on the
	// widest line — not padded independently.
	input := "##\n####\n##\n"
	got := alignOutput(input, AlignCenter, 10)
	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")

	// Widest line is "####" (4 chars). Center pad = (10-4)/2 = 3.
	wantPad := (10 - 4) / 2
	for i, l := range lines {
		if !strings.HasPrefix(l, strings.Repeat(" ", wantPad)) {
			t.Errorf("line %d: want %d leading spaces, got %q", i, wantPad, l)
		}
		// Must not have more than wantPad leading spaces.
		if strings.HasPrefix(l, strings.Repeat(" ", wantPad+1)) {
			t.Errorf("line %d: too many leading spaces in %q", i, l)
		}
	}
}
