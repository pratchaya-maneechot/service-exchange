package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/command"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/observability"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type UpdateUserProfileCommandHandler struct {
	userRepo user.UserRepository
	logger   *slog.Logger
	config   *config.Config
}

func NewUpdateUserProfileCommandHandler(
	userRepo user.UserRepository,
	logger *slog.Logger,
	cfg *config.Config,
) *UpdateUserProfileCommandHandler {
	handlerLogger := logger.With(slog.String("component", "UpdateUserProfileCommandHandler"))
	return &UpdateUserProfileCommandHandler{
		userRepo: userRepo,
		logger:   handlerLogger,
		config:   cfg,
	}
}

func (h *UpdateUserProfileCommandHandler) Handle(ctx context.Context, cmd command.UpdateUserProfileCommand) (*user.User, error) {
	logger := observability.GetLoggerFromContext(ctx).With(
		slog.String("method", "HandleUpdateUserProfileCommand"),
		slog.String("user_id", string(cmd.UserID)),
	)
	tracer := otel.Tracer(fmt.Sprintf("%s.command-handler", h.config.Name))
	ctx, span := tracer.Start(ctx, "UpdateUserProfileCommandHandler.Handle", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()
	span.SetAttributes(
		attribute.String("user.id", string(cmd.UserID)),
		attribute.String("profile.display_name", *cmd.DisplayName),
	)
	logger.Debug("Attempting to update user profile.")
	domUser, err := h.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			span.SetStatus(codes.Error, "User for profile update not found")
			span.SetAttributes(attribute.String("error.type", "user_not_found"))
			logger.Warn("User for profile update not found", "user_id", string(cmd.UserID), slog.Any("error", err))
			return nil, user.ErrUserNotFound
		}
		span.SetStatus(codes.Error, "Failed to find user for profile update")
		span.RecordError(err)
		logger.Error("Failed to find user for profile update", "user_id", string(cmd.UserID), slog.Any("error", err))
		return nil, err
	}
	logger.Debug("User found for profile update", "user_id", string(domUser.ID))
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
	logger.Debug("User profile updated in domain model", "user_id", string(domUser.ID))
	if err = h.userRepo.Save(ctx, domUser); err != nil {
		span.SetStatus(codes.Error, "Failed to save updated user profile to repository")
		span.RecordError(err)
		logger.Error("Failed to save updated user profile to repository", "user_id", string(domUser.ID), slog.Any("error", err))
		return nil, err
	}
	span.SetStatus(codes.Ok, "User profile updated successfully")
	logger.Info("User profile updated successfully", "user_id", string(domUser.ID))
	return domUser, nil
}
