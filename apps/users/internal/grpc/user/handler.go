package handler

import (
	"context"
	"errors"
	"log/slog"

	pb "github.com/pratchaya-maneechot/service-exchange/apps/users/api/proto/user"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/command"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/query"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/pkg/bus"
	errs "github.com/pratchaya-maneechot/service-exchange/apps/users/pkg/errors"
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

// --- Helper for Error Mapping ---
func mapErrorToGRPCCode(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, errs.ErrNotFound) {
		return status.Errorf(codes.NotFound, "%s", err.Error())
	}
	if errors.Is(err, errs.ErrValidation) || errors.Is(err, errs.ErrInvalidArgument) {
		return status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}
	if errors.Is(err, errs.ErrAlreadyExists) {
		return status.Errorf(codes.AlreadyExists, "%s", err.Error())
	}
	// Add other error mappings as needed (e.g., Unauthorized, Forbidden)
	return status.Errorf(codes.Internal, "internal server error: %v", err)
}

// --- Command Handlers ---

func (h *UserGRPCHandler) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	cmd := command.RegisterUserCommand{
		LineUserID:  req.GetLineUserId(),
		Email:       req.GetEmail().GetValue(),
		Password:    req.GetPassword().GetValue(),
		DisplayName: req.GetDisplayName(),
		AvatarURL:   req.GetAvatarUrl().GetValue(),
	}

	result, err := h.commandBus.Dispatch(ctx, cmd)
	if err != nil {
		h.logger.Error("Failed to dispatch RegisterUserCommand", "error", err)
		return nil, mapErrorToGRPCCode(err)
	}

	// Type assertion for the result from the bus
	res, ok := result.(struct {
		UserID    string
		JWT       string
		ExpiresIn int
	})
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal server error: unexpected response from RegisterUserCommand handler")
	}

	return &pb.RegisterUserResponse{
		UserId:           res.UserID,
		JwtToken:         res.JWT,
		ExpiresInSeconds: int32(res.ExpiresIn),
	}, nil
}

func (h *UserGRPCHandler) UserLogin(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	cmd := command.UserLoginCommand{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	result, err := h.commandBus.Dispatch(ctx, cmd)
	if err != nil {
		h.logger.Error("Failed to dispatch UserLoginCommand", "error", err)
		return nil, mapErrorToGRPCCode(err)
	}

	res, ok := result.(struct {
		JWT       string
		ExpiresIn int
	})
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal server error: unexpected response from UserLoginCommand handler")
	}

	return &pb.UserLoginResponse{
		JwtToken:         res.JWT,
		ExpiresInSeconds: int32(res.ExpiresIn),
	}, nil
}

func (h *UserGRPCHandler) LineLogin(ctx context.Context, req *pb.LineLoginRequest) (*pb.LineLoginResponse, error) {
	cmd := command.LineLoginCommand{
		LineUserID: req.GetLineUserId(),
		IDToken:    req.GetIdToken(),
	}

	result, err := h.commandBus.Dispatch(ctx, cmd)
	if err != nil {
		h.logger.Error("Failed to dispatch LineLoginCommand", "error", err)
		return nil, mapErrorToGRPCCode(err)
	}

	res, ok := result.(struct {
		UserID    string
		JWT       string
		ExpiresIn int
	})
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal server error: unexpected response from LineLoginCommand handler")
	}

	return &pb.LineLoginResponse{
		UserId:           res.UserID,
		JwtToken:         res.JWT,
		ExpiresInSeconds: int32(res.ExpiresIn),
	}, nil
}

func (h *UserGRPCHandler) UpdateUserProfile(ctx context.Context, req *pb.UpdateUserProfileRequest) (*emptypb.Empty, error) {
	userID := ids.UserID(req.GetUserId())
	cmd := command.UpdateUserProfileCommand{
		UserID:      userID,
		DisplayName: GetValuePointer(req.GetDisplayName()),
		FirstName:   GetValuePointer(req.GetFirstName()),
		LastName:    GetValuePointer(req.GetLastName()),
		Bio:         GetValuePointer(req.GetBio()),
		AvatarURL:   GetValuePointer(req.GetAvatarUrl()),
		PhoneNumber: GetValuePointer(req.GetPhoneNumber()),
		Address:     GetValuePointer(req.GetAddress()),
		Preferences: ToStringInterfaceMap(req.GetPreferences()),
	}

	_, err := h.commandBus.Dispatch(ctx, cmd)
	if err != nil {
		h.logger.Error("Failed to dispatch UpdateUserProfileCommand", "error", err)
		return nil, mapErrorToGRPCCode(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *UserGRPCHandler) SubmitIdentityVerification(ctx context.Context, req *pb.SubmitIdentityVerificationRequest) (*emptypb.Empty, error) {
	userID := ids.UserID(req.GetUserId())

	cmd := command.SubmitIdentityVerificationCommand{
		UserID:         userID,
		DocumentType:   MapProtoDocumentTypeToDomain(req.GetDocumentType()),
		DocumentURLs:   req.GetDocumentUrls(),
		DocumentNumber: req.GetDocumentNumber().GetValue(),
	}

	_, err := h.commandBus.Dispatch(ctx, cmd)
	if err != nil {
		h.logger.Error("Failed to dispatch SubmitIdentityVerificationCommand", "error", err)
		return nil, mapErrorToGRPCCode(err)
	}

	return &emptypb.Empty{}, nil
}

// --- Query Handlers ---

func (h *UserGRPCHandler) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.UserProfileDTO, error) {
	userID := ids.UserID(req.GetUserId())
	qry := query.GetUserProfileQuery{
		UserID: userID,
	}

	result, err := h.queryBus.Dispatch(ctx, qry)
	if err != nil {
		h.logger.Error("Failed to dispatch GetUserProfileQuery", "error", err)
		return nil, mapErrorToGRPCCode(err)
	}

	internalDTO, ok := result.(*query.UserProfileDTO)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal server error: unexpected response from GetUserProfileQuery handler")
	}

	return h.MapUserProfileDTOToProto(internalDTO), nil
}

func (h *UserGRPCHandler) GetUserIdentityVerification(ctx context.Context, req *pb.GetUserIdentityVerificationRequest) (*pb.IdentityVerificationDTO, error) {
	userID := ids.UserID(req.GetUserId())
	qry := query.GetUserIdentityVerificationQuery{
		UserID: userID,
	}

	result, err := h.queryBus.Dispatch(ctx, qry)
	if err != nil {
		h.logger.Error("Failed to dispatch GetUserIdentityVerificationQuery", "error", err)
		return nil, mapErrorToGRPCCode(err)
	}

	internalDTO, ok := result.(*query.IdentityVerificationDTO)
	if !ok {
		return nil, status.Errorf(codes.Internal, "internal server error: unexpected response from GetUserIdentityVerificationQuery handler")
	}

	return h.MapIdentityVerificationDTOToProto(internalDTO), nil
}
