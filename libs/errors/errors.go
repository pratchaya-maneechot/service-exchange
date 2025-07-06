package errs

import (
	"errors"
)

var (
	ErrNotFound        = New(CodeNotFound, "not found")
	ErrAlreadyExists   = New(CodeAlreadyExists, "already exists")
	ErrInternal        = New(CodeInternal, "internal server error")
	ErrUnauthorized    = New(CodeUnauthorized, "unauthorized")
	ErrForbidden       = New(CodeForbidden, "forbidden")
	ErrInvalidArgument = New(CodeInvalidArgument, "invalid argument")
)

type ErrorInternal struct {
	Code    Code
	Message string
}

func (e *ErrorInternal) Error() string {
	return e.Message
}

func New(code Code, msg string) *ErrorInternal {
	return &ErrorInternal{
		Code:    code,
		Message: msg,
	}
}

func IsErrorInternal(err error) bool {
	var de *ErrorInternal
	return errors.As(err, &de)
}

func GetErrorInternalCode(err error) (Code, bool) {
	var de *ErrorInternal
	if errors.As(err, &de) {
		return de.Code, true
	}
	return "", false
}
