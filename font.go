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
	Enabled      bool
	EqualChar    bool // bit 0
	Underscore   bool // bit 1
	Hierarchy    bool // bit 2
	OppositePair bool // bit 3
	BigX         bool // bit 4
	Hardblank    bool // bit 5
	VEqualChar   bool // bit 8
	VUnderscore  bool // bit 9
	VHierarchy   bool // bit 10
	HLine        bool // bit 11
	Vline        bool // bit 12
}

type LayoutMode struct {
	FullWidth bool // no smushing
	Kerning   bool // overlap by 1
	Smushing  bool // smushing enabled
	Universal bool // universal smushing
	VKerning  bool // vertical kerning
	VSmushing bool // vertical smushing
}

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
	layoutMode     LayoutMode
}

func (f *FigFont) Render(text string) string {
	renderer := New(f)
	return renderer.Render(text)
}
