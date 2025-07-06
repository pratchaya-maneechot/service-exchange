package handlers

import (
	"context"
	"log/slog"

	pb "github.com/pratchaya-maneechot/service-exchange/apps/users/api/proto/user"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/command"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/query"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/grpc/mappers"
	libCommand "github.com/pratchaya-maneechot/service-exchange/libs/bus/command"
	libQuery "github.com/pratchaya-maneechot/service-exchange/libs/bus/query"
	libGrpc "github.com/pratchaya-maneechot/service-exchange/libs/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserGRPCHandler struct {
	pb.UnimplementedUserServiceServer
	commandBus libCommand.CommandBus
	queryBus   libQuery.QueryBus
	logger     *slog.Logger
}

func RegisUserGRPCHandler(
	gs *grpc.Server,
	cb libCommand.CommandBus, qb libQuery.QueryBus, logger *slog.Logger,
) {
	server := &UserGRPCHandler{
		commandBus: cb,
		queryBus:   qb,
		logger:     logger,
	}
	pb.RegisterUserServiceServer(gs, server)
}

func (h *UserGRPCHandler) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	cmd := command.RegisterUserCommand{
		LineUserID:  req.GetLineUserId(),
		Email:       libGrpc.StringValueToPtr(req.GetEmail()),
		Password:    libGrpc.StringValueToPtr(req.GetPassword()),
		DisplayName: req.GetDisplayName(),
		AvatarURL:   req.GetAvatarUrl().GetValue(),
	}
	result, err := h.commandBus.Dispatch(ctx, cmd)
	if err != nil {
		return nil, libGrpc.NewGRPCErrCode(err)
	}
	usr, ok := result.(*command.RegisterUserDto)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal server error: unexpected response from RegisterUserCommand handler")
	}
	return &pb.RegisterUserResponse{
		UserId:           string(usr.UserID),
		JwtToken:         string(usr.UserID),
		ExpiresInSeconds: int32(0),
	}, nil
}

func (h *UserGRPCHandler) UserLogin(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	panic("un implement")
}

func (h *UserGRPCHandler) LineLogin(ctx context.Context, req *pb.LineLoginRequest) (*pb.LineLoginResponse, error) {
	panic("un implement")
}

func (h *UserGRPCHandler) UpdateUserProfile(ctx context.Context, req *pb.UpdateUserProfileRequest) (*emptypb.Empty, error) {
	userID := ids.UserID(req.GetUserId())
	cmd := command.UpdateUserProfileCommand{
		UserID:      userID,
		DisplayName: libGrpc.StringValueToPtr(req.GetDisplayName()),
		FirstName:   libGrpc.StringValueToPtr(req.GetFirstName()),
		LastName:    libGrpc.StringValueToPtr(req.GetLastName()),
		Bio:         libGrpc.StringValueToPtr(req.GetBio()),
		AvatarURL:   libGrpc.StringValueToPtr(req.GetAvatarUrl()),
		PhoneNumber: libGrpc.StringValueToPtr(req.GetPhoneNumber()),
		Address:     libGrpc.StringValueToPtr(req.GetAddress()),
		Preferences: libGrpc.StringMapToAnyMap(req.GetPreferences()),
	}
	if _, err := h.commandBus.Dispatch(ctx, cmd); err != nil {
		h.logger.Error("Failed to dispatch UpdateUserProfileCommand", "error", err)
		return nil, libGrpc.NewGRPCErrCode(err)
	}
	return &emptypb.Empty{}, nil
}

func (h *UserGRPCHandler) SubmitIdentityVerification(ctx context.Context, req *pb.SubmitIdentityVerificationRequest) (*emptypb.Empty, error) {
	panic("un implement")
}

func (h *UserGRPCHandler) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.UserProfileDTO, error) {
	userID := ids.UserID(req.GetUserId())
	qry := query.GetUserProfileQuery{
		UserID: userID,
	}
	result, err := h.queryBus.Dispatch(ctx, qry)
	if err != nil {
		return nil, libGrpc.NewGRPCErrCode(err)
	}
	internalDTO, ok := result.(*query.UserProfileDTO)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal server error: unexpected response from GetUserProfileQuery handler")
	}
	return mappers.MapUserProfileDTOToProto(internalDTO), nil
}

func (h *UserGRPCHandler) GetUserIdentityVerification(ctx context.Context, req *pb.GetUserIdentityVerificationRequest) (*pb.IdentityVerificationDTO, error) {
	panic("un implement")
}
