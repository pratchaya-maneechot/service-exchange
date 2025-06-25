package repositories

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/role"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user" // Domain models
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/observability"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/postgres"
	db "github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/postgres/generated"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/pkg/utils"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db     *db.Queries
	pool   *pgxpool.Pool
	logger *slog.Logger
	config *config.Config
}

func NewPostgresUserRepository(cfg *config.Config, dbPool *postgres.DBPool, logger *slog.Logger) user.UserRepository {
	// Enrich the logger with a common attribute for this repository
	repoLogger := logger.With(slog.String("component", "userRepository"))
	return &userRepository{
		db:     db.New(dbPool.Pool),
		pool:   dbPool.Pool,
		logger: repoLogger, // Use the enriched logger
		config: cfg,
	}
}

// FindByID retrieves a user by their ID.
func (r *userRepository) FindByID(ctx context.Context, id ids.UserID) (*user.User, error) {
	// Always get the logger from context first to ensure trace correlation
	logger := observability.GetLoggerFromContext(ctx).With(slog.String("method", "FindByID"))
	tracer := otel.Tracer(fmt.Sprintf("%s.repository", r.config.Name))
	ctx, span := tracer.Start(ctx, "UserRepository.FindByID", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	span.SetAttributes(attribute.String("db.user_id", string(id)))
	logger.Debug("Attempting to find user by ID", "user_id", id)

	raw, err := r.db.FindUserByID(ctx, encodeUID(string(id)))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			span.SetStatus(codes.Error, "User not found in database")
			logger.Warn("User not found", "user_id", id, slog.Any("error", err)) // Log actual error for debugging
			return nil, user.ErrUserNotFound
		}
		span.SetStatus(codes.Error, "Failed to query user by ID")
		span.RecordError(err)
		logger.Error("Failed to query user by ID", "user_id", id, slog.Any("error", err))
		return nil, fmt.Errorf("failed to query user full aggregate by ID: %w", err)
	}

	span.SetStatus(codes.Ok, "User found successfully")
	logger.Debug("Successfully found user by ID", "user_id", id)

	preferencesJSON, err := utils.ByteToMap(raw.Preferences)
	if err != nil {
		// This is a data integrity/marshalling error, usually critical
		logger.Error("Failed to unmarshal user preferences", "user_id", id, slog.Any("error", err))
		return nil, fmt.Errorf("failed to unmarshal preferences to JSON: %w", err)
	}

	uRoles, err := r.db.GetUserRoles(ctx, raw.ID)
	if err != nil {
		logger.Error("Failed to query user roles", "user_id", id, slog.Any("error", err))
		return nil, fmt.Errorf("failed to query get user role: %w", err)
	}

	var roles []role.Role
	for _, ur := range uRoles {
		idVal := uint(ur.ID)
		roles = append(roles, role.Role{
			ID:          &idVal,
			Name:        role.RoleName(ur.Name),
			Description: *ur.Description,
		})
	}

	var lastLoginAt *time.Time
	if raw.LastLoginAt.Valid {
		lastLoginAt = &raw.LastLoginAt.Time
	}

	createdAt := raw.CreatedAt.Time
	updatedAt := raw.UpdatedAt.Time

	return &user.User{
		ID:           ids.UserID(raw.ID.String()),
		LineUserID:   raw.LineUserID,
		Email:        *raw.Email,
		PasswordHash: *raw.PasswordHash,
		Status:       user.UserStatus(raw.Status),
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		LastLoginAt:  lastLoginAt,
		Profile: user.Profile{
			UserID:      ids.UserID(raw.ID.String()),
			DisplayName: *raw.DisplayName,
			FirstName:   raw.FirstName,
			LastName:    raw.LastName,
			Bio:         raw.Bio,
			AvatarURL:   raw.AvatarUrl,
			PhoneNumber: raw.PhoneNumber,
			Address:     raw.Address,
			Preferences: *preferencesJSON,
		},
		Roles: roles,
	}, nil
}

