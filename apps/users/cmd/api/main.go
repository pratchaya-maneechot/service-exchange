package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal"
)

func main() {
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	app, err := internal.InitializeApp(rootCtx)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	app.Logger.Info("Application initialized successfully.")

	go func() {
		if err := app.Infra.Metric.Start(); err != nil {
			// Log only, don't os.Exit() here, as it's a background goroutine
			app.Logger.Error("failed to start metrics server", "error", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-sigChan
		app.Logger.Info("Received OS signal, initiating graceful shutdown", "signal", s.String())
		rootCancel()
	}()

	app.Logger.Info("Starting gRPC server",
		"address", app.Config.Server.Address,
		"environment", app.Config.Environment,
		"version", app.Config.Version)

	if err := app.Server.Start(rootCtx); err != nil {
		// If the error is due to context cancellation (normal shutdown), don't treat as fatal
		if err == context.Canceled || err == context.DeadlineExceeded {
			app.Logger.Info("gRPC server stopped gracefully due to context cancellation.")
		} else {
			app.Logger.Error("gRPC server encountered an unrecoverable error", "error", err)
			os.Exit(1)
		}
	}

	app.Logger.Info("gRPC server stopped. Proceeding with application shutdown handler.")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), app.Config.Server.ShutdownTimeout)
	defer shutdownCancel()

	if err := app.Shutdown.Handle(shutdownCtx); err != nil {
		app.Logger.Error("Application graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	app.Logger.Info("Application shutdown completed successfully.")
	os.Exit(0)
}
