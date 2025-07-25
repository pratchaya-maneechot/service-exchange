// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: apps/tasks/api/proto/task/task.proto

package task

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Task represents a single task in the system.
type Task struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Task) Reset() {
	*x = Task{}
	mi := &file_apps_tasks_api_proto_task_task_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Task) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Task) ProtoMessage() {}

func (x *Task) ProtoReflect() protoreflect.Message {
	mi := &file_apps_tasks_api_proto_task_task_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Task.ProtoReflect.Descriptor instead.
func (*Task) Descriptor() ([]byte, []int) {
	return file_apps_tasks_api_proto_task_task_proto_rawDescGZIP(), []int{0}
}

func (x *Task) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Task) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Task) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Task) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

// Request message for creating a task.
type CreateTaskRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"` // Field 1 is more conventional for the primary field
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateTaskRequest) Reset() {
	*x = CreateTaskRequest{}
	mi := &file_apps_tasks_api_proto_task_task_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTaskRequest) ProtoMessage() {}

func (x *CreateTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_apps_tasks_api_proto_task_task_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTaskRequest.ProtoReflect.Descriptor instead.
func (*CreateTaskRequest) Descriptor() ([]byte, []int) {
	return file_apps_tasks_api_proto_task_task_proto_rawDescGZIP(), []int{1}
}

func (x *CreateTaskRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// Request message for retrieving a single task.
type GetTaskRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"` // The ID of the task to retrieve.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetTaskRequest) Reset() {
	*x = GetTaskRequest{}
	mi := &file_apps_tasks_api_proto_task_task_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTaskRequest) ProtoMessage() {}

func (x *GetTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_apps_tasks_api_proto_task_task_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTaskRequest.ProtoReflect.Descriptor instead.
func (*GetTaskRequest) Descriptor() ([]byte, []int) {
	return file_apps_tasks_api_proto_task_task_proto_rawDescGZIP(), []int{2}
}

func (x *GetTaskRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// Request message for listing tasks.
// This can be extended with pagination, filters, etc.
type ListTasksRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListTasksRequest) Reset() {
	*x = ListTasksRequest{}
	mi := &file_apps_tasks_api_proto_task_task_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListTasksRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListTasksRequest) ProtoMessage() {}

func (x *ListTasksRequest) ProtoReflect() protoreflect.Message {
	mi := &file_apps_tasks_api_proto_task_task_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListTasksRequest.ProtoReflect.Descriptor instead.
func (*ListTasksRequest) Descriptor() ([]byte, []int) {
	return file_apps_tasks_api_proto_task_task_proto_rawDescGZIP(), []int{3}
}

// Response message for listing tasks.
type ListTasksResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Tasks         []*Task                `protobuf:"bytes,1,rep,name=tasks,proto3" json:"tasks,omitempty"` // 'repeated' indicates a list of Task objects.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListTasksResponse) Reset() {
	*x = ListTasksResponse{}
	mi := &file_apps_tasks_api_proto_task_task_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListTasksResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListTasksResponse) ProtoMessage() {}

func (x *ListTasksResponse) ProtoReflect() protoreflect.Message {
	mi := &file_apps_tasks_api_proto_task_task_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListTasksResponse.ProtoReflect.Descriptor instead.
func (*ListTasksResponse) Descriptor() ([]byte, []int) {
	return file_apps_tasks_api_proto_task_task_proto_rawDescGZIP(), []int{4}
}

func (x *ListTasksResponse) GetTasks() []*Task {
	if x != nil {
		return x.Tasks
	}
	return nil
}

// Request message for updating a task.
type UpdateTaskRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`     // The ID of the task to update.
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"` // The new name for the task.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateTaskRequest) Reset() {
	*x = UpdateTaskRequest{}
	mi := &file_apps_tasks_api_proto_task_task_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateTaskRequest) ProtoMessage() {}

func (x *UpdateTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_apps_tasks_api_proto_task_task_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateTaskRequest.ProtoReflect.Descriptor instead.
func (*UpdateTaskRequest) Descriptor() ([]byte, []int) {
	return file_apps_tasks_api_proto_task_task_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateTaskRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateTaskRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// Request message for deleting a task.
type DeleteTaskRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"` // The ID of the task to delete.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteTaskRequest) Reset() {
	*x = DeleteTaskRequest{}
	mi := &file_apps_tasks_api_proto_task_task_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteTaskRequest) ProtoMessage() {}

func (x *DeleteTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_apps_tasks_api_proto_task_task_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteTaskRequest.ProtoReflect.Descriptor instead.
func (*DeleteTaskRequest) Descriptor() ([]byte, []int) {
	return file_apps_tasks_api_proto_task_task_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteTaskRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_apps_tasks_api_proto_task_task_proto protoreflect.FileDescriptor

const file_apps_tasks_api_proto_task_task_proto_rawDesc = "" +
	"\n" +
	"$apps/tasks/api/proto/task/task.proto\x12\x04task\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x1bgoogle/protobuf/empty.proto\"\xa0\x01\n" +
	"\x04Task\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x129\n" +
	"\n" +
	"created_at\x18\a \x01(\v2\x1a.google.protobuf.TimestampR\tcreatedAt\x129\n" +
	"\n" +
	"updated_at\x18\b \x01(\v2\x1a.google.protobuf.TimestampR\tupdatedAt\"'\n" +
	"\x11CreateTaskRequest\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\" \n" +
	"\x0eGetTaskRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\"\x12\n" +
	"\x10ListTasksRequest\"5\n" +
	"\x11ListTasksResponse\x12 \n" +
	"\x05tasks\x18\x01 \x03(\v2\n" +
	".task.TaskR\x05tasks\"7\n" +
	"\x11UpdateTaskRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\"#\n" +
	"\x11DeleteTaskRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id2\x9d\x02\n" +
	"\vTaskService\x121\n" +
	"\n" +
	"CreateTask\x12\x17.task.CreateTaskRequest\x1a\n" +
	".task.Task\x12+\n" +
	"\aGetTask\x12\x14.task.GetTaskRequest\x1a\n" +
	".task.Task\x12<\n" +
	"\tListTasks\x12\x16.task.ListTasksRequest\x1a\x17.task.ListTasksResponse\x121\n" +
	"\n" +
	"UpdateTask\x12\x17.task.UpdateTaskRequest\x1a\n" +
	".task.Task\x12=\n" +
	"\n" +
	"DeleteTask\x12\x17.task.DeleteTaskRequest\x1a\x16.google.protobuf.EmptyBKZIgithub.com/pratchaya-maneechot/service-exchange/apps/tasks/api/proto/taskb\x06proto3"

var (
	file_apps_tasks_api_proto_task_task_proto_rawDescOnce sync.Once
	file_apps_tasks_api_proto_task_task_proto_rawDescData []byte
)

func file_apps_tasks_api_proto_task_task_proto_rawDescGZIP() []byte {
	file_apps_tasks_api_proto_task_task_proto_rawDescOnce.Do(func() {
		file_apps_tasks_api_proto_task_task_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_apps_tasks_api_proto_task_task_proto_rawDesc), len(file_apps_tasks_api_proto_task_task_proto_rawDesc)))
	})
	return file_apps_tasks_api_proto_task_task_proto_rawDescData
}

var file_apps_tasks_api_proto_task_task_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_apps_tasks_api_proto_task_task_proto_goTypes = []any{
	(*Task)(nil),                  // 0: task.Task
	(*CreateTaskRequest)(nil),     // 1: task.CreateTaskRequest
	(*GetTaskRequest)(nil),        // 2: task.GetTaskRequest
	(*ListTasksRequest)(nil),      // 3: task.ListTasksRequest
	(*ListTasksResponse)(nil),     // 4: task.ListTasksResponse
	(*UpdateTaskRequest)(nil),     // 5: task.UpdateTaskRequest
	(*DeleteTaskRequest)(nil),     // 6: task.DeleteTaskRequest
	(*timestamppb.Timestamp)(nil), // 7: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),         // 8: google.protobuf.Empty
}
var file_apps_tasks_api_proto_task_task_proto_depIdxs = []int32{
	7, // 0: task.Task.created_at:type_name -> google.protobuf.Timestamp
	7, // 1: task.Task.updated_at:type_name -> google.protobuf.Timestamp
	0, // 2: task.ListTasksResponse.tasks:type_name -> task.Task
	1, // 3: task.TaskService.CreateTask:input_type -> task.CreateTaskRequest
	2, // 4: task.TaskService.GetTask:input_type -> task.GetTaskRequest
	3, // 5: task.TaskService.ListTasks:input_type -> task.ListTasksRequest
	5, // 6: task.TaskService.UpdateTask:input_type -> task.UpdateTaskRequest
	6, // 7: task.TaskService.DeleteTask:input_type -> task.DeleteTaskRequest
	0, // 8: task.TaskService.CreateTask:output_type -> task.Task
	0, // 9: task.TaskService.GetTask:output_type -> task.Task
	4, // 10: task.TaskService.ListTasks:output_type -> task.ListTasksResponse
	0, // 11: task.TaskService.UpdateTask:output_type -> task.Task
	8, // 12: task.TaskService.DeleteTask:output_type -> google.protobuf.Empty
	8, // [8:13] is the sub-list for method output_type
	3, // [3:8] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_apps_tasks_api_proto_task_task_proto_init() }
func file_apps_tasks_api_proto_task_task_proto_init() {
	if File_apps_tasks_api_proto_task_task_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_apps_tasks_api_proto_task_task_proto_rawDesc), len(file_apps_tasks_api_proto_task_task_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_apps_tasks_api_proto_task_task_proto_goTypes,
		DependencyIndexes: file_apps_tasks_api_proto_task_task_proto_depIdxs,
		MessageInfos:      file_apps_tasks_api_proto_task_task_proto_msgTypes,
	}.Build()
	File_apps_tasks_api_proto_task_task_proto = out.File
	file_apps_tasks_api_proto_task_task_proto_goTypes = nil
	file_apps_tasks_api_proto_task_task_proto_depIdxs = nil
}
