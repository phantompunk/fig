package font

import (
	"bytes"
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/phantompunk/fig/assets"
)

// Format discriminates between supported font file formats.
type Format uint8

const (
	FormatFLF Format = iota // .flf — standard Figlet
	FormatTLF               // .tlf — Toilet (future)
)

// FontLoader abstracts the source of font data.
// New sources (HTTP, zip bundle) can implement this interface.
type FontLoader interface {
	// Load returns the raw bytes and format of a font file by name.
	// name has no extension — the loader resolves the extension.
	Load(name string) ([]byte, Format, error)

	// List returns all font names available from this source, without extensions.
	List() ([]string, error)
}

type EmbedLoader struct {
	FS  embed.FS
	Dir string
}

func (e EmbedLoader) Load(name string) ([]byte, Format, error) {
	candidates := []struct {
		ext    string
		format Format
	}{
		{".flf", FormatFLF},
		{".tlf", FormatTLF},
	}

	for _, c := range candidates {
		path := filepath.Join(e.Dir, name+c.ext)
		data, err := e.FS.ReadFile(path)
		if err == nil {
			return data, c.format, nil
		}
	}

	return nil, 0, fmt.Errorf("font %q not found in embedded filesystem", name)

}

// List returns all font names in the embedded directory, without extensions.
func (e EmbedLoader) List() ([]string, error) {
	entries, err := e.FS.ReadDir(e.Dir)
	if err != nil {
		return nil, fmt.Errorf("listing embedded fonts in %q: %w", e.Dir, err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if isFontFile(name) {
			names = append(names, strings.TrimSuffix(name, filepath.Ext(name)))
		}
	}
	return names, nil
}

// isFontFile reports whether the filename has a recognised font extension.
func isFontFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".flf" || ext == ".tlf"
}

// ListFonts returns a list of available font names.
func ListFonts() []string {
	files, err := assets.FontFS.ReadDir(".")
	if err != nil {
		return nil
	}

	names := []string{}
	for _, file := range files {
		name := file.Name()[:len(file.Name())-4]
		names = append(names, name)
	}

	return names
}

func loadFont(name string) (*Font, error) {
	fileName := filepath.Join(".", name+".flf")

	data, err := assets.FontFS.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return parseFont(name, bytes.NewReader(data))
}

// BundledLoader returns a FontLoader for the fonts embedded at build time.
func BundledLoader() FontLoader {
	return EmbedLoader{FS: assets.FontFS, Dir: "."}
}
