package handles

import (
	"context"
	"errors"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/command"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/common"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/entities"
)

type UpdateUserProfileHandler struct {
	repo domain.Repository
}

func NewUpdateUserProfileHandler(
	repo domain.Repository,
) *UpdateUserProfileHandler {
	return &UpdateUserProfileHandler{
		repo: repo,
	}
}

func (h UpdateUserProfileHandler) Handle(ctx context.Context, cmd common.Command) (interface{}, error) {
	input := cmd.(command.UpdateUserProfile)
	prev, err := h.repo.User.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if prev == nil {
		return nil, errors.New(`user not found`)
	}
	err = h.repo.User.Save(ctx, &entities.User{
		Name:     input.Name,
		Phone:    input.Phone,
		Password: input.Password,
		Email:    input.Email,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
