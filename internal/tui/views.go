package tui

import (
	"fmt"
	"strings"

	gloss "charm.land/lipgloss/v2"
)

func (m model) topBar() string {
	appName := gloss.NewStyle().Bold(true).Foreground(gloss.Color("#C4C7D4")).Render("fig")

	var tags []string
	for _, tag := range tagCycle {
		if tag == m.activeTag {
			tags = append(tags, gloss.NewStyle().Foreground(gloss.Color("#C4C7D4")).Bold(true).Render("["+tag+"]"))
		} else {
			tags = append(tags, gloss.NewStyle().Foreground(gloss.Color("#626784")).Render(tag))
		}
	}
	tagSection := strings.Join(tags, " ")

	sep := gloss.NewStyle().Foreground(gloss.Color("#626784")).Render("|")
	var textLabel string
	if m.focusState == focusTextInput {
		textLabel = gloss.NewStyle().Foreground(gloss.Color("#C4C7D4")).Bold(true).Render("text")
	} else {
		textLabel = gloss.NewStyle().Foreground(gloss.Color("#626784")).Render("text")
	}
	right := sep + "   " + textLabel + m.textInput.View()

	leftPart := appName + "   " + tagSection
	// spacingWidth := max(m.width-gloss.Width(leftPart)-gloss.Width(right)-2, 1)
	spacing := strings.Repeat(" ", 3)

	content := leftPart + spacing + right
	bar := gloss.NewStyle().Padding(0, 1, 0, 1).Render(content)
	underline := gloss.NewStyle().Foreground(gloss.Color("#626784")).Render(strings.Repeat("‾", m.width))
	return gloss.JoinVertical(gloss.Left, bar, underline)
}

func (m model) helpBox() string {
	var controls string
	switch m.focusState {
	case focusFilter:
		controls = "Enter apply   Esc cancel"
	case focusTextInput:
		controls = "Enter apply   Esc cancel   ^u clear"
	default:
		controls = "↑/k ↓/j navigate   / filter   i edit text   f favorite   c copy   q quit" // f favorite   ? help
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
		// Border(gloss.RoundedBorder()).
		BorderForeground(gloss.Color("#C4C7D4")).
		Foreground(gloss.Color("#C4C7D4")).
		Bold(true)
}

func (m model) boxStyle() gloss.Style {
	return gloss.NewStyle().
		Width(m.width-4).
		// Border(gloss.HiddenBorder()).
		// BorderForeground(gloss.Color("#626784")).
		Padding(0, 1, 0, 1).
		Foreground(gloss.Color("#626784"))
}

func (m model) statusView() string {
	var status string
	if m.focusState == focusFilter {
		status = " filter" + m.filterInput.View()
	} else if m.copyMsg != "" {
		status = " " + m.copyMsg
	} else if m.filterQuery != "" {
		status = fmt.Sprintf(" filter: %q — %d/%d fonts", m.filterQuery, len(m.filteredFonts), len(m.fonts))
	}
	return gloss.NewStyle().Foreground(gloss.Color("#626784")).Render(status)
}
