package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SendCmdMsg struct {
	ID    string
	Input string
}
type RecvCmdMsg struct {
	Output string
	Err    error
}
type MenuMsg struct {
	Menu MenuModel
}

func (a *ClientApp) SendCommand() tea.Msg {
	var err error = nil

	a.Cmd, err = a.Client.RunCommand(a.Ctx, a.Cmd)
	if err != nil {
		return RecvCmdMsg{
			Output: "",
			Err:    err,
		}
	}
	return RecvCmdMsg{
		Output: a.Cmd.Out,
		Err:    nil,
	}
}

// TODO: Debe de devolver un mensaje para que la funcion update actualice el estado
func (a *ClientApp) FetchImplants() tea.Msg {
	var items []string
	availableImplants, err := a.Client.GetImplants(a.Ctx, nil)
	if err != nil {
		logger.Fatal(err)
	}
	for _, v := range availableImplants.Implants {
		// logger.Debug(v.Id)
		items = append(items, v.Id)
	}
	return MenuMsg{
		Menu: NewMenu(items, 0),
	}
}
