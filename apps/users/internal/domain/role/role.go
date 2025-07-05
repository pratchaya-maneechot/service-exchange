package role

type RoleName string

const (
	RoleNamePoster RoleName = "POSTER"
	RoleNameTasker RoleName = "TASKER"
	RoleNameAdmin  RoleName = "ADMIN"
)

type Role struct {
	ID          uint
	Name        RoleName
	Description *string
}

func NewRoleFromRepository(id uint, name string, description *string) *Role {
	return &Role{
		ID:          id,
		Name:        RoleName(name),
		Description: description,
	}
}
