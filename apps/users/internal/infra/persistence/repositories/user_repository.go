package repositories

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/role"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
	db "github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/persistence/postgres/generated"
	"github.com/pratchaya-maneechot/service-exchange/libs/infra/observability"
	lp "github.com/pratchaya-maneechot/service-exchange/libs/infra/postgres"
	"github.com/pratchaya-maneechot/service-exchange/libs/utils"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db     *db.Queries
	pool   *pgxpool.Pool
	logger *slog.Logger
	config *config.Config
	tracer trace.Tracer
}

func NewPostgresUserRepository(cfg *config.Config, dbPool *lp.DBPool, logger *slog.Logger) user.UserRepository {
	repoLogger := logger.With(slog.String("component", "userRepository"))
	return &userRepository{
		db:     db.New(dbPool.Pool),
		pool:   dbPool.Pool,
		logger: repoLogger,
		config: cfg,
		tracer: otel.Tracer(fmt.Sprintf("%s.repository", cfg.Name)),
	}
}

func (r *userRepository) FindByID(ctx context.Context, id ids.UserID) (*user.User, error) {
	logger := observability.LoggerFromCtx(ctx).With(slog.String("user_id", string(id)))

	ctx, span := r.tracer.Start(ctx, "UserRepository.FindByID", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.operation", "read_by_id"),
		attribute.String("db.user_id", string(id)),
	)

	raw, err := r.db.FindUserByID(ctx, lp.ToUUID(string(id)))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			span.SetStatus(codes.Ok, "User not found in DB")
			span.SetAttributes(attribute.Bool("user.found", false))
			logger.Debug("User not found in DB.", "user_id", string(id))
			return nil, user.ErrUserNotFound
		}

		span.SetStatus(codes.Error, "Failed to query user by ID from DB")
		span.RecordError(err)
		span.SetAttributes(attribute.String("error.type", "db_read_error"))
		logger.Error("Failed to query user by ID from DB", "user_id", string(id), slog.Any("error", err))
		return nil, fmt.Errorf("failed to query user by ID: %w", err)
	}

	uRoles, err := r.db.GetUserRoles(ctx, raw.ID)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to query user roles from DB")
		span.RecordError(err)
		span.SetAttributes(attribute.String("error.type", "db_read_error"))
		logger.Error("Failed to query user roles from DB", "user_id", string(id), slog.Any("error", err))
		return nil, fmt.Errorf("failed to query user roles: %w", err)
	}
	var roles = utils.ArrayMap(uRoles, func(ur db.Role) role.Role { return *role.NewRoleFromRepository(uint(ur.ID), ur.Name, ur.Description) })

	span.SetStatus(codes.Ok, "User loaded from DB")
	span.SetAttributes(attribute.Bool("user.found", true))
	logger.Debug("User successfully loaded from DB.", "user_id", string(id))

	preferencesJSON, err := utils.ByteToMap(raw.Preferences)
	if err != nil {
		logger.Error("Failed to unmarshal user preferences", "user_id", raw.ID, slog.Any("error", err))
		return nil, fmt.Errorf("failed to marshal preferences to JSON: %w", err)
	}

	resp, err := user.NewUserFromRepository(
		raw.ID.String(),
		raw.LineUserID,
		raw.Email,
		raw.PasswordHash,
		raw.Status,
		*lp.ToTime(raw.CreatedAt),
		*lp.ToTime(raw.UpdatedAt),
		lp.ToTime(raw.LastLoginAt),
		user.NewProfileFromRepository(raw.ID.String(), raw.DisplayName, *preferencesJSON),
		roles,
	)
	if err != nil {
		logger.Error("Failed to transform json to user model", "user_id", raw.ID, slog.Any("error", err))
		return nil, fmt.Errorf("failed to transform json to user model: %w", err)
	}

	return resp, nil
}

