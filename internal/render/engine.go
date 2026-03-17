package render

import (
	"sync"

	"github.com/phantompunk/fig/internal/font"
)

type FilterFunc func(*Canvas) *Canvas

type RenderOptions struct {
	FontName   string
	FilterFunc []FilterFunc
	Align      Alignment // AlignLeft (default), AlignCenter, or AlignRight
	Width      int       // terminal width override; 0 means detect at render time
}

// renderCacheKey is the unique identity of a rendered output. FilterFunc is
// intentionally excluded — function values are not comparable. Callers that
// supply FilterFunc bypass the cache entirely (see Render).
type renderCacheKey struct {
	text     string
	fontName string
	align    Alignment
	width    int
}

type Engine struct {
	registry  *font.FontRegistry
	TermWidth func() int // injectable for tests; defaults to terminalWidth
	cacheMu   sync.RWMutex
	cache     map[renderCacheKey]string
}

func New(loaders ...font.FontLoader) *Engine {
	return &Engine{
		registry:  font.NewRegistry(loaders...),
		TermWidth: terminalWidth,
		cache:     make(map[renderCacheKey]string),
	}
}

func NewEngine() *Engine {
	return &Engine{
		registry:  font.NewRegistry(),
		TermWidth: terminalWidth,
		cache:     make(map[renderCacheKey]string),
	}
}

// CacheLen returns the number of entries currently in the render cache.
// Intended for testing and diagnostics.
func (e *Engine) CacheLen() int {
	e.cacheMu.RLock()
	defer e.cacheMu.RUnlock()
	return len(e.cache)
}

func (e *Engine) Render(text string, opts RenderOptions) (string, error) {
	// Resolve terminal width once so it is consistent across both the cache
	// key and the alignment call. Width=0 means "detect now".
	effectiveWidth := opts.Width
	if opts.Align != AlignLeft && effectiveWidth == 0 {
		effectiveWidth = e.TermWidth()
	}

	// FilterFunc values are not comparable, so skip the cache when any are set.
	if len(opts.FilterFunc) == 0 {
		key := renderCacheKey{text: text, fontName: opts.FontName, align: opts.Align, width: effectiveWidth}
		e.cacheMu.RLock()
		if result, ok := e.cache[key]; ok {
			e.cacheMu.RUnlock()
			return result, nil
		}
		e.cacheMu.RUnlock()

		out, err := e.render(text, opts, effectiveWidth)
		if err != nil {
			return "", err
		}
		e.cacheMu.Lock()
		e.cache[key] = out
		e.cacheMu.Unlock()
		return out, nil
	}

	return e.render(text, opts, effectiveWidth)
}

func (e *Engine) render(text string, opts RenderOptions, effectiveWidth int) (string, error) {
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
		out = alignOutput(out, opts.Align, effectiveWidth)
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
