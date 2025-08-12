package query

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/role"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
	"github.com/pratchaya-maneechot/service-exchange/libs/infra/observability"
	"github.com/pratchaya-maneechot/service-exchange/libs/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type GetUserProfileQuery struct {
	UserID ids.UserID `json:"userId"`
}

type UserProfileDTO struct {
	UserID      string          `json:"userId"`
	LineUserID  string          `json:"lineUserId"`
	Email       *string         `json:"email"`
	DisplayName string          `json:"displayName"`
	FirstName   *string         `json:"firstName,omitempty"`
	LastName    *string         `json:"lastName,omitempty"`
	Bio         *string         `json:"bio,omitempty"`
	AvatarURL   *string         `json:"avatarUrl,omitempty"`
	PhoneNumber *string         `json:"phoneNumber,omitempty"`
	Address     *string         `json:"address,omitempty"`
	Preferences map[string]any  `json:"preferences"`
	Status      user.UserStatus `json:"status"`
	IsVerified  bool            `json:"isVerified"`
	LastLoginAt *time.Time      `json:"lastLoginAt,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
	Roles       []string        `json:"roles"`
}

type GetUserProfileQueryHandler struct {
	userRepo user.UserRepository
	logger   *slog.Logger
	config   *config.Config
	tracer   trace.Tracer
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
		tracer:   otel.Tracer(fmt.Sprintf("%s.query-handler", cfg.Name)),
	}
}

func (h *GetUserProfileQueryHandler) Handle(ctx context.Context, qry GetUserProfileQuery) (*UserProfileDTO, error) {
	logger := observability.LoggerFromCtx(ctx).With(slog.String("user_id", string(qry.UserID)))

	ctx, span := h.tracer.Start(ctx, "GetUserProfileQueryHandler.Handle", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	span.SetAttributes(attribute.String("query.user_id", string(qry.UserID)))

	u, err := h.userRepo.FindByID(ctx, qry.UserID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			span.SetStatus(codes.Ok, "User profile not found")
			span.SetAttributes(attribute.String("error.type", string(user.ErrUserNotFound.Code)))
			span.SetAttributes(attribute.Bool("user.found", false))
			logger.Warn(user.ErrUserNotFound.Message, "user_id", string(qry.UserID))
			return nil, user.ErrUserNotFound
		}

		span.SetStatus(codes.Error, "Failed to retrieve user profile")
		span.RecordError(err)
		span.SetAttributes(attribute.String("error.type", "repository_read_error"))
		logger.Error("Failed to retrieve user profile from repository", "user_id", string(qry.UserID), slog.Any("error", err))
		return nil, err
	}

	span.SetStatus(codes.Ok, "User profile retrieved")
	span.SetAttributes(attribute.Bool("user.found", true))
	logger.Info("User profile retrieved successfully.", "user_id", string(u.ID))

	resp := &UserProfileDTO{
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
	return resp, nil
}
