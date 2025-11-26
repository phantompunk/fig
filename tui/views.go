package tui

import (
	"fmt"
	"strings"

	gloss "github.com/charmbracelet/lipgloss"
)

func (m model) textInputBox() string {
	var b strings.Builder
	b.WriteString("Input")
	b.WriteString(m.textInput.View())

	if m.focusState == focusTextInput {
		return m.selectedBoxStyle().Render(b.String())
	}
	return m.boxStyle().Render(b.String())
}

func (m model) helpBox() string {
	controls := "↑/k:up ↓/j:down i:text ctrl+c:quit"
	list := fmt.Sprintf("%d/%d", m.cursor+1, len(m.fonts))
	spacingWidth := m.width - gloss.Width(controls) - gloss.Width(list) - 2
	if spacingWidth < 0 {
		spacingWidth = 0
	}
	spacing := strings.Repeat(" ", spacingWidth)
	content := gloss.JoinHorizontal(
		gloss.Top,
		controls,
		spacing,
		list,
	)
	return gloss.NewStyle().Padding(0, 1, 0, 1).Render(content)
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
		return fmt.Sprintf(" Count %d, selected: %d", len(m.fonts), m.cursor)
	}
	return fmt.Sprintf(" Count %d, selected: %d, %s, height: %d, vh: %d", len(m.fonts), m.cursor, m.fonts[m.cursor].font.Name(), m.fonts[m.cursor].height, m.viewHeight)
}