// FindByLineUserID retrieves a User aggregate by their LINE User ID, using sqlc generated query with JOINs.
func (r *userRepository) FindByLineUserID(ctx context.Context, lineUserID string) (*user.User, error) {
	logger := observability.GetLoggerFromContext(ctx).With(slog.String("method", "FindByLineUserID"))
	tracer := otel.Tracer(fmt.Sprintf("%s.repository", r.config.Name))
	ctx, span := tracer.Start(ctx, "UserRepository.FindByLineUserID", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	span.SetAttributes(attribute.String("db.line_user_id", lineUserID))
	logger.Debug("Attempting to find user by Line User ID", "line_user_id", lineUserID)

	raw, err := r.db.FindUserByLineUserID(ctx, lineUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			span.SetStatus(codes.Error, "User not found by Line User ID")
			logger.Warn("User not found by Line User ID", "line_user_id", lineUserID, slog.Any("error", err))
			return nil, user.ErrUserNotFound
		}
		span.SetStatus(codes.Error, "Failed to query user by Line User ID")
		span.RecordError(err)
		logger.Error("Failed to query user by Line User ID", "line_user_id", lineUserID, slog.Any("error", err))
		return nil, fmt.Errorf("failed to query user full aggregate by LINE User ID: %w", err)
	}

	span.SetStatus(codes.Ok, "User found by Line User ID successfully")
	logger.Debug("Successfully found user by Line User ID", "line_user_id", lineUserID)

	preferencesJSON, err := utils.ByteToMap(raw.Preferences)
	if err != nil {
		logger.Error("Failed to unmarshal user preferences for Line User ID", "line_user_id", lineUserID, slog.Any("error", err))
		return nil, fmt.Errorf("failed to marshal preferences to JSON: %w", err)
	}
	userID := ids.UserID(raw.ID.String())
	return &user.User{
		ID:           userID,
		LineUserID:   lineUserID,
		Email:        *raw.Email,
		PasswordHash: "", // Password hash should generally not be returned in read operations
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
	}, nil
}

