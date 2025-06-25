package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/query"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/role"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/observability"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/pkg/utils"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type GetUserProfileQueryHandler struct {
	userRepo user.UserRepository
	logger   *slog.Logger
	config   *config.Config
}

func NewGetUserProfileQueryHandler(
	userRepo user.UserRepository,
	logger *slog.Logger,
	cfg *config.Config,
) *GetUserProfileQueryHandler {
	handlerLogger := logger.With(slog.String("component", "GetUserProfileQueryHandler"))
	return &GetUserProfileQueryHandler{
		userRepo: userRepo,
		logger:   handlerLogger,
		config:   cfg,
	}
}

func (h *GetUserProfileQueryHandler) Handle(ctx context.Context, qry query.GetUserProfileQuery) (*query.UserProfileDTO, error) {
	logger := observability.GetLoggerFromContext(ctx).With(
		slog.String("method", "HandleGetUserProfileQuery"),
		slog.String("user_id", string(qry.UserID)),
	)
	tracer := otel.Tracer(fmt.Sprintf("%s.query-handler", h.config.Name))
	ctx, span := tracer.Start(ctx, "GetUserProfileQueryHandler.Handle", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()
	span.SetAttributes(
		attribute.String("query.user_id", string(qry.UserID)),
	)
	logger.Debug("Attempting to retrieve user profile.")
	u, err := h.userRepo.FindByID(ctx, qry.UserID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			span.SetStatus(codes.Error, "User profile not found")
			logger.Warn("User profile not found for ID", "user_id", string(qry.UserID), slog.Any("error", err))
			return nil, user.ErrUserNotFound
		}
		span.SetStatus(codes.Error, "Failed to retrieve user profile from repository")
		span.RecordError(err)
		logger.Error("Failed to retrieve user profile from repository", "user_id", string(qry.UserID), slog.Any("error", err))
		return nil, err
	}
	span.SetStatus(codes.Ok, "User profile retrieved successfully")
	logger.Info("User profile retrieved successfully", "user_id", string(u.ID))
	resp := &query.UserProfileDTO{
		UserID:      string(u.ID),
		LineUserID:  u.LineUserID,
		Email:       u.Email,
		DisplayName: u.Profile.DisplayName,
		FirstName:   u.Profile.FirstName,
		LastName:    u.Profile.LastName,
		Bio:         u.Profile.Bio,
		AvatarURL:   u.Profile.AvatarURL,
		PhoneNumber: u.Profile.PhoneNumber,
		Address:     u.Profile.Address,
		Preferences: u.Profile.Preferences,
		Status:      u.Status,
		IsVerified:  u.IsVerified(),
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
		Roles: utils.ArrayMap(u.Roles, func(r role.Role) string {
			return string(r.Name)
		}),
	}
	logger.Debug("User profile DTO constructed successfully", "user_id", string(u.ID))
	return resp, nil
}
