package handlers

import (
	"context"
	"log/slog"

	pb "github.com/pratchaya-maneechot/service-exchange/apps/users/api/proto/user"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/command"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/query"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/grpc/mappers"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/grpc/utils"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/pkg/bus"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/pkg/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserGRPCHandler struct {
	pb.UnimplementedUserServiceServer
	commandBus bus.CommandBus
	queryBus   bus.QueryBus
	logger     *slog.Logger
}

func RegisUserGRPCHandler(
	gs *grpc.Server,
	cb bus.CommandBus, qb bus.QueryBus, logger *slog.Logger,
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
		Email:       utils.GetValuePointer(req.GetEmail()),
		Password:    utils.GetValuePointer(req.GetPassword()),
		DisplayName: req.GetDisplayName(),
		AvatarURL:   req.GetAvatarUrl().GetValue(),
	}
	result, err := h.commandBus.Dispatch(ctx, cmd)
	if err != nil {
		return nil, grpcutil.NewGRPCErrCode(err)
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
		DisplayName: utils.GetValuePointer(req.GetDisplayName()),
		FirstName:   utils.GetValuePointer(req.GetFirstName()),
		LastName:    utils.GetValuePointer(req.GetLastName()),
		Bio:         utils.GetValuePointer(req.GetBio()),
		AvatarURL:   utils.GetValuePointer(req.GetAvatarUrl()),
		PhoneNumber: utils.GetValuePointer(req.GetPhoneNumber()),
		Address:     utils.GetValuePointer(req.GetAddress()),
		Preferences: utils.GetStringInterface(req.GetPreferences()),
	}
	if _, err := h.commandBus.Dispatch(ctx, cmd); err != nil {
		h.logger.Error("Failed to dispatch UpdateUserProfileCommand", "error", err)
		return nil, grpcutil.NewGRPCErrCode(err)
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
		return nil, grpcutil.NewGRPCErrCode(err)
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
