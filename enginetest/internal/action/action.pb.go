// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.1
// source: github.com/dogmatiq/enginekit/enginetest/internal/action/action.proto

package action

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Action struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Behavior:
	//
	//	*Action_Fail
	//	*Action_Log
	//	*Action_ExecuteCommand
	//	*Action_RecordEvent
	//	*Action_ScheduleTimeout
	//	*Action_Destroy
	//	*Action_End
	Behavior isAction_Behavior `protobuf_oneof:"behavior"`
}

func (x *Action) Reset() {
	*x = Action{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Action) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Action) ProtoMessage() {}

func (x *Action) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Action.ProtoReflect.Descriptor instead.
func (*Action) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_rawDescGZIP(), []int{0}
}

func (m *Action) GetBehavior() isAction_Behavior {
	if m != nil {
		return m.Behavior
	}
	return nil
}

func (x *Action) GetFail() string {
	if x, ok := x.GetBehavior().(*Action_Fail); ok {
		return x.Fail
	}
	return ""
}

func (x *Action) GetLog() string {
	if x, ok := x.GetBehavior().(*Action_Log); ok {
		return x.Log
	}
	return ""
}

func (x *Action) GetExecuteCommand() *anypb.Any {
	if x, ok := x.GetBehavior().(*Action_ExecuteCommand); ok {
		return x.ExecuteCommand
	}
	return nil
}

func (x *Action) GetRecordEvent() *anypb.Any {
	if x, ok := x.GetBehavior().(*Action_RecordEvent); ok {
		return x.RecordEvent
	}
	return nil
}

func (x *Action) GetScheduleTimeout() *ScheduleTimeoutDetails {
	if x, ok := x.GetBehavior().(*Action_ScheduleTimeout); ok {
		return x.ScheduleTimeout
	}
	return nil
}

func (x *Action) GetDestroy() *Empty {
	if x, ok := x.GetBehavior().(*Action_Destroy); ok {
		return x.Destroy
	}
	return nil
}

func (x *Action) GetEnd() *Empty {
	if x, ok := x.GetBehavior().(*Action_End); ok {
		return x.End
	}
	return nil
}

type isAction_Behavior interface {
	isAction_Behavior()
}

type Action_Fail struct {
	Fail string `protobuf:"bytes,1,opt,name=fail,proto3,oneof"`
}

type Action_Log struct {
	Log string `protobuf:"bytes,2,opt,name=log,proto3,oneof"`
}

type Action_ExecuteCommand struct {
	ExecuteCommand *anypb.Any `protobuf:"bytes,3,opt,name=execute_command,json=executeCommand,proto3,oneof"`
}

type Action_RecordEvent struct {
	RecordEvent *anypb.Any `protobuf:"bytes,4,opt,name=record_event,json=recordEvent,proto3,oneof"`
}

type Action_ScheduleTimeout struct {
	ScheduleTimeout *ScheduleTimeoutDetails `protobuf:"bytes,5,opt,name=schedule_timeout,json=scheduleTimeout,proto3,oneof"`
}

type Action_Destroy struct {
	Destroy *Empty `protobuf:"bytes,6,opt,name=destroy,proto3,oneof"`
}

type Action_End struct {
	End *Empty `protobuf:"bytes,7,opt,name=end,proto3,oneof"`
}

func (*Action_Fail) isAction_Behavior() {}

func (*Action_Log) isAction_Behavior() {}

func (*Action_ExecuteCommand) isAction_Behavior() {}

func (*Action_RecordEvent) isAction_Behavior() {}

func (*Action_ScheduleTimeout) isAction_Behavior() {}

func (*Action_Destroy) isAction_Behavior() {}

func (*Action_End) isAction_Behavior() {}

type ScheduleTimeoutDetails struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timeout *anypb.Any             `protobuf:"bytes,1,opt,name=timeout,proto3" json:"timeout,omitempty"`
	At      *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=at,proto3" json:"at,omitempty"`
}

func (x *ScheduleTimeoutDetails) Reset() {
	*x = ScheduleTimeoutDetails{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScheduleTimeoutDetails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScheduleTimeoutDetails) ProtoMessage() {}

func (x *ScheduleTimeoutDetails) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScheduleTimeoutDetails.ProtoReflect.Descriptor instead.
func (*ScheduleTimeoutDetails) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_rawDescGZIP(), []int{1}
}

func (x *ScheduleTimeoutDetails) GetTimeout() *anypb.Any {
	if x != nil {
		return x.Timeout
	}
	return nil
}

func (x *ScheduleTimeoutDetails) GetAt() *timestamppb.Timestamp {
	if x != nil {
		return x.At
	}
	return nil
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_rawDescGZIP(), []int{2}
}

var File_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto protoreflect.FileDescriptor

