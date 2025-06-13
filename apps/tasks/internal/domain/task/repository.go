package task

import (
	"context"

	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/domain/shared"
)

type TaskRepository interface {
	Save(ctx context.Context, entity *Task) error
	Delete(ctx context.Context, id shared.TaskID) error
	Detail(ctx context.Context, id shared.TaskID) (*Task, error)
}
