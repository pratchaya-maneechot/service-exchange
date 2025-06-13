package task

import "github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/domain/shared"

type Task struct {
	ID   shared.TaskID
	Name string
}

func NewTask(id shared.TaskID, name string) *Task {
	return &Task{
		ID:   id,
		Name: name,
	}
}
