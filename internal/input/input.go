package input

import (
	"io"
	"os"
	"strings"
	"unicode"
)

type Source interface {
	Read() (string, error)
}

type Args struct {
	Args []string
}

func (a Args) Read() (string, error) {
	if len(a.Args) == 0 {
		return "", nil
	}
	return normalizeInput(strings.Join(a.Args, " ")), nil
}

type Stdin struct {
	Reader io.Reader
}

func (s Stdin) Read() (string, error) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, s.Reader)
	if err != nil {
		return "", err
	}
	return normalizeInput(buf.String()), nil
}

func Resolve(args []string) Source {
	if stdinHasData() {
		return Stdin{Reader: os.Stdin}
	}
	return Args{Args: args}
}

func stdinHasData() bool {
	data, _ := os.Stdin.Stat()
	return (data.Mode() & os.ModeCharDevice) == 0
}

func normalizeInput(s string) string {
	// Trim surrounding whitespace.
	s = strings.TrimSpace(s)

	// Normalize line endings.
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")

	// Trim trailing whitespace from each line.
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRightFunc(line, unicode.IsSpace)
	}

	return strings.Join(lines, "\n")
}
