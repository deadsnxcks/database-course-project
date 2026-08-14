package main

import (
	"github.com/deadsnxcks/dbcp/gateway/internal/config"
	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/handlers/cargo"
	cargotype "github.com/deadsnxcks/dbcp/gateway/internal/http-server/handlers/cargo-type"
	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/handlers/operations"
	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/handlers/opercargo"
	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/handlers/report"
	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/handlers/storageloc"
	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/handlers/vessel"
	"github.com/deadsnxcks/dbcp/gateway/internal/lib/logger/sl"
	"github.com/deadsnxcks/dbcp/gateway/internal/lib/logger/slogpretty"
	cargov1 "github.com/deadsnxcks/dbcp/protos/gen/go/cargo"
	cargotypev1 "github.com/deadsnxcks/dbcp/protos/gen/go/cargotype"
	operationv1 "github.com/deadsnxcks/dbcp/protos/gen/go/operation"
	opercargov1 "github.com/deadsnxcks/dbcp/protos/gen/go/opercargo"
	reportv1 "github.com/deadsnxcks/dbcp/protos/gen/go/report"
	storagelocv1 "github.com/deadsnxcks/dbcp/protos/gen/go/storageloc"
	vesselv1 "github.com/deadsnxcks/dbcp/protos/gen/go/vessel"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	envLocal 	= "local"
	envDev 		= "dev"
	envProd 	= "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting gateway", slog.String("address", cfg.HTTPServer.Address))

	router := chi.NewRouter()

	corsMiddleware := cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(corsMiddleware)

	conn, err := grpc.NewClient(cfg.GRPC.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to connect grpc", sl.Err(err))
	}
	defer conn.Close()

	vesselClient := vesselv1.NewVesselServiceClient(conn)
	vesselHandler := vessel.New(log, vesselClient)

	cargoTypeClient := cargotypev1.NewCargoTypeServiceClient(conn)
	cargoTypeHandler := cargotype.New(log, cargoTypeClient)

	cargoClient := cargov1.NewCargoServiceClient(conn)
	cargoHandler := cargo.New(log, cargoClient)

	operationsClient := operationv1.NewOperationServiceClient(conn)
	operationsHandler := operations.New(log, operationsClient)

	operCargoClient := opercargov1.NewOperationCargoServiceClient(conn)
	operCargoHandler := opercargo.New(log, operCargoClient)

	storageLocClient := storagelocv1.NewStorageLocationServiceClient(conn)
	storageLocHandler := storageloc.New(log, storageLocClient)

	reportClient := reportv1.NewReportServiceClient(conn)
	reportHandler := report.New(log, reportClient)

	router.Route("/vessel", func(r chi.Router) {
		r.Get("/", vesselHandler.List())
		r.Post("/", vesselHandler.Create())
		r.Get("/{id}", vesselHandler.Get())
		r.Put("/{id}", vesselHandler.Update())
		r.Delete("/{id}", vesselHandler.Delete())
	})

	router.Route("/cargotype", func(r chi.Router) {
		r.Get("/", cargoTypeHandler.List())
		r.Post("/", cargoTypeHandler.Create())
		r.Get("/{id}", cargoTypeHandler.Get())
		r.Put("/{id}", cargoTypeHandler.Update())
		r.Delete("/{id}", cargoTypeHandler.Delete())
	})

	router.Route("/cargo", func(r chi.Router) {
		r.Get("/", cargoHandler.List())
		r.Post("/", cargoHandler.Create())
		r.Get("/{id}", cargoHandler.Get())
		r.Put("/{id}", cargoHandler.Update())
		r.Delete("/{id}", cargoHandler.Delete())
	})

	router.Route("/operation", func(r chi.Router) {
		r.Get("/", operationsHandler.List())
		r.Post("/", operationsHandler.Create())
		r.Get("/{id}", operationsHandler.Get())
		r.Put("/{id}", operationsHandler.Update())
		r.Delete("/{id}", operationsHandler.Delete())
	})

	router.Route("/opercargo", func(r chi.Router) {
		r.Get("/", operCargoHandler.List())
		r.Post("/", operCargoHandler.Create())
		r.Delete("/", operCargoHandler.Delete())
	})

	router.Route("/storageloc", func(r chi.Router) {
		r.Get("/", storageLocHandler.List())
		r.Post("/", storageLocHandler.Create())
		r.Get("/{id}", storageLocHandler.Get())
		r.Put("/{id}", storageLocHandler.Update())
		r.Delete("/{id}", storageLocHandler.Delete())
		r.Post("/{id}/use", storageLocHandler.Use())
		r.Post("/{id}/reset", storageLocHandler.Reset())
	})

	router.Method("Get", "/report-cargo-detail", reportHandler.CargoDetailReport())
	router.Method("Get", "/report-cargo-type", reportHandler.CargoTypeReport())

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}