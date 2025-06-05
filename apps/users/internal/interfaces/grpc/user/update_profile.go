package server

import (
	"context"

	userpb "github.com/pratchaya-maneechot/service-exchange/apps/users/api/proto/user"
)

func (s *UserServer) UpdateProfile(ctx context.Context, req *userpb.UpdateProfileRequest) (*userpb.UpdateProfileResponse, error) {
	// TODO: Implement UpdateProfile
	return &userpb.UpdateProfileResponse{}, nil
}
