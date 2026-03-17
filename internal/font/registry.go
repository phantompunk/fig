package font

import (
	"fmt"
	"sync"

)

// entry is the value stored in the registry cache.
// once ensures the font is parsed exactly once even under concurrent access.
type entry struct {
	once sync.Once
	font *Font
	err  error
}

// FontRegistry is a concurrency-safe, lazily-populated font cache.
// Fonts are loaded and parsed once then reused across all callers.
type FontRegistry struct {
	loaders []FontLoader
	cache   sync.Map // map[string]*entry
}

// NewRegistry returns a FontRegistry that searches the given loaders in order.
// The first loader that successfully returns data for a name wins.
func NewRegistry(loaders ...FontLoader) *FontRegistry {
	return &FontRegistry{loaders: loaders}
}

// Get returns the named font, loading and parsing it on first access.
// Subsequent calls for the same name are served from cache.
// Safe for concurrent use: if multiple goroutines request the same uncached
// font simultaneously the file is loaded and parsed exactly once.
func (r *FontRegistry) Get(name string) (*Font, error) {
	actual, _ := r.cache.LoadOrStore(name, &entry{})
	e := actual.(*entry)
	e.once.Do(func() {
		e.font, e.err = r.load(name)
	})
	return e.font, e.err
}

// Preload eagerly loads the named fonts into the cache.
// Errors for individual fonts are collected and returned as a combined error;
// successfully loaded fonts are still cached.
func (r *FontRegistry) Preload(names []string) error {
	var errs []error
	for _, name := range names {
		if _, err := r.Get(name); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("preload: %d font(s) failed: %v", len(errs), errs)
	}
	return nil
}

// Available returns the union of all font names across all loaders.
// Names are deduplicated; order is loaders-first then alphabetical within
// each loader's output.
func (r *FontRegistry) Available() ([]string, error) {
	seen := make(map[string]struct{})
	var names []string

	for _, l := range r.loaders {
		list, err := l.List()
		if err != nil {
			continue // best-effort: skip unavailable loaders
		}
		for _, name := range list {
			if _, ok := seen[name]; !ok {
				seen[name] = struct{}{}
				names = append(names, name)
			}
		}
	}

	if len(names) == 0 {
		return nil, fmt.Errorf("no fonts available across %d loader(s)", len(r.loaders))
	}
	return names, nil
}

// load performs the actual file read and parse for a single font name.
// It is called at most once per name (guarded by sync.Once in the entry).
func (r *FontRegistry) load(name string) (*Font, error) {
	for _, l := range r.loaders {
		data, _, err := l.Load(name)
		if err != nil {
			continue
		}
		return Parse(data, name, fmt.Sprintf("loader:%T", l))
	}
	return nil, fmt.Errorf("font %q not found in any loader", name)
}

