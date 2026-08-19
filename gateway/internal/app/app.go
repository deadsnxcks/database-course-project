package app

import (
	"context"
	"fmt"
	"log/slog"

	httpserverapp "github.com/deadsnxcks/dbcp/gateway/internal/app/http_server"
	"github.com/deadsnxcks/dbcp/gateway/internal/config"
	grpcclient "github.com/deadsnxcks/dbcp/gateway/internal/grpc"
	"github.com/deadsnxcks/dbcp/gateway/internal/lib/logger/sl"
	cargov1 "github.com/deadsnxcks/dbcp/protos/gen/go/cargo"
	cargotypev1 "github.com/deadsnxcks/dbcp/protos/gen/go/cargotype"
	operationv1 "github.com/deadsnxcks/dbcp/protos/gen/go/operation"
	opercargov1 "github.com/deadsnxcks/dbcp/protos/gen/go/opercargo"
	reportv1 "github.com/deadsnxcks/dbcp/protos/gen/go/report"
	storagelocv1 "github.com/deadsnxcks/dbcp/protos/gen/go/storageloc"
	vesselv1 "github.com/deadsnxcks/dbcp/protos/gen/go/vessel"
	"google.golang.org/grpc"
)

type App struct {
	HTTPServer *httpserverapp.App
	conn *grpc.ClientConn
}

func New(
	log *slog.Logger,
	cfg config.Config,
) (*App, error) {
	const op = "app.New"
	conn, err := grpcclient.New(cfg.GRPC.Address, log)
	if err != nil {
		log.Error("failed to create grpc client", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	vesselClient := vesselv1.NewVesselServiceClient(conn)
	cargoTypeClient := cargotypev1.NewCargoTypeServiceClient(conn)
	cargoClient := cargov1.NewCargoServiceClient(conn)
	operationsClient := operationv1.NewOperationServiceClient(conn)
	operCargoClient := opercargov1.NewOperationCargoServiceClient(conn)
	storageLocClient := storagelocv1.NewStorageLocationServiceClient(conn)
	reportClient := reportv1.NewReportServiceClient(conn)

	httpserver := httpserverapp.New(
		log,
		vesselClient,
		cargoClient,
		cargoTypeClient,
		operationsClient,
		operCargoClient,
		storageLocClient,
		reportClient,
		cfg.HTTPServer,
	)

	return &App{
		HTTPServer: httpserver,
		conn: conn,
	}, nil
}

func (a *App) Stop(ctx context.Context) error {
	if err := a.HTTPServer.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop HTTP server: %w", err)
	}

	if err := a.conn.Close(); err != nil {
		return fmt.Errorf("failed to close gRPC connection: %w", err)
	}

	return nil
}