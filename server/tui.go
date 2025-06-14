package main

import (
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	windowStyle = lipgloss.NewStyle().
			BorderForeground(lipgloss.Color("6")).
			Width(0).
			Padding(0, 1).
			Align(lipgloss.Center).
			Border(lipgloss.ThickBorder())
	keymap = mainKeymap{
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
		Start: key.NewBinding(
			key.WithKeys("s,r"),
			key.WithHelp("s/r", "Encender servicio."),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "Ocultar/Mostrar ayuda."),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q/esc", "Salir."),
		),
	}
	ProgColors = ColorPalette{
		Dark: Palette{
			Primary:   "#42B2D7",
			Secondary: "#8F7EE7",
			Tertiary:  "#B6C2CF",
		},
		Light: Palette{
			Primary:   "#227D9B",
			Secondary: "#352C63",
			Tertiary:  "#172B4D",
		},
	}
)

type Palette struct {
	Primary     string
	Secondary   string
	Tertiary    string
	FocusFg     string
	FocusBg     string
	FocusBorder string
}
type ColorPalette struct {
	Dark  Palette
	Light Palette
}

type mainKeymap struct {
	Start key.Binding
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Help  key.Binding
	Quit  key.Binding
	CtrlC key.Binding
	CtrlL key.Binding
}

func (k mainKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.CtrlL}
}

func (k mainKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Start, k.Enter},
		{k.Help, k.CtrlC, k.CtrlL, k.Quit},
	}
}

type MainModel struct {
	// Element []string
	// Cursor  int
	// _keymap mainKeymap
	tabs []TabModel
	help help.Model
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch _msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(_msg, keymap.Down):
			var (
				i int
				t TabModel
			)
			for i, t = range m.tabs {
				if t.isActive {
					// t.isActive = false
					break
				}
			}
			m.tabs[i].isActive = false
			if i == 0 {
				m.tabs[len(m.tabs)-1].isActive = true
			} else {
				i -= 1
				m.tabs[i].isActive = true
			}
		case key.Matches(_msg, keymap.Up):
			var (
				i int
				t TabModel
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
		case key.Matches(_msg, keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case key.Matches(_msg, keymap.Quit):
			return m, tea.Quit
			// return m, tea.Printf("Testing here!")
		// case key.Matches(_msg, keymap.Start):

		case key.Matches(_msg, keymap.CtrlC):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *MainModel) View() string {
	var tabs []string
	for _, t := range m.tabs {
		tabs = append(tabs, t.View())
	}

	helpView := m.help.View(keymap)
	mainWindow := windowStyle.Render(lipgloss.JoinHorizontal(lipgloss.Center, tabs...))

	height := 8 - strings.Count(mainWindow, "\n") - strings.Count(helpView, "\n")
	return mainWindow + strings.Repeat("\n", height) + helpView
}

func main() {
	tabs := []TabModel{
		{
			title:    "uno",
			isActive: true,
		},
		{
			title: "dos",
		},
		{
			title: "tres",
		},
	}
	m := MainModel{
		tabs: tabs,
		help: help.New(),
	}
	go func() {
		for i := range 20 {
			time.Sleep(2 * time.Second)
			tea.Printf("[%d] Testing here!", i)
			// log.Info("Testing Here!", "attempt", i)
		}
	}()
	if _, err := tea.NewProgram(&m).Run(); err != nil {
		os.Exit(1)
	}
}
