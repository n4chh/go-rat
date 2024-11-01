package main

import (
	"errors"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/grpc/status"
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

type ErrorMsg struct {
	Err           error
	GrpcStatusErr *status.Status
}

func (a *ClientApp) SendCommand() tea.Msg {
	var err error = nil

	a.Cmd, err = a.Client.RunCommand(a.Ctx, a.Cmd)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if !ok {
			log.Fatal("cant get status error")
			// handle error
		} else if statusErr.Message() == "implant closed" {
			return a.FetchImplants()
		}
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
		statusErr, ok := status.FromError(err)
		if !ok {
			return ErrorMsg{
				Err:           err,
				GrpcStatusErr: nil,
			}
		}
		return ErrorMsg{
			Err:           err,
			GrpcStatusErr: statusErr,
		}
	}
	for _, v := range availableImplants.Implants {
		// logger.Debug(v.Id)
		items = append(items, v.Id)
	}
	if len(items) == 0 {
		return ErrorMsg{
			Err:           errors.New("no implants"),
			GrpcStatusErr: nil,
		}
	}

	return MenuMsg{
		Menu: NewMenu(items, 0),
	}
}
