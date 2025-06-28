/**
 * DeleteUserCommand - Command
 */

package command

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/infra/observability"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// DeleteUserCommand represents the input for the DeleteUser operation.
type DeleteUserCommand struct {
	// Add command-specific fields here, e.g.,
	// ID         ids.UserID `json:"id"`
	// Name       string                   `json:"name" validate:"required"`
	// Description string                   `json:"description,omitempty"`
}

// DeleteUserDto represents the output after the DeleteUser operation.
type DeleteUserDto struct {
	ID string `json:"id"` // Example: ID of the created/updated resource
}

// DeleteUserCommandHandler handles the DeleteUserCommand.
type DeleteUserCommandHandler struct {
	userRepo user.UserRepository
	logger   *slog.Logger
	config   *config.Config
	tracer   trace.Tracer
}

// NewDeleteUserCommandHandler creates a new instance of DeleteUserCommandHandler.
func NewDeleteUserCommandHandler(
	userRepo user.UserRepository,
	logger *slog.Logger,
	cfg *config.Config,
) *DeleteUserCommandHandler {
	return &DeleteUserCommandHandler{
		userRepo: userRepo,
		logger:   logger.With(slog.String("component", "DeleteUserCommandHandler")),
		config:   cfg,
		tracer:   otel.Tracer(fmt.Sprintf("%s.command-handler", cfg.Name)),
	}
}

// Handle executes the DeleteUserCommand.
func (h *DeleteUserCommandHandler) Handle(ctx context.Context, cmd DeleteUserCommand) (*DeleteUserDto, error) {
	logger := observability.GetLoggerFromContext(ctx).With(slog.String("command", "DeleteUser"))

	ctx, span := h.tracer.Start(ctx, "DeleteUserCommand.Handle", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	// Example: span.SetAttributes(attribute.String("user.name", cmd.Name))

	// Implement your business logic here.
	// For example:
	// 1. Validate command input.
	// 2. Check for existing resource if needed.
	// 3. Create/load domain entity.
	//    newUser, err := user.NewUser(
	//        ids.NewUserID(), // Generate a new ID for new entities
	//        cmd.Name,
	//        cmd.Description,
	//    )
	//    if err != nil {
	//        span.SetStatus(codes.Error, "Failed to create user domain model")
	//        span.RecordError(err)
	//        logger.Error("Failed to create domain user", slog.Any("error", err))
	//        return nil, err
	//    }
	// 4. Save the domain entity using the repository.
	//    if err := h.userRepo.Save(ctx, newUser); err != nil {
	//        span.SetStatus(codes.Error, "Failed to save user")
	//        span.RecordError(err)
	//        logger.Error("Failed to save user to repository", slog.Any("error", err))
	//        return nil, err
	//    }

	span.SetStatus(codes.Ok, "DeleteUser handled successfully")
	logger.Info("DeleteUser executed successfully.")

	return &DeleteUserDto{
		ID: "example-id", // Replace with actual ID if generated
	}, nil
}
