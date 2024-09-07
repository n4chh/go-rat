package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/iortego42/go-rat/log"
)

var logger = log.InitLogger()

var ITEMSTYLE lipgloss.Style = lipgloss.NewStyle().
	// Border(lipgloss.NormalBorder()).
	// BorderForeground(lipgloss.Color("#9fe0f0")).
	Foreground(lipgloss.NoColor{})

var SELECTORSTYLE lipgloss.Style = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#9fe0f0")).
	Bold(true)

	// Padding(1)

// TYPES
type SelectMsg struct{ Implant string }

type MenuModel struct {
	Elements []string
	Cursor   int
	List     *list.List
}

// ----
func initMenu(elements []string, cursor int) *MenuModel {
	return &MenuModel{
		Elements: elements,
		Cursor:   cursor,
		List:     list.New(elements).ItemStyle(ITEMSTYLE).EnumeratorStyle(SELECTORSTYLE),
	}
}

func (m *MenuModel) Init() tea.Cmd {
	return nil
}

func (m *MenuModel) Quit() tea.Msg {
	var msg GoBackMsg = true
	return msg
}

// TODO: Funcion de recarga de tipo tea.Cmd para actualizar los implants,
// implementar spinners de carga o algo similar mientras
// espera la respuesta del servidor

func (m *MenuModel) Choose() tea.Msg {
	var msg SelectMsg
	msg.Implant = m.Elements[m.Cursor]
	return msg
}

func (m *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil
	switch msg := msg.(type) {
	case tea.KeyMsg:
		s := msg.String()
		switch s {
		case "q", "esq":
			cmd = m.Quit
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Elements)-1 {
				m.Cursor++
			}
		case "enter":
			cmd = m.Choose
		}
	}
	return m, cmd
}

func (m *MenuModel) View() string {
	m.List = m.List.Enumerator(m.selector)
	return fmt.Sprintln(m.List)
}

func (m *MenuModel) selector(items list.Items, i int) string {
	if i == m.Cursor {
		return "|> "
	}
	return ""
}
