package tui

import (
	"fmt"
	"strings"

	gloss "charm.land/lipgloss/v2"
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

func (m model) filterBox() string {
	var b strings.Builder
	b.WriteString("Filter /")
	b.WriteString(m.filterInput.View())
	return m.selectedBoxStyle().Render(b.String())
}

func (m model) helpBox() string {
	var controls string
	if m.focusState == focusFilter {
		controls = "enter:confirm  esc:clear"
	} else {
		controls = "↑/k:up ↓/j:down gg:start G:end ^U:pgup ^D:pgdn tab:tag /:filter i:text a:align c:copy q:quit"
	}
	list := fmt.Sprintf("%d/%d  ", m.cursor+1, len(m.filteredFonts))
	spacingWidth := max(m.width-gloss.Width(controls)-gloss.Width(list)-2, 0)
	spacing := strings.Repeat(" ", spacingWidth)
	content := gloss.JoinHorizontal(
		gloss.Top,
		controls,
		spacing,
		list,
	)
	return gloss.NewStyle().Foreground(gloss.Color("#626784")).Padding(0, 1, 0, 1).Render(content)
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
	var status string
	if m.copyMsg != "" {
		status = " " + m.copyMsg
	} else if m.filterQuery != "" && m.activeTag != "all" {
		status = fmt.Sprintf(" [%s]  Filter: %q — %d/%d fonts", m.activeTag, m.filterQuery, len(m.filteredFonts), len(m.fonts))
	} else if m.activeTag != "all" {
		status = fmt.Sprintf(" [%s] — %d/%d fonts", m.activeTag, len(m.filteredFonts), len(m.fonts))
	} else if m.filterQuery != "" {
		status = fmt.Sprintf(" Filter: %q — %d/%d fonts", m.filterQuery, len(m.filteredFonts), len(m.fonts))
	} else if len(m.filteredFonts) == 0 || m.cursor >= len(m.filteredFonts) {
		status = fmt.Sprintf(" Count %d, selected: %d", len(m.filteredFonts), m.cursor)
	} else {
		status = fmt.Sprintf(" Count %d, selected: %d, %s, height: %d, vh: %d", len(m.filteredFonts), m.cursor, m.filteredFonts[m.cursor].name, m.filteredFonts[m.cursor].height, m.viewHeight)
	}
	return gloss.NewStyle().Foreground(gloss.Color("#626784")).Render(status)
}
