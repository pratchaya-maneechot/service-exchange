package server

import (
	"context"
	"log/slog"

	taskpb "github.com/pratchaya-maneechot/service-exchange/apps/tasks/api/proto/task"
	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/pkg/bus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TaskServiceServer struct {
	taskpb.UnimplementedTaskServiceServer
	bus    *bus.Bus
	logger *slog.Logger
}

func (t *TaskServiceServer) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.Task, error) {
	panic("unimplemented")
}

func (t *TaskServiceServer) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

func (t *TaskServiceServer) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.Task, error) {
	panic("unimplemented")
}

func (t *TaskServiceServer) ListTasks(ctx context.Context, req *taskpb.ListTasksRequest) (*taskpb.ListTasksResponse, error) {
	panic("unimplemented")
}

func (t *TaskServiceServer) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.Task, error) {
	panic("unimplemented")
}

func Register(
	gs *grpc.Server,
	bus *bus.Bus,
	logger *slog.Logger,
) {
	server := &TaskServiceServer{
		bus:    bus,
		logger: logger,
	}
	taskpb.RegisterTaskServiceServer(gs, server)
}
