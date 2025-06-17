package command

import (
	"context"

	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/domain/shared/ids"
	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/domain/task"
)

type CreateTaskCommand struct {
	Name string
}

type UpdateTaskCommand struct {
	TaskID string
	Name   string
}

type TaskCommandService struct {
	taskRepo task.TaskRepository
}

func NewTaskCommandService(repo task.TaskRepository) *TaskCommandService {
	return &TaskCommandService{
		taskRepo: repo,
	}
}

func (s *TaskCommandService) HandleCreateTaskCommand(ctx context.Context, cmd CreateTaskCommand) (ids.TaskID, error) {
	newTaskID := ids.NewTaskID()
	p := task.NewTask(newTaskID, cmd.Name)
	if err := s.taskRepo.Save(ctx, p); err != nil {
		return "", err
	}
	return p.ID, nil
}
