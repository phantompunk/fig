package tui

import (
	"fmt"
	"strings"

	gloss "github.com/charmbracelet/lipgloss"
)

var (
	fname  = gloss.NewStyle().Italic(true)
	output = gloss.NewStyle().PaddingLeft(4)
)

func (m model) newTextInputView() string {
	var b strings.Builder
	b.WriteString("Input")
	b.WriteString(m.textInput.View())

	if m.focusState == focusTextInput {
		return m.selectedBoxStyle().Render(b.String())
	}
	return m.boxStyle().Render(b.String())
}

func (m model) fontViewOG(index int) string {
	preview := m.PreviewFont(index)
	if index == m.cursor {
		return m.selectedBoxStyle().Render(preview)
	}
	return m.boxStyle().Render(preview)
}

func (m model) fontView(index int) string {
	preview := m.PreviewFont(index)
	if index == m.cursor {
		return m.selectedBoxStyle().Render(preview)
	}
	return m.boxStyle().Render(preview)
}

func (m model) selectedBoxStyle() gloss.Style {
	return m.boxStyle().
		BorderForeground(gloss.Color("#C4C7D4")).
		Foreground(gloss.Color("#C4C7D4")).
		Bold(true)
}

func (m model) boxStyle() gloss.Style {
	return gloss.NewStyle().
		Width(m.width-4).
		Border(gloss.RoundedBorder()).
		BorderForeground(gloss.Color("#626784")).
		Padding(0, 1, 0, 1).
		Foreground(gloss.Color("#626784"))
}

func (m model) statusView() string {
	if len(m.fonts) == 0 || m.cursor >= len(m.fonts) {
		return fmt.Sprintf("Count %d, selected: %d", len(m.fonts), m.cursor)
	}
	return fmt.Sprintf("Count %d, selected: %d, %s, height: %d", len(m.fonts), m.cursor, m.fonts[m.cursor].font.Name(), m.fonts[m.cursor].height)
}
