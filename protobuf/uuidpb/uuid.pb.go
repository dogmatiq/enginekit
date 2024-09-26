// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: github.com/dogmatiq/enginekit/protobuf/uuidpb/uuid.proto

package uuidpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// UUID is a size-optimized representation of an RFC 9562 UUID.
type UUID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Upper uint64 `protobuf:"varint,1,opt,name=upper,proto3" json:"upper,omitempty"`
	Lower uint64 `protobuf:"varint,2,opt,name=lower,proto3" json:"lower,omitempty"`
}

func (x *UUID) Reset() {
	*x = UUID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UUID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UUID) ProtoMessage() {}

func (x *UUID) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UUID.ProtoReflect.Descriptor instead.
func (*UUID) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_rawDescGZIP(), []int{0}
}

func (x *UUID) GetUpper() uint64 {
	if x != nil {
		return x.Upper
	}
	return 0
}

func (x *UUID) GetLower() uint64 {
	if x != nil {
		return x.Lower
	}
	return 0
}

var File_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto protoreflect.FileDescriptor

var file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_rawDesc = []byte{
	0x0a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x67,
	0x6d, 0x61, 0x74, 0x69, 0x71, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x6b, 0x69, 0x74, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x75, 0x75, 0x69, 0x64, 0x70, 0x62, 0x2f,
	0x75, 0x75, 0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x64, 0x6f, 0x67, 0x6d,
	0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x22, 0x32, 0x0a, 0x04, 0x55, 0x55,
	0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x75, 0x70, 0x70, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x05, 0x75, 0x70, 0x70, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x6f, 0x77, 0x65,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x42, 0x2f,
	0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x67,
	0x6d, 0x61, 0x74, 0x69, 0x71, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x6b, 0x69, 0x74, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x75, 0x75, 0x69, 0x64, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_rawDescOnce sync.Once
	file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_rawDescData = file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_rawDesc
)

func file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_rawDescGZIP() []byte {
	file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_rawDescOnce.Do(func() {
		file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_rawDescData)
	})
	return file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_rawDescData
}

var file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_goTypes = []any{
	(*UUID)(nil), // 0: dogma.protobuf.UUID
}
var file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_init() }
func file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_init() {
	if File_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*UUID); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_goTypes,
		DependencyIndexes: file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_depIdxs,
		MessageInfos:      file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_msgTypes,
	}.Build()
	File_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto = out.File
	file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_rawDesc = nil
	file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_goTypes = nil
	file_github_com_dogmatiq_enginekit_protobuf_uuidpb_uuid_proto_depIdxs = nil
}
