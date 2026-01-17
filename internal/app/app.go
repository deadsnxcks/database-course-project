package app

import (
	"context"
	grpcapp "dbcp/internal/app/grpc"
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

	grpcApp := grpcapp.New(log, vesselService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}