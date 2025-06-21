package handler

import (
	"context"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/command"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
)

type UpdateUserProfileCommandHandler struct {
	userRepo user.UserRepository
}

func NewUpdateUserProfileCommandHandler(
	ctx context.Context,
	userRepo user.UserRepository,
) *UpdateUserProfileCommandHandler {
	return &UpdateUserProfileCommandHandler{
		userRepo: userRepo,
	}
}

func (h *UpdateUserProfileCommandHandler) Handle(ctx context.Context, cmd command.UpdateUserProfileCommand) (*user.User, error) {
	domUser, err := h.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return domUser, nil
}
