package repository

import (
	"context"
	"sync"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
)

type impl struct {
	users map[string]*user.User
	mutex sync.RWMutex
}

// AddRoleToUser implements user.UserRepository.
func (i *impl) AddRoleToUser(ctx context.Context, userID ids.UserID, roleID uint) error {
	panic("unimplemented")
}

// ExistsByLineUserID implements user.UserRepository.
func (i *impl) ExistsByLineUserID(ctx context.Context, lineUserID string) (bool, error) {
	panic("unimplemented")
}

// FindByID implements user.UserRepository.
func (i *impl) FindByID(ctx context.Context, id ids.UserID) (*user.User, error) {
	panic("unimplemented")
}

// FindByLineUserID implements user.UserRepository.
func (i *impl) FindByLineUserID(ctx context.Context, lineUserID string) (*user.User, error) {
	panic("unimplemented")
}

// GetRoleByID implements user.UserRepository.
func (i *impl) GetRoleByID(ctx context.Context, roleID uint) (*user.Role, error) {
	panic("unimplemented")
}

// Save implements user.UserRepository.
func (i *impl) Save(ctx context.Context, user *user.User) error {
	panic("unimplemented")
}

func NewUserRepository() user.UserRepository {
	return &impl{
		users: make(map[string]*user.User),
		mutex: sync.RWMutex{},
	}
}
