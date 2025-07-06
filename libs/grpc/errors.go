package grpc

import (
	"context"
	"errors"
	"fmt"

	errs "github.com/pratchaya-maneechot/service-exchange/apps/users/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewGRPCErrCode(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, context.Canceled) {
		return status.Errorf(codes.Canceled, "request was cancelled: %s", err.Error())
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return status.Errorf(codes.DeadlineExceeded, "request timed out: %s", err.Error())
	}

	if domainErrorCode, ok := errs.GetErrorInternalCode(err); ok {
		switch domainErrorCode {
		case errs.CodeAlreadyExists:
			return status.Errorf(codes.AlreadyExists, "%s", err.Error())

		case errs.CodeNotFound:
			return status.Errorf(codes.NotFound, "%s", err.Error())

		// INVALID_ARGUMENT
		case errs.CodeInvalidArgument:
			return status.Errorf(codes.InvalidArgument, "%s", err.Error())

		default:
			fmt.Printf("Unhandled domain error code: %s - %s\n", domainErrorCode, err.Error())
			return status.Errorf(codes.Internal, "internal server error: an unmapped domain error occurred")
		}
	}

	fmt.Printf("Unhandled error type: %T - %v\n", err, err)
	return status.Errorf(codes.Internal, "internal server error: an unexpected error occurred")
}
