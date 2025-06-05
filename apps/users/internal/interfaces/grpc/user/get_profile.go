package server

import (
	"context"

	userpb "github.com/pratchaya-maneechot/service-exchange/apps/users/api/proto/user"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *UserServer) GetProfile(ctx context.Context, req *userpb.GetProfileRequest) (*userpb.GetProfileResponse, error) {
	user := &userpb.User{
		Id:        req.UserId,
		Name:      "John Doe",
		Phone:     wrapperspb.String("123-456-7890"),
		Email:     wrapperspb.String("john.doe@example.com"),
		Password:  wrapperspb.String("hashedpassword123"), // In practice, this would be a hashed value
		Roles:     []userpb.UserRole{userpb.UserRole_POSTER, userpb.UserRole_TASKER},
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}

	lineProfile := &userpb.LineProfile{
		Id:            "U4af4980629abc",
		UserId:        req.UserId,
		DisplayName:   "Johnny",
		PictureUrl:    "https://example.com/profile.jpg",
		StatusMessage: wrapperspb.String("Available for tasks!"),
	}

	response := &userpb.GetProfileResponse{
		User:        user,
		LineProfile: lineProfile,
	}
	return response, nil
}
