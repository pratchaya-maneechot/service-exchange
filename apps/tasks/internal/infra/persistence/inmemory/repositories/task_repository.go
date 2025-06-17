package repository

import (
	"context"
	"sync"

	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/domain/shared/ids"
	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/domain/task"
)

type impl struct {
	tasks map[string]*task.Task
	mutex sync.RWMutex
}

func (r *impl) Delete(ctx context.Context, id ids.TaskID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	panic("unimplemented")
}

func (r *impl) Detail(ctx context.Context, id ids.TaskID) (*task.Task, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	panic("unimplemented")
}

func (r *impl) Save(ctx context.Context, entity *task.Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	panic("unimplemented")
}

func NewTaskRepository() task.TaskRepository {
	return &impl{
		tasks: make(map[string]*task.Task),
		mutex: sync.RWMutex{},
	}
}
