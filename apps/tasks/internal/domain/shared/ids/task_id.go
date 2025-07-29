package ids

import "github.com/pratchaya-maneechot/service-exchange/apps/tasks/pkg/util"

type TaskID string

func NewTaskID() TaskID {
	return TaskID(util.UUID())
}
