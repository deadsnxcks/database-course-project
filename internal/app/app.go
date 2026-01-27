package app

import (
	"context"
	grpcapp "dbcp/internal/app/grpc"
	cargoservice "dbcp/internal/services/cargo"
	cargotypeservice "dbcp/internal/services/cargo-type"
	operationservice "dbcp/internal/services/operation"
	opercargoservice "dbcp/internal/services/opercargo"
	reportservice "dbcp/internal/services/report"
	storagelocservice "dbcp/internal/services/storageloc"
	vesselservice "dbcp/internal/services/vessel"
	"dbcp/internal/storage/postgresql"
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