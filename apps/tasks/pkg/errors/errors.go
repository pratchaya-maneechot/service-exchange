package errors

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrAlreadyExists   = errors.New("already exists")
	ErrValidation      = errors.New("validation error")
	ErrInternal        = errors.New("internal server error")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrInvalidArgument = errors.New("invalid argument")
)

func New(msg string) error {
	return errors.New(msg)
}

func Wrap(err error, msg string) error {
	return errors.Join(errors.New(msg), err)
}
