package role

import errs "github.com/pratchaya-maneechot/service-exchange/libs/errors"

var (
	ErrRoleNotFound = errs.New(errs.CodeNotFound, "role not found")
)
