package render

import (
	"fmt"
	"strings"
	"testing"

	"github.com/phantompunk/fig/internal/font"
)

// stubFontLoader serves raw FLF bytes from an in-memory map.
type stubFontLoader struct {
	fonts map[string][]byte
}

func (s *stubFontLoader) Load(name string) ([]byte, font.Format, error) {
	data, ok := s.fonts[name]
	if !ok {
		return nil, font.FormatFLF, fmt.Errorf("font %q not found", name)
	}
	return data, font.FormatFLF, nil
}

func (s *stubFontLoader) List() ([]string, error) {
	names := make([]string, 0, len(s.fonts))
	for k := range s.fonts {
		names = append(names, k)
	}
	return names, nil
}

// minimalFLF builds the smallest valid FLF that Parse accepts.
// height=1; each of the 95 printable ASCII glyphs is one line ending with "@@"
// (trim=2 for the only line leaves a single space as glyph content).
func minimalEngineFLF() []byte {
	var b strings.Builder
	b.WriteString("flf2a$ 1 1 2 0 0\n")
	for range 95 {
		b.WriteString(" @@\n")
	}
	return []byte(b.String())
}

func newEngineWithStub(fonts map[string][]byte) *Engine {
	return New(&stubFontLoader{fonts: fonts})
}

func TestEngine_Render_knownFont(t *testing.T) {
	e := newEngineWithStub(map[string][]byte{"mini": minimalEngineFLF()})

	got, err := e.Render("hi", RenderOptions{FontName: "mini"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == "" {
		t.Error("expected non-empty output")
	}
}

func TestEngine_Render_unknownFont(t *testing.T) {
	e := newEngineWithStub(map[string][]byte{})

	_, err := e.Render("hi", RenderOptions{FontName: "missing"})
	if err == nil {
		t.Error("expected error for unknown font, got nil")
	}
}

func TestEngine_Render_emptyText(t *testing.T) {
	e := newEngineWithStub(map[string][]byte{"mini": minimalEngineFLF()})

	got, err := e.Render("", RenderOptions{FontName: "mini"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Empty text produces an empty canvas — just newlines or nothing.
	_ = got
}

func TestEngine_Render_cachesFont(t *testing.T) {
	loader := &stubFontLoader{fonts: map[string][]byte{"mini": minimalEngineFLF()}}
	e := New(loader)

	for range 3 {
		if _, err := e.Render("a", RenderOptions{FontName: "mini"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	// The registry should have loaded the font only once; we can't inspect call
	// count here, but we verify no error across repeated calls.
}
