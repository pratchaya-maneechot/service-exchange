package handler

import (
	"context"
	"errors"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewGRPCErrCode(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, context.Canceled):
		return status.Errorf(codes.Canceled, "request was cancelled: %v", err)
	case errors.Is(err, context.DeadlineExceeded):
		return status.Errorf(codes.DeadlineExceeded, "request timed out: %v", err)
	case errors.Is(err, user.ErrUserNotFound):
		return status.Errorf(codes.NotFound, "%s", err.Error())
	case errors.Is(err, user.ErrLineUserIDAlreadyExists), errors.Is(err, user.ErrLineUserAlreadyExists):
		return status.Errorf(codes.AlreadyExists, "%s", err.Error())
	case errors.Is(err, user.ErrEmailAlreadyExists):
		return status.Errorf(codes.AlreadyExists, "%s", err.Error())
	case errors.Is(err, user.ErrMissingLineIDOrEmail):
		return status.Errorf(codes.InvalidArgument, "%s", err.Error())
	case errors.Is(err, user.ErrInvalidCredentials):
		return status.Errorf(codes.Unauthenticated, "%s", err.Error())
	case errors.Is(err, user.ErrRoleAlreadyAssigned):
		return status.Errorf(codes.FailedPrecondition, "%s", err.Error())
	case errors.Is(err, user.ErrInvalidVerificationStatusTransition):
		return status.Errorf(codes.FailedPrecondition, "%s", err.Error())
	case errors.Is(err, user.ErrMissingDocumentURLs), errors.Is(err, user.ErrMissingDocumentType):
		return status.Errorf(codes.InvalidArgument, "%s", err.Error())
	default:
		return status.Errorf(codes.Internal, "internal server error: %s", err.Error())
	}
}
