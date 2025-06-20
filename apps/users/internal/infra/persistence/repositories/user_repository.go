package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/role"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user" // Domain models
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/postgres"
	db "github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/postgres/generated"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresUserRepository implements the user.UserRepository interface.
type userRepository struct {
	db   *db.Queries
	pool *pgxpool.Pool // Keep pool for transaction management
}

// NewPostgresUserRepository creates a new PostgresUserRepository.
func NewPostgresUserRepository(dbPool *postgres.DBPool) user.UserRepository {
	return &userRepository{
		db:   db.New(dbPool.Pool),
		pool: dbPool.Pool,
	}
}

// FindByID retrieves a User aggregate by its ID, using sqlc generated query with JOINs.
func (r *userRepository) FindByID(ctx context.Context, id ids.UserID) (*user.User, error) {
	raw, err := r.db.FindUserByID(ctx, encodeUID(uuid.MustParse(string(id))))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to query user full aggregate by ID: %w", err)
	}
	userID := ids.UserID(raw.ID.String())
	preferencesJSON, err := utils.ByteToMap(raw.Preferences)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal preferences to JSON: %w", err)
	}
	return &user.User{
		ID:           ids.UserID(userID),
		LineUserID:   raw.LineUserID,
		Email:        *raw.Email,
		PasswordHash: *raw.PasswordHash,
		Status:       user.UserStatus(raw.Status),
		CreatedAt:    raw.CreatedAt.Time,
		UpdatedAt:    raw.UpdatedAt.Time,
		LastLoginAt:  &raw.LastLoginAt.Time,
		Profile: user.Profile{
			UserID:      userID,
			DisplayName: *raw.DisplayName,
			FirstName:   raw.FirstName,
			LastName:    raw.LastName,
			Bio:         raw.Bio,
			AvatarURL:   raw.AvatarUrl,
			PhoneNumber: raw.PhoneNumber,
			Address:     raw.Address,
			Preferences: *preferencesJSON,
		},
	}, err
}

// FindByLineUserID retrieves a User aggregate by their LINE User ID, using sqlc generated query with JOINs.
func (r *userRepository) FindByLineUserID(ctx context.Context, lineUserID string) (*user.User, error) {
	raw, err := r.db.FindUserByLineUserID(ctx, lineUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to query user full aggregate by LINE User ID: %w", err)
	}

	preferencesJSON, err := utils.ByteToMap(raw.Preferences)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal preferences to JSON: %w", err)
	}
	userID := ids.UserID(raw.ID.String())
	return &user.User{
		ID:           userID,
		LineUserID:   lineUserID,
		Email:        *raw.Email,
		PasswordHash: "",
		Status:       user.UserStatus(raw.Status),
		CreatedAt:    raw.CreatedAt.Time,
		UpdatedAt:    raw.UpdatedAt.Time,
		LastLoginAt:  &raw.LastLoginAt.Time,
		Profile: user.Profile{
			UserID:      userID,
			DisplayName: *raw.DisplayName,
			FirstName:   raw.FirstName,
			LastName:    raw.LastName,
			Bio:         raw.Bio,
			AvatarURL:   raw.AvatarUrl,
			PhoneNumber: raw.PhoneNumber,
			Address:     raw.Address,
			Preferences: *preferencesJSON,
		},
	}, err
}

