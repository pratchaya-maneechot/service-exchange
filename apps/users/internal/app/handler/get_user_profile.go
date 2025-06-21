package handler

import (
	"context"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/query"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/role"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/pkg/utils"
)

type GetUserProfileQueryHandler struct {
	userRepo user.UserRepository
}

func NewGetUserProfileQueryHandler(
	ctx context.Context,
	userRepo user.UserRepository,
) *GetUserProfileQueryHandler {
	return &GetUserProfileQueryHandler{
		userRepo: userRepo,
	}
}

// HandleGetUserProfile handles the query to retrieve a user's profile.
func (h *GetUserProfileQueryHandler) Handle(ctx context.Context, qry query.GetUserProfileQuery) (*query.UserProfileDTO, error) {
	u, err := h.userRepo.FindByID(ctx, qry.UserID)
	if err != nil {
		return nil, err
	}
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
	return resp, nil
}
