package main

import (
	"context"
	"errors"
	"log" // Using standard log for fatal errors before slog is initialized
	"net/http"
	"os"
	"os/signal"
	"syscall"

	// For context.WithTimeout
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal"
)

func main() {
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()
	app, err := internal.InitializeApp(rootCtx)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	defer app.Cleanup()
	logger := app.Logger
	logger.Info("Application initialized successfully.")
	if app.MetricServer != nil {
		go func() {
			if err := app.MetricServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) { // Add http.ErrServerClosed check
				logger.Error("failed to start metrics server", "error", err)
			}
		}()
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-sigChan
		logger.Info("Received OS signal, initiating graceful shutdown", "signal", s.String())
		rootCancel()
	}()
	logger.Info("Starting users server",
		"address", app.Config.Server.Address,
		"environment", app.Config.Environment,
		"version", app.Config.Version)

	if err := app.Server.Start(rootCtx); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			logger.Info("gRPC server stopped gracefully due to context cancellation.")
		} else {
			logger.Error("gRPC server encountered an unrecoverable error", "error", err)
			os.Exit(1)
		}
	}
	logger.Info("gRPC server stopped. Proceeding with application shutdown handler.")
	logger.Info("Application shutdown completed successfully.")
}
