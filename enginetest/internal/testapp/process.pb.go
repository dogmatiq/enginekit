// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.1
// source: github.com/dogmatiq/enginekit/enginetest/internal/testapp/process.proto

package testapp

import (
	action "github.com/dogmatiq/enginekit/enginetest/internal/action"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

type ProcessEventA struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	InstanceId    string                 `protobuf:"bytes,1,opt,name=instance_id,json=instanceId,proto3" json:"instance_id,omitempty"`
	Actions       []*action.Action       `protobuf:"bytes,2,rep,name=actions,proto3" json:"actions,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProcessEventA) Reset() {
	*x = ProcessEventA{}
	mi := &file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProcessEventA) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessEventA) ProtoMessage() {}

func (x *ProcessEventA) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessEventA.ProtoReflect.Descriptor instead.
func (*ProcessEventA) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_rawDescGZIP(), []int{0}
}

func (x *ProcessEventA) GetInstanceId() string {
	if x != nil {
		return x.InstanceId
	}
	return ""
}

func (x *ProcessEventA) GetActions() []*action.Action {
	if x != nil {
		return x.Actions
	}
	return nil
}

var File_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto protoreflect.FileDescriptor

const file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_rawDesc = "" +
	"\n" +
	"Ggithub.com/dogmatiq/enginekit/enginetest/internal/testapp/process.proto\x12\x1ddogmatiq.enginekit.enginetest\x1aEgithub.com/dogmatiq/enginekit/enginetest/internal/action/action.proto\"q\n" +
	"\rProcessEventA\x12\x1f\n" +
	"\vinstance_id\x18\x01 \x01(\tR\n" +
	"instanceId\x12?\n" +
	"\aactions\x18\x02 \x03(\v2%.dogmatiq.enginekit.enginetest.ActionR\aactionsB;Z9github.com/dogmatiq/enginekit/enginetest/internal/testappb\x06proto3"

var (
	file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_rawDescOnce sync.Once
	file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_rawDescData []byte
)

func file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_rawDescGZIP() []byte {
	file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_rawDescOnce.Do(func() {
		file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_rawDesc), len(file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_rawDesc)))
	})
	return file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_rawDescData
}

var file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_goTypes = []any{
	(*ProcessEventA)(nil), // 0: dogmatiq.enginekit.enginetest.ProcessEventA
	(*action.Action)(nil), // 1: dogmatiq.enginekit.enginetest.Action
}
var file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_depIdxs = []int32{
	1, // 0: dogmatiq.enginekit.enginetest.ProcessEventA.actions:type_name -> dogmatiq.enginekit.enginetest.Action
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_init() }
func file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_init() {
	if File_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_rawDesc), len(file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_goTypes,
		DependencyIndexes: file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_depIdxs,
		MessageInfos:      file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_msgTypes,
	}.Build()
	File_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto = out.File
	file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_goTypes = nil
	file_github_com_dogmatiq_enginekit_enginetest_internal_testapp_process_proto_depIdxs = nil
}
