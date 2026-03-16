package canvas

import (
	"testing"

	"github.com/phantompunk/fig/assert"
	fig "github.com/phantompunk/fig/internal/font"
)

func TestNewCanvas(t *testing.T) {
	c := NewCanvas(4, 10)
	if c.height != 4 {
		t.Errorf("want height 4, got %d", c.height)
	}
	if len(c.cells) != 4 {
		t.Errorf("want 4 rows, got %d", len(c.cells))
	}
	if len(c.cells[0]) != 10 {
		t.Errorf("want width 10, got %d", len(c.cells[0]))
	}
}

func TestStampWriteGlyphAtOffset(t *testing.T) {
	c := NewCanvas(3, 10)
	glyph := [][]rune{
		[]rune("##"),
		[]rune("##"),
		[]rune("##"),
	}
	c.Stamp(glyph, 2)

	want := [][]rune{
		{0, 0, '#', '#', 0, 0, 0, 0, 0, 0},
		{0, 0, '#', '#', 0, 0, 0, 0, 0, 0},
		{0, 0, '#', '#', 0, 0, 0, 0, 0, 0},
	}
	assert.CanvasEqual(t, want, c.cells)
}

func TestStampOutOfBounds(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("stamp panicked: %v", r)
		}
	}()

	c := NewCanvas(3, 4)
	glyph := [][]rune{[]rune("ABCDE")}
	c.Stamp(glyph[:1], 0)
}

func TestSmushCell(t *testing.T) {
	hb := '$'
	tests := []struct {
		name  string
		left  rune
		right rune
		rule  fig.SmushRule
		want  rune
		ok    bool
	}{
		{"equal chars: same letters", 'X', 'X', fig.EqualChars, 'X', true},
		{"equal chars: diff letters", 'X', 'Y', fig.EqualChars, 0, false},
		{"underscore with pipe", '_', '|', fig.Underscore, '|', true},
		{"underscore with letter", 'l', '_', fig.Underscore, 0, false},
		{"underscore with self", '_', '_', fig.Underscore, 0, false},
		{"hierachy: pipe vs slash", '|', '/', fig.Heirarchy, '/', true},
		{"hierachy: bracket vs pipe", '[', '|', fig.Heirarchy, '[', true},
		{"hierachy: paren vs bracket", '|', '/', fig.Heirarchy, '/', true},
		{"hierachy: same char", '[', ']', fig.Heirarchy, 0, false},
		{"hierachy: letters", 'A', 'B', fig.Heirarchy, 0, false},
		{"opposite: brackets", '[', ']', fig.OppositePair, '|', true},
		{"opposite: rev brackets", ']', '[', fig.OppositePair, '|', true},
		{"opposite: braces", '{', '}', fig.OppositePair, '|', true},
		{"opposite: same bracket", '[', '[', fig.OppositePair, 0, false},
		{"bigX: forward+back", '/', '\\', fig.BigX, '|', true},
		{"bigX: back+forward", '\\', '/', fig.BigX, 'Y', true},
		{"bigX: angles", '>', '<', fig.BigX, 'X', true},
		{"hardblank: two $", '$', '$', fig.HardblankRule(hb), hb, true},
		{"hardblank: $+space", '$', ' ', fig.HardblankRule(hb), 0, false},
		{"hardblank: $+letter", '$', 'A', fig.HardblankRule(hb), 0, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := smushCell(tc.left, tc.right, []fig.SmushRule{tc.rule}, hb)
			if ok != tc.ok || got != tc.want {
				t.Errorf("smushCell(%q, %q)=(%q, %v), want (%q,%v)", tc.left, tc.right, got, ok, tc.want, tc.ok)
			}
		})
	}
}

func TestFindOverlap_noSmushPossible(t *testing.T) {
	// Two glyphs with incompatible chars in any overlap position
	// Using rule1 only: chars must match to smush
	existing := [][]rune{
		[]rune("AB"),
	}
	incoming := [][]rune{
		[]rune("CD"),
	}
	c := NewCanvas(1, 4)
	c.Stamp(existing, 0)
	// Max overlap would be 2, but A≠C and B≠D — should back off to 0
	got := c.FindOverlap(incoming, 2, []fig.SmushRule{fig.EqualCharsRule('$')}, '$')
	if got != 0 {
		t.Errorf("want overlap 0, got %d", got)
	}
}

func TestFindOverlap_fullOverlapClean(t *testing.T) {
	existing := [][]rune{[]rune("XX")}
	incoming := [][]rune{[]rune("XX")}
	c := NewCanvas(1, 4)
	c.Stamp(existing, 0)
	got := c.FindOverlap(incoming, 2, []fig.SmushRule{fig.EqualCharsRule('$')}, '$')
	if got != 2 {
		t.Errorf("want overlap 2, got %d", got)
	}
}

func TestFindOverlap_partialOverlap(t *testing.T) {
	// Right col matches (X==X), left col doesn't (A≠B)
	// Full overlap of 2 fails, back off to 1 — right col only
	existing := [][]rune{[]rune("AX")}
	incoming := [][]rune{[]rune("XB")}
	c := NewCanvas(1, 4)
	c.Stamp(existing, 0)
	got := c.FindOverlap(incoming, 2, []fig.SmushRule{fig.EqualChars}, '$')
	if got != 1 {
		t.Errorf("want overlap 1, got %d", got)
	}
}

func TestStampSmush_mergesOverlapZone(t *testing.T) {
	// "AB" already on canvas, stamping "BC" with 1-col overlap
	// Overlap col: B == B (rule1 match) → B survives
	c := NewCanvas(1, 6)
	c.Stamp([][]rune{[]rune("AB")}, 0)

	incoming := [][]rune{[]rune("BC")}
	c.StampSmush(incoming, 1, []fig.SmushRule{fig.EqualCharsRule('$')}, '$')

	want := []rune{'A', 'B', 'C', 0, 0, 0}
	assert.RowEqual(t, want, c.cells[0])
}

func TestStampSmush_noOverwriteOnFailedSmush(t *testing.T) {
	// If smush fails (shouldn't happen post-FindOverlap, but be defensive)
	// existing cell wins
	c := NewCanvas(1, 4)
	c.Stamp([][]rune{[]rune("AB")}, 0)
	// Stamp something that FindOverlap would resolve to 0 overlap
	c.StampSmush([][]rune{[]rune("CD")}, 2, []fig.SmushRule{fig.EqualChars}, '$')
	want := []rune{'A', 'B', 'C', 'D'}
	assert.RowEqual(t, want, c.cells[0])
}

func TestString_replacesHardblanksAndZeros(t *testing.T) {
    c := NewCanvas(2, 4)
    c.cells[0] = []rune{'A', '$', 0, 'B'}
    c.cells[1] = []rune{0, 'C', '$', 0}

    got := c.String('$')
    want := "A  B\n C \n"
    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}
