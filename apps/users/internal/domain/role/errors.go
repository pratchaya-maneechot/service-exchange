package role

import errs "github.com/pratchaya-maneechot/service-exchange/apps/users/pkg/errors"

var (
	ErrRoleNotFound = errs.New(errs.CodeNotFound, "role not found")
)
