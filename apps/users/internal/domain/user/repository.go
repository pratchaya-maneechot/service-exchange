package user

import (
	"context"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/role"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"
)

// UserRepository is the interface that provides access to User aggregates.
// It defines the contract for persisting and retrieving User domain objects.
type UserRepository interface {
	// FindByID retrieves a User aggregate by its ID.
	FindByID(ctx context.Context, id ids.UserID) (*User, error)

	// FindByLineUserID retrieves a User aggregate by their LINE User ID.
	FindByLineUserID(ctx context.Context, lineUserID string) (*User, error)

	// Save persists a User aggregate (either creating or updating).
	Save(ctx context.Context, user *User) error

	// ExistsByLineUserID checks if a user with the given LINE User ID already exists.
	ExistsByLineUserID(ctx context.Context, lineUserID string) (bool, error)

	// CreateUserRole adds a role to an existing user.
	// This might be a separate method if roles are managed outside the main User aggregate persistence.
	CreateUserRole(ctx context.Context, userID ids.UserID, roleID uint) error

	// GetRoleByID retrieves a Role by its ID.
	GetRoleByID(ctx context.Context, roleID uint) (*role.Role, error)
}
