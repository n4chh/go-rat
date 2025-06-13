package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var PROMPTSTYLE = lipgloss.NewStyle().
	SetString("ID:").
	Foreground(lipgloss.Color("5")).
	Bold(true)

type PromptReadyMsg bool

type PromptKeyMap struct {
	Quit, Enter key.Binding
}

type PromptModel struct {
	Implant string
	Ti      textinput.Model
	KeyMap  PromptKeyMap
}

func (p *PromptModel) PromptReady() tea.Msg {
	var ret PromptReadyMsg = true
	return ret
}

func NewPromptModel(implant string) PromptModel {
	_keymap := PromptKeyMap{
		Quit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "Salir de la interfaz del implant y volver al menú."),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithKeys("enter", "Manda un comando al implant seleccionado"),
		),
	}
	p := PromptModel{}
	p.Implant = implant
	p.Ti = textinput.New()
	p.Ti.Prompt = ""
	p.Ti.PromptStyle = lipgloss.NewStyle().
		SetString("|>").
		Foreground(lipgloss.Color("5"))
	p.Ti.Placeholder = "Introduzca un comando."
	p.Ti.PlaceholderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("default"))
	p.KeyMap = _keymap
	return p
}

func (p *PromptModel) Quit() tea.Msg {
	var msg GoBackMsg = true
	p.Implant = ""
	return msg
}

func (p PromptModel) Init() tea.Cmd {
	return textinput.Blink
}

func (p PromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil
	switch _msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(_msg, p.KeyMap.Quit):
			cmd = p.Quit
		case key.Matches(_msg, p.KeyMap.Enter):
			value := p.Ti.Value()
			p.Ti.Reset()
			cmd = func() tea.Msg {
				return SendCmdMsg{
					ID:    p.Implant,
					Input: value,
				}
			}
		default:
			p.Ti, cmd = p.Ti.Update(msg)
		}
	}
	return p, cmd
}

func (p PromptModel) View() string {
	s := ""
	s += PROMPTSTYLE.Render("[" + p.Implant + "]")
	s += "\n"
	s += p.Ti.View()
	return s
}

// func mainLoop(client grpcapi.AdminClient, implant uuid.UUID) {
// 	var (
// 		ctx context.Context
// 		cmd *grpcapi.Command
// 		err error
// 	)
// 	ctx = context.Background()
// 	cmd = new(grpcapi.Command)
// 	PROMPT.SetString("[", implant.String(), "]\n|>")
//
// 	cmd.Id = implant.String()
// 	s := bufio.NewScanner(os.Stdin)
// 	for {
// 		cmd.Out = ""
// 		fmt.Print(PROMPT.SetString("ID: ").Bold(true).Render())
// 		fmt.Println(PROMPT.SetString("").UnsetForeground().Render(implant.String()))
// 		fmt.Print(PROMPT.Render())
// 		s.Scan()
// 		strCmd := strings.Trim(s.Text(), " \n")
// 		if strCmd == "" {
// 			continue
// 		}
// 		if strCmd == "exit" {
// 			fmt.Println(PROMPT.SetString("[+]").Render(), "Bye")
// 			return
// 		}
// 		cmd.In = strCmd
// 		cmd, err = client.RunCommand(ctx, cmd)
// 		if cmd == nil && err == nil {
// 			log.Info("Implant cerrado", "id", implant.String())
// 			return
// 		}
// 		if err != nil {
// 			statusErr, ok := status.FromError(err)
// 			if !ok {
// 				log.Fatal("Error creando statusErr")
// 			}
// 			if statusErr.Message() == "implant closed" {
// 				logger.Info("Implant Cerrado", "id", implant.String())
// 				return
// 			}
// 			log.Fatal(err)
// 		}
// 		fmt.Print(cmd.Out)
// 		if cmd.Out != "" && cmd.Out[len(cmd.Out)-1] != '\n' && cmd.In != "clear" {
// 			fmt.Println(PROMPT.Bold(true).Background(lipgloss.Color("#9fe0f0")).Foreground(lipgloss.Color("#1a1a1a")).SetString("%").Render())
// 		}
// 	}
// }
