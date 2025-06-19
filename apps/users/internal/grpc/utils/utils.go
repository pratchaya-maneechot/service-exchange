package utils

import (
	"errors"

	errs "github.com/pratchaya-maneechot/service-exchange/apps/users/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// GetValuePointer converts google.protobuf.StringValue to *string
func GetValuePointer(sv *wrapperspb.StringValue) *string {
	if sv == nil {
		return nil
	}
	val := sv.GetValue()
	return &val
}

// ToStringInterfaceMap converts map<string, string> from proto to map[string]any
func ToStringInterfaceMap(m map[string]string) map[string]any {
	if m == nil {
		return nil
	}
	res := make(map[string]any)
	for k, v := range m {
		res[k] = v
	}
	return res
}

func MapErrorToGRPCCode(err error) error {
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
