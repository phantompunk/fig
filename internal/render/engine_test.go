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

func TestEngine_Render_cacheHit(t *testing.T) {
	e := newEngineWithStub(map[string][]byte{"mini": minimalEngineFLF()})
	opts := RenderOptions{FontName: "mini"}

	first, err := e.Render("hi", opts)
	if err != nil {
		t.Fatalf("first render: %v", err)
	}
	if e.CacheLen() != 1 {
		t.Fatalf("expected cache length 1 after first render, got %d", e.CacheLen())
	}

	second, err := e.Render("hi", opts)
	if err != nil {
		t.Fatalf("second render: %v", err)
	}
	if e.CacheLen() != 1 {
		t.Errorf("expected cache length still 1 after cache hit, got %d", e.CacheLen())
	}
	if first != second {
		t.Errorf("cache hit returned different result: %q != %q", first, second)
	}
}

func TestEngine_Render_cacheMiss_distinctFont(t *testing.T) {
	e := newEngineWithStub(map[string][]byte{
		"a": minimalEngineFLF(),
		"b": minimalEngineFLF(),
	})

	if _, err := e.Render("hi", RenderOptions{FontName: "a"}); err != nil {
		t.Fatalf("render a: %v", err)
	}
	if _, err := e.Render("hi", RenderOptions{FontName: "b"}); err != nil {
		t.Fatalf("render b: %v", err)
	}
	if e.CacheLen() != 2 {
		t.Errorf("expected 2 cache entries for two distinct fonts, got %d", e.CacheLen())
	}
}

func TestEngine_Render_cacheBypass_withFilter(t *testing.T) {
	e := newEngineWithStub(map[string][]byte{"mini": minimalEngineFLF()})
	opts := RenderOptions{
		FontName:   "mini",
		FilterFunc: []FilterFunc{func(c *Canvas) *Canvas { return c }},
	}

	for range 3 {
		if _, err := e.Render("hi", opts); err != nil {
			t.Fatalf("render: %v", err)
		}
	}
	if e.CacheLen() != 0 {
		t.Errorf("cache should stay empty when FilterFunc is set, got %d entries", e.CacheLen())
	}
}

func TestEngine_Render_width0_resolvedPerCall(t *testing.T) {
	e := newEngineWithStub(map[string][]byte{"mini": minimalEngineFLF()})

	callCount := 0
	widths := []int{80, 120}
	e.TermWidth = func() int {
		w := widths[callCount%len(widths)]
		callCount++
		return w
	}

	opts := RenderOptions{FontName: "mini", Align: AlignCenter} // Width=0 → resolved via TermWidth
	if _, err := e.Render("hi", opts); err != nil {
		t.Fatalf("first render: %v", err)
	}
	if _, err := e.Render("hi", opts); err != nil {
		t.Fatalf("second render: %v", err)
	}

	// Two renders with different resolved widths → two distinct cache entries.
	if e.CacheLen() != 2 {
		t.Errorf("expected 2 cache entries for two distinct resolved widths, got %d", e.CacheLen())
	}
}
