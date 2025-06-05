package repositories

import (
	"context"
	"errors"
	"sync"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/entities"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/repositories"
)

type UserRepository struct {
	users map[string]*entities.User
	mutex sync.RWMutex
}

func NewUserRepository() repositories.UserRepository {
	return &UserRepository{
		users: make(map[string]*entities.User),
		mutex: sync.RWMutex{},
	}
}

func (r *UserRepository) Save(ctx context.Context, entity *entities.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if entity == nil {
		return errors.New("nil user entity")
	}

	r.users[entity.ID] = entity
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(r.users, id)
	return nil
}
