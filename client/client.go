package main

import (
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/iortego42/go-rat/grpcapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

func main() {
	var (
		opts         []grpc.DialOption
		conn         *grpc.ClientConn
		admin_client grpcapi.AdminClient
		err          error
	)
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err = grpc.NewClient("127.0.0.1:9090", opts...)
	if err != nil {
		log.Fatal("[!] No se pudo establecer conexión con el servidor principal.", "ERROR", err)
	}
	defer conn.Close()
	admin_client = grpcapi.NewAdminClient(conn)
	ctx := context.Background()
	var cmd = new(grpcapi.Command)
	cmd.In = os.Args[1]
	cmd.Id = os.Args[2]
	cmd, err = admin_client.RunCommand(ctx, cmd)
	log.Debug("[*] Resultado recibido.")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cmd.Out)
}
