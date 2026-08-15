package grpcclient

import (
	"log/slog"

	vesselv1 "github.com/deadsnxcks/dbcp/protos/gen/go/vessel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	Vessel vesselv1.VesselServiceClient
}

func New(addr string, log *slog.Logger) (*Clients, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Clients{
		Vessel: vesselv1.NewVesselServiceClient(conn),
	}, nil
}
