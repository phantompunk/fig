package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/phantompunk/stencil/pkg/font"
	"github.com/phantompunk/stencil/pkg/stencil"
)

var docStyle = lipgloss.NewStyle().Margin(0, 2)
type item struct {
	title string
	desc  string
}

func (i item) Title() string { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
  case tea.WindowSizeMsg:
    h, v := docStyle.GetFrameSize()
    m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func Render() {
	items := []list.Item{}

	for _, fname := range font.ListFonts() {
		st, _ := stencil.NewStencil(fname, fname)
		items = append(items, item{title: fname, desc: st.Draw()})
	}

  d := list.NewDefaultDelegate()
  d.SetHeight(10)
  m := model{list: list.New(items, d, 0, 0)}
  m.list.Title = "Fonts"

  p := tea.NewProgram(m, tea.WithAltScreen())
  if _, err := p.Run(); err != nil {
    fmt.Println("Error running program", err)
    os.Exit(1)
  } 
}
