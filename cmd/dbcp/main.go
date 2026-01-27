package main

import (
	"context"
	"dbcp/internal/app"
	"dbcp/internal/config"
	"dbcp/internal/lib/logger/slogpretty"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal 	= "local"
	envDev 		= "dev"
	envProd 	= "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	ctx := context.Background()
	application := app.New(log, cfg.GRPC.Port, cfg.DBConnString, ctx)

	go func () {
		application.GRPCServer.MustRun()
	}()

	log.Info("application start!")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<- stop
	application.GRPCServer.Stop()
	log.Info("Gracefully stopped")	
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