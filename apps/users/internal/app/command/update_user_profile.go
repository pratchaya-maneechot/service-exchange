package command

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
	"github.com/pratchaya-maneechot/service-exchange/libs/infra/observability"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type UpdateUserProfileCommand struct {
	UserID      ids.UserID     `json:"-"`
	DisplayName *string        `json:"displayName,omitempty"`
	FirstName   *string        `json:"firstName,omitempty"`
	LastName    *string        `json:"lastName,omitempty"`
	Bio         *string        `json:"bio,omitempty"`
	AvatarURL   *string        `json:"avatarUrl,omitempty"`
	PhoneNumber *string        `json:"phoneNumber,omitempty" validate:"omitempty,e164"`
	Address     *string        `json:"address,omitempty"`
	Preferences map[string]any `json:"preferences,omitempty"`
}

type UpdateUserProfileDto struct {
	UserID string `json:"UserId" validate:"required"`
}

type UpdateUserProfileCommandHandler struct {
	userRepo user.UserRepository
	logger   *slog.Logger
	config   *config.Config
	tracer   trace.Tracer
}

func NewUpdateUserProfileCommandHandler(
	userRepo user.UserRepository,
	logger *slog.Logger,
	cfg *config.Config,
) *UpdateUserProfileCommandHandler {
	return &UpdateUserProfileCommandHandler{
		userRepo: userRepo,
		logger:   logger.With(slog.String("component", "UpdateUserProfileCommandHandler")),
		config:   cfg,
		tracer:   otel.Tracer(fmt.Sprintf("%s.command-handler", cfg.Name)),
	}
}
func (h *UpdateUserProfileCommandHandler) Handle(ctx context.Context, cmd UpdateUserProfileCommand) (*UpdateUserProfileDto, error) {
	logger := observability.LoggerFromCtx(ctx).With(slog.String("user_id", string(cmd.UserID)))

	ctx, span := h.tracer.Start(ctx, "UpdateUserProfileCommand.Handle", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	span.SetAttributes(attribute.String("user.id", string(cmd.UserID)))

	if cmd.DisplayName != nil {
		span.SetAttributes(attribute.String("profile.display_name", *cmd.DisplayName))
	}

	domUser, err := h.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to retrieve user")
		span.RecordError(err)
		span.SetAttributes(attribute.String("error.type", "repository_read_error"))
		logger.Error("Failed to retrieve user for profile update", "user_id", string(cmd.UserID), slog.Any("error", err))
		return nil, err
	}
	if domUser == nil {
		span.SetStatus(codes.Error, "User not found")
		span.SetAttributes(attribute.String("error.type", string(user.ErrUserNotFound.Code)))
		logger.Warn(user.ErrUserNotFound.Message, "user_id", cmd.UserID)
		return nil, user.ErrUserNotFound
	}

	domUser.UpdateProfile(
		cmd.DisplayName,
		cmd.FirstName,
		cmd.LastName,
		cmd.Bio,
		cmd.AvatarURL,
		cmd.PhoneNumber,
		cmd.Address,
		cmd.Preferences,
	)

	if err = h.userRepo.Save(ctx, domUser); err != nil {
		span.SetStatus(codes.Error, "Failed to save user profile")
		span.RecordError(err)
		span.SetAttributes(attribute.String("error.type", "repository_write_error"))
		logger.Error("Failed to save updated user profile to repository", "user_id", domUser.ID, slog.Any("error", err))
		return nil, err
	}

	span.SetStatus(codes.Ok, "User profile updated")
	logger.Info("User profile updated successfully.", "user_id", domUser.ID)

	return &UpdateUserProfileDto{
		UserID: string(domUser.ID),
	}, nil
}
