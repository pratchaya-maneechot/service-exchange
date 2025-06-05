package handles

import (
	"context"
	"errors"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/command"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/common"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/entities"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/events"
	uid "github.com/pratchaya-maneechot/service-exchange/apps/users/pkg/uuid"
)

type LineRegisterHandler struct {
	repo     domain.Repository
	eventBus common.EventBus
}

func NewLineRegisterHandler(
	repo domain.Repository,
	eventBus common.EventBus,
) *LineRegisterHandler {
	return &LineRegisterHandler{
		repo:     repo,
		eventBus: eventBus,
	}
}

func (h LineRegisterHandler) Handle(ctx context.Context, cmd common.Command) (interface{}, error) {
	input := cmd.(command.LineRegister)
	profile, err := h.repo.LineProfile.GetByID(ctx, input.LineRefID)
	if err != nil {
		return nil, err
	}
	if profile != nil {
		return nil, errors.New(`user already have line account`)
	}
	userModel := &entities.User{
		ID:    uid.Generate(),
		Name:  input.DisplayName,
		Roles: []entities.UserRole{entities.RolePoster},
	}
	err = h.repo.User.Save(ctx, userModel)
	if err != nil {
		return nil, err
	}
	newProfile := &entities.LineProfile{
		ID:            input.LineRefID,
		UserID:        userModel.ID,
		DisplayName:   input.DisplayName,
		PictureURL:    input.PictureURL,
		StatusMessage: input.StatusMessage,
	}
	err = h.repo.LineProfile.Save(ctx, newProfile)
	if err != nil {
		return nil, err
	}
	h.eventBus.Publish(events.UserCreated{
		ID:        userModel.ID,
		Name:      userModel.Name,
		Phone:     userModel.Phone,
		Password:  userModel.Password,
		Email:     userModel.Email,
		Roles:     userModel.Roles,
		CreatedAt: userModel.CreatedAt,
		UpdatedAt: userModel.UpdatedAt,
	})
	return newProfile, nil
}
