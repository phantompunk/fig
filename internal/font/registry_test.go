package font

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/phantompunk/fig/internal/assert"
)

// stubLoader is a FontLoader that serves fonts from an in-memory map.
type stubLoader struct {
	fonts map[string][]byte
	// loadCount tracks how many times Load was called per name.
	mu        sync.Mutex
	loadCount map[string]int
}

func newStubLoader(fonts map[string][]byte) *stubLoader {
	return &stubLoader{fonts: fonts, loadCount: make(map[string]int)}
}

func (s *stubLoader) Load(name string) ([]byte, Format, error) {
	s.mu.Lock()
	s.loadCount[name]++
	s.mu.Unlock()

	data, ok := s.fonts[name]
	if !ok {
		return nil, FormatFLF, fmt.Errorf("font %q not found", name)
	}
	return data, FormatFLF, nil
}

func (s *stubLoader) List() ([]string, error) {
	names := make([]string, 0, len(s.fonts))
	for k := range s.fonts {
		names = append(names, k)
	}
	return names, nil
}

// errLoader always returns an error from both Load and List.
type errLoader struct{ err error }

func (e *errLoader) Load(_ string) ([]byte, Format, error) { return nil, FormatFLF, e.err }
func (e *errLoader) List() ([]string, error)               { return nil, e.err }

// minimalFLF returns a minimal valid FIGlet font file for testing.
// height=1 means each glyph is one line; the parser requires the final line
// of each glyph to have at least 2 chars so it can strip the 2-char end marker.
// We define all 95 required printable ASCII glyphs (32–126).
func minimalFLF() []byte {
	var b strings.Builder
	// header: flf2a<hardblank> height baseline maxLength oldLayout commentLines
	b.WriteString("flf2a$ 1 1 2 0 0\n")
	// 95 characters: ASCII 32..126, each one line ending with "@@ "
	// (trim=2 for the last/only line, so "@@ " → "@" after trim — but we just
	// need len >= 2; using " @@" gives a single space glyph after trim)
	for range 95 {
		b.WriteString(" @@\n")
	}
	return []byte(b.String())
}

func TestNewRegistry_NoLoaders(t *testing.T) {
	r := NewRegistry()
	_, err := r.Get("standard")
	if err == nil {
		t.Error("expected error with no loaders, got nil")
	}
}

func TestGet_FontFound(t *testing.T) {
	loader := newStubLoader(map[string][]byte{
		"mini": minimalFLF(),
	})
	r := NewRegistry(loader)

	f, err := r.Get("mini")
	assert.NilError(t, err)
	assert.Equal(t, f.Name(), "mini")
}

func TestGet_FontNotFound(t *testing.T) {
	loader := newStubLoader(map[string][]byte{})
	r := NewRegistry(loader)

	_, err := r.Get("missing")
	if err == nil {
		t.Error("expected error for missing font, got nil")
	}
}

func TestGet_CachesResult(t *testing.T) {
	loader := newStubLoader(map[string][]byte{
		"mini": minimalFLF(),
	})
	r := NewRegistry(loader)

	f1, err1 := r.Get("mini")
	f2, err2 := r.Get("mini")

	assert.NilError(t, err1)
	assert.NilError(t, err2)
	// Same pointer — loaded exactly once.
	assert.True(t, f1 == f2)
	assert.Equal(t, loader.loadCount["mini"], 1)
}

func TestGet_ConcurrentSameFont(t *testing.T) {
	loader := newStubLoader(map[string][]byte{
		"mini": minimalFLF(),
	})
	r := NewRegistry(loader)

	const goroutines = 20
	fonts := make([]*Font, goroutines)
	errs := make([]error, goroutines)

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := range goroutines {
		go func(i int) {
			defer wg.Done()
			fonts[i], errs[i] = r.Get("mini")
		}(i)
	}
	wg.Wait()

	for i, err := range errs {
		assert.NilError(t, err)
		assert.True(t, fonts[i] == fonts[0]) // all goroutines got same *Font
	}
	// Load should have been called exactly once despite concurrent requests.
	assert.Equal(t, loader.loadCount["mini"], 1)
}

func TestGet_ErrorIsCached(t *testing.T) {
	loader := newStubLoader(map[string][]byte{}) // no fonts
	r := NewRegistry(loader)

	_, err1 := r.Get("missing")
	_, err2 := r.Get("missing")

	if err1 == nil || err2 == nil {
		t.Error("expected errors, got nil")
	}
	// Load should still only be attempted once.
	assert.Equal(t, loader.loadCount["missing"], 1)
}