var file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_rawDesc = []byte{
	0x0a, 0x45, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x67,
	0x6d, 0x61, 0x74, 0x69, 0x71, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x6b, 0x69, 0x74, 0x2f,
	0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x74, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1d, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x74, 0x69,
	0x71, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x6b, 0x69, 0x74, 0x2e, 0x65, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0x74, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x9a, 0x03, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a,
	0x04, 0x66, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x66,
	0x61, 0x69, 0x6c, 0x12, 0x12, 0x0a, 0x03, 0x6c, 0x6f, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x00, 0x52, 0x03, 0x6c, 0x6f, 0x67, 0x12, 0x3f, 0x0a, 0x0f, 0x65, 0x78, 0x65, 0x63, 0x75,
	0x74, 0x65, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x48, 0x00, 0x52, 0x0e, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74,
	0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x39, 0x0a, 0x0c, 0x72, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x41, 0x6e, 0x79, 0x48, 0x00, 0x52, 0x0b, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x12, 0x62, 0x0a, 0x10, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x5f,
	0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x35, 0x2e,
	0x64, 0x6f, 0x67, 0x6d, 0x61, 0x74, 0x69, 0x71, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x6b,
	0x69, 0x74, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x53, 0x63,
	0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x44, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x73, 0x48, 0x00, 0x52, 0x0f, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65,
	0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x40, 0x0a, 0x07, 0x64, 0x65, 0x73, 0x74, 0x72,
	0x6f, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61,
	0x74, 0x69, 0x71, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x6b, 0x69, 0x74, 0x2e, 0x65, 0x6e,
	0x67, 0x69, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x48, 0x00,
	0x52, 0x07, 0x64, 0x65, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x12, 0x38, 0x0a, 0x03, 0x65, 0x6e, 0x64,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x74, 0x69,
	0x71, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x6b, 0x69, 0x74, 0x2e, 0x65, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x48, 0x00, 0x52, 0x03,
	0x65, 0x6e, 0x64, 0x42, 0x0a, 0x0a, 0x08, 0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x22,
	0x74, 0x0a, 0x16, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x6f,
	0x75, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x2e, 0x0a, 0x07, 0x74, 0x69, 0x6d,
	0x65, 0x6f, 0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79,
	0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x2a, 0x0a, 0x02, 0x61, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x02, 0x61, 0x74, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x3a,
	0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x67,
	0x6d, 0x61, 0x74, 0x69, 0x71, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x6b, 0x69, 0x74, 0x2f,
	0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x74, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_rawDescOnce sync.Once
	file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_rawDescData = file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_rawDesc
)

func file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_rawDescGZIP() []byte {
	file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_rawDescOnce.Do(func() {
		file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_rawDescData)
	})
	return file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_rawDescData
}

var file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_goTypes = []interface{}{
	(*Action)(nil),                 // 0: dogmatiq.enginekit.enginetest.Action
	(*ScheduleTimeoutDetails)(nil), // 1: dogmatiq.enginekit.enginetest.ScheduleTimeoutDetails
	(*Empty)(nil),                  // 2: dogmatiq.enginekit.enginetest.Empty
	(*anypb.Any)(nil),              // 3: google.protobuf.Any
	(*timestamppb.Timestamp)(nil),  // 4: google.protobuf.Timestamp
}
var file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_depIdxs = []int32{
	3, // 0: dogmatiq.enginekit.enginetest.Action.execute_command:type_name -> google.protobuf.Any
	3, // 1: dogmatiq.enginekit.enginetest.Action.record_event:type_name -> google.protobuf.Any
	1, // 2: dogmatiq.enginekit.enginetest.Action.schedule_timeout:type_name -> dogmatiq.enginekit.enginetest.ScheduleTimeoutDetails
	2, // 3: dogmatiq.enginekit.enginetest.Action.destroy:type_name -> dogmatiq.enginekit.enginetest.Empty
	2, // 4: dogmatiq.enginekit.enginetest.Action.end:type_name -> dogmatiq.enginekit.enginetest.Empty
	3, // 5: dogmatiq.enginekit.enginetest.ScheduleTimeoutDetails.timeout:type_name -> google.protobuf.Any
	4, // 6: dogmatiq.enginekit.enginetest.ScheduleTimeoutDetails.at:type_name -> google.protobuf.Timestamp
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_init() }
func file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_init() {
	if File_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Action); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScheduleTimeoutDetails); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Action_Fail)(nil),
		(*Action_Log)(nil),
		(*Action_ExecuteCommand)(nil),
		(*Action_RecordEvent)(nil),
		(*Action_ScheduleTimeout)(nil),
		(*Action_Destroy)(nil),
		(*Action_End)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_goTypes,
		DependencyIndexes: file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_depIdxs,
		MessageInfos:      file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_msgTypes,
	}.Build()
	File_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto = out.File
	file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_rawDesc = nil
	file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_goTypes = nil
	file_github_com_dogmatiq_enginekit_enginetest_internal_action_action_proto_depIdxs = nil
}
