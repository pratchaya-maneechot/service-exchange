package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/role"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/postgres"
)

type ShutdownHandler struct {
	logger *slog.Logger

	dbPool       *postgres.DBPool
	roleCacheSvc *role.RoleCacheService
}

func NewShutdownHandler(
	logger *slog.Logger,
	dbPool *postgres.DBPool,
	roleCacheSvc *role.RoleCacheService,
) *ShutdownHandler {
	return &ShutdownHandler{
		logger:       logger,
		dbPool:       dbPool,
		roleCacheSvc: roleCacheSvc,
	}
}

func (sh *ShutdownHandler) Handle(shutdownCtx context.Context) error {
	errChan := make(chan error, 1)

	go func() {
		sh.logger.Info("Starting graceful shutdown of application components...")
		var shutdownErrors []error

		if sh.dbPool != nil {
			sh.logger.Info("Closing database connections...")
			sh.dbPool.Close()
			sh.logger.Info("Database connections closed.")
		}

		if sh.roleCacheSvc != nil {
			// sh.roleCacheSvc.Stop()
			//  if err := sh.roleCacheSvc.Stop(ctx); err != nil {
			//     shutdownErrors = append(shutdownErrors, fmt.Errorf("failed to stop role cache service: %w", err))
			// } else {
			//     sh.logger.Info("Role cache manager stopped.")
			// }
			sh.logger.Info("Role cache manager stopping (if it has background goroutines)...")
		}

		if len(shutdownErrors) > 0 {
			errChan <- fmt.Errorf("errors during shutdown: %v", shutdownErrors)
		} else {
			errChan <- nil
		}
	}()

	select {
	case <-shutdownCtx.Done():
		return fmt.Errorf("graceful shutdown timed out or cancelled: %w", shutdownCtx.Err())
	case err := <-errChan:
		return err
	}
}
