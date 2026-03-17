package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	gloss "charm.land/lipgloss/v2"
	"github.com/phantompunk/fig/internal/font"
	"github.com/phantompunk/fig/internal/render"
)

type focusState int

const (
	focusFontList focusState = iota
	focusTextInput
)

type item struct {
	name   string
	index  int
	height int
}

type model struct {
	engine     *render.Engine
	textInput  textinput.Model
	fonts      []item
	text       string
	cursor     int
	width      int
	height     int
	viewHeight int
	ready      bool
	offset     int
	focusState focusState
	align      render.Alignment
}

func newModel() *model {
	textInput := textinput.New()
	return &model{
		textInput: textInput,
		engine:    render.New(font.BundledLoader()),
	}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		return loadFontsWithEngine(m.engine)
	}
}

func (m model) View() tea.View {
	inputBox := m.textInputBox()
	status := m.statusView()
	helpBox := m.helpBox()

	if !m.ready {
		v := tea.NewView("loading fonts")
		v.AltScreen = true
		return v
	}

	previewContent := m.renderPreviews()

	// Limit preview to viewHeight lines to ensure we don't overflow
	lines := strings.Split(previewContent, "\n")
	if len(lines) > m.viewHeight {
		lines = lines[:m.viewHeight]
		previewContent = strings.Join(lines, "\n")
	}

	footer := gloss.JoinVertical(gloss.Left, status, helpBox)

	const (
		inputBoxHeight = 3
		statusHeight   = 1
		helpBoxHeight  = 1
	)
	footerHeight := statusHeight + helpBoxHeight
	contentHeight := m.height - inputBoxHeight - footerHeight

	mainContent := gloss.NewStyle().
		Height(contentHeight).
		Render(previewContent)

	content := gloss.JoinVertical(gloss.Left, inputBox, mainContent, footer)
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
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

		case "g":
			m.cursor = 0
			m.ensureSelectedVisible()

		case "G":
			m.cursor = len(m.fonts) - 1
			m.ensureSelectedVisible()

		case "a":
			if m.focusState == focusFontList {
				m.align = (m.align + 1) % 3
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
		m.height = msg.Height

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
	if len(m.fonts) == 0 {
		return
	}

	if m.cursor >= len(m.fonts) {
		m.cursor = len(m.fonts) - 1
	}

	if m.cursor < 0 {
		m.cursor = 0
	}

	heightUpToCursor := 0
	for i := 0; i <= m.cursor; i++ {
		heightUpToCursor += m.fonts[i].height
	}

	if heightUpToCursor <= m.viewHeight {
		m.offset = 0
		return
	}

	safeMargin := 3
	safeViewHeight := m.viewHeight - safeMargin
	if safeViewHeight < 5 {
		safeViewHeight = 5
	}

	m.offset = heightUpToCursor - safeViewHeight
	if m.offset < 0 {
		m.offset = 0
	}
}

type fontsLoadedMsg struct {
	fonts []item
}

func loadFontsWithEngine(engine *render.Engine) tea.Msg {
	fontNames, err := engine.ListFonts()
	if err != nil {
		return fontsLoadedMsg{fonts: nil}
	}

	items := make([]item, 0, len(fontNames))
	for i, name := range fontNames {
		h, err := engine.FontHeight(name)
		if err != nil {
			continue
		}
		items = append(items, item{
			name:   name,
			index:  i,
			height: h + 3,
		})
	}

	return fontsLoadedMsg{fonts: items}
}

func (m model) PreviewFont(index int) string {
	if index < 0 || index >= len(m.fonts) {
		return "Invalid font index"
	}

	name := m.fonts[index].name
	opts := render.RenderOptions{
		FontName: name,
		Align:    m.align,
		Width:    m.width,
	}

	text := m.text
	if text == "" {
		text = name
	}

	rendered, err := m.engine.Render(text, opts)
	if err != nil {
		rendered = fmt.Sprintf("[render error: %v]", err)
	}

	fname := gloss.NewStyle().Italic(true)
	output := gloss.NewStyle().PaddingLeft(4)
	return fmt.Sprintf("%s\n%s", fname.Render(name), output.Render(rendered))
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

	if _, err := tea.NewProgram(m).Run(); err != nil {
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

			if i == start && startOff > 0 {
				lines := strings.Split(preview, "\n")
				safeOffset := startOff
				if safeOffset > len(lines) {
					safeOffset = len(lines)
				}
				lines = lines[safeOffset:]
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
