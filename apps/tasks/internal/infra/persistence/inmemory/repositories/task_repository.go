package repository

import (
	"context"
	"sync"

	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/domain/shared"
	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/domain/task"
)

type TaskRepository struct {
	tasks map[string]*task.Task
	mutex sync.RWMutex
}

// Delete implements task.TaskRepository.
func (r *TaskRepository) Delete(ctx context.Context, id shared.TaskID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	panic("unimplemented")
}

// Detail implements task.TaskRepository.
func (r *TaskRepository) Detail(ctx context.Context, id shared.TaskID) (*task.Task, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	panic("unimplemented")
}

// Save implements task.TaskRepository.
func (r *TaskRepository) Save(ctx context.Context, entity *task.Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	panic("unimplemented")
}

func NewTaskRepository() task.TaskRepository {
	return &TaskRepository{
		tasks: make(map[string]*task.Task),
		mutex: sync.RWMutex{},
	}
}
