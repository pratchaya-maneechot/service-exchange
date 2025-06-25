package handler

import (
	"context"
	"fmt"

	"log/slog"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/command"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/role"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/observability"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type RegisterUserCommandHandler struct {
	userRepo     user.UserRepository
	roleCacheSvc *role.RoleCacheService
	logger       *slog.Logger
	config       *config.Config
}

func NewRegisterUserCommandHandler(
	userRepo user.UserRepository,
	rcs *role.RoleCacheService,
	logger *slog.Logger,
	cfg *config.Config,
) *RegisterUserCommandHandler {

	handlerLogger := logger.With(slog.String("component", "RegisterUserCommandHandler"))
	return &RegisterUserCommandHandler{
		userRepo:     userRepo,
		roleCacheSvc: rcs,
		logger:       handlerLogger,
		config:       cfg,
	}
}

func (h *RegisterUserCommandHandler) Handle(ctx context.Context, cmd command.RegisterUserCommand) (*command.RegisterUserDto, error) {
	logger := observability.GetLoggerFromContext(ctx).With(
		slog.String("method", "HandleRegisterUserCommand"),
		slog.String("line_user_id", cmd.LineUserID),
		slog.String("display_name", cmd.DisplayName),
	)
	tracer := otel.Tracer(fmt.Sprintf("%s.command-handler", h.config.Name))
	ctx, span := tracer.Start(ctx, "RegisterUserCommandHandler.Handle", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()
	span.SetAttributes(
		attribute.String("user.line_user_id", cmd.LineUserID),
		attribute.String("user.display_name", cmd.DisplayName),
	)
	logger.Debug("Attempting to register new user.")
	existing, err := h.userRepo.ExistsByLineUserID(ctx, cmd.LineUserID)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to check existing user by Line User ID")
		span.RecordError(err)
		logger.Error("Failed to check existing user by Line User ID", slog.Any("error", err))
		return nil, err
	}
	if existing {
		span.SetStatus(codes.Error, "Line User ID already exists")
		logger.Warn("User registration failed: Line User ID already exists.", "line_user_id", cmd.LineUserID)
		return nil, user.ErrLineUserAlreadyExists
	}
	logger.Debug("Line User ID is unique.")
	domUser, err := user.NewUser(ids.NewUserID(), cmd.LineUserID, cmd.Email, cmd.Password)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to create new domain user")
		span.RecordError(err)
		logger.Error("Failed to create new domain user", slog.Any("error", err))
		return nil, err
	}
	span.SetAttributes(attribute.String("user.id", string(domUser.ID)))
	logger.Debug("New domain user created", "new_user_id", domUser.ID)
	defaultRole, err := h.roleCacheSvc.GetRoleByName(role.RoleNamePoster)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to get default role")
		span.RecordError(err)
		logger.Error("Failed to get default role", "role_name", role.RoleNamePoster, slog.Any("error", err))
		return nil, err
	}
	logger.Debug("Default role retrieved", "default_role_id", defaultRole.ID, "default_role_name", defaultRole.Name)
	if err = domUser.AddRole(defaultRole); err != nil {
		span.SetStatus(codes.Error, "Failed to add default role to user")
		span.RecordError(err)
		logger.Error("Failed to add default role to user", "default_role_id", defaultRole.ID, slog.Any("error", err))
		return nil, err
	}
	logger.Debug("Default role added to user", "new_user_id", domUser.ID)
	if err = h.userRepo.Save(ctx, domUser); err != nil {
		span.SetStatus(codes.Error, "Failed to save user to repository")
		span.RecordError(err)
		logger.Error("Failed to save user to repository", "new_user_id", domUser.ID, slog.Any("error", err))
		return nil, err
	}
	span.SetStatus(codes.Ok, "User registered successfully")
	logger.Info("User registered successfully", "user_id", domUser.ID)
	return &command.RegisterUserDto{
		UserID: string(domUser.ID),
	}, nil
}
