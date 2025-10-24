package fig

// ListFonts returns a list of available font names.
func ListFonts() []string {
	files, err := figFS.ReadDir("fonts")
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