func (r *userRepository) FindByLineUserID(ctx context.Context, lineUserID string) (*user.User, error) {
	logger := observability.LoggerFromCtx(ctx).With(slog.String("method", "FindByLineUserID"))
	ctx, span := r.tracer.Start(ctx, "UserRepository.FindByLineUserID", trace.WithSpanKind(trace.SpanKindClient))
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

	uRoles, err := r.db.GetUserRoles(ctx, raw.ID)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to query user roles from DB")
		span.RecordError(err)
		span.SetAttributes(attribute.String("error.type", "db_read_error"))
		logger.Error("Failed to query user roles from DB", "user_id", raw.ID, slog.Any("error", err))
		return nil, fmt.Errorf("failed to query user roles: %w", err)
	}
	var roles = utils.ArrayMap(uRoles, func(ur db.Role) role.Role { return *role.NewRoleFromRepository(uint(ur.ID), ur.Name, ur.Description) })

	span.SetStatus(codes.Ok, "User found by Line User ID successfully")
	logger.Debug("Successfully found user by Line User ID", "line_user_id", lineUserID)

	preferencesJSON, err := utils.ByteToMap(raw.Preferences)
	if err != nil {
		logger.Error("Failed to unmarshal user preferences for Line User ID", "line_user_id", lineUserID, slog.Any("error", err))
		return nil, fmt.Errorf("failed to marshal preferences to JSON: %w", err)
	}
	resp, err := user.NewUserFromRepository(
		raw.ID.String(),
		raw.LineUserID,
		raw.Email,
		raw.PasswordHash,
		raw.Status,
		*lp.ToTime(raw.CreatedAt),
		*lp.ToTime(raw.UpdatedAt),
		lp.ToTime(raw.LastLoginAt),
		user.NewProfileFromRepository(raw.ID.String(), raw.DisplayName, *preferencesJSON),
		roles,
	)
	if err != nil {
		logger.Error("Failed to transform json to user model", "user_id", raw.ID, slog.Any("error", err))
		return nil, fmt.Errorf("failed to transform json to user model: %w", err)
	}

	return resp, nil
}

func (r *userRepository) Save(ctx context.Context, u *user.User) (err error) {
	logger := observability.LoggerFromCtx(ctx).With(slog.String("user_id", string(u.ID)))

	ctx, span := r.tracer.Start(ctx, "UserRepository.Save", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.user_id", string(u.ID)),
	)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to begin transaction")
		span.RecordError(err)
		span.SetAttributes(attribute.String("error.type", "db_transaction_error"))
		logger.Error("Failed to begin DB transaction", slog.Any("error", err))
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			panic(r)
		} else if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				logger.Error("Failed to rollback DB transaction", slog.Any("error", rollbackErr))
			}
		} else {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				span.SetStatus(codes.Error, "Failed to commit transaction")
				span.RecordError(commitErr)
				span.SetAttributes(attribute.String("error.type", "db_transaction_error"))
				logger.Error("Failed to commit DB transaction", slog.Any("error", commitErr))
				err = fmt.Errorf("failed to commit transaction: %w", commitErr)
			}
		}
	}()

	qtx := r.db.WithTx(tx)
	userID := lp.ToUUID(string(u.ID))

	userExists, err := qtx.UserExistsByID(ctx, userID)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to check user existence in DB")
		span.RecordError(err)
		span.SetAttributes(attribute.String("error.type", "db_read_error"))
		logger.Error("Failed to check user existence in DB", slog.Any("error", err))
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	if userExists {
		logger.Debug("Updating existing user in DB.")
		input := db.UpdateUserParams{
			ID:           userID,
			LineUserID:   u.LineUserID,
			Email:        u.Email,
			PasswordHash: u.PasswordHash,
			Status:       string(u.Status),
			LastLoginAt:  lp.ToTimestamp(u.LastLoginAt),
		}
		if _, err = qtx.UpdateUser(ctx, input); err != nil {
			span.SetStatus(codes.Error, "Failed to update user in DB")
			span.RecordError(err)
			span.SetAttributes(attribute.String("error.type", "db_write_error"))
			logger.Error("Failed to update user in DB", slog.Any("error", err))
			return fmt.Errorf("failed to update user: %w", err)
		}
	} else {
		logger.Debug("Creating new user in DB.")
		input := db.CreateUserParams{
			ID:           userID,
			LineUserID:   u.LineUserID,
			Email:        u.Email,
			PasswordHash: u.PasswordHash,
			Status:       string(u.Status),
		}
		if _, err = qtx.CreateUser(ctx, input); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == "23505" {
					logger.Warn("Duplicate entry during user creation in DB.", slog.Any("error", err))
					span.SetStatus(codes.Error, "Duplicate user entry in DB")
					span.SetAttributes(attribute.String("error.type", "db_unique_violation"))
					return user.ErrLineUserAlreadyExists
				}
			}
			span.SetStatus(codes.Error, "Failed to create user in DB")
			span.RecordError(err)
			span.SetAttributes(attribute.String("error.type", "db_write_error"))
			logger.Error("Failed to insert user in DB", slog.Any("error", err))
			return fmt.Errorf("failed to insert user: %w", err)
		}
	}

	preferencesByte, marshalErr := utils.MapToByte(u.Profile.Preferences)
	if marshalErr != nil {
		span.SetStatus(codes.Error, "Failed to marshal preferences")
		span.RecordError(marshalErr)
		span.SetAttributes(attribute.String("error.type", "data_serialization_error"))
		logger.Error("Failed to marshal user preferences for DB", slog.Any("error", marshalErr))
		return marshalErr
	}
	logger.Debug("Saving user profile in DB.")
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
		span.SetStatus(codes.Error, "Failed to upsert user profile in DB")
		span.RecordError(err)
		span.SetAttributes(attribute.String("error.type", "db_write_error"))
		logger.Error("Failed to upsert user profile in DB", slog.Any("error", err))
		return fmt.Errorf("failed to upsert user profile: %w", err)
	}

	logger.Debug("Updating user roles in DB.")
	if err := qtx.DeleteUserRoles(ctx, userID); err != nil {
		span.SetStatus(codes.Error, "Failed to delete old roles in DB")
		span.RecordError(err)
		span.SetAttributes(attribute.String("error.type", "db_write_error"))
		logger.Error("Failed to delete old user roles in DB", slog.Any("error", err))
		return fmt.Errorf("failed to delete old user roles: %w", err)
	}
	for _, r := range u.Roles {
		roleInput := db.CreateUserRoleParams{
			UserID: userID,
			RoleID: int32(r.ID),
		}
		if _, err := qtx.CreateUserRole(ctx, roleInput); err != nil {
			span.SetStatus(codes.Error, "Failed to add role in DB")
			span.RecordError(err)
			span.SetAttributes(attribute.String("error.type", "db_write_error"))
			logger.Error("Failed to add user role in DB", slog.Any("error", err), "role_id", r.ID)
			return fmt.Errorf("failed to add role %d to user %s: %w", r.ID, u.ID, err)
		}
	}

	span.SetStatus(codes.Ok, "User saved to DB")
	logger.Info("User saved successfully to DB.")
	return nil
}

