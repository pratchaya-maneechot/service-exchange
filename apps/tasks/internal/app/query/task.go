package query

import (
	"context"

	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/domain/shared/ids"
	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/domain/task"
)

type GetTaskDetail struct {
	TaskID string
}

type TaskQueryService struct {
	taskRepo task.TaskRepository
}

func NewTaskQueryService(repo task.TaskRepository) *TaskQueryService {
	return &TaskQueryService{
		taskRepo: repo,
	}
}

func (s *TaskQueryService) HandleGetTaskDetail(ctx context.Context, query GetTaskDetail) (*task.Task, error) {
	taskID := ids.TaskID(query.TaskID)
	t, err := s.taskRepo.Detail(ctx, taskID)
	if err != nil {
		return nil, err
	}
	return t, nil
}
