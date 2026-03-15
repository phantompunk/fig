package fig

import (
	"flag"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/phantompunk/fig/assert"
)

var update = flag.Bool("update", false, "regenerate golden files using figlet")

// fontDir returns the absolute path to the embedded fonts directory so figlet
// can be pointed at exactly the same font files the renderer uses.
func fontDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Join(filepath.Dir(file), "fonts")
}

func figletRender(t *testing.T, fontName, text string) string {
	t.Helper()
	out, err := exec.Command("figlet", "-d", fontDir(t), "-f", fontName, text).Output()
	if err != nil {
		t.Fatalf("figlet -f %s %q: %v", fontName, text, err)
	}
	return string(out)
}

func TestRenderGolden(t *testing.T) {
	tests := []struct {
		name string
		font string
		text string
	}{
		// standard: EqualChar + Underscore + Hierarchy + OppositePair smush rules
		{"standard_Hi", "standard", "Hi"},
		{"standard_fig", "standard", "fig"},
		{"standard_Hello", "standard", "Hello"},
		// 3x5: full-width mode, no smushing or kerning
		{"3x5_Hi", "3x5", "Hi"},
		// banner: different smush rule set
		{"banner_Hi", "banner", "Hi"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if *update {
				assert.Golden(t, tt.name, figletRender(t, tt.font, tt.text), true)
				return
			}

			f, err := Font(tt.font)
			if err != nil {
				t.Fatalf("load font %q: %v", tt.font, err)
			}
			assert.Golden(t, tt.name, f.Render(tt.text), false)
		})
	}
}
