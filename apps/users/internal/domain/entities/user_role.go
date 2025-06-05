package entities

import "slices"

type UserRole string

const (
	RolePoster UserRole = "Poster"
	RoleTasker UserRole = "Tasker"
)

var ValidRoles = []UserRole{RolePoster, RoleTasker}

func (r UserRole) IsValid() bool {
	return slices.Contains(ValidRoles, r)
}

func (r UserRole) String() string {
	return string(r)
}
