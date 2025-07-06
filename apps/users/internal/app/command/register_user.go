package command

import (
	"context"
	"errors"
	"fmt"

	"log/slog"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/role"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
	"github.com/pratchaya-maneechot/service-exchange/libs/infra/observability"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type RegisterUserCommand struct {
	LineUserID  string  `json:"lineUserId" validate:"required"`
	Email       *string `json:"email,omitempty" validate:"omitempty,email"`
	Password    *string `json:"password,omitempty"`
	DisplayName string  `json:"displayName" validate:"required"`
	AvatarURL   string  `json:"avatarUrl,omitempty"`
}

type RegisterUserDto struct {
	UserID string `json:"UserId" validate:"required"`
}

type RegisterUserCommandHandler struct {
	userRepo     user.UserRepository
	roleCacheSvc *role.RoleCacheService
	logger       *slog.Logger
	config       *config.Config
	tracer       trace.Tracer
}

func NewRegisterUserCommandHandler(
	userRepo user.UserRepository,
	rcs *role.RoleCacheService,
	logger *slog.Logger,
	cfg *config.Config,
) *RegisterUserCommandHandler {
	return &RegisterUserCommandHandler{
		userRepo:     userRepo,
		roleCacheSvc: rcs,
		logger:       logger.With(slog.String("component", "RegisterUserCommandHandler")),
		config:       cfg,
		tracer:       otel.Tracer(fmt.Sprintf("%s.command-handler", cfg.Name)),
	}
}

func (h *RegisterUserCommandHandler) Handle(ctx context.Context, cmd RegisterUserCommand) (*RegisterUserDto, error) {
	logger := observability.LoggerFromCtx(ctx).With(slog.String("line_user_id", cmd.LineUserID))

	ctx, span := h.tracer.Start(ctx, "RegisterUserCommand.Handle", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	span.SetAttributes(attribute.String("user.line_user_id", cmd.LineUserID))

	existing, err := h.userRepo.ExistsByLineUserID(ctx, cmd.LineUserID)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to check existing user")
		span.RecordError(err)
		span.SetAttributes(attribute.String("error.type", "repository_read_error"))
		logger.Error("Failed to check existing user", slog.Any("error", err))
		return nil, err
	}
	if existing {
		span.SetStatus(codes.Error, "User already exists")
		span.SetAttributes(attribute.String("error.type", string(user.ErrLineUserAlreadyExists.Code)))
		logger.Warn(user.ErrLineUserAlreadyExists.Message, "line_user_id", cmd.LineUserID)
		return nil, user.ErrLineUserAlreadyExists
	}

	domUser, err := user.NewUser(ids.NewUserID(), cmd.LineUserID, cmd.Email, cmd.Password)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to create user domain model")
		span.RecordError(err)
		span.SetAttributes(attribute.String("error.type", "domain_creation_error"))
		logger.Error("Failed to create domain user", slog.Any("error", err))
		return nil, err
	}
	span.SetAttributes(attribute.String("user.id", string(domUser.ID)))

	defaultRole, err := h.roleCacheSvc.GetRoleByName(role.RoleNamePoster)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to get default role")
		span.RecordError(err)
		span.SetAttributes(attribute.String("error.type", "role_service_error"))
		logger.Error("Failed to get default role", "role_name", role.RoleNamePoster, slog.Any("error", err))
		return nil, err
	}

	if err = domUser.AddRole(defaultRole); err != nil {
		span.SetStatus(codes.Error, "Failed to assign role to user model")
		span.RecordError(err)
		span.SetAttributes(attribute.String("error.type", "domain_role_assignment_error"))
		logger.Error("Failed to add role to domain user", slog.Any("error", err))
		return nil, err
	}

	if err = h.userRepo.Save(ctx, domUser); err != nil {
		span.SetStatus(codes.Error, "Failed to save user")
		span.RecordError(err)

		if errors.Is(err, user.ErrLineUserAlreadyExists) {
			span.SetAttributes(attribute.String("error.type", string(user.ErrLineUserAlreadyExists.Code)))
		} else {
			span.SetAttributes(attribute.String("error.type", "repository_write_error"))
		}
		logger.Error("Failed to save user to repository", "user_id", domUser.ID, slog.Any("error", err))
		return nil, err
	}

	span.SetStatus(codes.Ok, "User registered")
	logger.Info("User registered successfully.", "user_id", domUser.ID)

	return &RegisterUserDto{
		UserID: string(domUser.ID),
	}, nil
}
