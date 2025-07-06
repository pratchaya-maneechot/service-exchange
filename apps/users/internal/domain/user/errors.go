package user

import errs "github.com/pratchaya-maneechot/service-exchange/libs/errors"

var (
	ErrUserNotFound                        = errs.New(errs.CodeNotFound, "user not found")
	ErrLineUserIDAlreadyExists             = errs.New(errs.CodeAlreadyExists, "LINE user ID already exists")
	ErrEmailAlreadyExists                  = errs.New(errs.CodeAlreadyExists, "email already exists")
	ErrLineUserAlreadyExists               = errs.New(errs.CodeAlreadyExists, "line user already exists")
	ErrMissingLineIDOrEmail                = errs.New(errs.CodeInvalidArgument, "either LINE user ID or email must be provided")
	ErrInvalidCredentials                  = errs.New(errs.CodeUnauthorized, "invalid credentials")
	ErrRoleAlreadyAssigned                 = errs.New(errs.CodeAlreadyExists, "role already assigned to user")
	ErrInvalidVerificationStatusTransition = errs.New(errs.CodeInvalidArgument, "invalid identity verification status transition")
	ErrMissingDocumentURLs                 = errs.New(errs.CodeInvalidArgument, "document URLs are required for identity verification")
	ErrMissingDocumentType                 = errs.New(errs.CodeInvalidArgument, "document type is required for identity verification")
)
