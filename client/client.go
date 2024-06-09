package main

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/iortego42/go-rat/grpcapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"time"
)

var prompt = lipgloss.NewStyle().
	SetString("|> ").
	Foreground(lipgloss.Color("#9fef00"))

var viewportStyle = lipgloss.NewStyle().
	//Foreground(lipgloss.Color("#ffffa0")).
	Padding(0, 1).
	Border(lipgloss.RoundedBorder())

var bye = lipgloss.NewStyle().
	SetString("Bye").
	Border(lipgloss.RoundedBorder()).
	Padding(0, 1).
	Bold(true)

var outputStyle = lipgloss.NewStyle().
	Padding(0).
	TabWidth(0)

type model struct {
	input    textinput.Model
	viewport viewport.Model
	client   grpcapi.AdminClient
	ctx      context.Context
	cmd      *grpcapi.Command
	id       *grpcapi.Identity
}

func initModel() *model {
	var m *model = &model{}
	m.input = textinput.New()
	m.input.Prompt = ""
	m.input.PromptStyle = prompt
	m.input.Focus()
	m.viewport = viewport.New(80, 20)
	m.viewport.Style = viewportStyle
	m.cmd = new(grpcapi.Command)
	m.cmd.Id = os.Args[1]
	m.ctx = context.Background()
	return m
}

func (m model) exec() tea.Msg {
	var _exec = func() tea.Msg {
		var err error
		m.cmd, err = m.client.RunCommand(m.ctx, m.cmd)
		if err != nil {
			return err.Error()
		}
		return m.cmd.Out
	}
	return _exec()
}

func (m model) _run() string {
	var err error

	m.cmd, err = m.client.RunCommand(m.ctx, m.cmd)
	if err != nil {
		return err.Error()
	}
	return m.cmd.Out
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var iCmd, vpCmd tea.Cmd

	m.input, iCmd = m.input.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.cmd.In = m.input.Value()
			m.input.Reset()
			m.exec()
			m.viewport.SetContent(m._run())
			m.viewport.GotoBottom()
			//return m, m.exec
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.viewport.SetContent("")
		}
		//default:
		//	return m, nil
	}
	return m, tea.Batch(iCmd, vpCmd)
}

func (m model) View() string {
	//if m.cmd.Out == "" {
	//	return m.input.View()
	//} else {
	//	out := outputStyle.Render(m.cmd.Out)
	//	m.cmd.Out = ""
	//	return out
	//}
	return fmt.Sprintf(
		"%s\n%s\n%s",
		"R.A.T. Console",
		m.viewport.View(),
		m.input.View())
}

func main() {
	var (
		opts []grpc.DialOption
		err  error
		conn *grpc.ClientConn
		//client grpcapi.AdminClient
		m *model
	)
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err = grpc.NewClient("127.0.0.1:9090", opts...)
	if err != nil {
		log.Fatal("[!] No se pudo establecer conexión con el servidor principal.", "ERROR", err)
	}
	defer conn.Close()
	m = initModel()
	m.client = grpcapi.NewAdminClient(conn)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(bye.Render())
	time.Sleep(1000)
	//client = grpcapi.NewAdminClient(conn)
	//mainLoop(client)
}
