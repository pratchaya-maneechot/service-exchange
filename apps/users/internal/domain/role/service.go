package role

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/observability"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

type RoleCacheService struct {
	roleReader      RoleReader
	logger          *slog.Logger
	metricsRecorder observability.MetricsRecorder
	refreshInterval time.Duration

	rolesByName map[RoleName]Role
	rolesByID   map[uint]Role
	rolesMutex  sync.RWMutex
	initOnce    sync.Once

	stopRefreshChan chan struct{}
	wg              sync.WaitGroup
}

func NewRoleCacheService(
	reader RoleReader,
	logger *slog.Logger,
	metricsRecorder observability.MetricsRecorder,
	refreshInterval time.Duration,
) *RoleCacheService {
	return &RoleCacheService{
		roleReader:      reader,
		logger:          logger,
		metricsRecorder: metricsRecorder,
		refreshInterval: refreshInterval,
		rolesByName:     make(map[RoleName]Role),
		rolesByID:       make(map[uint]Role),
		stopRefreshChan: make(chan struct{}),
	}
}

func (rcm *RoleCacheService) InitAndStartRefresh(parentCtx context.Context) error { // Renamed param for clarity
	var initErr error
	rcm.initOnce.Do(func() {
		initErr = rcm.loadRolesFromDB(parentCtx) // Initial load can use parentCtx
		if initErr != nil {
			rcm.logger.Error("Failed initial load of roles from DB", "error", initErr)
			return
		}
		rcm.logger.Info("Initial roles loaded successfully.")
		rcm.metricsRecorder.RecordRoleCacheInitSuccess()

		rcm.wg.Add(1)
		go func() {
			defer rcm.wg.Done()
			// Create a cancellable child context for the refresh loop
			refreshCtx, refreshCancel := context.WithCancel(context.Background())
			defer refreshCancel() // Ensure this child context is cancelled when goroutine exits

			ticker := time.NewTicker(rcm.refreshInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					rcm.logger.Info("Refreshing roles cache...")
					// Pass the refreshCtx to loadRolesFromDB
					if err := rcm.loadRolesFromDB(refreshCtx); err != nil {
						rcm.logger.Error("Failed to refresh roles cache", "error", err)
						rcm.metricsRecorder.RecordRoleCacheRefreshFailed()
					} else {
						rcm.logger.Info("Roles cache refreshed successfully.")
						rcm.metricsRecorder.RecordRoleCacheRefreshSuccess()
					}
				case <-rcm.stopRefreshChan: // Explicit stop signal
					rcm.logger.Info("Stopping roles cache refresh goroutine by stop signal.")
					return
				case <-parentCtx.Done(): // Listen to parent context cancellation as well
					rcm.logger.Info("Stopping roles cache refresh goroutine due to parent context cancellation.")
					return
				}
			}
		}()
	})
	return initErr
}

func (rcm *RoleCacheService) Stop() {
	rcm.logger.Info("Signaling roles cache refresh goroutine to stop.")
	close(rcm.stopRefreshChan)
	rcm.wg.Wait()
	rcm.logger.Info("Roles cache refresh goroutine stopped.")
}

func (rcm *RoleCacheService) loadRolesFromDB(ctx context.Context) error {
	tr := otel.Tracer("domain-role-cache")
	ctx, span := tr.Start(ctx, "RoleCacheService.loadRolesFromDB")
	defer span.End()

	allRoles, err := rcm.roleReader.GetAllRoles(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch roles from DB")
		return fmt.Errorf("failed to fetch all roles from source: %w", err)
	}

	newRolesByName := make(map[RoleName]Role)
	newRolesByID := make(map[uint]Role)

	for _, r := range allRoles {
		newRolesByID[*r.ID] = r
		newRolesByName[r.Name] = r
	}

	rcm.rolesMutex.Lock()
	rcm.rolesByName = newRolesByName
	rcm.rolesByID = newRolesByID
	rcm.rolesMutex.Unlock()

	span.SetStatus(codes.Ok, "Roles loaded from DB successfully")
	rcm.metricsRecorder.RecordRoleCacheLoadCount(float64(len(allRoles)))

	return nil
}

func (rcm *RoleCacheService) GetRoleByName(name RoleName) (Role, error) {
	rcm.rolesMutex.RLock()
	defer rcm.rolesMutex.RUnlock()
	if rcm.rolesByName == nil {
		rcm.metricsRecorder.RecordRoleCacheMiss("uninitialized")
		return Role{}, fmt.Errorf("role cache not initialized")
	}
	role, ok := rcm.rolesByName[name]
	if !ok {
		rcm.metricsRecorder.RecordRoleCacheMiss("not_found")
		return Role{}, ErrRoleNotFound
	}
	rcm.metricsRecorder.RecordRoleCacheHit()
	return role, nil
}

func (rcm *RoleCacheService) GetRoleByID(id uint) (Role, error) {
	rcm.rolesMutex.RLock()
	defer rcm.rolesMutex.RUnlock()

	if rcm.rolesByID == nil {
		rcm.metricsRecorder.RecordRoleCacheMiss("uninitialized")
		return Role{}, fmt.Errorf("role cache not initialized")
	}
	role, ok := rcm.rolesByID[id]
	if !ok {
		rcm.metricsRecorder.RecordRoleCacheMiss("not_found")
		return Role{}, ErrRoleNotFound
	}
	rcm.metricsRecorder.RecordRoleCacheHit()
	return role, nil
}
