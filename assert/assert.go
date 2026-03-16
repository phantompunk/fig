package assert

import (
	"os"
	"path/filepath"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

func True(t *testing.T, actual bool) {
	t.Helper()
	if !actual {
		t.Errorf("got %v; expected: True", actual)
	}
}

func False(t *testing.T, actual bool) {
	t.Helper()
	if actual {
		t.Errorf("got %v; expected: True", actual)
	}
}

func NilError(t *testing.T, actual error) {
	t.Helper()

	if actual != nil {
		t.Errorf("got %v; expected: nil", actual)
	}
}

func NotNil(t *testing.T, actual any) {
	t.Helper()

	if actual == nil {
		t.Errorf("got nil; expected: not nil")
	}
}

// Golden compares actual against the contents of the named golden file under
// testdata/. When update is true the file is written instead of read, allowing
// callers to regenerate baselines with -update.
func Golden(t *testing.T, name, actual string, update bool) {
	t.Helper()

	path := filepath.Join("testdata", name+".golden")

	if update {
		if err := os.MkdirAll("testdata", 0755); err != nil {
			t.Fatalf("create testdata dir: %v", err)
		}
		if err := os.WriteFile(path, []byte(actual), 0644); err != nil {
			t.Fatalf("write golden file %s: %v", path, err)
		}
		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("missing golden file %s — run with -update to create it", path)
	}
	Equal(t, string(want), actual)
}

func CanvasEqual(t *testing.T, want, got [][]rune) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("canvas height mismatch: want %d, got %d", len(want), len(got))
	}

	for y := range want {
		if len(want[y]) != len(got[y]) {
			t.Errorf("canvas row %d width mismatch: want %d, got %d", y, len(want[y]), len(got[y]))
			return
		}

		// 2. Check individual runes
		for x := 0; x < len(want[y]); x++ {
			if want[y][x] != got[y][x] {
				t.Errorf("mismatch at [%d, %d]: want %q, got %q", y, x, want[y][x], got[y][x])
				return
			}
		}
	}
}

func RowEqual(t *testing.T, want, got []rune) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("row width mismatch: want %d, got %d", len(want), len(got))
	}

	for idx := range want {
		if want[idx] != got[idx] {
			t.Errorf("unexpected cell at %d: want %q, got %q", idx, want[idx], got[idx])
		}
	}
}
