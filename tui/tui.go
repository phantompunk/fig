package tui

import (
	"fmt"

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
	top    int
	bottom int
}

type model struct {
	textInput  textinput.Model
	fonts      []item
	text       string
	cursor     int
	width      int
	viewHeight int
	ready      bool
	start      int
	end        int
	focusState focusState
}

func newModel() *model {
	textInput := textinput.New()
	return &model{textInput: textInput}
}

func (m model) Init() tea.Cmd { return loadFonts }

func (m model) View() string {
	var elements []string
	elements = append(elements, m.newTextInputView())

	if m.ready {
		for i := m.start; i < m.end && i < len(m.fonts); i++ {
			elements = append(elements, m.fontView(i))
		}
		// elements = append(elements, m.statusView())
	}

	elements = append(elements, "Press Ctrl+C to quit")
	return gloss.JoinVertical(gloss.Left, elements...)
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
				m.updateVisibleRange()
			}

		case "j", "down":
			if m.cursor < len(m.fonts)-1 && m.focusState == focusFontList {
				m.cursor++
				m.updateVisibleRange()
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
		m.updateVisibleRange()
		return m, nil

	case tea.WindowSizeMsg:
		m.width, m.viewHeight = msg.Width, msg.Height-8
		return m, nil
	}

	if m.textInput, cmd = m.textInput.Update(msg); cmd != nil {
		m.text = m.textInput.Value()
		return m, cmd
	}
	return m, nil
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
	totalHeight = 3
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

	for i, name := range fontNames {
		font, err := fig.Font(name)
		if err != nil {
			continue
		}

		items = append(items, item{
			name:   name,
			font:   font,
			index:  i,
			height: font.Height(),
		})
	}

	return fontsLoadedMsg{fonts: items}

	// all := []item{}
	// fonts := fig.ListFonts()
	// pos := 3
	// for i, name := range fonts {
	// 	itm := item{font: fig.Must(fig.Font(name)), index: i}
	// 	itm.height = itm.font.Height()
	// 	itm.top = pos
	// 	pos += itm.font.Height() + 3
	// 	itm.bottom = pos
	// 	all = append(all, itm)
	// }
	// return fontsLoadedMsg{fonts: all}
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
