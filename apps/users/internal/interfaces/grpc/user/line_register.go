package server

import (
	"context"

	"github.com/pkg/errors"
	userpb "github.com/pratchaya-maneechot/service-exchange/apps/users/api/proto/user"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/command"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/entities"
)

func (s *UserServer) LineRegister(ctx context.Context, req *userpb.LineRegisterRequest) (*userpb.LineRegisterResponse, error) {
	result, err := s.commandBus.Send(ctx, command.LineRegister{
		LineRefID:   req.GetLineId(),
		DisplayName: req.GetDisplayName(),
		PictureURL:  req.GetPictureUrl(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to process LineRegister")
	}
	profile := result.(*entities.LineProfile)
	return &userpb.LineRegisterResponse{
		UserId: profile.UserID,
	}, nil
}
