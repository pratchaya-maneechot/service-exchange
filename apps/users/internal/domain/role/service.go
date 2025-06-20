package role

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type RoleCacheService struct {
	roleReader      RoleReader
	logger          *slog.Logger
	refreshInterval time.Duration

	rolesByName map[RoleName]Role
	rolesByID   map[uint]Role
	rolesMutex  sync.RWMutex
	initOnce    sync.Once
}

// NewRoleCacheService creates a new RoleCacheService.
func NewRoleCacheService(reader RoleReader, logger *slog.Logger, refreshInterval time.Duration) *RoleCacheService {
	return &RoleCacheService{
		roleReader:      reader,
		logger:          logger,
		refreshInterval: refreshInterval,
		rolesByName:     make(map[RoleName]Role),
		rolesByID:       make(map[uint]Role),
	}
}
func (rcm *RoleCacheService) InitAndStartRefresh(ctx context.Context) error {
	var initErr error
	rcm.initOnce.Do(func() {
		initErr = rcm.loadRolesFromDB(ctx)
		if initErr != nil {
			return
		}
		rcm.logger.Info("Initial roles loaded successfully.")

		go func() {
			ticker := time.NewTicker(rcm.refreshInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					rcm.logger.Info("Refreshing roles cache...")
					if err := rcm.loadRolesFromDB(ctx); err != nil {
						rcm.logger.Error("Failed to refresh roles cache", "error", err)
					} else {
						rcm.logger.Info("Roles cache refreshed successfully.")
					}
				case <-ctx.Done():
					rcm.logger.Info("Stopping roles cache refresh goroutine.")
					return
				}
			}
		}()
	})
	return initErr
}

func (rcm *RoleCacheService) loadRolesFromDB(ctx context.Context) error {
	allRoles, err := rcm.roleReader.GetAllRoles(ctx) // ใช้ roleReader และส่ง ctx
	if err != nil {
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

	return nil
}

// GetRoleByName retrieves a role from the cache.
func (rcm *RoleCacheService) GetRoleByName(name RoleName) (Role, error) {
	rcm.rolesMutex.RLock()
	defer rcm.rolesMutex.RUnlock()
	if rcm.rolesByName == nil {
		return Role{}, fmt.Errorf("role cache not initialized")
	}
	role, ok := rcm.rolesByName[name]
	if !ok {
		return Role{}, ErrRoleNotFound
	}
	return role, nil
}

// GetRoleByID retrieves a role from the cache.
func (rcm *RoleCacheService) GetRoleByID(id uint) (Role, error) {
	rcm.rolesMutex.RLock()
	defer rcm.rolesMutex.RUnlock()

	if rcm.rolesByID == nil {
		return Role{}, fmt.Errorf("role cache not initialized")
	}
	role, ok := rcm.rolesByID[id]
	if !ok {
		return Role{}, ErrRoleNotFound
	}
	return role, nil
}
