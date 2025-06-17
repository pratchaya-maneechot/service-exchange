package user

type RoleName string

const (
	RoleNamePoster RoleName = "POSTER"
	RoleNameTasker RoleName = "TASKER"
	RoleNameAdmin  RoleName = "ADMIN"
)

type Role struct {
	ID          uint // Internal ID, managed by system (e.g., from a static lookup table)
	Name        RoleName
	Description string
}

// NewRole creates a new Role instance.
func NewRole(id uint, name RoleName, description string) Role {
	return Role{
		ID:          id,
		Name:        name,
		Description: description,
	}
}
