package repositories

import (
	"context"
	"errors"
	"sync"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/entities"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/repositories"
)

type LineProfileRepository struct {
	profiles map[string]*entities.LineProfile
	mutex    sync.RWMutex
}

func NewLineProfileRepository() repositories.LineProfileRepository {
	return &LineProfileRepository{
		profiles: make(map[string]*entities.LineProfile),
		mutex:    sync.RWMutex{},
	}
}

func (r *LineProfileRepository) Save(ctx context.Context, entity *entities.LineProfile) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if entity == nil {
		return errors.New("nil line profile entity")
	}

	r.profiles[entity.UserID] = entity
	return nil
}

func (r *LineProfileRepository) GetByID(ctx context.Context, id string) (*entities.LineProfile, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	profile, exists := r.profiles[id]
	if !exists {
		return nil, nil
	}
	return profile, nil
}
