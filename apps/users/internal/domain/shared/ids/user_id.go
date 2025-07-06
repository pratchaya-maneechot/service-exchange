package ids

import "github.com/pratchaya-maneechot/service-exchange/libs/utils"

type UserID string

func NewUserID() UserID {
	return UserID(utils.GenerateUID())
}
