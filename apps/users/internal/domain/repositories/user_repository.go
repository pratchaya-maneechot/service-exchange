package repositories

import (
	"context"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/entities"
)

type UserRepository interface {
	Save(ctx context.Context, entity *entities.User) error
	GetByID(ctx context.Context, id string) (*entities.User, error)
	Delete(ctx context.Context, id string) error
}

type LineProfileRepository interface {
	Save(ctx context.Context, entity *entities.LineProfile) error
	GetByID(ctx context.Context, id string) (*entities.LineProfile, error)
}
