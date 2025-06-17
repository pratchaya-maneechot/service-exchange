package task

import "github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/domain/shared/ids"

type Task struct {
	ID   ids.TaskID
	Name string
}

func NewTask(id ids.TaskID, name string) *Task {
	return &Task{
		ID:   id,
		Name: name,
	}
}
