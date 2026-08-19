package httpserverapp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/deadsnxcks/dbcp/gateway/internal/config"
	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/router"
	"github.com/deadsnxcks/dbcp/gateway/internal/lib/logger/sl"
	cargov1 "github.com/deadsnxcks/dbcp/protos/gen/go/cargo"
	cargotypev1 "github.com/deadsnxcks/dbcp/protos/gen/go/cargotype"
	operationv1 "github.com/deadsnxcks/dbcp/protos/gen/go/operation"
	opercargov1 "github.com/deadsnxcks/dbcp/protos/gen/go/opercargo"
	reportv1 "github.com/deadsnxcks/dbcp/protos/gen/go/report"
	storagelocv1 "github.com/deadsnxcks/dbcp/protos/gen/go/storageloc"
	vesselv1 "github.com/deadsnxcks/dbcp/protos/gen/go/vessel"
)

type App struct {
	log        *slog.Logger
	server *http.Server
}

func New(
	log *slog.Logger,
	vesselClient vesselv1.VesselServiceClient,
	cargoClient cargov1.CargoServiceClient,
	cargoTypeClient cargotypev1.CargoTypeServiceClient,
	operationClient operationv1.OperationServiceClient,
	opercargoClient opercargov1.OperationCargoServiceClient,
	storagelocClient storagelocv1.StorageLocationServiceClient,
	reportClient reportv1.ReportServiceClient,
	cfg config.HTTPServer,
) *App {
	router := router.New(
		log,
		vesselClient,
		cargoClient,
		cargoTypeClient,
		operationClient,
		opercargoClient,
		storagelocClient,
		reportClient,
	)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &App{
		log: log,
		server: srv,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		a.log.Error("failed to run HTTP server", sl.Err(err))
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "httpserverapp.Run"

	a.log.Info("starting HTTP server", slog.String("address", a.server.Addr))

	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	const op = "httpserverapp.Stop"

	a.log.Info("stopping HTTP server", slog.String("address", a.server.Addr))

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}