// Save persists a User aggregate. It handles creation/update for user, profile, and related entities.
// This method orchestrates multiple sqlc-generated inserts/updates within a transaction.
func (r *userRepository) Save(ctx context.Context, u *user.User) error {
	logger := observability.GetLoggerFromContext(ctx).With(slog.String("method", "Save"), slog.String("user_id", string(u.ID)))
	tracer := otel.Tracer(fmt.Sprintf("%s.repository", r.config.Name))
	ctx, span := tracer.Start(ctx, "UserRepository.Save", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	span.SetAttributes(attribute.String("db.user_id", string(u.ID)))
	logger.Debug("Starting transaction to save user")

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to begin transaction")
		span.RecordError(err)
		logger.Error("Failed to begin transaction", slog.Any("error", err))
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			panic(r) // Re-throw panic after rollback
		} else if err != nil { // Only rollback if an error occurred
			tx.Rollback(ctx)
		}
	}()

	qtx := r.db.WithTx(tx)
	userID := encodeUID(string(u.ID))

	userExists, err := qtx.UserExistsByID(ctx, userID)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to check user existence")
		span.RecordError(err)
		logger.Error("Failed to check user existence", slog.Any("error", err))
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	if userExists {
		logger.Debug("Updating existing user")
		input := db.UpdateUserParams{
			ID:           userID,
			LineUserID:   u.LineUserID,
			Email:        &u.Email,
			PasswordHash: &u.PasswordHash,
			Status:       string(u.Status),
			LastLoginAt:  pgtype.Timestamptz{Time: *u.LastLoginAt},
		}
		if _, err = qtx.UpdateUser(ctx, input); err != nil {
			span.SetStatus(codes.Error, "Failed to update user")
			span.RecordError(err)
			logger.Error("Failed to update user", slog.Any("error", err))
			return fmt.Errorf("failed to update user: %w", err)
		}
	} else {
		logger.Debug("Creating new user")
		input := db.CreateUserParams{
			ID:           userID,
			LineUserID:   u.LineUserID,
			Email:        &u.Email,
			PasswordHash: &u.PasswordHash,
			Status:       string(u.Status),
		}
		if _, err = qtx.CreateUser(ctx, input); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == "23505" { // Unique violation
					logger.Warn("Duplicate user ID or Line User ID", slog.Any("error", err))
					span.SetStatus(codes.Error, "Duplicate user entry")
					return user.ErrLineUserAlreadyExists
				}
			}
			span.SetStatus(codes.Error, "Failed to insert user")
			span.RecordError(err)
			logger.Error("Failed to insert user", slog.Any("error", err))
			return fmt.Errorf("failed to insert user: %w", err)
		}
	}

	preferencesByte, err := utils.MapToByte(u.Profile.Preferences)
	if err != nil {
		logger.Error("Failed to marshal user preferences for saving", slog.Any("error", err))
		return err
	}
	logger.Debug("Saving user profile")
	profileInput := db.UpsertUserProfileParams{
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
	if _, err = qtx.UpsertUserProfile(ctx, profileInput); err != nil {
		span.SetStatus(codes.Error, "Failed to upsert user profile")
		span.RecordError(err)
		logger.Error("Failed to upsert user profile", slog.Any("error", err))
		return fmt.Errorf("failed to upsert user profile: %w", err)
	}

	logger.Debug("Updating user roles", "roles_count", len(u.Roles))
	if err := qtx.DeleteUserRoles(ctx, userID); err != nil {
		span.SetStatus(codes.Error, "Failed to delete old user roles")
		span.RecordError(err)
		logger.Error("Failed to delete old user roles", slog.Any("error", err))
		return fmt.Errorf("failed to delete old user roles: %w", err)
	}
	for _, r := range u.Roles {
		roleInput := db.CreateUserRoleParams{
			UserID: userID,
			RoleID: int32(*r.ID),
		}
		if _, err := qtx.CreateUserRole(ctx, roleInput); err != nil {
			span.SetStatus(codes.Error, "Failed to add user role")
			span.RecordError(err)
			logger.Error("Failed to add user role", slog.Any("error", err), "role_id", r.ID)
			return fmt.Errorf("failed to add role %d to user %s: %w", *r.ID, u.ID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		span.SetStatus(codes.Error, "Failed to commit transaction")
		span.RecordError(err)
		logger.Error("Failed to commit transaction", slog.Any("error", err))
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	span.SetStatus(codes.Ok, "User saved successfully")
	logger.Info("User saved successfully") // user_id is already in the logger's context
	return nil
}

// ExistsByLineUserID checks if a user with the given LINE User ID already exists.
func (r *userRepository) ExistsByLineUserID(ctx context.Context, lineUserID string) (bool, error) {
	logger := observability.GetLoggerFromContext(ctx).With(slog.String("method", "ExistsByLineUserID"))
	tracer := otel.Tracer(fmt.Sprintf("%s.repository", r.config.Name))
	ctx, span := tracer.Start(ctx, "UserRepository.ExistsByLineUserID", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	span.SetAttributes(attribute.String("db.line_user_id", lineUserID))
	logger.Debug("Checking if Line User ID exists in database", "line_user_id", lineUserID)

	exists, err := r.db.UserExistsByLineUserID(ctx, lineUserID)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to check Line User ID existence")
		span.RecordError(err)
		logger.Error("Failed to check Line User ID existence", slog.Any("error", err), "line_user_id", lineUserID)
		return false, fmt.Errorf("failed to check line user ID existence: %w", err)
	}
	span.SetStatus(codes.Ok, "Line User ID existence checked")
	span.SetAttributes(attribute.Bool("user.exists", exists))
	logger.Debug("Line User ID existence check complete", "line_user_id", lineUserID, "exists", exists)
	return exists, nil
}

// CreateUserRole adds a role to an existing user.
func (r *userRepository) CreateUserRole(ctx context.Context, usrId ids.UserID, roleID uint) error {
	logger := observability.GetLoggerFromContext(ctx).With(
		slog.String("method", "CreateUserRole"),
		slog.String("user_id", string(usrId)),
		slog.Uint64("role_id", uint64(roleID)), // Use Uint64 for uint
	)
	tracer := otel.Tracer(fmt.Sprintf("%s.repository", r.config.Name))
	ctx, span := tracer.Start(ctx, "UserRepository.CreateUserRole", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	span.SetAttributes(attribute.String("db.user_id", string(usrId)), attribute.Int("db.role_id", int(roleID)))
	logger.Debug("Attempting to create user role")

	userID := encodeUID(string(usrId))

	// Check if role exists
	roleExists, err := r.db.RoleExistsByID(ctx, int32(roleID))
	if err != nil {
		span.SetStatus(codes.Error, "Failed to check role existence")
		span.RecordError(err)
		logger.Error("Failed to check role existence", slog.Any("error", err))
		return fmt.Errorf("failed to check role existence: %w", err)
	}
	if !roleExists {
		span.SetStatus(codes.Error, "Role not found")
		logger.Warn("Role not found")
		return role.ErrRoleNotFound
	}

	alreadyHasRole, err := r.db.UserRoleExists(ctx, db.UserRoleExistsParams{
		UserID: userID,
		RoleID: int32(roleID),
	})
	if err != nil {
		span.SetStatus(codes.Error, "Failed to check if user already has role")
		span.RecordError(err)
		logger.Error("Failed to check if user already has role", slog.Any("error", err))
		return fmt.Errorf("failed to check if user already has role: %w", err)
	}
	if alreadyHasRole {
		span.SetStatus(codes.Error, "User already has role")
		logger.Warn("User already has role")
		return user.ErrRoleAlreadyAssigned
	}

	_, err = r.db.CreateUserRole(ctx, db.CreateUserRoleParams{
		UserID: userID,
		RoleID: int32(roleID),
	})
	if err != nil {
		span.SetStatus(codes.Error, "Failed to add role to user")
		span.RecordError(err)
		logger.Error("Failed to add role to user", slog.Any("error", err))
		return fmt.Errorf("failed to add role to user: %w", err)
	}

	span.SetStatus(codes.Ok, "User role created successfully")
	logger.Info("Successfully created user role")
	return nil
}

// GetRoleByID retrieves a Role by its ID.
func (r *userRepository) GetRoleByID(ctx context.Context, roleID uint) (*role.Role, error) {
	logger := observability.GetLoggerFromContext(ctx).With(
		slog.String("method", "GetRoleByID"),
		slog.Uint64("role_id", uint64(roleID)),
	)
	tracer := otel.Tracer(fmt.Sprintf("%s.repository", r.config.Name))
	ctx, span := tracer.Start(ctx, "UserRepository.GetRoleByID", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	span.SetAttributes(attribute.Int("db.role_id", int(roleID)))
	logger.Debug("Attempting to get role by ID")

	raw, err := r.db.FindRoleByID(ctx, int32(roleID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			span.SetStatus(codes.Error, "Role not found")
			logger.Warn("Role not found")
			return nil, role.ErrRoleNotFound
		}
		span.SetStatus(codes.Error, "Failed to get role by ID")
		span.RecordError(err)
		logger.Error("Failed to get role by ID", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get role by ID: %w", err)
	}

	span.SetStatus(codes.Ok, "Role found successfully")
	logger.Debug("Successfully found role by ID")

	id := uint(raw.ID)
	return &role.Role{
		ID:          &id,
		Name:        role.RoleName(raw.Name),
		Description: *raw.Description,
	}, nil
}

func encodeUID(str string) pgtype.UUID {
	return pgtype.UUID{
		Bytes: [16]byte(uuid.MustParse(str)),
		Valid: true,
	}
}
