package role

import (
	"context"
)

// RoleReader is the interface that provides read access to Role entity.
// It defines the contract for retrieving Role domain objects.
type RoleReader interface {
	// GetAllRoles retrieves all Roles entity.
	GetAllRoles(ctx context.Context) ([]Role, error)
}
