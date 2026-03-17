package font

type Font struct {
	name     string
	metadata Metadata
	glyphs   GlyphDict
	rules    []SmushRule
}

// LoadFont loads a FIGlet font by name.
func LoadFont(name string) (*Font, error) {
	return loadFont(name)
}

// Must is a helper method that panics if there is an error loading the font.
func Must(t *Font, err error) *Font {
	if err != nil {
		panic(err)
	}
	return t
}

// Name returns the name of the font.
func (f *Font) Name() string     { return f.name }
func (f *Font) Height() int      { return f.metadata.height }
func (f *Font) Hardblank() rune  { return f.metadata.hardBlank }
func (f *Font) MaxLength() int   { return f.metadata.maxLength }
func (f *Font) Rules() []SmushRule { return f.rules }
func (f *Font) IsFullWidth() bool  { return f.metadata.layoutMode.FullWidth }

// GlyphRunes returns the glyph for char as a 2D rune slice ready for the
// canvas. Regular spaces are converted to 0 (transparent); hardblanks and all
// other characters are kept as-is.
func (f *Font) GlyphRunes(char rune) [][]rune {
	g := f.getGlyph(char)
	rows := make([][]rune, len(g.lines))
	for i, line := range g.lines {
		row := make([]rune, len(line))
		for j, ch := range line {
			if ch == ' ' {
				row[j] = 0
			} else {
				row[j] = ch
			}
		}
		rows[i] = row
	}
	return rows
}

func (f *Font) Render(text string) string {
	renderer := New(f)
	return renderer.Render(text)
}

func (f *Font) Lines(text string) []string {
	renderer := New(f)
	return renderer.Lines(text)
}

func (f *Font) getGlyph(char rune) Glyph {
	return f.glyphs[char]
}

func NewFigFont(name string, meta Metadata) *Font {
	return &Font{
		name:     name,
		metadata: meta,
		glyphs:   make(GlyphDict),
		rules:    setRules(meta),
	}
}

func setRules(meta Metadata) []SmushRule {
	rules := []SmushRule{}
	if meta.smushMode.EqualChar {
		rules = append(rules, EqualCharsRule(meta.hardBlank))
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
		rules = append(rules, HardblankRule(meta.hardBlank))
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
