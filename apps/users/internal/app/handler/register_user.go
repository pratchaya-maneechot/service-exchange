package handler

import (
	"context"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/command"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/role"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
)

type RegisterUserCommandHandler struct {
	userRepo     user.UserRepository
	roleCacheSvc *role.RoleCacheService
	// userService   user.UserService          // A domain service if you want to encapsulate complex domain logic
	// eventProducer message.UserEventProducer // Infrastructure layer dependency for publishing events
}

func NewRegisterUserCommandHandler(
	userRepo user.UserRepository,
	rcs *role.RoleCacheService,
) *RegisterUserCommandHandler {
	return &RegisterUserCommandHandler{
		userRepo:     userRepo,
		roleCacheSvc: rcs,
	}
}

func (h *RegisterUserCommandHandler) Handle(ctx context.Context, cmd command.RegisterUserCommand) (*user.User, error) {
	existing, err := h.userRepo.ExistsByLineUserID(ctx, cmd.LineUserID)
	if err != nil {
		return nil, err
	}
	if existing {
		return nil, user.ErrLineUserAlreadyExists
	}
	domUser, err := user.NewUser(ids.NewUserID(), cmd.LineUserID, cmd.Email, cmd.Password)
	if err != nil {
		return nil, err
	}
	defaultRole, err := h.roleCacheSvc.GetRoleByName(role.RoleNamePoster)
	if err != nil {
		return nil, err
	}
	if err = domUser.AddRole(defaultRole); err != nil {
		return nil, err
	}
	if err = h.userRepo.Save(ctx, domUser); err != nil {
		return nil, err
	}
	return domUser, nil
}
