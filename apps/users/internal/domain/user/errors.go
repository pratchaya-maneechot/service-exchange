package user

import "errors"

var (
	ErrUserNotFound                        = errors.New("user not found")
	ErrLineUserIDAlreadyExists             = errors.New("LINE user ID already exists")
	ErrEmailAlreadyExists                  = errors.New("email already exists")
	ErrLineUserAlreadyExists               = errors.New("line user already exists")
	ErrMissingLineIDOrEmail                = errors.New("either LINE user ID or email must be provided")
	ErrInvalidCredentials                  = errors.New("invalid credentials")
	ErrRoleAlreadyAssigned                 = errors.New("role already assigned to user")
	ErrInvalidVerificationStatusTransition = errors.New("invalid identity verification status transition")
	ErrMissingDocumentURLs                 = errors.New("document URLs are required for identity verification")
	ErrMissingDocumentType                 = errors.New("document type is required for identity verification")
)
