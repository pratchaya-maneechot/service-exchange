package query

import (
	"context"
	"errors"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
)

// userQueryHandler implements the UserQueryHandler interface.
type userQueryHandler struct {
	userRepo user.UserRepository
}

// HandleGetUserIdentityVerification implements UserQueryHandler.
func (h *userQueryHandler) HandleGetUserIdentityVerification(query GetUserIdentityVerificationQuery) (*IdentityVerificationDTO, error) {
	panic("unimplemented")
}

// NewUserQueryHandler creates a new UserQueryHandler.
func NewUserQueryHandler(userRepo user.UserRepository) UserQueryHandler {
	return &userQueryHandler{userRepo: userRepo}
}

// HandleGetUserProfile handles the query to retrieve a user's profile.
func (h *userQueryHandler) HandleGetUserProfile(query GetUserProfileQuery) (*UserProfileDTO, error) {
	ctx := context.Background()

	u, err := h.userRepo.FindByID(ctx, query.UserID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	profileDTO := &UserProfileDTO{
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
	}

	return profileDTO, nil
}
