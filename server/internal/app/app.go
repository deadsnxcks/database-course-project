package app

import (
	"context"
	grpcapp "github.com/deadsnxcks/dbcp/server/internal/app/grpc"
	cargoservice "github.com/deadsnxcks/dbcp/server/internal/services/cargo"
	cargotypeservice "github.com/deadsnxcks/dbcp/server/internal/services/cargo-type"
	operationservice "github.com/deadsnxcks/dbcp/server/internal/services/operation"
	opercargoservice "github.com/deadsnxcks/dbcp/server/internal/services/opercargo"
	reportservice "github.com/deadsnxcks/dbcp/server/internal/services/report"
	storagelocservice "github.com/deadsnxcks/dbcp/server/internal/services/storageloc"
	vesselservice "github.com/deadsnxcks/dbcp/server/internal/services/vessel"
	"github.com/deadsnxcks/dbcp/server/internal/storage/postgresql"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	connString string,
	ctx context.Context,
) *App {
	storage, err := postgresql.New(ctx, connString)
	if err != nil {
		panic(err)
	}

	vesselService := vesselservice.New(log, storage)
	cargoTypeService := cargotypeservice.New(log, storage)
	cargoService := cargoservice.New(log, storage)
	storageLocService := storagelocservice.New(log, storage)
	operationService := operationservice.New(log, storage)
	operCargoService := opercargoservice.New(log, storage)
	reportService := reportservice.New(log, storage)

	grpcApp := grpcapp.New(
		log, 
		vesselService, 
		cargoTypeService, 
		cargoService,
		storageLocService, 
		operationService,
		operCargoService,
		reportService,
		grpcPort,
	)

	return &App{
		GRPCServer: grpcApp,
	}
}