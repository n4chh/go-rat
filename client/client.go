package main

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/iortego42/go-rat/grpcapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO: Revisar si merece la pena
// func handleStatusCodes(s *status.Status) int{
// 	if s == nil {
// 		return -1
// 	}
// 	if s.Message() == "implant closed" {
// 		return 0
// 	}
// }

// global types

type (
	GoBackMsg bool
	// Modes
	Mode int
)

const (
	// PromptImplant Mode = iota
	Default Mode = iota
	SelectImplant
	PromptImplant
	// PromptReady
	CmdOutRecived
	Error
)

func (m Mode) String() string {
	switch m {
	case PromptImplant:
		return "Introducir comando para implant"
		// return "Generando linea de comandos para implant"
	// case PromptReady:
	case SelectImplant:
		return "Seleccionar Implant"
	case CmdOutRecived:
		return "Resultado Recibido"
	case Default:
		return "Por defecto"
	}
	return "desconocido"
}

// app

type AppKeyMap struct {
	CtrlC key.Binding
}

type ClientApp struct {
	Menu   MenuModel
	Prompt PromptModel
	State  Mode
	KeyMap AppKeyMap
	Client grpcapi.AdminClient
	Ctx    context.Context
	Cmd    *grpcapi.Command
	Err    error
}

func (a ClientApp) Init() tea.Cmd {
	return a.FetchImplants
}

func (a ClientApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch _msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(_msg, a.KeyMap.CtrlC):
			return a, tea.Quit
		}
	case MenuMsg:
		a.Menu = _msg.Menu
		a.State = SelectImplant
	case SelectMsg:
		a.Prompt = NewPromptModel(_msg.Implant)
		cmds = append(cmds, PromptReady)
		return a, tea.Batch(cmds...)
	case PromptReadyMsg:
		a.State = PromptImplant
		return a, a.Prompt.Ti.Focus()
	case GoBackMsg:
		if a.State == SelectImplant {
			return a, tea.Quit
		}
		a.State = SelectImplant
	case SendCmdMsg:
		a.Cmd.In = _msg.Input
		a.Cmd.Id = _msg.ID
		cmds = append(cmds, a.SendCommand)
		return a, tea.Batch(cmds...)
	case RecvCmdMsg:
		a.State = CmdOutRecived
		if _msg.Err != nil {
			a.Err = _msg.Err
			a.State = Error
		}
	}
	switch a.State {
	case SelectImplant:
		m, cmd := a.Menu.Update(msg)
		newMenu, ok := m.(MenuModel)
		if !ok {
			logger.Fatal("Bad assertion", "menu", m)
		}
		a.Menu = newMenu
		cmds = append(cmds, cmd)
		return a, tea.Batch(cmds...)
	case PromptImplant:
		p, cmd := a.Prompt.Update(msg)
		newPrompt, ok := p.(PromptModel)
		if !ok {
			logger.Fatal("Bad assertion", "prompt", p)
		}
		a.Prompt = newPrompt
		cmds = append(cmds, cmd)
		return a, tea.Batch(cmds...)
	case CmdOutRecived:
		cmds = append(cmds, PromptReady)
		return a, tea.Batch(cmds...)
	case Error:
		a.State = Default
		return a, nil
	}
	return a, nil
}

func (a ClientApp) View() string {
	if a.Err != nil {
		return a.Err.Error()
	}
	switch a.State {
	case SelectImplant:
		return a.Menu.View()
	case PromptImplant:
		return a.Prompt.View()
	case CmdOutRecived:
		return a.Prompt.View() + "\n" + a.Cmd.Out
	}
	return ""
}

func NewClientApp(conn *grpc.ClientConn) ClientApp {
	_keymap := AppKeyMap{
		CtrlC: key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("Ctrl C", "Force Quit")),
	}
	a := ClientApp{
		Client: grpcapi.NewAdminClient(conn),
		Ctx:    context.Background(),
		Cmd:    new(grpcapi.Command),
		State:  Default,
		KeyMap: _keymap,
	}
	return a
}

func main() {
	var (
		opts []grpc.DialOption
		conn *grpc.ClientConn
		err  error
		app  ClientApp
	)
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err = grpc.NewClient("127.0.0.1:9090", opts...)
	if err != nil {
		logger.Fatal("[!] No se pudo establecer conexión con el servidor principal.", "ERROR", err)
	}
	defer conn.Close()
	app = NewClientApp(conn)
	p := tea.NewProgram(app)
	if _, err = p.Run(); err != nil {
		logger.Fatal("Ocurrió un error en la ejecución del programa.", "error", err)
	}
}
