package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TabModel struct {
	title    string
	isActive bool
}

var (
	inactiveTabStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				Padding(0, 1).
				BorderForeground(lipgloss.AdaptiveColor{Dark: ProgColors.Dark.Tertiary, Light: ProgColors.Light.Tertiary})
	activeTabStyle = inactiveTabStyle.Border(lipgloss.ThickBorder()).BorderForeground(lipgloss.AdaptiveColor{Light: ProgColors.Light.Primary, Dark: ProgColors.Dark.Primary})
)

func (t TabModel) Init() tea.Cmd {
	return nil
}

func (t TabModel) Update(msg tea.Msg) (TabModel, tea.Cmd) {
	return t, nil
}

func (t TabModel) View() string {
	if t.isActive {
		return activeTabStyle.Render(t.title)
	}
	return inactiveTabStyle.Render(t.title)
}
