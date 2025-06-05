package domain

import "github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/repositories"

type Repository struct {
	User        repositories.UserRepository
	LineProfile repositories.LineProfileRepository
}

func NewRepository(
	user repositories.UserRepository,
	lineProfile repositories.LineProfileRepository,
) *Repository {
	return &Repository{
		User:        user,
		LineProfile: lineProfile,
	}
}
