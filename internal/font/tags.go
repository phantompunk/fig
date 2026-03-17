package font

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type fontEntry struct {
	Name string   `yaml:"name"`
	Tags []string `yaml:"tags"`
}

type fontsYAML struct {
	Fonts []fontEntry `yaml:"fonts"`
}

// TagMap maps a font name to its slice of tags.
type TagMap map[string][]string

// LoadTagMap parses the embedded YAML and returns a TagMap.
func LoadTagMap(data []byte) (TagMap, error) {
	var doc fontsYAML
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parsing fonts.yaml: %w", err)
	}
	m := make(TagMap, len(doc.Fonts))
	for _, e := range doc.Fonts {
		m[e.Name] = e.Tags
	}
	return m, nil
}

// HasTag reports whether the named font has the given tag.
// Always returns true when tag is "all".
func (tm TagMap) HasTag(name, tag string) bool {
	if tag == "all" {
		return true
	}
	for _, t := range tm[name] {
		if t == tag {
			return true
		}
	}
	return false
}
