package main

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
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
	mode int8
)

const (
	PromptImplant mode = iota
	SelectImplant
)

func (m mode) String() string {
	switch m {
	case PromptImplant:
		return "Implant Prompt"
	case SelectImplant:
		return "Select Implant"
	}
	return "unknown"
}

// app

type AppKeyMap struct {
	CtrlC key.Binding
}

type ClientApp struct {
	Menu   *MenuModel
	Prompt textinput.Model
	State  mode
	KeyMap AppKeyMap
	Client grpcapi.AdminClient
}

func (a *ClientApp) Init() tea.Cmd {
	a.FetchImplants()
	return nil
}

func (a *ClientApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch _msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(_msg, a.KeyMap.CtrlC):
			return a, tea.Quit
		}
	case SelectMsg:
		a.State = PromptImplant

	case GoBackMsg:
		if a.State == SelectImplant {
			return a, tea.Quit
		}
		a.State = SelectImplant
	}
	switch a.State {
	case SelectImplant:
		m, cmd := a.Menu.Update(msg)
		newMenu, ok := m.(*MenuModel)
		if !ok {
			log.Fatal("Bad assertion", "menu", m)
		}
		a.Menu = newMenu
		cmds = append(cmds, cmd)
		return a, tea.Batch(cmds...)
		// case PromptImplant:
		// 	m, cmd := a.Menu.Update(msg)
		// 	newMenu, ok := m.(MenuModel)
		// 	if !ok {
		// 		log.Fatal("Bad assertion", "menu", m)
		// 	}
		// 	a.Menu = newMenu
		// 	cmds = append(cmds, cmd)
		// 	return a, tea.Batch(cmds...)
	default:
		a.FetchImplants()
		logger.Info(a)
	}
	return a, nil
}

func (a *ClientApp) View() string {
	switch a.State {
	case SelectImplant:
		return a.Menu.View()
	case PromptImplant:
		return a.Prompt.View()
	}
	return ""
}

func (a *ClientApp) FetchImplants() tea.Msg {
	var (
		ctx   = context.Background()
		items []string
	)
	availableImplants, err := a.Client.GetImplants(ctx, nil)
	if err != nil {
		logger.Fatal(err)
	}
	for _, v := range availableImplants.Implants {
		logger.Debug(v.Id)
		items = append(items, v.Id)
	}
	a.Menu = initMenu(items, 0)
	a.State = SelectImplant
	return nil
}

func NewClientApp(conn *grpc.ClientConn) *ClientApp {
	_keymap := AppKeyMap{
		CtrlC: key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("Ctrl C", "Force Quit")),
	}
	a := &ClientApp{
		Prompt: textinput.Model{},
		Client: grpcapi.NewAdminClient(conn),
		KeyMap: _keymap,
		// State:  SelectImplant,
	}
	return a
}

func main() {
	var (
		opts []grpc.DialOption
		conn *grpc.ClientConn
		err  error
		app  *ClientApp
	)
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err = grpc.NewClient("127.0.0.1:9090", opts...)
	if err != nil {
		log.Fatal("[!] No se pudo establecer conexión con el servidor principal.", "ERROR", err)
	}
	defer conn.Close()
	app = NewClientApp(conn)
	p := tea.NewProgram(app)
	if _, err = p.Run(); err != nil {
		log.Fatal("Ocurrió un error en la ejecución del programa.", "error", err)
	}
}