// Save persists a User aggregate. It handles creation/update for user, profile, and related entities.
// This method orchestrates multiple sqlc-generated inserts/updates within a transaction.
func (r *userRepository) Save(ctx context.Context, u *user.User) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.db.WithTx(tx)
	userID := encodeUID(uuid.MustParse(string(u.ID)))
	userExists, err := qtx.UserExistsByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	if userExists {
		input := db.UpdateUserParams{
			ID:           userID,
			LineUserID:   u.LineUserID,
			Email:        &u.Email,
			PasswordHash: &u.PasswordHash,
			Status:       string(u.Status),
			LastLoginAt:  pgtype.Timestamptz{Time: *u.LastLoginAt},
		}
		if _, err = qtx.UpdateUser(ctx, input); err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	} else {
		input := db.CreateUserParams{
			ID:           userID,
			LineUserID:   u.LineUserID,
			Email:        &u.Email,
			PasswordHash: &u.PasswordHash,
			Status:       string(u.Status),
		}
		if _, err = qtx.CreateUser(ctx, input); err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
				if pgErr.ConstraintName == "users_line_user_id_key" {
					return user.ErrLineUserIDAlreadyExists
				}
				if pgErr.ConstraintName == "users_email_key" {
					return user.ErrEmailAlreadyExists
				}
			}
			return fmt.Errorf("failed to insert user: %w", err)
		}
	}

	profileExists, err := qtx.ProfileExistsByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check profile existence: %w", err)
	}

	preferencesByte, err := utils.MapToByte(u.Profile.Preferences)
	if err != nil {
		return err
	}

	if profileExists {
		input := db.UpdateProfileParams{
			UserID:      userID,
			DisplayName: u.Profile.DisplayName,
			FirstName:   u.Profile.FirstName,
			LastName:    u.Profile.LastName,
			Bio:         u.Profile.Bio,
			AvatarUrl:   u.Profile.AvatarURL,
			PhoneNumber: u.Profile.PhoneNumber,
			Address:     u.Profile.Address,
			Preferences: preferencesByte,
		}
		if _, err = qtx.UpdateProfile(ctx, input); err != nil {
			return fmt.Errorf("failed to update profile: %w", err)
		}
	} else {
		input := db.CreateProfileParams{
			UserID:      userID,
			DisplayName: u.Profile.DisplayName,
			FirstName:   u.Profile.FirstName,
			LastName:    u.Profile.LastName,
			Bio:         u.Profile.Bio,
			AvatarUrl:   u.Profile.AvatarURL,
			PhoneNumber: u.Profile.PhoneNumber,
			Address:     u.Profile.Address,
			Preferences: preferencesByte,
		}
		if _, err = qtx.CreateProfile(ctx, input); err != nil {
			return fmt.Errorf("failed to insert profile: %w", err)
		}
	}

	if len(u.Roles) > 0 {
		for _, role := range u.Roles {
			if err = qtx.DeleteUserRole(ctx, db.DeleteUserRoleParams{
				UserID: userID,
				RoleID: int32(*role.ID),
			}); err != nil {
				return fmt.Errorf("failed to delete user role: %w", err)
			}
			if _, err = qtx.CreateUserRole(ctx, db.CreateUserRoleParams{
				UserID: userID,
				RoleID: int32(*role.ID),
			}); err != nil {
				return fmt.Errorf("failed to create user role: %w", err)
			}
		}
	}

	return tx.Commit(ctx)
}

// ExistsByLineUserID checks if a user with the given LINE User ID already exists.
func (r *userRepository) ExistsByLineUserID(ctx context.Context, lineUserID string) (bool, error) {
	exists, err := r.db.UserExistsByLineUserID(ctx, lineUserID)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence by line user ID: %w", err)
	}
	return exists, nil
}

// AddRoleToUser adds a role to an existing user.
func (r *userRepository) CreateUserRole(ctx context.Context, usrId ids.UserID, roleID uint) error {
	userID := encodeUID(uuid.MustParse(string(usrId)))
	// Check if role exists
	roleExists, err := r.db.RoleExistsByID(ctx, int32(roleID))
	if err != nil {
		return fmt.Errorf("failed to check role existence: %w", err)
	}
	if !roleExists {
		return role.ErrRoleNotFound
	}

	alreadyHasRole, err := r.db.UserRoleExists(ctx, db.UserRoleExistsParams{
		UserID: userID,
		RoleID: int32(roleID),
	})
	if err != nil {
		return fmt.Errorf("failed to check if user already has role: %w", err)
	}
	if alreadyHasRole {
		return user.ErrRoleAlreadyAssigned
	}

	_, err = r.db.CreateUserRole(ctx, db.CreateUserRoleParams{
		UserID: userID,
		RoleID: int32(roleID),
	})
	if err != nil {
		return fmt.Errorf("failed to add role to user: %w", err)
	}
	return nil
}

// GetRoleByID retrieves a Role by its ID.
func (r *userRepository) GetRoleByID(ctx context.Context, roleID uint) (*role.Role, error) {
	raw, err := r.db.FindRoleByID(ctx, int32(roleID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, role.ErrRoleNotFound
		}
		return nil, fmt.Errorf("failed to get role by ID: %w", err)
	}
	id := uint(raw.ID)
	return &role.Role{
		ID:          &id,
		Name:        role.RoleName(raw.Name),
		Description: *raw.Description,
	}, nil
}

// Utility functions (from your example)
func encodeUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: [16]byte(id),
		Valid: true,
	}
}
