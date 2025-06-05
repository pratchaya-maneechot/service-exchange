package command

type UpdateUserProfile struct {
	ID       string
	Name     string
	Phone    *string
	Password *string
	Email    *string
}

var UpdateUserProfileCommand = "UpdateUserProfile"

func (c UpdateUserProfile) CommandName() string {
	return UpdateUserProfileCommand
}
