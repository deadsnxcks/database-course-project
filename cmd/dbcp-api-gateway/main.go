package main

import (
	"dbcp-api-gateway/internal/config"
	cargotype "dbcp-api-gateway/internal/http-server/handlers/cargo-type"
	"dbcp-api-gateway/internal/http-server/handlers/vessel"
	"dbcp-api-gateway/internal/lib/logger/slogpretty"
	cargotypev1 "dbcp-api-gateway/protos/gen/go/cargotype"
	vesselv1 "dbcp-api-gateway/protos/gen/go/vessel"
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
		AllowedOrigins:   []string{"*"}, // или конкретный фронт: "http://localhost:3000"
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // 5 минут
	})
	router.Use(corsMiddleware)

	conn, err := grpc.NewClient(cfg.GRPC.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to connect grpc", err)
	}
	defer conn.Close()

	vesselClient := vesselv1.NewVesselServiceClient(conn)
	vesselHandler := vessel.New(log, vesselClient)

	cargoTypeClient := cargotypev1.NewCargoTypeServiceClient(conn)
	cargoTypeHandler := cargotype.New(log, cargoTypeClient)

	router.Route("/vessels", func(r chi.Router) {
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

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", err.Error())
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