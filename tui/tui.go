package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/phantompunk/fig"
)

type FontPreview struct {
	Name   string
	Font   *fig.FigFont
	Output string
}

func (f FontPreview) Preview(text string) string {
	return f.Font.Render(text)
}

var allFonts = []FontPreview{{
	Name: "standard",
	Font: fig.Must(fig.Font("standard")),
}, {
	Name: "slant",
	Font: fig.Must(fig.Font("slant")),
}, {
	Name: "lean",
	Font: fig.Must(fig.Font("lean")),
}, {
	Name: "big",
	Font: fig.Must(fig.Font("big")),
}}

type fontModel struct {
	text          string
	fonts         []FontPreview
	selectedFont  int
	width, height int
	focusState    focusState
	newTextInput  textinput.Model
}

type focusState int

const (
	focusStateList focusState = iota
	focusStateTextInput
)

func (a fontModel) Init() tea.Cmd { return loadFonts }

func (a fontModel) View() string {
	var elements []string
	elements = append(elements, a.newTextInputView())
	for i := range a.fonts {
		elements = append(elements, a.fontView(i))
	}

	elements = append(elements, "Press Ctrl+C to quit")
	return gloss.JoinVertical(gloss.Left, elements...)
}

func (a fontModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case fontsLoadedMsg:
		a.fonts = msg.fonts
		return a, nil
	case tea.WindowSizeMsg:
		a.width, a.height = msg.Width, msg.Width
		return a, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit
		case "up":
			if a.selectedFont > 0 && a.focusState == focusStateList {
				a.selectedFont--
			}
		case "down":
			if a.selectedFont < len(a.fonts)-1 && a.focusState == focusStateList {
				a.selectedFont++
			}
		case "i":
			if a.focusState == focusStateList {
				a.toggleFocusState()
			}
		case "esc", "enter":
			if a.focusState == focusStateTextInput {
				a.text = a.newTextInput.Value()
			}
			a.toggleFocusState()
		}
	}

	if a.newTextInput, cmd = a.newTextInput.Update(msg); cmd != nil {
		a.text = a.newTextInput.Value()
		return a, cmd
	}
	return a, nil
}

type fontsLoadedMsg struct {
	fonts []FontPreview
}

func loadFonts() tea.Msg {
	all := []FontPreview{}
	fonts := fig.ListFonts()
	for _, name := range fonts[5:7] {
		all = append(all, FontPreview{Name: name, Font: fig.Must(fig.Font(name))})
	}
	return fontsLoadedMsg{fonts: all}
}

func (a fontModel) selectedBoxStyle() gloss.Style {
	return a.boxStyle().
		BorderForeground(gloss.Color("#00FF00")).
		Bold(true)
}

func (a fontModel) boxStyle() gloss.Style {
	return gloss.NewStyle().
		Width(a.width-4).
		Border(gloss.RoundedBorder()).
		BorderForeground(gloss.Color("#FFF")).
		Padding(0, 1, 0, 1).
		Foreground(gloss.Color("#FFF"))
}

func (a fontModel) PreviewFont(index int) string {
	tmpl := "%s\n%s"
	message := fmt.Sprintf(tmpl, a.fonts[index].Name, a.fonts[index].Font.Render(a.text))
	return message
}

func (a fontModel) fontView(index int) string {
	preview := a.PreviewFont(index)
	if index == a.selectedFont {
		return a.selectedBoxStyle().Render(preview)
	}
	return a.boxStyle().Render(preview)
}

func (a fontModel) newTextInputView() string {
	var b strings.Builder
	b.WriteString("Input ")
	b.WriteString(a.newTextInput.View())

	if a.focusState == focusStateTextInput {
		return a.selectedBoxStyle().Render(b.String())
	}
	return a.boxStyle().Render(b.String())
}

func (a *fontModel) toggleFocusState() {
	switch a.focusState {
	case focusStateTextInput:
		a.focusState = focusStateList
		a.selectedFont = 0
		a.newTextInput.Focus()
	case focusStateList:
		a.focusState = focusStateTextInput
		a.selectedFont = -1
		a.newTextInput.Focus()
	}
}

func Start() error {
	// m := fontModel{fonts: allFonts, text: "demo"}
	m := fontModel{text: "demo", newTextInput: textinput.New()}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		return fmt.Errorf("running interactive table: %w", err)
	}

	return nil
}
