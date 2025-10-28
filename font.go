package fig

type FigFont struct {
	name     string
	metadata Metadata
	glyphs   GlyphDict
	rules    []SmushRule
}

func Font(name string) (*FigFont, error) {
	return loadFont(name)
}

func (f *FigFont) getGlyph(char rune) Glyph {
	return f.glyphs[char]
}

func NewFigFont(name string, meta Metadata) *FigFont {
	return &FigFont{
		name:     name,
		metadata: meta,
		glyphs:   make(GlyphDict),
		rules:    setRules(meta),
	}
}

func setRules(meta Metadata) []SmushRule {
	rules := []SmushRule{}
	if meta.smushMode.EqualChar {
		rules = append(rules, EqualChars)
	}

	if meta.smushMode.Underscore {
		rules = append(rules, Underscore)
	}

	if meta.smushMode.Underscore {
		rules = append(rules, Underscore)
	}

	if meta.smushMode.Hierarchy {
		rules = append(rules, Heirarchy)
	}

	if meta.smushMode.OppositePair {
		rules = append(rules, OppositePair)
	}

	if meta.smushMode.BigX {
		rules = append(rules, BigX)
	}

	if meta.smushMode.Hardblank {
		rules = append(rules, Hardblank)
	}
	return rules
}

type SmushMode struct {
	Enabled           bool
	HorizontalLayout  bool
	HorizontalFit     bool // enable kerning
	VerticalLayout    bool
	VerticalFit       bool
	HorizontalSmush   bool
	VerticalSmush     bool
	EqualChar         bool
	Underscore        bool
	Hierarchy         bool
	OppositePair      bool
	BigX              bool
	Hardblank         bool
	HorizontalKerning bool
	HorizontalFull    bool
	FullLayoutActive  bool
}

type VerticalMode struct{}

type Glyph struct {
	lines []string
	width int
}

type GlyphDict map[rune]Glyph

type Metadata struct {
	signature      string
	hardBlank      rune
	height         int
	baseline       int
	maxLength      int
	oldLayout      int
	commentLines   int
	printDirection int
	fullLayout     int
	codeTag        int
	comments       string
	smushMode      SmushMode
	verticalMode   VerticalMode
}

func (f *FigFont) Render(text string) string {
	renderer := New(f)
	return renderer.Render(text)
}

func (f *FigFont) RenderExp(text string) string {
	renderer := New(f)
	return renderer.RenderExp(text)
}