func (r *userRepository) ExistsByLineUserID(ctx context.Context, lineUserID string) (bool, error) {
	logger := observability.LoggerFromCtx(ctx).With(slog.String("line_user_id", lineUserID))

	ctx, span := r.tracer.Start(ctx, "UserRepository.ExistsByLineUserID", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.operation", "check_existence"),
		attribute.String("db.line_user_id", lineUserID),
	)

	exists, err := r.db.UserExistsByLineUserID(ctx, lineUserID)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to query DB for existence")
		span.RecordError(err)
		span.SetAttributes(attribute.String("error.type", "db_read_error"))
		logger.Error("Failed to check Line User ID existence in DB", slog.Any("error", err))
		return false, fmt.Errorf("failed to check line user ID existence: %w", err)
	}
	span.SetStatus(codes.Ok, "Existence checked in DB")
	span.SetAttributes(attribute.Bool("user.exists", exists))
	logger.Debug("Line User ID existence check complete in DB.", "exists", exists)
	return exists, nil
}

func (r *userRepository) CreateUserRole(ctx context.Context, usrId ids.UserID, roleID uint) error {
	logger := observability.LoggerFromCtx(ctx).With(
		slog.String("method", "CreateUserRole"),
		slog.String("user_id", string(usrId)),
		slog.Uint64("role_id", uint64(roleID)),
	)
	ctx, span := r.tracer.Start(ctx, "UserRepository.CreateUserRole", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	span.SetAttributes(attribute.String("db.user_id", string(usrId)), attribute.Int("db.role_id", int(roleID)))
	logger.Debug("Attempting to create user role")

	userID := lp.ToUUID(string(usrId))

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

func (r *userRepository) GetRoleByID(ctx context.Context, roleID uint) (*role.Role, error) {
	logger := observability.LoggerFromCtx(ctx).With(
		slog.String("method", "GetRoleByID"),
		slog.Uint64("role_id", uint64(roleID)),
	)
	ctx, span := r.tracer.Start(ctx, "UserRepository.GetRoleByID", trace.WithSpanKind(trace.SpanKindClient))
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

	return role.NewRoleFromRepository(roleID, raw.Name, raw.Description), nil
}
