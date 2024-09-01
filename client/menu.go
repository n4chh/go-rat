package main

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/iortego42/go-rat/grpcapi"
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

type menu struct {
	elements []string
	cursor   int
}

func initMenu(elements []string, cursor int) menu {
	return menu{
		elements: elements,
		cursor:   cursor,
	}
}

func (m menu) Init() tea.Cmd {
	return nil
}

func (m menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.elements)-1 {
				m.cursor++
			}
		case "enter":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m menu) View() string {
	l := list.New(m.elements)
	l.ItemStyle(ITEMSTYLE)
	l.EnumeratorStyle(SELECTORSTYLE)
	l = l.Enumerator(m.selector)
	return fmt.Sprintln(l)
}

func (m menu) selector(items list.Items, i int) string {
	if i == m.cursor {
		return "|> "
	}
	return ""
}

func implantsMenu(client grpcapi.AdminClient) {
	var (
		m     menu
		ctx   = context.Background()
		items []string
	)

	availableImplants, err := client.GetImplants(ctx, nil)
	if err != nil {
		logger.Fatal(err)
	}
	for _, v := range availableImplants.Implants {
		logger.Debug(v.Id)
		items = append(items, v.Id)
	}
	logger.Debug(items)
	m = initMenu(items, 0)
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		logger.Fatal("There was an error on menu", "menu", err)
	}
}
