package main

import tea "github.com/charmbracelet/bubbletea"

// func newPromptModel() textinput.Model {
// 	promptModel.PromptStyle = lipgloss.NewStyle().
// 		SetString("|> ").
// 		Foreground(lipgloss.Color("#9fe0f0"))
// 	keyMap := keymap{
// 		next: key.NewBinding(
// 			key.WithKeys("tab"),
// 			key.WithHelp("tab", "next"),
// 		),
// 		prev: key.NewBinding(
// 			key.WithKeys("shift+tab"),
// 			key.WithHelp("shift+tab", "prev"),
// 		),
// 		add: key.NewBinding(
// 			key.WithKeys("ctrl+n"),
// 			key.WithHelp("ctrl+n", "add an editor"),
// 		),
// 		remove: key.NewBinding(
// 			key.WithKeys("ctrl+w"),
// 			key.WithHelp("ctrl+w", "remove an editor"),
// 		),
// 		quit: key.NewBinding(
// 			key.WithKeys("esc", "ctrl+c"),
// 			key.WithHelp("esc", "quit"),
// 		),
// 	}
// 	promptModel := textinput.Model{
// 	  keymap: keyMap
// 	}
// 	return promptModel
// }

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

func (m MenuModel) promptUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}