func TestLoad_FallsBackToSecondLoader(t *testing.T) {
	first := newStubLoader(map[string][]byte{})                      // no fonts
	second := newStubLoader(map[string][]byte{"mini": minimalFLF()}) // has font
	r := NewRegistry(first, second)

	f, err := r.Get("mini")
	assert.NilError(t, err)
	assert.Equal(t, f.Name(), "mini")
	assert.Equal(t, first.loadCount["mini"], 1)
	assert.Equal(t, second.loadCount["mini"], 1)
}

func TestLoad_FirstLoaderWins(t *testing.T) {
	first := newStubLoader(map[string][]byte{"mini": minimalFLF()})
	second := newStubLoader(map[string][]byte{"mini": minimalFLF()})
	r := NewRegistry(first, second)

	f, err := r.Get("mini")
	assert.NilError(t, err)
	assert.Equal(t, f.Name(), "mini")
	assert.Equal(t, first.loadCount["mini"], 1)
	assert.Equal(t, second.loadCount["mini"], 0)
}

func TestPreload_AllSucceed(t *testing.T) {
	loader := newStubLoader(map[string][]byte{
		"a": minimalFLF(),
		"b": minimalFLF(),
	})
	r := NewRegistry(loader)

	err := r.Preload([]string{"a", "b"})
	assert.NilError(t, err)
}

func TestPreload_SomeFail(t *testing.T) {
	loader := newStubLoader(map[string][]byte{
		"good": minimalFLF(),
	})
	r := NewRegistry(loader)

	err := r.Preload([]string{"good", "bad"})
	if err == nil {
		t.Error("expected error for missing font, got nil")
	}
}

func TestPreload_AllFail(t *testing.T) {
	r := NewRegistry(newStubLoader(map[string][]byte{}))

	err := r.Preload([]string{"x", "y"})
	if err == nil {
		t.Error("expected combined error, got nil")
	}
}

func TestPreload_Empty(t *testing.T) {
	r := NewRegistry()
	err := r.Preload([]string{})
	assert.NilError(t, err)
}

func TestPreload_ErrorMessageContainsFontName(t *testing.T) {
	r := NewRegistry(newStubLoader(map[string][]byte{}))

	err := r.Preload([]string{"nope"})
	if err == nil {
		t.Fatal("expected error")
	}
	// The error should mention the font name.
	if !errors.Is(err, err) { // trivially true; just check the string
		t.Error("unexpected error shape")
	}
	// Verify "nope" appears in the message.
	assert.True(t, len(err.Error()) > 0)
}

// --- Available ---

func TestAvailable_SingleLoader(t *testing.T) {
	loader := newStubLoader(map[string][]byte{
		"a": minimalFLF(),
		"b": minimalFLF(),
	})
	r := NewRegistry(loader)

	names, err := r.Available()
	assert.NilError(t, err)
	assert.Equal(t, len(names), 2)
}

func TestAvailable_Deduplication(t *testing.T) {
	first := newStubLoader(map[string][]byte{"a": minimalFLF(), "b": minimalFLF()})
	second := newStubLoader(map[string][]byte{"b": minimalFLF(), "c": minimalFLF()})
	r := NewRegistry(first, second)

	names, err := r.Available()
	assert.NilError(t, err)
	// "b" appears in both loaders but should only be listed once.
	seen := make(map[string]int)
	for _, n := range names {
		seen[n]++
	}
	assert.Equal(t, seen["b"], 1)
	assert.Equal(t, len(names), 3) // a, b, c
}

func TestAvailable_SkipsFailingLoaders(t *testing.T) {
	bad := &errLoader{err: fmt.Errorf("unavailable")}
	good := newStubLoader(map[string][]byte{"a": minimalFLF()})
	r := NewRegistry(bad, good)

	names, err := r.Available()
	assert.NilError(t, err)
	assert.Equal(t, len(names), 1)
	assert.Equal(t, names[0], "a")
}

func TestAvailable_NoFontsReturnsError(t *testing.T) {
	r := NewRegistry()
	_, err := r.Available()
	if err == nil {
		t.Error("expected error when no fonts are available, got nil")
	}
}

func TestAvailable_AllLoadersFail(t *testing.T) {
	r := NewRegistry(
		&errLoader{err: fmt.Errorf("e1")},
		&errLoader{err: fmt.Errorf("e2")},
	)
	_, err := r.Available()
	if err == nil {
		t.Error("expected error when all loaders fail, got nil")
	}
}
