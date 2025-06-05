package handles

import (
	"context"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/query"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/common"
)

type GetUserProfileHandler struct {
	repo domain.Repository
}

func NewGetUserProfileHandler(
	repo domain.Repository,
) *GetUserProfileHandler {
	return &GetUserProfileHandler{
		repo: repo,
	}
}
func (h GetUserProfileHandler) Handle(ctx context.Context, qry common.Query) (interface{}, error) {
	params := qry.(query.GetUserProfile)
	profile, err := h.repo.User.GetByID(ctx, params.Id)
	if err != nil {
		return nil, err
	}
	return profile, nil
}
