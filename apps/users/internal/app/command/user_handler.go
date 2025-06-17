package command

import (

	// Or a structured logger

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
)

// userCommandHandler implements the UserCommandHandler interface.
type userCommandHandler struct {
	userRepo user.UserRepository
	// userService   user.UserService          // A domain service if you want to encapsulate complex domain logic
	// eventProducer message.UserEventProducer // Infrastructure layer dependency for publishing events
}

// HandleLineLogin implements UserCommandHandler.
func (h *userCommandHandler) HandleLineLogin(cmd LineLoginCommand) (ids.UserID, string, int, error) {
	panic("unimplemented")
}

// HandleRegisterUser implements UserCommandHandler.
func (h *userCommandHandler) HandleRegisterUser(cmd RegisterUserCommand) (ids.UserID, error) {
	panic("unimplemented")
}

// HandleSubmitIdentityVerification implements UserCommandHandler.
func (h *userCommandHandler) HandleSubmitIdentityVerification(cmd SubmitIdentityVerificationCommand) error {
	panic("unimplemented")
}

// HandleUpdateUserProfile implements UserCommandHandler.
func (h *userCommandHandler) HandleUpdateUserProfile(cmd UpdateUserProfileCommand) error {
	panic("unimplemented")
}

// HandleUserLogin implements UserCommandHandler.
func (h *userCommandHandler) HandleUserLogin(cmd UserLoginCommand) (string, int, error) {
	panic("unimplemented")
}

// NewUserCommandHandler creates a new UserCommandHandler.
func NewUserCommandHandler(
	userRepo user.UserRepository,
) UserCommandHandler {
	return &userCommandHandler{
		userRepo: userRepo,
	}
}
