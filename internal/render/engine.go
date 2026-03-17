package render

import (
	"github.com/phantompunk/fig/internal/font"
)

type FilterFunc func(*Canvas) *Canvas

type RenderOptions struct {
	FontName   string
	FilterFunc []FilterFunc
	Align      Alignment // AlignLeft (default), AlignCenter, or AlignRight
	Width      int       // terminal width override; 0 means detect at render time
}

type Engine struct {
	registry  *font.FontRegistry
	TermWidth func() int // injectable for tests; defaults to terminalWidth
}

func New(loaders ...font.FontLoader) *Engine {
	return &Engine{registry: font.NewRegistry(loaders...), TermWidth: terminalWidth}
}

func NewEngine() *Engine {
	return &Engine{registry: font.NewRegistry(), TermWidth: terminalWidth}
}

func (e *Engine) Render(text string, opts RenderOptions) (string, error) {
	f, err := e.registry.Get(opts.FontName)
	if err != nil {
		return "", err
	}

	runes := []rune(text)
	canvas := NewCanvas(f.Height(), len(runes)*f.MaxLength())

	rules := f.Rules()
	hb := f.Hardblank()

	cursor := 0
	for i, char := range runes {
		glyph := f.GlyphRunes(char)
		w := glyphWidth(glyph)
		if i == 0 {
			leftEdge := glyphLeftEdge(glyph)
			canvas.Stamp(glyph, -leftEdge)
			cursor = w - leftEdge
			continue
		}
		if f.IsFullWidth() {
			canvas.Stamp(glyph, cursor)
			cursor += w
		} else {
			overlap := canvas.FindOverlap(glyph, f.MaxLength(), rules, hb)
			canvas.StampSmush(glyph, overlap, rules, hb)
		}
	}

	out := canvas.String(hb, cursor)
	if opts.Align != AlignLeft {
		w := opts.Width
		if w == 0 {
			w = e.TermWidth()
		}
		out = alignOutput(out, opts.Align, w)
	}
	return out, nil
}

// ListFonts returns the names of all fonts available across the engine's loaders.
func (e *Engine) ListFonts() ([]string, error) {
	return e.registry.Available()
}

// FontHeight returns the line height of the named font. The font is loaded and
// cached on first call; subsequent calls are served from the registry cache.
func (e *Engine) FontHeight(name string) (int, error) {
	f, err := e.registry.Get(name)
	if err != nil {
		return 0, err
	}
	return f.Height(), nil
}

// glyphWidth returns the natural width of a glyph (max row length, including
// trailing zeros), which is used for full-width placement.
func glyphWidth(glyph [][]rune) int {
	w := 0
	for _, row := range glyph {
		if len(row) > w {
			w = len(row)
		}
	}
	return w
}

// glyphLeftEdge returns the minimum leading-zero count across all non-empty
// rows of the glyph. This is the amount to left-flush when rendering the first
// character so the output has no unnecessary leading whitespace.
func glyphLeftEdge(glyph [][]rune) int {
	minLeft := -1
	for _, row := range glyph {
		for x, ch := range row {
			if ch != 0 {
				if minLeft < 0 || x < minLeft {
					minLeft = x
				}
				break
			}
		}
	}
	if minLeft < 0 {
		return 0
	}
	return minLeft
}
