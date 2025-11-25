package tui

import (
	"fmt"
	"strings"

	gloss "github.com/charmbracelet/lipgloss"
)

// func (m model) renderPreviews() string {
// 	var b strings.Builder
//
// 	if m.ready {
// 		start, startOff, end := m.visibleRange()
// 		for i := start; i <= end; i++ {
// 			preview := m.fontViewOG(i)
//
// 			// Clip the first item if needed
// 			if i == start && startOff > 0 {
// 				lines := strings.Split(preview, "\n")
// 				lines = lines[startOff:]
// 				preview = strings.Join(lines, "\n")
// 			}
//
// 			b.WriteString(preview)
// 			if i < end {
// 				b.WriteString("\n")
// 			}
// 		}
// 	}
//
// 	return b.String()
// }

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
	helpBox := "Press Ctrl+C to quit"
	return m.helpBoxStyle().Render(helpBox)
}

// func (m model) fontViewOG(index int) string {
// 	preview := m.PreviewFont(index)
// 	if index == m.cursor {
// 		return m.selectedBoxStyle().Render(preview)
// 	}
// 	return m.boxStyle().Render(preview)
// }

func (m model) fontView(index int) string {
	preview := m.PreviewFont(index)
	if index == m.cursor {
		return m.selectedBoxStyle().Render(preview)
	}
	return m.boxStyle().Render(preview)
}

func (m model) helpBoxStyle() gloss.Style {
	return gloss.NewStyle().
		Width(m.width-4).
		Padding(0, 1, 0, 1)
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
