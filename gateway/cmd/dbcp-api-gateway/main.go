package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/deadsnxcks/dbcp/gateway/internal/app"
	"github.com/deadsnxcks/dbcp/gateway/internal/config"
	"github.com/deadsnxcks/dbcp/gateway/internal/lib/logger/sl"
	"github.com/deadsnxcks/dbcp/gateway/internal/lib/logger/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting gateway", slog.String("address", cfg.HTTPServer.Address))

	application, err := app.New(log, *cfg)
	if err != nil {
		log.Error("failed to create application", sl.Err(err))
		os.Exit(1)
	}
	
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go application.HTTPServer.MustRun()

	<-ctx.Done()
	log.Info("shutdown signal received")

	shutDownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTPServer.ShutdownTimeout)
	defer cancel()

	if err := application.Stop(shutDownCtx); err != nil {
		log.Error("failed to stop application", sl.Err(err))
		os.Exit(1)
	}
	
	log.Info("gateway stopped")
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
