package main

import (
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	windowStyle = lipgloss.NewStyle().
			BorderForeground(lipgloss.Color("6")).
			Padding(2).
			Align(lipgloss.Center).
			Border(lipgloss.NormalBorder())
	inactiveTabStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				Padding(0, 1).
				BorderForeground(lipgloss.Color("6"))
	activeTabStyle = inactiveTabStyle.BorderForeground(lipgloss.Color("3"))
	keymap         = mainKeymap{
		CtrlC: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("Ctrl C", "Forzar salida"),
		),
		CtrlL: key.NewBinding(
			key.WithKeys("ctrl+l"),
			key.WithHelp("Ctrl L", "Limpiar la pantalla"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "Mover el cursor hacia arriba."),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "Mover el cursor hacia abajo."),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Selecionar implant."),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q,esc", "Salir."),
		),
	}
)

type mainKeymap struct {
	Up, Down, Enter, Quit, CtrlC, CtrlL key.Binding
}
type tab struct {
	title    string
	isActive bool
}

func (t tab) View() string {
	if t.isActive {
		return activeTabStyle.Render(t.title)
	}
	return inactiveTabStyle.Render(t.title)
}

type mainModel struct {
	// Element []string
	// Cursor  int
	// _keymap mainKeymap
	tabs []tab
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch _msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(_msg, keymap.Up):
			var (
				i int
				t tab
			)
			for i, t = range m.tabs {
				if t.isActive {
					// t.isActive = false
					break
				}
			}
			m.tabs[i].isActive = false
			if i == len(m.tabs)-1 {
				m.tabs[0].isActive = true
			} else {
				i += 1
				m.tabs[i].isActive = true
			}
		case key.Matches(_msg, keymap.CtrlC):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *mainModel) View() string {
	var tabs []string
	for _, t := range m.tabs {
		tabs = append(tabs, t.View())
	}
	return windowStyle.Render(lipgloss.JoinHorizontal(lipgloss.Center, tabs...))
}

func main() {
	tabs := []tab{
		{
			title: "uno",
		},
		{
			title: "dos",
		},
		{
			title: "tres",
		},
	}
	m := mainModel{
		tabs: tabs,
	}
	if _, err := tea.NewProgram(&m).Run(); err != nil {
		os.Exit(1)
	}
}
