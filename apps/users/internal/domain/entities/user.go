package entities

import (
	"slices"
	"time"
)

type User struct {
	ID        string
	Name      string
	Phone     *string
	Password  *string
	Email     *string
	Roles     []UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) IsPoster() bool {
	return slices.Contains(u.Roles, RolePoster)
}

func (u *User) IsTasker() bool {
	return slices.Contains(u.Roles, RoleTasker)
}
