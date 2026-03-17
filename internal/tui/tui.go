package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	gloss "charm.land/lipgloss/v2"
	"github.com/atotto/clipboard"
	"github.com/phantompunk/fig/internal/font"
	"github.com/phantompunk/fig/internal/render"
)

type focusState int

const (
	focusFontList focusState = iota
	focusTextInput
	focusFilter
)

type item struct {
	name   string
	index  int
	height int
}

type model struct {
	engine        *render.Engine
	textInput     textinput.Model
	filterInput   textinput.Model
	fonts         []item
	filteredFonts []item
	filterQuery   string
	text          string
	cursor        int
	width         int
	height        int
	viewHeight    int
	ready         bool
	offset        int
	focusState    focusState
	align         render.Alignment
	copyMsg       string
}

func newModel() *model {
	textInput := textinput.New()
	filterInput := textinput.New()
	filterInput.Placeholder = "filter fonts..."
	return &model{
		textInput:   textInput,
		filterInput: filterInput,
		engine:      render.New(font.BundledLoader()),
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
		inputBoxHeight  = 3
		filterBoxHeight = 3
		statusHeight    = 1
		helpBoxHeight   = 1
	)
	footerHeight := statusHeight + helpBoxHeight
	contentHeight := m.height - inputBoxHeight - footerHeight
	if m.focusState == focusFilter {
		contentHeight -= filterBoxHeight
	}

	mainContent := gloss.NewStyle().
		Height(contentHeight).
		Render(previewContent)

	var parts []string
	parts = append(parts, inputBox)
	if m.focusState == focusFilter {
		parts = append(parts, m.filterBox())
	}
	parts = append(parts, mainContent, footer)
	content := gloss.JoinVertical(gloss.Left, parts...)

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		m.copyMsg = ""
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			if m.focusState != focusFilter && m.focusState != focusTextInput {
				return m, tea.Quit
			}

		case "esc":
			if m.focusState == focusFilter {
				m.focusState = focusFontList
				m.filterInput.Blur()
				m.filterQuery = ""
				m.applyFilter()
				return m, nil
			}
			return m, tea.Quit

		case "k", "up":
			if m.cursor > 0 && m.focusState == focusFontList {
				m.cursor--
				m.ensureSelectedVisible()
			}

		case "j", "down":
			if m.cursor < len(m.filteredFonts)-1 && m.focusState == focusFontList {
				m.cursor++
				m.ensureSelectedVisible()
			}

		case "g":
			if m.focusState == focusFontList {
				m.cursor = 0
				m.ensureSelectedVisible()
			}

		case "G":
			if m.focusState == focusFontList {
				m.cursor = len(m.filteredFonts) - 1
				m.ensureSelectedVisible()
			}

		case "ctrl+u":
			if m.focusState == focusFontList {
				m.cursor = m.pageMove(-m.viewHeight / 2)
				m.ensureSelectedVisible()
			}

		case "ctrl+d":
			if m.focusState == focusFontList {
				m.cursor = m.pageMove(m.viewHeight / 2)
				m.ensureSelectedVisible()
			}

		case "/":
			if m.focusState == focusFontList {
				m.focusState = focusFilter
				m.filterInput.SetValue("")
				m.filterQuery = ""
				m.filterInput.Focus()
				return m, nil
			}

		case "c":
			if m.focusState == focusFontList && len(m.filteredFonts) > 0 {
				rendered := m.renderedOutput()
				if err := clipboard.WriteAll(rendered); err != nil {
					m.copyMsg = "copy failed"
				} else {
					m.copyMsg = "copied!"
				}
			}

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
			} else if m.focusState == focusFilter {
				m.focusState = focusFontList
				m.filterInput.Blur()
				return m, nil
			}
		}

	case fontsLoadedMsg:
		m.fonts = msg.fonts
		m.ready = true
		m.cursor = 0
		m.applyFilter()
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

	if m.focusState == focusTextInput {
		if m.textInput, cmd = m.textInput.Update(msg); cmd != nil {
			m.text = m.textInput.Value()
			return m, cmd
		}
	} else if m.focusState == focusFilter {
		if m.filterInput, cmd = m.filterInput.Update(msg); cmd != nil {
			newQuery := m.filterInput.Value()
			if newQuery != m.filterQuery {
				m.filterQuery = newQuery
				m.applyFilter()
				m.cursor = 0
				m.ensureSelectedVisible()
			}
			return m, cmd
		}
	}

	return m, nil
}

func (m *model) applyFilter() {
	if m.filterQuery == "" {
		m.filteredFonts = m.fonts
		return
	}
	q := strings.ToLower(m.filterQuery)
	result := make([]item, 0)
	for _, f := range m.fonts {
		if strings.Contains(strings.ToLower(f.name), q) {
			result = append(result, f)
		}
	}
	m.filteredFonts = result
}

func (m *model) visibleRange() (start, startOffset, end int) {
	if len(m.filteredFonts) == 0 {
		return 0, 0, 0
	}

	off := m.offset
	start = 0

	for start < len(m.filteredFonts) && off >= m.filteredFonts[start].height {
		off -= m.filteredFonts[start].height
		start++
	}

	if start >= len(m.filteredFonts) {
		return len(m.filteredFonts) - 1, 0, len(m.filteredFonts) - 1
	}

	startOffset = off
	remaining := m.viewHeight
	remaining -= (m.filteredFonts[start].height - startOffset)

	end = start

	for remaining > 0 && end+1 < len(m.filteredFonts) {
		end++
		remaining -= m.filteredFonts[end].height
	}

	return
}

func (m *model) ensureSelectedVisible() {
	if len(m.filteredFonts) == 0 {
		m.offset = 0
		return
	}

	if m.cursor >= len(m.filteredFonts) {
		m.cursor = len(m.filteredFonts) - 1
	}

	if m.cursor < 0 {
		m.cursor = 0
	}

	heightUpToCursor := 0
	for i := 0; i <= m.cursor; i++ {
		heightUpToCursor += m.filteredFonts[i].height
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

// pageMove advances the cursor by delta lines worth of font entries (negative = up).
// It walks through fonts accumulating heights until the target line count is reached.
func (m *model) pageMove(delta int) int {
	if len(m.filteredFonts) == 0 {
		return 0
	}
	if delta < 0 {
		remaining := -delta
		c := m.cursor
		for c > 0 && remaining > 0 {
			c--
			remaining -= m.filteredFonts[c].height
		}
		return c
	}
	remaining := delta
	c := m.cursor
	for c < len(m.filteredFonts)-1 && remaining > 0 {
		remaining -= m.filteredFonts[c].height
		c++
	}
	return c
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
	if index < 0 || index >= len(m.filteredFonts) {
		return "Invalid font index"
	}

	name := m.filteredFonts[index].name
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

// renderedOutput returns the plain rendered text for the currently selected font.
func (m model) renderedOutput() string {
	if len(m.filteredFonts) == 0 || m.cursor >= len(m.filteredFonts) {
		return ""
	}
	name := m.filteredFonts[m.cursor].name
	text := m.text
	if text == "" {
		text = name
	}
	out, err := m.engine.Render(text, render.RenderOptions{FontName: name})
	if err != nil {
		return ""
	}
	return out
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
	if len(m.filteredFonts) == 0 {
		return gloss.NewStyle().Foreground(gloss.Color("#626784")).Render("  no fonts match")
	}

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
