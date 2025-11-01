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

type model struct {
	textInput    textinput.Model
	fonts        []*fig.FigFont
	text         string
	selectedFont int
	width        int
	height       int
	ready        bool
	focusState   focusState
}

func newModel() *model {
	textInput := textinput.New()
	return &model{textInput: textInput}
}

func (a model) Init() tea.Cmd { return loadFonts }

func (a model) View() string {
	var elements []string
	elements = append(elements, a.newTextInputView())
	for i := range a.fonts {
		elements = append(elements, a.fontView(i))
	}

	elements = append(elements, "Press Ctrl+C to quit")
	return gloss.JoinVertical(gloss.Left, elements...)
}

func (a model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return a, tea.Quit

		case "k", "up":
			if a.selectedFont > 0 && a.focusState == focusFontList {
				a.selectedFont--
			}

		case "j", "down":
			if a.selectedFont < len(a.fonts)-1 && a.focusState == focusFontList {
				a.selectedFont++
			}

		case "i":
			if a.focusState == focusFontList {
				a.toggleFocusState()
				return a, nil
			}

		case "enter":
			if a.focusState == focusTextInput {
				a.toggleFocusState()
			}
		}

	case fontsLoadedMsg:
		a.fonts = msg.fonts
		return a, nil

	case tea.WindowSizeMsg:
		a.width, a.height = msg.Width, msg.Width
		return a, nil
	}

	if a.textInput, cmd = a.textInput.Update(msg); cmd != nil {
		a.text = a.textInput.Value()
		return a, cmd
	}
	return a, nil
}

type fontsLoadedMsg struct {
	fonts []*fig.FigFont
}

func loadFonts() tea.Msg {
	all := []*fig.FigFont{}
	fonts := fig.ListFonts()
	for _, name := range fonts[5:7] {
		all = append(all, fig.Must(fig.Font(name)))
	}
	return fontsLoadedMsg{fonts: all}
}

func (a model) PreviewFont(index int) string {
	tmpl := "%s\nn\t%s"
	fname := gloss.NewStyle().Italic(true)
	output := gloss.NewStyle().PaddingLeft(4)
	if a.text == "" {
		return fmt.Sprintf(tmpl, fname.Render(a.fonts[index].Name()), output.Render(a.fonts[index].Render(a.fonts[index].Name())))
	}
	return fmt.Sprintf(tmpl, fname.Render(a.fonts[index].Name()), output.Render(a.fonts[index].Render(a.text)))
}

func (a *model) toggleFocusState() {
	switch a.focusState {
	case focusTextInput:
		a.focusState = focusFontList
		a.selectedFont = 0
		a.textInput.Blur()
	case focusFontList:
		a.focusState = focusTextInput
		a.selectedFont = -1
		a.textInput.Focus()
	}
}

func Start() error {
	m := newModel()

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		return fmt.Errorf("running tui: %w", err)
	}

	return nil
}
