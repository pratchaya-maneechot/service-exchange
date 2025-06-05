package main

import (
	"context"
	"log"
	"os"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal"
)

func main() {
	app, err := internal.InitializeApp()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := app.Infra.Metric.Start(); err != nil {
			app.Logger.Error("failed to start metrics server", "error", err)
		}
	}()

	go app.Shutdown.HandleSignals(cancel)

	app.Logger.Info("starting gRPC server",
		"address", app.Config.Server.Address,
		"environment", app.Config.Environment,
		"version", app.Config.Version)

	if err := app.Server.Start(ctx); err != nil {
		app.Logger.Error("gRPC server error", "error", err)
		os.Exit(1)
	}

	app.Logger.Info("server shutdown completed successfully")
}
