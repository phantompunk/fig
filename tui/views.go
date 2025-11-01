package tui

import (
	"strings"

	gloss "github.com/charmbracelet/lipgloss"
)

var (
	fname  = gloss.NewStyle().Italic(true)
	output = gloss.NewStyle().PaddingLeft(4)
)

func (a model) newTextInputView() string {
	var b strings.Builder
	b.WriteString("Input")
	b.WriteString(a.textInput.View())

	if a.focusState == focusTextInput {
		return a.selectedBoxStyle().Render(b.String())
	}
	return a.boxStyle().Render(b.String())
}

func (a model) fontView(index int) string {
	preview := a.PreviewFont(index)
	if index == a.selectedFont {
		return a.selectedBoxStyle().Render(preview)
	}
	return a.boxStyle().Render(preview)
}

func (a model) selectedBoxStyle() gloss.Style {
	return a.boxStyle().
		BorderForeground(gloss.Color("#C4C7D4")).
		Foreground(gloss.Color("#C4C7D4")).
		Bold(true)
}

func (a model) boxStyle() gloss.Style {
	return gloss.NewStyle().
		Width(a.width-4).
		Border(gloss.RoundedBorder()).
		BorderForeground(gloss.Color("#626784")).
		Padding(0, 1, 0, 1).
		Foreground(gloss.Color("#626784"))
}
