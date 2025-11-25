package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/phantompunk/fig"
)

type focusState int

const (
	focusFontList focusState = iota
	focusTextInput
)

type item struct {
	name   string
	font   *fig.FigFont
	index  int
	height int
}

type model struct {
	textInput   textinput.Model
	fonts       []item
	text        string
	cursor      int
	width       int
	viewHeight  int
	ready       bool
	start       int
	end         int
	totalHeight int
	offset      int
	focusState  focusState
}

func newModel() *model {
	textInput := textinput.New()
	return &model{textInput: textInput}
}

func (m model) Init() tea.Cmd { return loadFonts }

func (m model) View() string {
	inputBox := m.textInputBox()
	visible := m.renderPreviews()
	status := m.statusView()
	helpBox := m.helpBox()
	return gloss.JoinVertical(gloss.Left, inputBox, visible, status, helpBox)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "k", "up":
			if m.cursor > 0 && m.focusState == focusFontList {
				m.cursor--
				m.ensureSelectedVisible()
			}

		case "j", "down":
			if m.cursor < len(m.fonts)-1 && m.focusState == focusFontList {
				m.cursor++
				m.ensureSelectedVisible()
			}

		case "i":
			if m.focusState == focusFontList {
				m.toggleFocusState()
				return m, nil
			}

		case "enter":
			if m.focusState == focusTextInput {
				m.toggleFocusState()
			}
		}

	case fontsLoadedMsg:
		m.fonts = msg.fonts
		m.ready = true
		m.cursor = 0
		m.ensureSelectedVisible()
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width

		const (
			inputBoxHeight = 3
			statusHeight   = 1
			helpBoxHeight  = 2
		)
		m.viewHeight = msg.Height - inputBoxHeight - statusHeight - helpBoxHeight
		if m.viewHeight < 5 {
			m.viewHeight = 5
		}
		return m, nil
	}

	if m.textInput, cmd = m.textInput.Update(msg); cmd != nil {
		m.text = m.textInput.Value()
		return m, cmd
	}
	return m, nil
}

func (m *model) visibleRange() (start, startOffset, end int) {
	if len(m.fonts) == 0 {
		return 0, 0, 0
	}

	off := m.offset
	start = 0

	for start < len(m.fonts) && off >= m.fonts[start].height {
		off -= m.fonts[start].height
		start++
	}

	if start >= len(m.fonts) {
		return len(m.fonts) - 1, 0, len(m.fonts) - 1
	}

	startOffset = off
	remaining := m.viewHeight
	remaining -= (m.fonts[start].height - startOffset)

	end = start

	for remaining > 0 && end+1 < len(m.fonts) {
		end++
		remaining -= m.fonts[end].height
	}

	return
}

func (m *model) ensureSelectedVisible() {
	start, _, end := m.visibleRange()

	if m.cursor < start {
		m.offset = 0
		for i := 0; i < m.cursor; i++ {
			m.offset += m.fonts[i].height
		}
		return
	}

	if m.cursor > end {
		offset := 0
		for i := 0; i <= m.cursor; i++ {
			offset += m.fonts[i].height
		}
		m.offset = offset - m.viewHeight

		if m.offset < 0 {
			m.offset = 0
		}
		return
	}
}

func (m *model) updateVisibleRange() {
	if len(m.fonts) == 0 {
		m.start = 0
		m.end = 0
		return
	}

	// Ensure cursor is within bounds
	if m.cursor >= len(m.fonts) {
		m.cursor = len(m.fonts) - 1
	}

	// Find the range of fonts that fit in the visible area
	totalHeight := 0
	start := 0
	end := 0

	// Find start index - scroll to show cursor at top if needed
	for i := 0; i < len(m.fonts); i++ {
		if i == m.cursor {
			start = i
			break
		}
	}

	// Find end index - fill viewport from start position
	totalHeight = 0
	for i := start; i < len(m.fonts) && totalHeight < m.viewHeight; i++ {
		totalHeight += m.fonts[i].height + 3
		end = i + 1
	}

	// Ensure end doesn't exceed array bounds
	if end > len(m.fonts) {
		end = len(m.fonts)
	}

	m.start = start
	m.end = end
}

type fontsLoadedMsg struct {
	fonts []item
}

func loadFonts() tea.Msg {
	fontNames := fig.ListFonts()
	items := make([]item, 0, len(fontNames))

	for i, name := range fontNames[:20] {
		font, err := fig.Font(name)
		if err != nil {
			continue
		}

		items = append(items, item{
			name:   name,
			font:   font,
			index:  i,
			height: font.Height() + 3,
		})
	}

	return fontsLoadedMsg{fonts: items}
}

func (m model) PreviewFont(index int) string {
	if index < 0 || index >= len(m.fonts) {
		return "Invalid font index"
	}
	tmpl := "%s\n%s"
	fname := gloss.NewStyle().Italic(true)
	output := gloss.NewStyle().PaddingLeft(4)
	if m.text == "" {
		return fmt.Sprintf(tmpl, fname.Render(m.fonts[index].font.Name()), output.Render(m.fonts[index].font.Render(m.fonts[index].font.Name())))
	}
	return fmt.Sprintf(tmpl, fname.Render(m.fonts[index].font.Name()), output.Render(m.fonts[index].font.Render(m.text)))
}

func (m *model) toggleFocusState() {
	switch m.focusState {
	case focusTextInput:
		m.focusState = focusFontList
		m.textInput.Blur()
	case focusFontList:
		m.focusState = focusTextInput
		m.textInput.Focus()
	}
}

func Start() error {
	m := newModel()

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		return fmt.Errorf("running tui: %w", err)
	}

	return nil
}

func (m model) renderPreviews() string {
	var b strings.Builder

	if m.ready {
		start, startOff, end := m.visibleRange()
		for i := start; i <= end; i++ {
			preview := m.fontViewOG(i)

			// Clip the first item if needed
			if i == start && startOff > 0 {
				lines := strings.Split(preview, "\n")
				lines = lines[startOff:]
				preview = strings.Join(lines, "\n")
			}

			b.WriteString(preview)
			if i < end {
				b.WriteString("\n")
			}
		}
	}

	return b.String()
}

func (m model) fontViewOG(index int) string {
	preview := m.PreviewFont(index)
	if index == m.cursor {
		return m.selectedBoxStyle().Render(preview)
	}
	return m.boxStyle().Render(preview)
}
