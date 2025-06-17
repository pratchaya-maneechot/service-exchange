package ids

import "github.com/pratchaya-maneechot/service-exchange/apps/users/pkg/utils"

type UserID string

func NewUserID() UserID {
	return UserID(utils.GenerateUID())
}
