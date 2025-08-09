package handlers

import (
	"context"

	pb "github.com/pratchaya-maneechot/service-exchange/apps/users/api/proto/user"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/command"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/query"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/grpc/views"
	lg "github.com/pratchaya-maneechot/service-exchange/libs/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserGRPCHandler struct {
	pb.UnimplementedUserServiceServer
	lg.GrpcHandlerOption
}

func RegisUserGRPCHandler(
	gs *grpc.Server,
	opt lg.GrpcHandlerOption,
) {
	pb.RegisterUserServiceServer(gs, &UserGRPCHandler{
		GrpcHandlerOption: opt,
	})
}

func (h *UserGRPCHandler) LineRegister(ctx context.Context, req *pb.LineRegisterRequest) (*pb.LineRegisterResponse, error) {
	cmd := command.RegisterUserCommand{
		LineUserID:  req.GetLineUserId(),
		Email:       lg.StringValueToPtr(req.GetEmail()),
		Password:    lg.StringValueToPtr(req.GetPassword()),
		DisplayName: req.GetDisplayName(),
		AvatarURL:   lg.StringValueToPtr(req.GetAvatarUrl()),
	}
	if err := h.Validator.Struct(cmd); err != nil {
		return nil, h.ValidationErrors(err)
	}

	result, err := h.Command.Dispatch(ctx, cmd)
	if err != nil {
		return nil, lg.NewGRPCErrCode(err)
	}

	usr, ok := result.(*command.RegisterUserDto)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal server error: unexpected response from RegisterUserCommand handler")
	}

	return &pb.LineRegisterResponse{
		UserId: string(usr.UserID),
	}, nil
}

func (h *UserGRPCHandler) UpdateUserProfile(ctx context.Context, req *pb.UpdateUserProfileRequest) (*emptypb.Empty, error) {
	userID := ids.UserID(req.GetUserId())
	cmd := command.UpdateUserProfileCommand{
		UserID:      userID,
		DisplayName: lg.StringValueToPtr(req.GetDisplayName()),
		FirstName:   lg.StringValueToPtr(req.GetFirstName()),
		LastName:    lg.StringValueToPtr(req.GetLastName()),
		Bio:         lg.StringValueToPtr(req.GetBio()),
		AvatarURL:   lg.StringValueToPtr(req.GetAvatarUrl()),
		PhoneNumber: lg.StringValueToPtr(req.GetPhoneNumber()),
		Address:     lg.StringValueToPtr(req.GetAddress()),
		Preferences: lg.StringMapToAnyMap(req.GetPreferences()),
	}
	if _, err := h.Command.Dispatch(ctx, cmd); err != nil {
		h.Logger.Error("Failed to dispatch UpdateUserProfileCommand", "error", err)
		return nil, lg.NewGRPCErrCode(err)
	}
	return &emptypb.Empty{}, nil
}

func (h *UserGRPCHandler) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.UserProfile, error) {
	userID := ids.UserID(req.UserId)
	qry := query.GetUserProfileQuery{
		UserID: userID,
	}
	result, err := h.Query.Dispatch(ctx, qry)
	if err != nil {
		return nil, lg.NewGRPCErrCode(err)
	}
	internalDTO, ok := result.(*query.UserProfileDTO)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal server error: unexpected response from GetUserProfileQuery handler")
	}
	return views.UserProfile(internalDTO), nil
}
