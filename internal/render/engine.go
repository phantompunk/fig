package render

import (
	"github.com/phantompunk/fig/internal/font"
)

type FilterFunc func(*Canvas) *Canvas

type RenderOptions struct {
	FontName   string
	FilterFunc []FilterFunc
}

type Engine struct {
	registry *font.FontRegistry
}

func New(loaders ...font.FontLoader) *Engine {
	return &Engine{registry: font.NewRegistry(loaders...)}
}

func NewEngine() *Engine {
	return &Engine{registry: font.NewRegistry()}
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

	return canvas.String(hb, cursor), nil
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
